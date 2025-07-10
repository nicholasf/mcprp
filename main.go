// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var httpAddr = flag.String("http", ":8080", "use SSE HTTP at this address")

type SayHiParams struct {
	Name string `json:"name"`
}

func SayHi(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[SayHiParams]) (*mcp.CallToolResultFor[any], error) {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Hi " + params.Arguments.Name},
		},
	}, nil
}

func main() {
	server1 := mcp.NewServer("greeter1", "v0.0.1", nil)
	mcp.AddTool(server1, &mcp.Tool{Name: "greet1", Description: "say hi"}, SayHi)

	server2 := mcp.NewServer("greeter2", "v0.0.1", nil)
	mcp.AddTool(server2, &mcp.Tool{Name: "greet2", Description: "say hello"}, SayHi)

	log.Printf("MCP servers serving at %s", ":8080")
	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		url := request.URL.Path
		log.Printf("Handling request for URL %s\n", url)
		switch url {
		case "/greeter1":
			return server1
		case "/greeter2":
			return server2
		default:
			return nil
		}
	})
	http.ListenAndServe(*httpAddr, handler)
}
