fmt:
	go fmt ./...
lint:
	golangci-lint run
protoc:
	protoc -I protos/proto protos/proto/* --go_out=./protos/gen --go_opt=paths=source_relative --go-grpc_out=./protos/gen --go-grpc_opt=paths=source_relative