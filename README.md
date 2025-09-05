# gcp-security-mcp-server
A test mcp server that we can use for the hackathon.

# Gemini CLI Settings for MCP server
In your `~/.gemini/settings.json` file, make sure you have the following settings:

```
{ ...file contains other config objects
  "mcpServers": {
    "security-tools": {
      "command": "path/to/server/binary", // should be home/user/gcp-security-mcp-server/gcp-security-mcp
      "cwd": "path/to/server/directory", // should be home/user/gcp-security-mcp-server
      "timeout": 30000,
      "trust": false
    }
  }
}
```

# Development instructions
1. Ask Manraj to share the firebase workspace we're working in

2. Open the terminal by clicking the bottom bar near where the error and warning icons are.

3. Authorize with github by running 
```gh auth login```

4. Start making changes to the MCP server code.