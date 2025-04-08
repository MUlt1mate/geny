package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var testInput = `
protoc -I=. -I=./vendor --go_out=. --go_opt=paths=source_relative proto/somepackage/somepackage.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --go_out=. --go_opt=paths=source_relative proto/example.proto proto/example2.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=paths=source_relative,marshaller=easyjson proto/example.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=only=client,paths=source_relative proto/example2.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=paths=source_relative,autoURI=true proto/no_url.proto
easyjson -all proto/example.pb.go
`

func main() {
	var (
		batch *CommandBatch
		err   error
	)
	g := &geny{}
	if batch, err = g.ParseText(testInput); err != nil {
		log.Fatal(err)
	}

	output, _ := yaml.Marshal(batch)

	if batch, err = g.ParseYAML(output); err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("output.yaml", output, 0666)
}
