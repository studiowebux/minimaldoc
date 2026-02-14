package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/studiowebux/minimaldoc/internal/core"
)

// BaseFuncMap returns the common template functions used across all generators.
func BaseFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":           dictFunc,
		"safeHTML":       safeHTMLFunc,
		"json":           jsonFunc,
		"lower":          lowerFunc,
		"upper":          upperFunc,
		"hasPrefix":      strings.HasPrefix,
		"add":            addFunc,
		"join":           strings.Join,
		"replace":        replaceFunc,
		"hasCustomTheme": hasCustomThemeFunc,
	}
}

// hasCustomThemeFunc checks if custom theme configuration is present.
func hasCustomThemeFunc(cfg core.ThemeConfig) bool {
	return cfg.HasCustomColors() || cfg.HasCustomFonts() || cfg.HasHeroBackground()
}

// dictFunc creates a map from key-value pairs for passing multiple values to templates.
func dictFunc(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict requires an even number of arguments")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// safeHTMLFunc marks a string as safe HTML that should not be escaped.
func safeHTMLFunc(s string) template.HTML {
	return template.HTML(s)
}

// jsonFunc marshals a value to JSON for use in templates.
func jsonFunc(v any) (template.JS, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return template.JS(bytes), nil
}

// lowerFunc converts any value to lowercase string.
func lowerFunc(v any) string {
	return strings.ToLower(fmt.Sprintf("%v", v))
}

// upperFunc converts any value to uppercase string.
func upperFunc(v any) string {
	return strings.ToUpper(fmt.Sprintf("%v", v))
}

// addFunc adds two integers.
func addFunc(a, b int) int {
	return a + b
}

// replaceFunc replaces all occurrences of old with new in input.
func replaceFunc(input, old, new string) string {
	return strings.ReplaceAll(input, old, new)
}

// ExtendFuncMap returns a new FuncMap that combines base functions with additional ones.
func ExtendFuncMap(additional template.FuncMap) template.FuncMap {
	base := BaseFuncMap()
	for k, v := range additional {
		base[k] = v
	}
	return base
}

// StatusFuncMap returns template functions specific to status pages.
func StatusFuncMap() template.FuncMap {
	return ExtendFuncMap(template.FuncMap{
		"statusColor": func(status core.ComponentStatus) string {
			switch status {
			case core.StatusOperational:
				return "green"
			case core.StatusDegraded:
				return "yellow"
			case core.StatusPartialOutage:
				return "orange"
			case core.StatusMajorOutage:
				return "red"
			case core.StatusMaintenance:
				return "blue"
			default:
				return "gray"
			}
		},
		"severityColor": func(severity core.IncidentSeverity) string {
			switch severity {
			case core.SeverityMinor:
				return "yellow"
			case core.SeverityMajor:
				return "orange"
			case core.SeverityCritical:
				return "red"
			default:
				return "gray"
			}
		},
		"incidentStatusLabel": func(status core.IncidentStatus) string {
			switch status {
			case core.IncidentInvestigating:
				return "Investigating"
			case core.IncidentIdentified:
				return "Identified"
			case core.IncidentMonitoring:
				return "Monitoring"
			case core.IncidentResolved:
				return "Resolved"
			default:
				return string(status)
			}
		},
		"maintenanceStatusLabel": func(status core.MaintenanceStatus) string {
			switch status {
			case core.MaintenanceScheduled:
				return "Scheduled"
			case core.MaintenanceInProgress:
				return "In Progress"
			case core.MaintenanceCompleted:
				return "Completed"
			default:
				return string(status)
			}
		},
		"formatTime": func(t time.Time) string {
			return t.Format("Jan 2, 2006 15:04 MST")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
	})
}

// ChangelogFuncMap returns template functions specific to changelog pages.
func ChangelogFuncMap() template.FuncMap {
	return ExtendFuncMap(template.FuncMap{
		"changeColor": func(ct core.ChangeType) string {
			return ct.Color()
		},
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatDateRFC": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
	})
}

// OpenAPIFuncMap returns template functions specific to OpenAPI pages.
func OpenAPIFuncMap() template.FuncMap {
	return ExtendFuncMap(template.FuncMap{
		"endpointID": func(endpoint *core.APIEndpoint) string {
			if endpoint.OperationID != "" {
				return endpoint.OperationID
			}
			id := endpoint.Method + "-" + endpoint.Path
			result := ""
			lastWasDash := false
			for _, r := range id {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
					result += string(r)
					lastWasDash = false
				} else if !lastWasDash {
					result += "-"
					lastWasDash = true
				}
			}
			return strings.TrimSuffix(result, "-")
		},
		"stripSpecExt": func(name string) string {
			name = strings.TrimSuffix(name, ".yaml")
			name = strings.TrimSuffix(name, ".yml")
			name = strings.TrimSuffix(name, ".json")
			return name
		},
	})
}

// PortfolioFuncMap returns template functions specific to portfolio pages.
func PortfolioFuncMap() template.FuncMap {
	return ExtendFuncMap(template.FuncMap{
		"formatDate": func(t any) string {
			if tm, ok := t.(interface{ Format(string) string }); ok {
				return tm.Format("January 2006")
			}
			return ""
		},
	})
}
