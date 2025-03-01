gen-cal:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       calculatorpb/calculator.proto
run-server:
	go run server/server.go
run-client:
	go run client/client.go
run-ssl:
	./ssl/SSLCmd.sh