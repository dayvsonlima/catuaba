package mcp

import (
	"log"
	"os"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates and configures the Catuaba MCP server with all resources and tools.
func NewServer(projectDir string) *server.MCPServer {
	s := server.NewMCPServer(
		"catuaba",
		"0.1.0",
		server.WithResourceCapabilities(false, false),
		server.WithToolCapabilities(false),
		server.WithLogging(),
		server.WithRecovery(),
		server.WithInstructions("Catuaba MCP server exposes project structure (models, routes, controllers, middleware, env, UI components) as compact JSON and provides code generation tools. Use get_component to discover available templ components and their signatures before creating views."),
	)

	cache := NewCache()

	registerResources(s, projectDir, cache)
	registerTools(s, projectDir, cache)

	return s
}

// Serve starts the MCP server on stdio transport.
func Serve(projectDir string) error {
	s := NewServer(projectDir)
	errLogger := log.New(os.Stderr, "[catuaba-mcp] ", log.LstdFlags)
	return server.ServeStdio(s, server.WithErrorLogger(errLogger))
}

func registerResources(s *server.MCPServer, projectDir string, cache *Cache) {
	addResource(s, "catuaba://project/overview", "Project Overview",
		"Compact JSON overview: module name, Go version, database, port, plugins, directories",
		projectDir, cache, handleOverviewResource)

	addResource(s, "catuaba://project/models", "Models",
		"All model structs with fields, types, and tags (JSON/GORM)",
		projectDir, cache, handleModelsResource)

	addResource(s, "catuaba://project/routes", "Routes",
		"All HTTP routes with method, path, handler, and middleware",
		projectDir, cache, handleRoutesResource)

	addResource(s, "catuaba://project/controllers", "Controllers",
		"All controller packages with their exported functions",
		projectDir, cache, handleControllersResource)

	addResource(s, "catuaba://project/middleware", "Middleware",
		"All middleware functions with file locations",
		projectDir, cache, handleMiddlewareResource)

	addResource(s, "catuaba://project/env", "Environment Variables",
		"List of environment variable names from .env (no values)",
		projectDir, cache, handleEnvResource)

	addResource(s, "catuaba://project/components", "UI Components",
		"All templ UI components with signatures, params, children support, and associated types",
		projectDir, cache, handleComponentsResource)
}

func registerTools(s *server.MCPServer, projectDir string, cache *Cache) {
	s.AddTool(
		mcplib.NewTool("get_model",
			mcplib.WithDescription("Get model struct details. Returns all models if name is omitted."),
			mcplib.WithString("name", mcplib.Description("Model name (optional, returns all if empty)")),
		),
		makeGetModelHandler(projectDir, cache),
	)

	s.AddTool(
		mcplib.NewTool("get_routes",
			mcplib.WithDescription("List HTTP routes. Optionally filter by path prefix."),
			mcplib.WithString("prefix", mcplib.Description("Path prefix filter (optional, e.g. /auth)")),
		),
		makeGetRoutesHandler(projectDir, cache),
	)

	s.AddTool(
		mcplib.NewTool("get_logs",
			mcplib.WithDescription("Get recent application logs, filtered and compact for debugging."),
			mcplib.WithNumber("lines", mcplib.Description("Number of lines to return (default 50)")),
			mcplib.WithString("level", mcplib.Description("Filter by log level: error, warn, info, debug")),
			mcplib.WithString("path", mcplib.Description("Filter by request path prefix")),
		),
		makeGetLogsHandler(projectDir),
	)

	s.AddTool(
		mcplib.NewTool("generate_scaffold",
			mcplib.WithDescription("Generate a full scaffold: model + CRUD controllers + routes"),
			mcplib.WithString("name", mcplib.Required(), mcplib.Description("Resource name (e.g. Product)")),
			mcplib.WithString("fields", mcplib.Required(), mcplib.Description("Comma-separated field definitions (e.g. title:string,price:float64)")),
		),
		makeGenerateScaffoldHandler(cache),
	)

	s.AddTool(
		mcplib.NewTool("generate_model",
			mcplib.WithDescription("Generate a model file"),
			mcplib.WithString("name", mcplib.Required(), mcplib.Description("Model name (e.g. Order)")),
			mcplib.WithString("fields", mcplib.Required(), mcplib.Description("Comma-separated field definitions (e.g. total:float64,status:string)")),
		),
		makeGenerateModelHandler(cache),
	)

	s.AddTool(
		mcplib.NewTool("install_plugin",
			mcplib.WithDescription("Install a Catuaba plugin by name, URL, or local path"),
			mcplib.WithString("source", mcplib.Required(), mcplib.Description("Plugin name, git URL, or local path")),
		),
		makeInstallPluginHandler(cache),
	)

	s.AddTool(
		mcplib.NewTool("get_component",
			mcplib.WithDescription("Get UI component details. Returns all components if name is omitted, or the full source of a specific component."),
			mcplib.WithString("name", mcplib.Description("Component name (optional, returns all if empty)")),
		),
		makeGetComponentHandler(projectDir, cache),
	)

	s.AddTool(
		mcplib.NewTool("run_command",
			mcplib.WithDescription("Run any catuaba generate subcommand"),
			mcplib.WithString("command", mcplib.Required(), mcplib.Description("Generate subcommand (e.g. controller, middleware, service, seed, migration)")),
			mcplib.WithString("args", mcplib.Required(), mcplib.Description("Comma-separated arguments (e.g. auth,login,register)")),
		),
		makeRunCommandHandler(cache),
	)
}

func addResource(s *server.MCPServer, uri, name, description, projectDir string, cache *Cache, handler resourceMaker) {
	s.AddResource(
		mcplib.NewResource(uri, name,
			mcplib.WithResourceDescription(description),
			mcplib.WithMIMEType("application/json"),
		),
		handler(projectDir, cache),
	)
}

type resourceMaker func(projectDir string, cache *Cache) server.ResourceHandlerFunc
