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
protoc -I=. -I=./vendor --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --validate_out=lang=go,paths=source_relative:. file.proto
protoc -I=. -I=./vendor --openapi_out=output_mode=source_relative:. --httpgo_out=. --httpgo_opt=paths=source_relative,only=server,context=native file.proto
protoc -I=. -I=./vendor --go_out=paths=source_relative:. --amqprpc4_out=. file.proto
find proto/. -type f -name "*.pb.go" -exec sed --in-place s/,omitempty// {} ;
find proto/. -type f -name "*.pb.go" -exec easyjson -all {} ;
gofmt -w ./proto/
// mock
mockgen -package=mocks -destination=internal/pkg/database/mock/pid_mock.go -source=internal/pkg/database/pid.go
mockgen -package=mocks -destination=internal/pkg/database/mock/tx_mock.go -source=internal/pkg/database/tx.go
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
	_ = os.WriteFile("geny.yaml", output, 0666)

	if batch, err = g.ParseYAML(output); err != nil {
		log.Fatal(err)
	}
	//if err = g.Exec(batch); err != nil {
	//	log.Fatal(err)
	//}

	text := g.FormatGoFile(batch)
	_ = os.WriteFile("generate.go", []byte(text), 0666)
}
