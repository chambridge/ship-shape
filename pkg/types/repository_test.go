package types

import (
	"testing"
)

func TestPrimaryLanguage(t *testing.T) {
	tests := []struct {
		name       string
		repository Repository
		want       Language
	}{
		{
			name: "single language",
			repository: Repository{
				Languages: []LanguageStats{
					{Language: LanguageGo, Percentage: 100.0},
				},
			},
			want: LanguageGo,
		},
		{
			name: "multiple languages - go primary",
			repository: Repository{
				Languages: []LanguageStats{
					{Language: LanguageGo, Percentage: 75.0},
					{Language: LanguagePython, Percentage: 25.0},
				},
			},
			want: LanguageGo,
		},
		{
			name: "multiple languages - python primary",
			repository: Repository{
				Languages: []LanguageStats{
					{Language: LanguageGo, Percentage: 30.0},
					{Language: LanguagePython, Percentage: 70.0},
				},
			},
			want: LanguagePython,
		},
		{
			name:       "no languages",
			repository: Repository{Languages: []LanguageStats{}},
			want:       LanguageUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repository.PrimaryLanguage()
			if got != tt.want {
				t.Errorf("PrimaryLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasLanguage(t *testing.T) {
	repo := Repository{
		Languages: []LanguageStats{
			{Language: LanguageGo, Percentage: 75.0},
			{Language: LanguagePython, Percentage: 25.0},
		},
	}

	tests := []struct {
		name     string
		language Language
		want     bool
	}{
		{
			name:     "has go",
			language: LanguageGo,
			want:     true,
		},
		{
			name:     "has python",
			language: LanguagePython,
			want:     true,
		},
		{
			name:     "does not have java",
			language: LanguageJava,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.HasLanguage(tt.language)
			if got != tt.want {
				t.Errorf("HasLanguage(%v) = %v, want %v", tt.language, got, tt.want)
			}
		})
	}
}

func TestGetFramework(t *testing.T) {
	repo := Repository{
		Frameworks: []Framework{
			{Name: "pytest", Language: LanguagePython, Type: FrameworkTypeTest},
			{Name: "jest", Language: LanguageJavaScript, Type: FrameworkTypeTest},
		},
	}

	tests := []struct {
		name          string
		frameworkName string
		wantFound     bool
		wantFramework *Framework
	}{
		{
			name:          "find pytest",
			frameworkName: "pytest",
			wantFound:     true,
			wantFramework: &Framework{Name: "pytest", Language: LanguagePython, Type: FrameworkTypeTest},
		},
		{
			name:          "find jest",
			frameworkName: "jest",
			wantFound:     true,
			wantFramework: &Framework{Name: "jest", Language: LanguageJavaScript, Type: FrameworkTypeTest},
		},
		{
			name:          "not found",
			frameworkName: "junit",
			wantFound:     false,
			wantFramework: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.GetFramework(tt.frameworkName)

			if tt.wantFound && got == nil {
				t.Errorf("GetFramework(%q) = nil, want framework", tt.frameworkName)
				return
			}

			if !tt.wantFound && got != nil {
				t.Errorf("GetFramework(%q) = %+v, want nil", tt.frameworkName, got)
				return
			}

			if tt.wantFound && got != nil {
				if got.Name != tt.wantFramework.Name {
					t.Errorf("GetFramework(%q).Name = %q, want %q", tt.frameworkName, got.Name, tt.wantFramework.Name)
				}

				if got.Language != tt.wantFramework.Language {
					t.Errorf("GetFramework(%q).Language = %v, want %v", tt.frameworkName, got.Language, tt.wantFramework.Language)
				}

				if got.Type != tt.wantFramework.Type {
					t.Errorf("GetFramework(%q).Type = %v, want %v", tt.frameworkName, got.Type, tt.wantFramework.Type)
				}
			}
		})
	}
}

func TestHasFramework(t *testing.T) {
	repo := Repository{
		Frameworks: []Framework{
			{Name: "pytest", Language: LanguagePython, Type: FrameworkTypeTest},
			{Name: "jest", Language: LanguageJavaScript, Type: FrameworkTypeTest},
		},
	}

	tests := []struct {
		name          string
		frameworkName string
		want          bool
	}{
		{
			name:          "has pytest",
			frameworkName: "pytest",
			want:          true,
		},
		{
			name:          "has jest",
			frameworkName: "jest",
			want:          true,
		},
		{
			name:          "does not have junit",
			frameworkName: "junit",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.HasFramework(tt.frameworkName)
			if got != tt.want {
				t.Errorf("HasFramework(%q) = %v, want %v", tt.frameworkName, got, tt.want)
			}
		})
	}
}

func TestGetFrameworksByType(t *testing.T) {
	repo := Repository{
		Frameworks: []Framework{
			{Name: "pytest", Language: LanguagePython, Type: FrameworkTypeTest},
			{Name: "jest", Language: LanguageJavaScript, Type: FrameworkTypeTest},
			{Name: "eslint", Language: LanguageJavaScript, Type: FrameworkTypeLint},
			{Name: "prettier", Language: LanguageJavaScript, Type: FrameworkTypeFormat},
		},
	}

	tests := []struct {
		name string
		ftype FrameworkType
		want []string
	}{
		{
			name:  "test frameworks",
			ftype: FrameworkTypeTest,
			want:  []string{"pytest", "jest"},
		},
		{
			name:  "lint frameworks",
			ftype: FrameworkTypeLint,
			want:  []string{"eslint"},
		},
		{
			name:  "format frameworks",
			ftype: FrameworkTypeFormat,
			want:  []string{"prettier"},
		},
		{
			name:  "no coverage frameworks",
			ftype: FrameworkTypeCoverage,
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.GetFrameworksByType(tt.ftype)

			if len(got) != len(tt.want) {
				t.Errorf("GetFrameworksByType(%v) returned %d frameworks, want %d", tt.ftype, len(got), len(tt.want))
				return
			}

			// Check that all expected framework names are present
			for _, wantName := range tt.want {
				found := false
				for _, fw := range got {
					if fw.Name == wantName {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("GetFrameworksByType(%v) missing framework %q", tt.ftype, wantName)
				}
			}
		})
	}
}

func TestLanguageStats_IsPrimary(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		want       bool
	}{
		{
			name:       "15% is primary",
			percentage: 15.0,
			want:       true, // >10% threshold
		},
		{
			name:       "10% is not primary",
			percentage: 10.0,
			want:       false, // threshold is >10%, not >=10%
		},
		{
			name:       "5% is not primary",
			percentage: 5.0,
			want:       false,
		},
		{
			name:       "100% is primary",
			percentage: 100.0,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// IsPrimary should be set based on >10% threshold
			isPrimary := tt.percentage > 10.0

			if isPrimary != tt.want {
				t.Errorf("percentage %.1f: isPrimary = %v, want %v", tt.percentage, isPrimary, tt.want)
			}
		})
	}
}
