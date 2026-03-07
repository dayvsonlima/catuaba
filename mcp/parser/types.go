package parser

// RouteInfo represents a single route in the project.
type RouteInfo struct {
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	Handler    string   `json:"handler"`
	Middleware []string `json:"middleware,omitempty"`
}

// FieldInfo represents a struct field in a model.
type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	JSON string `json:"json,omitempty"`
	GORM string `json:"gorm,omitempty"`
}

// ModelInfo represents a Go model struct.
type ModelInfo struct {
	Name   string      `json:"name"`
	Fields []FieldInfo `json:"fields"`
}

// ControllerFunc represents a function in a controller package.
type ControllerFunc struct {
	Name string `json:"name"`
	File string `json:"file"`
}

// ControllerInfo represents a controller package with its functions.
type ControllerInfo struct {
	Package   string           `json:"package"`
	Functions []ControllerFunc `json:"functions"`
}

// MiddlewareInfo represents a middleware function.
type MiddlewareInfo struct {
	Name string `json:"name"`
	File string `json:"file"`
}

// ProjectOverview holds a summary of the Catuaba project.
type ProjectOverview struct {
	Module  string   `json:"module"`
	Go      string   `json:"go"`
	DB      string   `json:"db,omitempty"`
	Port    string   `json:"port,omitempty"`
	Plugins []string `json:"plugins,omitempty"`
	Dirs    []string `json:"dirs,omitempty"`
}

// ComponentParam represents a parameter of a templ component.
type ComponentParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ComponentInfo represents a templ component with its signature.
type ComponentInfo struct {
	Name     string           `json:"name"`
	File     string           `json:"file"`
	Params   []ComponentParam `json:"params"`
	Children bool             `json:"children"`
}

// ComponentType represents a Go type defined alongside components.
type ComponentType struct {
	Name   string      `json:"name"`
	Fields []FieldInfo `json:"fields"`
}

// ComponentsResult holds all parsed components and their associated types.
type ComponentsResult struct {
	Components []ComponentInfo `json:"components"`
	Types      []ComponentType `json:"types,omitempty"`
}

// LogEntry represents a compact log entry for the get_logs tool.
type LogEntry struct {
	Time    string `json:"t,omitempty"`
	Level   string `json:"l,omitempty"`
	Status  int    `json:"s,omitempty"`
	Method  string `json:"m,omitempty"`
	Path    string `json:"p,omitempty"`
	Dur     string `json:"d,omitempty"`
	Error   string `json:"e,omitempty"`
	Message string `json:"msg,omitempty"`
}
