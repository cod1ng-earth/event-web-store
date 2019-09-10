//go:generate go-bindata -pkg simba -nometadata -nocompress context.go.tpl

package simba

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/emicklei/proto"
)

func UpdateCode(batch bool, readLock string, bridges []string) {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := path.Base(currentDir)

	desc := createContextDescription(name, batch, readLock, bridges)
	f, err := os.Create("context.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	renderTemplate(f, desc)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func createContextDescription(name string, batch bool, readLock string, bridges []string) contextDescription {

	desc := contextDescription{
		Name:     name,
		Batch:    batch,
		ReadLock: readLock,
	}

	messageNames, ok := findMessageNames(filepath.Join("..", name, "topic.proto"))
	if !ok {
		log.Fatalf("topic wrap message not defined in proto file")
	}
	desc.MessageNames = messageNames

	publicProtoPath := filepath.Join("public", "topic.proto")
	if fileExists(publicProtoPath) {
		publicMessageNames, ok := findMessageNames(publicProtoPath)
		if ok {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current dir: %v", err)
			}
			publicPkgPath := trimFirstRune(strings.TrimPrefix(dir, filepath.Join(os.Getenv("GOPATH"), "src")))
			desc.Publisher = publisher{
				MessageNames: publicMessageNames,
				PkgPath:      publicPkgPath,
			}
		}
	}

	subscribers := []subscriber{}
	for _, bridge := range bridges {
		name = path.Base(bridge)
		messageNames, ok := findMessageNames(filepath.Join("..", name, "public", "topic.proto"))
		if !ok {
			log.Fatalf("topic wrap message not defined in proto file for bridge %v", bridge)
		}

		subscribers = append(subscribers, subscriber{
			Name:         name,
			PkgPath:      bridge,
			MessageNames: messageNames,
		})
	}
	desc.Bridges = subscribers

	return desc
}

func findMessageNames(protoPath string) ([]string, bool) {
	reader, err := os.Open(protoPath)
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
		if m.Name != "messages" {
			return
		}

		parent := m.Parent.(*proto.Message)
		if parent.Name != "TopicMessage" {
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
