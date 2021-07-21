.PHONY : gen

gen:
	$(info gen go files...)
	@go generate ./...

	$(info gen protobuf files...)
	@protoc \
	  --proto_path=./f1/pb \
	  --go_opt=paths=source_relative \
	  --go_out=./f1/pb \
	  --go-grpc_out=./f1/pb \
	  --go-grpc_opt=paths=source_relative f1.proto

	$(info gen OK)
