package mcp

type Capabilities struct {
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`
	Prompts struct {
		ListChanged bool `json:"listChanged"`
	} `json:"prompts"`
	Resources struct {
		Subscribe   bool `json:"subscribe"`
		ListChanged bool `json:"listChanged"`
	} `json:"resources"`
	Tools struct {
		ListChanged bool `json:"listChanged"`
	} `json:"listChanged"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Response struct {
	JSONRPCProtocolVersion string `json:"jsonrpc"`
	ID                     int    `json:"id"`
	Result                 struct {
		ProtocolVersion string       `json:"protocolVersion"`
		Capabilities    Capabilities `json:"capabilities"`
		ServerInfo      ServerInfo   `json:"serverInfo"`
	} `json:"Result"`
}

func NewResponse() Response {
	response := Response{
		JSONRPCProtocolVersion: "2.0",
		ID:                     1,
	}

	response.Result.ProtocolVersion = "2024-11-05"
	response.Result.Capabilities.Logging.Level = "debug"
	response.Result.Capabilities.Prompts.ListChanged = true
	response.Result.Capabilities.Resources.Subscribe = true
	response.Result.Capabilities.Resources.ListChanged = true
	response.Result.Capabilities.Tools.ListChanged = true
	response.Result.ServerInfo.Name = "DND Need a name!"
	response.Result.ServerInfo.Version = "1.0.0"

	return response
}
