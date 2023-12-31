package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	connectgoexample "connect-go-example"
	greetv1 "connect-go-example/gen/greet/v1"        // generated by protoc-gen-go
	"connect-go-example/gen/greet/v1/greetv1connect" // generated by protoc-gen-connect-go
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	stream *connect.ClientStream[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", stream.RequestHeader())
	var greeting strings.Builder
	for stream.Receive() {
		log.Println("Receive: ", stream.Msg())
		g := fmt.Sprintf("Hello, %s!\n", stream.Msg().Name)
		if _, err := greeting.WriteString(g); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: greeting.String(),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	interceptors := connect.WithInterceptors(connectgoexample.NewAuthInterceptor())

	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter, interceptors)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
