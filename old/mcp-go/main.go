package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("StreamableHTTP API Server", "1.0.0",
		server.WithInstructions("Testing 1 2 3"),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true))

	// Add RESTful tools
	s.AddTool(
		mcp.NewTool("get_user",text.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[SayHiParams]) (*mcp.CallToolResultFor[any], error) {
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
			mcp.WithDescription("Get user information"),
			mcp.WithString("user_id", mcp.Required()),
		),
		handleGetUser,
	)

	s.AddTool(
		mcp.NewTool("create_user",
			mcp.WithDescription("Create a new user"),
			mcp.WithString("name", mcp.Required()),
			mcp.WithString("email", mcp.Required()),
			mcp.WithNumber("age", mcp.Min(0)),
		),
		handleCreateUser,
	)

	s.AddTool(
		mcp.NewTool("search_users",
			mcp.WithDescription("Search users with filters"),
			mcp.WithString("query", mcp.Description("Search query")),
			mcp.WithNumber("limit", mcp.DefaultNumber(10), mcp.Max(100)),
			mcp.WithNumber("offset", mcp.DefaultNumber(0), mcp.Min(0)),
		),
		handleSearchUsers,
	)

	// Add resources
	s.AddResource(
		mcp.NewResource(
			"users://{user_id}",
			"User Profile",
			mcp.WithResourceDescription("User profile data"),
			mcp.WithMIMEType("application/json"),
		),
		handleUserResource,
	)

	// Start StreamableHTTP server
	log.Println("Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func handleGetUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userID := req.GetString("user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Simulate database lookup
	user, err := getUserFromDB(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"id":"%s","name":"%s","email":"%s","age":%d}`,
		user.ID, user.Name, user.Email, user.Age)), nil
}

func handleCreateUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := req.GetString("name", "")
	email := req.GetString("email", "")
	age := req.GetInt("age", 0)

	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}

	// Validate input
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	// Create user
	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: time.Now(),
	}

	if err := saveUserToDB(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"id":"%s","message":"User created successfully","user":{"id":"%s","name":"%s","email":"%s","age":%d}}`,
		user.ID, user.ID, user.Name, user.Email, user.Age)), nil
}

// Helper functions and types for the examples
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

func getUserFromDB(userID string) (*User, error) {
	// Placeholder implementation
	return &User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func generateID() string {
	// Placeholder implementation
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func saveUserToDB(user *User) error {
	// Placeholder implementation
	return nil
}

func handleSearchUsers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")
	limit := req.GetInt("limit", 10)
	offset := req.GetInt("offset", 0)

	// Search users with pagination
	_, total, err := searchUsersInDB(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"users":[{"id":"1","name":"John Doe","email":"john@example.com","age":30},{"id":"2","name":"Jane Smith","email":"jane@example.com","age":25}],"total":%d,"limit":%d,"offset":%d,"query":"%s"}`,
		total, limit, offset, query)), nil
}

func handleUserResource(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	userID := extractUserIDFromURI(req.Params.URI)

	user, err := getUserFromDB(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     fmt.Sprintf(`{"id":"%s","name":"%s","email":"%s","age":%d}`, user.ID, user.Name, user.Email, user.Age),
		},
	}, nil
}

// Additional helper functions

func searchUsersInDB(query string, limit, offset int) ([]*User, int, error) {
	// Placeholder implementation
	users := []*User{
		{ID: "1", Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: "2", Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	}
	return users, len(users), nil
}

func extractUserIDFromURI(uri string) string {
	// Extract user ID from URI like "users://123"
	if len(uri) > 8 && uri[:8] == "users://" {
		return uri[8:]
	}
	return uri
}
