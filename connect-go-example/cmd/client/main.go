package main

import (
	"context"
	"log"
	"net/http"

	connectgoexample "connect-go-example"
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
	stream := client.Greet(
		context.Background(),
	)
	if err := stream.Send(&greetv1.GreetRequest{Name: "Jane"}); err != nil {
		log.Println(err)
		return
	}

	if err := stream.Send(&greetv1.GreetRequest{Name: "Jack"}); err != nil {
		log.Println(err)
		return
	}

	res, err := stream.CloseAndReceive()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res.Msg.Greeting)
}
