package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
)

func main() {

	name := flag.String("name", "", "name of the context")
	batch := flag.Bool("batch", false, "use batch+finalize in addition to modify")
	readLock := flag.String("readLock", "parallel", "lock to protect model: exclusive, parallel, wait-free")
	flag.Parse()

	if *name == "" {
		log.Fatalf("name can not be empty")
	}

	if *readLock != "exclusive" && *readLock != "parallel" && *readLock != "wait-free" {
		log.Fatalf("readLock '%s' unknown", *readLock)
	}

	messageNames, ok := findMessageNames(*name)
	if !ok {
		log.Fatalf("topic wrap message not defined in proto file")
	}

	renderTemplate(os.Stdout, *name, *batch, *readLock, messageNames)
}

func findMessageNames(name string) ([]string, bool) {
	reader, err := os.Open("pkg/" + name + "/" + name + ".proto")
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
	funcMap := template.FuncMap{
		"title": strings.Title,
	}

	tplString, err := ioutil.ReadFile("cmd/dev-tools/simba/context.go.tpl")
	if err != nil {
		log.Fatalf("failed to read template: %v", err)
	}

	tpl, err := template.New("context").Funcs(funcMap).Parse(string(tplString))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

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

	err = tpl.Execute(w, data)
	if err != nil {
		log.Fatalf("template execution failed: %v", err)
	}
}
