package main

import (
	"context"
	"log/slog"
	"net"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/stretchr/testify/require"
)

func TestCall(t *testing.T) {
	ctx := context.Background()
	connA, connB := net.Pipe()

	sender := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(connA), nil)
	defer sender.Close()

	s := NewMCPServer()
	receiver := jsonrpc2.NewConn(ctx, jsonrpc2.NewPlainObjectStream(connB), s)
	defer receiver.Close()

	t.Run("Testing a valid call", func(t *testing.T) {
		var result string
		if err := sender.Call(ctx, "sayHello", nil, &result); err != nil {
			slog.Error(err.Error())
			t.Fail()
		}

		slog.Info(result)
		require.NotNil(t, result)
		require.Equal(t, result, "hello world")
	})

	t.Run("Testing a 404", func(t *testing.T) {
		var result string
		if err := sender.Call(ctx, "sayHello", nil, &result); err != nil {
			slog.Error(err.Error())
			t.Fail()
		}

		slog.Info(result)
		require.NotNil(t, result)
		require.Equal(t, result, "hello world")
	})
}
