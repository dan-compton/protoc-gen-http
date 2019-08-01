package main

import (
	"context"
	"fmt"
	"github.com/flowup/protoc-gen-http/examples/hello"
)

func main() {
	var c hello.GreeterClient = hello.NewGreeterHttpClient("http://localhost:8080")

	out, err := c.SayHello(context.Background(), &hello.HelloRequest{
		Name: "Peto",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
