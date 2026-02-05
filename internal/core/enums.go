package core

// HTTPMethod represents HTTP methods for API endpoints
type HTTPMethod string

const (
	HTTPGet     HTTPMethod = "GET"
	HTTPPost    HTTPMethod = "POST"
	HTTPPut     HTTPMethod = "PUT"
	HTTPPatch   HTTPMethod = "PATCH"
	HTTPDelete  HTTPMethod = "DELETE"
	HTTPHead    HTTPMethod = "HEAD"
	HTTPOptions HTTPMethod = "OPTIONS"
	HTTPTrace   HTTPMethod = "TRACE"
)

// AllHTTPMethods returns all valid HTTP methods in standard order
func AllHTTPMethods() []HTTPMethod {
	return []HTTPMethod{
		HTTPGet,
		HTTPPost,
		HTTPPut,
		HTTPPatch,
		HTTPDelete,
		HTTPHead,
		HTTPOptions,
		HTTPTrace,
	}
}

// IsValid checks if the HTTP method is valid
func (m HTTPMethod) IsValid() bool {
	switch m {
	case HTTPGet, HTTPPost, HTTPPut, HTTPPatch, HTTPDelete, HTTPHead, HTTPOptions, HTTPTrace:
		return true
	}
	return false
}

// ParameterLocation represents where a parameter is located in an API request
type ParameterLocation string

const (
	ParamInPath   ParameterLocation = "path"
	ParamInQuery  ParameterLocation = "query"
	ParamInHeader ParameterLocation = "header"
	ParamInCookie ParameterLocation = "cookie"
)

// IsValid checks if the parameter location is valid
func (p ParameterLocation) IsValid() bool {
	switch p {
	case ParamInPath, ParamInQuery, ParamInHeader, ParamInCookie:
		return true
	}
	return false
}

// SecuritySchemeType represents the type of security scheme
type SecuritySchemeType string

const (
	SecurityAPIKey         SecuritySchemeType = "apiKey"
	SecurityHTTP           SecuritySchemeType = "http"
	SecurityOAuth2         SecuritySchemeType = "oauth2"
	SecurityOpenIDConnect  SecuritySchemeType = "openIdConnect"
)

// IsValid checks if the security scheme type is valid
func (s SecuritySchemeType) IsValid() bool {
	switch s {
	case SecurityAPIKey, SecurityHTTP, SecurityOAuth2, SecurityOpenIDConnect:
		return true
	}
	return false
}

// SchemaType represents JSON schema types
type SchemaType string

const (
	SchemaString  SchemaType = "string"
	SchemaNumber  SchemaType = "number"
	SchemaInteger SchemaType = "integer"
	SchemaBoolean SchemaType = "boolean"
	SchemaArray   SchemaType = "array"
	SchemaObject  SchemaType = "object"
	SchemaNull    SchemaType = "null"
)

// IsValid checks if the schema type is valid
func (s SchemaType) IsValid() bool {
	switch s {
	case SchemaString, SchemaNumber, SchemaInteger, SchemaBoolean, SchemaArray, SchemaObject, SchemaNull:
		return true
	}
	return false
}

// ComponentStatus represents the operational status of a service component
type ComponentStatus string

const (
	StatusOperational   ComponentStatus = "operational"
	StatusDegraded      ComponentStatus = "degraded"
	StatusPartialOutage ComponentStatus = "partial_outage"
	StatusMajorOutage   ComponentStatus = "major_outage"
	StatusMaintenance   ComponentStatus = "maintenance"
)

// IsValid checks if the component status is valid
func (c ComponentStatus) IsValid() bool {
	switch c {
	case StatusOperational, StatusDegraded, StatusPartialOutage, StatusMajorOutage, StatusMaintenance:
		return true
	}
	return false
}

// Severity returns a numeric severity for sorting (higher = worse)
func (c ComponentStatus) Severity() int {
	switch c {
	case StatusOperational:
		return 0
	case StatusMaintenance:
		return 1
	case StatusDegraded:
		return 2
	case StatusPartialOutage:
		return 3
	case StatusMajorOutage:
		return 4
	}
	return 0
}

// IncidentStatus represents the status of an incident
type IncidentStatus string

const (
	IncidentInvestigating IncidentStatus = "investigating"
	IncidentIdentified    IncidentStatus = "identified"
	IncidentMonitoring    IncidentStatus = "monitoring"
	IncidentResolved      IncidentStatus = "resolved"
)

// IsValid checks if the incident status is valid
func (i IncidentStatus) IsValid() bool {
	switch i {
	case IncidentInvestigating, IncidentIdentified, IncidentMonitoring, IncidentResolved:
		return true
	}
	return false
}

// IsActive returns true if the incident is not yet resolved
func (i IncidentStatus) IsActive() bool {
	return i != IncidentResolved
}

// IncidentSeverity represents the severity level of an incident
type IncidentSeverity string

const (
	SeverityMinor    IncidentSeverity = "minor"
	SeverityMajor    IncidentSeverity = "major"
	SeverityCritical IncidentSeverity = "critical"
)

// IsValid checks if the incident severity is valid
func (s IncidentSeverity) IsValid() bool {
	switch s {
	case SeverityMinor, SeverityMajor, SeverityCritical:
		return true
	}
	return false
}

// Level returns a numeric severity for sorting (higher = worse)
func (s IncidentSeverity) Level() int {
	switch s {
	case SeverityMinor:
		return 1
	case SeverityMajor:
		return 2
	case SeverityCritical:
		return 3
	}
	return 0
}

// MaintenanceStatus represents the status of scheduled maintenance
type MaintenanceStatus string

const (
	MaintenanceScheduled  MaintenanceStatus = "scheduled"
	MaintenanceInProgress MaintenanceStatus = "in_progress"
	MaintenanceCompleted  MaintenanceStatus = "completed"
)

// IsValid checks if the maintenance status is valid
func (m MaintenanceStatus) IsValid() bool {
	switch m {
	case MaintenanceScheduled, MaintenanceInProgress, MaintenanceCompleted:
		return true
	}
	return false
}
