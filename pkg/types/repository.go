// Package types defines core data types used across Ship Shape.
package types

// Language represents a programming language detected in the repository.
type Language string

// Supported languages
const (
	LanguageGo         Language = "Go"
	LanguagePython     Language = "Python"
	LanguageJavaScript Language = "JavaScript"
	LanguageTypeScript Language = "TypeScript"
	LanguageJava       Language = "Java"
	LanguageRust       Language = "Rust"
	LanguageCSharp     Language = "C#"
	LanguageRuby       Language = "Ruby"
	LanguageUnknown    Language = "Unknown"
)

// LanguageStats contains statistics about a language in the repository.
type LanguageStats struct {
	Language   Language `json:"language"`
	FileCount  int      `json:"file_count"`
	Percentage float64  `json:"percentage"`
	IsPrimary  bool     `json:"is_primary"` // >10% of codebase
}

// Repository represents the analyzed repository context.
type Repository struct {
	// Path is the absolute path to the repository root
	Path string `json:"path"`

	// Languages detected in the repository with statistics
	Languages []LanguageStats `json:"languages"`

	// Frameworks detected (test frameworks, build systems, etc.)
	Frameworks []Framework `json:"frameworks"`

	// IsMonorepo indicates if this is a monorepo structure
	IsMonorepo bool `json:"is_monorepo"`

	// Workspaces contains monorepo workspace information
	Workspaces []Workspace `json:"workspaces,omitempty"`

	// TotalFiles is the count of analyzed files (excluding excluded paths)
	TotalFiles int `json:"total_files"`

	// ExcludedPaths are the patterns that were excluded during discovery
	ExcludedPaths []string `json:"excluded_paths"`
}

// Framework represents a detected framework or tool in the repository.
type Framework struct {
	// Name is the framework name (e.g., "pytest", "jest", "testing")
	Name string `json:"name"`

	// Language is the associated programming language
	Language Language `json:"language"`

	// Type categorizes the framework (test, build, lint, format, coverage)
	Type FrameworkType `json:"type"`

	// Version is the detected version (if available)
	Version string `json:"version,omitempty"`

	// ConfigFiles are the configuration files where this framework was detected
	ConfigFiles []string `json:"config_files,omitempty"`
}

// FrameworkType categorizes different types of frameworks and tools.
type FrameworkType string

const (
	FrameworkTypeTest     FrameworkType = "test"
	FrameworkTypeBuild    FrameworkType = "build"
	FrameworkTypeLint     FrameworkType = "lint"
	FrameworkTypeFormat   FrameworkType = "format"
	FrameworkTypeCoverage FrameworkType = "coverage"
	FrameworkTypeOther    FrameworkType = "other"
)

// Workspace represents a package or workspace in a monorepo.
type Workspace struct {
	// Name is the workspace/package name
	Name string `json:"name"`

	// Path is the relative path from repository root
	Path string `json:"path"`

	// Language is the primary language of this workspace
	Language Language `json:"language"`

	// Type indicates the workspace manager (npm, yarn, pnpm, go, maven, etc.)
	Type WorkspaceType `json:"type"`
}

// WorkspaceType identifies the workspace management system.
type WorkspaceType string

const (
	WorkspaceTypeNpm    WorkspaceType = "npm"
	WorkspaceTypeYarn   WorkspaceType = "yarn"
	WorkspaceTypePnpm   WorkspaceType = "pnpm"
	WorkspaceTypeGo     WorkspaceType = "go"
	WorkspaceTypeMaven  WorkspaceType = "maven"
	WorkspaceTypeGradle WorkspaceType = "gradle"
	WorkspaceTypeLerna  WorkspaceType = "lerna"
)

// PrimaryLanguage returns the primary language (highest percentage) in the repository.
// Returns LanguageUnknown if no languages detected.
func (r *Repository) PrimaryLanguage() Language {
	if len(r.Languages) == 0 {
		return LanguageUnknown
	}

	var primary LanguageStats
	for _, lang := range r.Languages {
		if lang.Percentage > primary.Percentage {
			primary = lang
		}
	}

	return primary.Language
}

// HasLanguage checks if a specific language is present in the repository.
func (r *Repository) HasLanguage(lang Language) bool {
	for _, l := range r.Languages {
		if l.Language == lang {
			return true
		}
	}

	return false
}

// GetFramework returns the first framework matching the given name.
// Returns nil if not found.
func (r *Repository) GetFramework(name string) *Framework {
	for i := range r.Frameworks {
		if r.Frameworks[i].Name == name {
			return &r.Frameworks[i]
		}
	}

	return nil
}

// HasFramework checks if a framework with the given name exists.
func (r *Repository) HasFramework(name string) bool {
	return r.GetFramework(name) != nil
}

// GetFrameworksByType returns all frameworks of a specific type.
func (r *Repository) GetFrameworksByType(ftype FrameworkType) []Framework {
	var frameworks []Framework
	for _, fw := range r.Frameworks {
		if fw.Type == ftype {
			frameworks = append(frameworks, fw)
		}
	}

	return frameworks
}
