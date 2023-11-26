package main

import (
	connectgoexample "connect-go-example"
	"context"
	"log"
	"net/http"

	greetv1 "connect-go-example/gen/greet/v1"
	"connect-go-example/gen/greet/v1/greetv1connect"

	"connectrpc.com/connect"
)

func main() {
	interceptors := connect.WithInterceptors(connectgoexample.NewAuthInterceptor())

	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		interceptors,
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}
