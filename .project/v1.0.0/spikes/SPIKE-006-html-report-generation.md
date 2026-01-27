# SPIKE-006: Interactive HTML Report Generation

## Overview
This spike validates the technical approach for generating beautiful, interactive HTML reports with charts, visualizations, and detailed analysis results using Go templates and modern JavaScript charting libraries.

**Associated User Stories**: SS-090
**Risk Level**: HIGH (13 story points)
**Priority**: P1 (High)
**Target Completion**: Week 6-7 of implementation

## Problem Statement
Ship Shape needs to generate comprehensive, visually appealing HTML reports that:
- Display overall quality scores with grade badges
- Show interactive charts (radar, pie, line, bar charts)
- Provide detailed findings with code snippets
- Enable drill-down navigation
- Work as standalone HTML (offline viewing)
- Support dark/light themes
- Be mobile-responsive
- Include file-level coverage heatmaps

**Key Challenges**:
- Embedding interactive charts in Go-generated HTML
- Standalone HTML (no external dependencies)
- Chart data generation from analysis results
- Large reports (10,000+ findings) performance
- Code syntax highlighting
- Responsive design
- Accessibility (WCAG 2.1 AA compliance)

## Existing Tools and Libraries

### Go HTML Templating
- **html/template** - Go standard library (secure, fast)
- **templ** - Type-safe Go templating
- **Gorilla templates** - Template utilities

### Chart Libraries for Go + HTML

**go-echarts** ([GitHub](https://github.com/go-echarts/go-echarts)) - Primary Choice
- ✅ Over 6.4k stars on GitHub
- ✅ Apache ECharts wrapper for Go
- ✅ Interactive charts (radar, pie, line, bar, heatmap, etc.)
- ✅ Generates standalone HTML with embedded JavaScript
- ✅ Responsive and customizable
- ✅ Dark/light theme support
- ✅ [Documentation](https://blog.logrocket.com/visualizing-data-go-echarts/)

**Alternative Libraries**:
- **Chart.js** - Popular JavaScript library (can be embedded)
- **gonum/plot** - Go plotting library (generates static images)
- **go chart** - Basic charting (limited interactivity)

### Syntax Highlighting
- **highlight.js** - JavaScript syntax highlighter (172 languages)
- **Prism.js** - Lightweight syntax highlighter
- **chroma** - Go syntax highlighter (generates HTML)

### CSS Frameworks
- **Tailwind CSS** - Utility-first CSS (via CDN or embedded)
- **Bootstrap** - Comprehensive framework
- **Pure.css** - Minimal CSS framework

## Spike Objectives
- [ ] Design report page structure and layout
- [ ] Implement Go template system
- [ ] Integrate go-echarts for interactive charts
- [ ] Create responsive CSS design
- [ ] Add syntax highlighting for code snippets
- [ ] Implement dark/light theme toggle
- [ ] Optimize for large reports (10,000+ findings)
- [ ] Validate accessibility (WCAG 2.1 AA)
- [ ] Test cross-browser compatibility

## Technical Investigation Areas

### 1. Report Structure Design

```
Ship Shape HTML Report
├── Header
│   ├── Logo and Title
│   ├── Repository Info
│   └── Overall Score Badge
├── Executive Summary
│   ├── Score Card (Grade + Score)
│   ├── Key Metrics Table
│   └── Quick Stats
├── Interactive Charts Section
│   ├── Score Breakdown (Radar Chart)
│   ├── Language Distribution (Pie Chart)
│   ├── Coverage Trends (Line Chart)
│   ├── Test Pyramid (Bar Chart)
│   └── Coverage Heatmap (Heatmap)
├── Detailed Analysis Tabs
│   ├── Test Quality
│   ├── Coverage
│   ├── Tool Adoption
│   ├── Test Smells
│   └── CI/CD Analysis
├── Findings Table
│   ├── Filterable by severity/type
│   ├── Sortable columns
│   └── Expandable details
└── Footer
    ├── Generation timestamp
    ├── Ship Shape version
    └── Export options
```

### 2. Go Template Implementation

**Base Template Structure**:
```go
package report

import (
    "html/template"
    "embed"
)

//go:embed templates/*.html
var templateFS embed.FS

type HTMLReportGenerator struct {
    templates *template.Template
    chartGen  *ChartGenerator
}

func NewHTMLReportGenerator() (*HTMLReportGenerator, error) {
    tmpl, err := template.ParseFS(templateFS, "templates/*.html")
    if err != nil {
        return nil, err
    }

    return &HTMLReportGenerator{
        templates: tmpl,
        chartGen:  NewChartGenerator(),
    }, nil
}

type ReportData struct {
    Title            string
    Repository       *RepositoryInfo
    GeneratedAt      time.Time
    OverallScore     int
    Grade            string
    DimensionScores  map[string]float64
    LanguageDistribution []*LanguageStat
    CoverageMetrics  *CoverageMetrics
    TestSmells       []*TestSmell
    Findings         []*Finding
    Charts           *ChartData
}

func (hrg *HTMLReportGenerator) Generate(
    report *Report,
    outputPath string,
) error {
    // Generate chart data
    chartData, err := hrg.chartGen.GenerateCharts(report)
    if err != nil {
        return err
    }

    // Prepare template data
    data := &ReportData{
        Title:       "Ship Shape Analysis Report",
        Repository:  report.Repository,
        GeneratedAt: time.Now(),
        OverallScore: report.Score,
        Grade:        report.Grade,
        Charts:       chartData,
        // ... more fields
    }

    // Create output file
    f, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer f.Close()

    // Execute template
    return hrg.templates.ExecuteTemplate(f, "report.html", data)
}
```

### 3. Chart Generation with go-echarts

```go
package report

import (
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
)

type ChartGenerator struct{}

type ChartData struct {
    RadarChartHTML    template.HTML
    PieChartHTML      template.HTML
    LineChartHTML     template.HTML
    BarChartHTML      template.HTML
    HeatmapHTML       template.HTML
}

// Generate radar chart for score dimensions
func (cg *ChartGenerator) GenerateRadarChart(
    dimensions map[string]float64,
) template.HTML {
    radar := charts.NewRadar()
    radar.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "Quality Dimensions",
        }),
        charts.WithTooltipOpts(opts.Tooltip{Show: true}),
        charts.WithLegendOpts(opts.Legend{Show: true}),
    )

    // Radar indicator setup
    indicators := make([]*opts.Indicator, 0)
    values := make([]float64, 0)

    for dim, score := range dimensions {
        indicators = append(indicators, &opts.Indicator{
            Name: dim,
            Max:  100,
        })
        values = append(values, score)
    }

    radar.AddSeries("Scores", []opts.RadarData{
        {Value: values},
    }).SetRadarComponent(
        opts.RadarComponent{Indicator: indicators},
    )

    // Render to HTML string
    var buf bytes.Buffer
    radar.Render(&buf)
    return template.HTML(buf.String())
}

// Generate pie chart for language distribution
func (cg *ChartGenerator) GeneratePieChart(
    languages []*LanguageStat,
) template.HTML {
    pie := charts.NewPie()
    pie.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "Language Distribution",
        }),
        charts.WithLegendOpts(opts.Legend{
            Orient: "vertical",
            Left:   "left",
        }),
    )

    items := make([]opts.PieData, 0)
    for _, lang := range languages {
        items = append(items, opts.PieData{
            Name:  lang.Name,
            Value: lang.Percentage,
        })
    }

    pie.AddSeries("Languages", items).
        SetSeriesOptions(
            charts.WithLabelOpts(opts.Label{
                Show:      true,
                Formatter: "{b}: {d}%",
            }),
        )

    var buf bytes.Buffer
    pie.Render(&buf)
    return template.HTML(buf.String())
}

// Generate line chart for coverage trends
func (cg *ChartGenerator) GenerateLineChart(
    historicalData []*HistoricalCoverage,
) template.HTML {
    line := charts.NewLine()
    line.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "Coverage Trends",
        }),
        charts.WithXAxisOpts(opts.XAxis{
            Name: "Date",
        }),
        charts.WithYAxisOpts(opts.YAxis{
            Name: "Coverage %",
            Min:  0,
            Max:  100,
        }),
        charts.WithTooltipOpts(opts.Tooltip{
            Show:    true,
            Trigger: "axis",
        }),
        charts.WithDataZoomOpts(opts.DataZoom{
            Type:  "slider",
            Start: 0,
            End:   100,
        }),
    )

    dates := make([]string, 0)
    lineCov := make([]opts.LineData, 0)
    branchCov := make([]opts.LineData, 0)

    for _, data := range historicalData {
        dates = append(dates, data.Date.Format("2006-01-02"))
        lineCov = append(lineCov, opts.LineData{
            Value: data.Coverage.OverallMetrics.LineRate * 100,
        })
        branchCov = append(branchCov, opts.LineData{
            Value: data.Coverage.OverallMetrics.BranchRate * 100,
        })
    }

    line.SetXAxis(dates).
        AddSeries("Line Coverage", lineCov).
        AddSeries("Branch Coverage", branchCov).
        SetSeriesOptions(
            charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
        )

    var buf bytes.Buffer
    line.Render(&buf)
    return template.HTML(buf.String())
}

// Generate heatmap for file-level coverage
func (cg *ChartGenerator) GenerateHeatmap(
    files []*FileCoverage,
) template.HTML {
    heatmap := charts.NewHeatMap()

    // Prepare data
    xAxis := make([]string, 0)  // Files
    yAxis := []string{"Coverage"}
    data := make([]opts.HeatMapData, 0)

    for i, file := range files {
        xAxis = append(xAxis, filepath.Base(file.Path))
        coveragePct := file.Metrics.LineRate * 100

        data = append(data, opts.HeatMapData{
            Value: [3]interface{}{i, 0, coveragePct},
        })
    }

    heatmap.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "File Coverage Heatmap",
        }),
        charts.WithXAxisOpts(opts.XAxis{
            Type:      "category",
            Data:      xAxis,
            SplitArea: &opts.SplitArea{Show: true},
        }),
        charts.WithYAxisOpts(opts.YAxis{
            Type:      "category",
            Data:      yAxis,
            SplitArea: &opts.SplitArea{Show: true},
        }),
        charts.WithVisualMapOpts(opts.VisualMap{
            Calculable: true,
            Min:        0,
            Max:        100,
            InRange: &opts.VisualMapInRange{
                Color: []string{"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8", "#ffffbf", "#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026"},
            },
        }),
    )

    heatmap.AddSeries("coverage", data)

    var buf bytes.Buffer
    heatmap.Render(&buf)
    return template.HTML(buf.String())
}
```

### 4. HTML Template Structure

**Base Template** (`templates/report.html`):
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.Repository.Name}}</title>

    <!-- Embedded CSS (for standalone HTML) -->
    <style>
        {{template "styles.css"}}
    </style>

    <!-- Theme toggle script -->
    <script>
        const theme = localStorage.getItem('theme') || 'light';
        document.documentElement.setAttribute('data-theme', theme);
    </script>
</head>
<body>
    <header class="report-header">
        <div class="container">
            <h1>🚢 Ship Shape Analysis Report</h1>
            <div class="repo-info">
                <span class="repo-name">{{.Repository.Name}}</span>
                <span class="generated-at">Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</span>
            </div>
            <button id="theme-toggle" class="theme-toggle">🌙</button>
        </div>
    </header>

    <main class="container">
        <!-- Executive Summary -->
        <section class="executive-summary">
            <div class="score-card">
                <div class="score-badge grade-{{.Grade}}">
                    <span class="grade">{{.Grade}}</span>
                    <span class="score">{{.OverallScore}}/100</span>
                </div>
                <div class="quick-stats">
                    {{template "quick-stats" .}}
                </div>
            </div>
        </section>

        <!-- Interactive Charts -->
        <section class="charts-section">
            <h2>📊 Visual Analysis</h2>
            <div class="chart-grid">
                <div class="chart-container">
                    {{.Charts.RadarChartHTML}}
                </div>
                <div class="chart-container">
                    {{.Charts.PieChartHTML}}
                </div>
                <div class="chart-container full-width">
                    {{.Charts.LineChartHTML}}
                </div>
                <div class="chart-container full-width">
                    {{.Charts.HeatmapHTML}}
                </div>
            </div>
        </section>

        <!-- Detailed Findings -->
        <section class="findings-section">
            <h2>🔍 Detailed Findings</h2>
            {{template "findings-table" .Findings}}
        </section>

        <!-- Test Smells -->
        <section class="smells-section">
            <h2>⚠️ Test Smells Detected</h2>
            {{template "smells-list" .TestSmells}}
        </section>
    </main>

    <footer class="report-footer">
        <div class="container">
            <p>Generated by Ship Shape v{{.Version}}</p>
        </div>
    </footer>

    <!-- Embedded JavaScript -->
    <script>
        {{template "scripts.js"}}
    </script>
</body>
</html>
```

**Findings Table Template**:
```html
{{define "findings-table"}}
<div class="findings-table-container">
    <div class="table-controls">
        <input type="text" id="findings-search" placeholder="Search findings...">
        <select id="severity-filter">
            <option value="">All Severities</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
        </select>
    </div>

    <table class="findings-table" id="findings-table">
        <thead>
            <tr>
                <th data-sort="severity">Severity</th>
                <th data-sort="type">Type</th>
                <th data-sort="title">Title</th>
                <th data-sort="file">File</th>
                <th data-sort="line">Line</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr class="finding-row severity-{{.Severity}}" data-finding-id="{{.ID}}">
                <td><span class="severity-badge severity-{{.Severity}}">{{.Severity}}</span></td>
                <td>{{.Type}}</td>
                <td>{{.Title}}</td>
                <td>{{.Location.File}}</td>
                <td>{{.Location.Line}}</td>
                <td>
                    <button class="btn-expand" onclick="toggleDetails('{{.ID}}')">Details</button>
                </td>
            </tr>
            <tr class="finding-details" id="details-{{.ID}}" style="display: none;">
                <td colspan="6">
                    <div class="details-content">
                        <p><strong>Description:</strong> {{.Description}}</p>
                        <p><strong>Remediation:</strong> {{.Remediation}}</p>
                        {{if .Code}}
                        <pre><code class="language-{{.Language}}">{{.Code}}</code></pre>
                        {{end}}
                    </div>
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>
</div>
{{end}}
```

### 5. Responsive CSS with Theme Support

```css
/* CSS Variables for theming */
:root[data-theme="light"] {
    --bg-primary: #ffffff;
    --bg-secondary: #f5f5f5;
    --text-primary: #333333;
    --text-secondary: #666666;
    --border-color: #e0e0e0;
    --grade-a: #4caf50;
    --grade-b: #2196f3;
    --grade-c: #ff9800;
    --grade-d: #f44336;
}

:root[data-theme="dark"] {
    --bg-primary: #1e1e1e;
    --bg-secondary: #2d2d2d;
    --text-primary: #e0e0e0;
    --text-secondary: #b0b0b0;
    --border-color: #404040;
    --grade-a: #66bb6a;
    --grade-b: #42a5f5;
    --grade-c: #ffa726;
    --grade-d: #ef5350;
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background-color: var(--bg-primary);
    color: var(--text-primary);
    margin: 0;
    padding: 0;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

/* Responsive Grid for Charts */
.chart-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
    gap: 20px;
    margin: 20px 0;
}

.chart-container.full-width {
    grid-column: 1 / -1;
}

@media (max-width: 768px) {
    .chart-grid {
        grid-template-columns: 1fr;
    }
}

/* Score Badge */
.score-badge {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    padding: 30px;
    border-radius: 10px;
    font-size: 24px;
    font-weight: bold;
}

.score-badge.grade-A,
.score-badge.grade-A\+ {
    background: linear-gradient(135deg, var(--grade-a), #66bb6a);
    color: white;
}

/* Findings Table */
.findings-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
}

.findings-table th,
.findings-table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
}

.findings-table th {
    background-color: var(--bg-secondary);
    font-weight: 600;
    cursor: pointer;
    user-select: none;
}

.findings-table th:hover {
    background-color: var(--border-color);
}

.severity-badge {
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
}

.severity-critical {
    background-color: #f44336;
    color: white;
}

.severity-high {
    background-color: #ff9800;
    color: white;
}
```

### 6. Interactive JavaScript for Filtering and Sorting

```javascript
// Theme toggle
document.getElementById('theme-toggle').addEventListener('click', () => {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    document.getElementById('theme-toggle').textContent = next === 'light' ? '🌙' : '☀️';
});

// Findings table filtering
document.getElementById('findings-search').addEventListener('input', (e) => {
    const search = e.target.value.toLowerCase();
    document.querySelectorAll('.finding-row').forEach(row => {
        const text = row.textContent.toLowerCase();
        row.style.display = text.includes(search) ? '' : 'none';
    });
});

// Severity filtering
document.getElementById('severity-filter').addEventListener('change', (e) => {
    const severity = e.target.value;
    document.querySelectorAll('.finding-row').forEach(row => {
        if (!severity || row.classList.contains(`severity-${severity}`)) {
            row.style.display = '';
        } else {
            row.style.display = 'none';
        }
    });
});

// Table sorting
document.querySelectorAll('th[data-sort]').forEach(th => {
    th.addEventListener('click', () => {
        const table = th.closest('table');
        const tbody = table.querySelector('tbody');
        const rows = Array.from(tbody.querySelectorAll('.finding-row'));
        const column = th.cellIndex;
        const order = th.dataset.order === 'asc' ? 'desc' : 'asc';

        rows.sort((a, b) => {
            const aVal = a.cells[column].textContent;
            const bVal = b.cells[column].textContent;
            return order === 'asc'
                ? aVal.localeCompare(bVal)
                : bVal.localeCompare(aVal);
        });

        rows.forEach(row => tbody.appendChild(row));
        th.dataset.order = order;
    });
});

// Toggle finding details
function toggleDetails(id) {
    const details = document.getElementById(`details-${id}`);
    details.style.display = details.style.display === 'none' ? '' : 'none';
}

// Initialize syntax highlighting
document.querySelectorAll('pre code').forEach(block => {
    hljs.highlightElement(block);
});
```

## Prototype Requirements

### Deliverable 1: Template System
**Files**: `internal/report/html/templates/*.html`
- Base layout template
- Executive summary template
- Charts section template
- Findings table template
- Test smells template

### Deliverable 2: Chart Generation Module
**Files**: `internal/report/html/charts.go`
- Radar chart generator
- Pie chart generator
- Line chart generator (trends)
- Heatmap generator
- Bar chart generator (test pyramid)

### Deliverable 3: CSS Styling
**Files**: `internal/report/html/templates/styles.css`
- Responsive grid layout
- Theme variables (light/dark)
- Component styles
- Mobile responsiveness
- Print styles

### Deliverable 4: Interactive JavaScript
**Files**: `internal/report/html/templates/scripts.js`
- Table filtering and sorting
- Theme toggle
- Details expansion
- Export functionality

### Deliverable 5: Report Generator
**Files**: `internal/report/html/generator.go`
- Template orchestration
- Data preparation
- Output file generation
- Standalone HTML embedding

## Performance Benchmarks

| Report Size | Findings | Generation Time | HTML Size | Load Time |
|------------|---------|-----------------|-----------|-----------|
| Small      | 100     | <1s             | 500KB     | <1s       |
| Medium     | 1,000   | <3s             | 2MB       | <2s       |
| Large      | 10,000  | <10s            | 10MB      | <5s       |

## Validation Strategy

### 1. Visual Testing
- Manual review on multiple browsers
- Mobile device testing
- Accessibility audit (WAVE, axe)
- Print layout verification

### 2. Performance Testing
- Lighthouse scores (>90 target)
- Large report load time
- Chart rendering performance
- Memory usage monitoring

### 3. Compatibility Testing
- Chrome, Firefox, Safari, Edge
- iOS Safari, Android Chrome
- Different screen sizes

## Risk Mitigation

### Risk 1: Large HTML file size
**Mitigation**:
- Lazy load charts
- Paginate findings table
- Compress embedded assets
- Optional external dependencies mode

### Risk 2: Chart rendering performance
**Mitigation**:
- Limit data points
- Use canvas instead of SVG for large datasets
- Progressive rendering
- Virtual scrolling for tables

### Risk 3: Accessibility issues
**Mitigation**:
- Semantic HTML
- ARIA labels
- Keyboard navigation
- Screen reader testing

## Go/No-Go Decision Criteria

### GO if:
- ✅ Interactive charts render correctly
- ✅ Responsive design works on mobile
- ✅ Theme toggle functional
- ✅ Performance acceptable (Large <10s)
- ✅ Accessibility score >90%
- ✅ Cross-browser compatible

### NO-GO if:
- ❌ Charts don't render in major browsers
- ❌ Performance unacceptable (>30s for large)
- ❌ Accessibility score <70%
- ❌ Not mobile-responsive

## Spike Deliverables

1. **Template System**
   - HTML templates
   - Template data structures
   - Template execution

2. **Chart Generation**
   - go-echarts integration
   - All chart types implemented
   - Chart customization

3. **Styling and Theming**
   - CSS with theme support
   - Responsive design
   - Component library

4. **Interactive Features**
   - Filtering and sorting
   - Theme toggle
   - Details expansion

5. **Sample Reports**
   - Small project example
   - Large project example
   - Monorepo example

## Success Metrics
- [ ] Charts render interactively
- [ ] Responsive on all devices
- [ ] Theme toggle works
- [ ] Performance targets met
- [ ] Accessibility compliant
- [ ] Cross-browser compatible

## Timeline
- **Week 1**: Template system and basic layout
- **Week 2**: Chart integration (go-echarts)
- **Week 3**: Styling and responsive design
- **Week 4**: Interactive features and optimization

## Sources and References
- [go-echarts GitHub](https://github.com/go-echarts/go-echarts)
- [Creating Charts in Go with ECharts](https://zetcode.com/golang/echarts/)
- [Visualizing data in Go-echarts](https://blog.logrocket.com/visualizing-data-go-echarts/)
- [Chart.js Documentation](https://www.chartjs.org/)
- [HTML Template Package (Go)](https://pkg.go.dev/html/template)
- [highlight.js](https://highlightjs.org/)
