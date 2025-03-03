run-error:
	go get golang.org/x/tools/go/packages
	go get golang.org/x/tools/go/ast/astutil
	go get golang.org/x/tools/imports
	go get github.com/urfave/cli/v2
run-graphql:
	go run github.com/99designs/gqlgen generate
run-server:
	go run server/server.go
proto-auth:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       proto/auth/auth.proto
proto-equipment:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       proto/equipment/equipment.proto
service-auth:
	go run service/auth/serviceAuth.go