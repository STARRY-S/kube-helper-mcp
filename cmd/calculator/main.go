package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

var (
	sse  bool
	bind string
	port int
)

func init() {
	flag.BoolVar(&sse, "sse", false, "Use SSE for streaming output")
	flag.StringVar(&bind, "bind", "0.0.0.0", "Bind address")
	flag.IntVar(&port, "port", 8188, "Bind port")
	flag.Parse()
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Calculator",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add a calculator tool
	calculatorTool := mcp.NewTool(
		"calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	// Add tool handler
	s.AddTool(calculatorTool, calculatorHandler)

	// Start the stdio server
	var err error
	if sse {
		listen := fmt.Sprintf("%v:%v", bind, port)
		u := fmt.Sprintf("http://%v", listen)
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(u),
		)
		logrus.Infof("Listen on %q", u)
		err = sseServer.Start(listen)
	} else {
		err = server.ServeStdio(s)
	}
	if err != nil {
		logrus.Fatalf("%v", err)
	}
}

func calculatorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, ok := request.Params.Arguments["operation"].(string)
	if !ok {
		return nil, errors.New("operation not provided")
	}
	x := request.Params.Arguments["x"].(float64)
	y := request.Params.Arguments["y"].(float64)

	var result float64
	switch op {
	case "add":
		result = x + y
	case "subtract":
		result = x - y
	case "multiply":
		result = x * y
	case "divide":
		if y == 0 {
			return nil, errors.New("cannot divide by zero")
		}
		result = x / y
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.3f", result)), nil
}
