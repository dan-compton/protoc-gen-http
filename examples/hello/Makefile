all:
	protoc -I ../googleapis -I . hello.proto --go_out=plugins=grpc:. --http_out=. --grpc-gateway_out=.