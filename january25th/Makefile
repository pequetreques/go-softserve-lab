create:
	protoc --proto_path=proto proto/*.proto --go_out=generated/
	protoc --proto_path=proto proto/*.proto --go-grpc_out=generated/

clean:
	rm generated/proto/*.go
