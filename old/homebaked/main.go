package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	ctx := context.Background()
	r := gin.Default()
	connA, connB := net.Pipe()

	sender := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(connA), nil)
	defer sender.Close()

	s := NewMCPServer()
	receiver := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(connB), s)
	defer receiver.Close()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// r.POST("/oauth", oauth)
	// r.GET("/.well-known/oauth-protected-resource", oauth)
	// r.GET("/.well-known/oauth-authorization-server/mcp", oauth)
	r.GET("/authorize", oauth)
	r.POST("/authorize", oauth)
	r.GET("/register", oauth)
	r.POST("/register", oauth)
	r.GET("/token", oauth)
	r.POST("/token", oauth)

	// GET and POST for mcp as per https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#streamable-http
	r.GET("/mcp/*jsonRPCRequest", func(c *gin.Context) {
		slog.Info("Received request", "request", c.Request)
		data := c.Param("jsonRPCRequest")

		var result string
		if err := sender.Call(ctx, data, nil, &result); err != nil {
			slog.Info("Error: ", "err", result)
			c.Error(err)
		}

		slog.Info("Generated MCP Response", "mcpresponse", result)

		c.Header("Content-Type", "application/json")
		c.JSON(202, result)
	})

	r.POST("/mcp/*jsonRPCRequest", func(c *gin.Context) {
		slog.Info("Received request", "request", c.Request)
		data := c.Param("jsonRPCRequest")

		var result string
		if err := sender.Call(ctx, data, nil, &result); err != nil {
			slog.Info("Error: ", "err", result)
			c.Error(err)
		}

		slog.Info("Generated MCP Response", "mcpresponse", result)

		c.Header("Content-Type", "application/json")
		c.JSON(202, result)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func oauth(c *gin.Context) {
	code := c.Query("code")
	slog.Info("Oauth called!", "code", code)

	c.JSON(200, code)
}

// basic session handler, returns a session id
// func initializer() string {
// 	return uuid.NewString()
// }
