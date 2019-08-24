//go:generate go-bindata context.go.tpl

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
)

func main() {

	batch := flag.Bool("batch", false, "use batch+finalize in addition to modify")
	readLock := flag.String("readLock", "parallel", "lock to protect model: exclusive, parallel, wait-free")
	flag.Parse()

	if *readLock != "exclusive" && *readLock != "parallel" && *readLock != "wait-free" {
		log.Fatalf("readLock '%s' unknown", *readLock)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := path.Base(currentDir)

	messageNames, ok := findMessageNames(name)
	if !ok {
		log.Fatalf("topic wrap message not defined in proto file")
	}

	renderTemplate(os.Stdout, name, *batch, *readLock, messageNames)
}

func findMessageNames(name string) ([]string, bool) {
	reader, err := os.Open(name + ".proto")
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

func renderTemplate(w io.Writer, name string, batch bool, lock string, messageNames []string) {
	data := struct {
		Name         string
		Batch        bool
		Lock         string
		MessageNames []string
	}{
		Name:         name,
		Batch:        batch,
		Lock:         lock,
		MessageNames: messageNames,
	}

	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	tpl, err := template.New("context").Funcs(funcMap).Parse(string(MustAsset("context.go.tpl")))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Fatalf("template execution failed: %v", err)
	}
}
