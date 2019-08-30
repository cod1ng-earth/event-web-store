//go:generate go-bindata context.go.tpl

package main

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
	"gopkg.in/alecthomas/kingpin.v2"
)

type contextDescription struct {
	Name         string
	Batch        bool
	ReadLock     string
	MessageNames []string
	Bridges      []bridgeDescription
}

type bridgeDescription struct {
	Name         string
	PkgPath      string
	MessageNames []string
}

func main() {

	batch := kingpin.Flag("batch", "use batch+finalize in addition to modify").Bool()
	readLock := kingpin.Flag("lock", "lock to protect model: exclusive, parallel, wait-free").String()
	bridges := kingpin.Flag("bridge", "connections to other contexts").Strings()
	kingpin.Parse()

	if *readLock == "" {
		*readLock = "parallel"
	}
	if *readLock != "exclusive" && *readLock != "parallel" && *readLock != "wait-free" {
		log.Fatalf("readLock '%s' unknown", *readLock)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := path.Base(currentDir)

	desc := createContextDescription(name, *batch, *readLock, *bridges)
	f, err := os.Create("context.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	renderTemplate(f, desc)
}

func createContextDescription(name string, batch bool, readLock string, bridges []string) contextDescription {

	desc := contextDescription{
		Name:     name,
		Batch:    batch,
		ReadLock: readLock,
	}

	messageNames, ok := findMessageNames(name)
	if !ok {
		log.Fatalf("topic wrap message not defined in proto file")
	}
	desc.MessageNames = messageNames

	bridgeDescs := []bridgeDescription{}
	for _, bridge := range bridges {
		name = path.Base(bridge)
		messageNames, ok := findMessageNames(name)
		if !ok {
			log.Fatalf("topic wrap message not defined in proto file for bridge %v", bridge)
		}

		bridgeDescs = append(bridgeDescs, bridgeDescription{
			Name:         name,
			PkgPath:      bridge,
			MessageNames: messageNames,
		})
	}
	desc.Bridges = bridgeDescs

	return desc
}

func findMessageNames(name string) ([]string, bool) {
	reader, err := os.Open(filepath.Join("..", name, "topic.proto"))
	if err != nil {
		log.Printf("failed to read proto file: %s", err)
		os.Exit(1)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Printf("failed to parse proto definition: %s", err)
		os.Exit(1)
	}

	names := []string{}
	proto.Walk(definition, proto.WithOneof(func(m *proto.Oneof) {
		if m.Name != name+"Message" {
			return
		}

		parent := m.Parent.(*proto.Message)
		if parent.Name != strings.Title(name)+"Messages" {
			return
		}

		for _, child := range m.Elements {
			oneOf := child.(*proto.OneOfField)
			names = append(names, oneOf.Name)
		}
	}))

	return names, len(names) != 0
}

func renderTemplate(w io.Writer, desc contextDescription) {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	tpl, err := template.New("context").Funcs(funcMap).Parse(string(MustAsset("context.go.tpl")))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	err = tpl.Execute(w, desc)
	if err != nil {
		log.Fatalf("template execution failed: %v", err)
	}
}
