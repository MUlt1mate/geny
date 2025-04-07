package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var input = `
protoc -I=. -I=./vendor --go_out=. --go_opt=paths=source_relative proto/somepackage/somepackage.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --go_out=. --go_opt=paths=source_relative proto/example.proto proto/example2.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=paths=source_relative,marshaller=easyjson proto/example.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=only=client,paths=source_relative proto/example2.proto
protoc -I=. -I=./vendor -I=/usr/local/include -I=./proto --httpgo_out=. --httpgo_opt=paths=source_relative,autoURI=true proto/no_url.proto
easyjson -all proto/example.pb.go
`

func main() {
	var batch = &CommandBatch{}
	for line := range strings.Lines(input) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		command := &SimpleCommand{}
		_ = command.Parse(line)
		log.Println(command.String())
		batch.Simple = append(batch.Simple, command)
	}
	output, _ := yaml.Marshal(batch)
	_ = os.WriteFile("output.yaml", output, 0666)
}
