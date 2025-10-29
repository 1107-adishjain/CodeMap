package models

// Analysis represents the top-level structure of our analysis-output.json.
type Analysis struct {
	Files []File `json:"files"`
}

// File represents a single source code file.
type File struct {
	Path      string     `json:"path"`
	Language  string     `json:"language"`
	Classes   []Class    `json:"classes,omitempty"`
	Functions []Function `json:"functions,omitempty"`
	Imports   []Import   `json:"imports,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// Class represents a class definition.
type Class struct {
	Name       string   `json:"name"`
	IsExported bool     `json:"is_exported"`
	Properties []string `json:"properties,omitempty"`
	Methods    []string `json:"methods,omitempty"`
}

// Function represents a function or method.
type Function struct {
	Name        string   `json:"name"`
	IsExported  bool     `json:"is_exported"`
	Params      []string `json:"params,omitempty"`
	ReturnTypes []string `json:"return_types,omitempty"`
	Calls       []string `json:"calls,omitempty"`
	IsMethodOf  string   `json:"is_method_of,omitempty"`
}

// Import represents an import statement.
type Import struct {
	Source string `json:"source"`
}
