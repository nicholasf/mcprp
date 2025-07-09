This project is mostly for experimenting with an MCP integration with Claude, the final goal being a RAG like application modelling DnD campaigns that allows Claude to read and write campaign information into a flat file and/or SQL format.

mcp-go: https://github.com/mark3labs/mcp-go?tab=readme-ov-file


The `capabilities` property in the MCP initialize response declares what features and functionality the server supports. This allows clients to adapt their behavior based on what's available.

## Structure
The capabilities object contains nested objects for different feature categories:



Yes, there are standard endpoints that MCP (Model Context Protocol) HTTP servers typically implement. The core endpoints are:

**Required endpoints:**
- `POST /` - The main endpoint that handles all MCP JSON-RPC requests
- `GET /sse` - Server-Sent Events endpoint for real-time communication (optional but recommended)

**Standard MCP JSON-RPC methods** (sent to the POST / endpoint):
- `initialize` - Establishes the connection and exchanges capabilities
- `ping` - Health check/keepalive
- `notifications/initialized` - Confirms initialization is complete

**Resource-related methods:**
- `resources/list` - List available resources
- `resources/read` - Read specific resource content
- `resources/subscribe` - Subscribe to resource changes
- `resources/unsubscribe` - Unsubscribe from resource changes

**Tool-related methods:**
- `tools/list` - List available tools
- `tools/call` - Execute a specific tool

**Prompt-related methods:**
- `prompts/list` - List available prompt templates
- `prompts/get` - Retrieve a specific prompt template

**Logging methods:**
- `logging/setLevel` - Set logging level

The server should also handle standard JSON-RPC error responses and support both request/response and notification message types. Most MCP servers will implement a subset of these based on their specific functionality - for example, a server that only provides tools might not implement the resources or prompts endpoints.

The exact implementation can vary, but these form the standard MCP protocol surface area for HTTP transport.

```json
{
  "capabilities": {
    "logging": { /* logging capabilities */ },
    "prompts": { /* prompt capabilities */ },
    "resources": { /* resource capabilities */ },
    "tools": { /* tool capabilities */ }
  }
}
```

## Logging Capabilities
```json
"logging": {
  "level": "info"  // Minimum log level the server accepts
}
```

## Prompts Capabilities
```json
"prompts": {
  "listChanged": true  // Server can notify when prompt list changes
}
```

## Resources Capabilities
```json
"resources": {
  "subscribe": true,      // Server supports resource subscriptions
  "listChanged": true     // Server can notify when resource list changes
}
```

## Tools Capabilities
```json
"tools": {
  "listChanged": true     // Server can notify when tool list changes
}
```

## Example Full Response
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "logging": {
        "level": "debug"
      },
      "prompts": {
        "listChanged": true
      },
      "resources": {
        "subscribe": true,
        "listChanged": true
      },
      "tools": {
        "listChanged": true
      }
    },
    "serverInfo": {
      "name": "example-server",
      "version": "1.0.0"
    }
  }
}
```

## Usage by Clients
Based on capabilities, clients can:
- **Skip unsupported features**: Don't try to subscribe to resources if `subscribe: false`
- **Enable notifications**: Listen for list changes if `listChanged: true`
- **Adjust logging**: Set appropriate log levels based on server's minimum level
- **UI adaptation**: Show/hide features based on what's available

## Key Points
- **Optional fields**: Servers only include capabilities they support
- **Boolean flags**: Most capabilities are simple true/false declarations
- **Extensible**: New capabilities can be added without breaking existing clients
- **Defensive programming**: Clients should assume capabilities are false/unavailable if not explicitly declared

This capability negotiation ensures robust interoperability between different MCP client and server implementations.



