package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOAS30(t *testing.T) {
	specContent := `openapi: "3.0.3"
info:
  title: Test API 3.0
  version: "1.0.0"
paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
      tags: [Users]
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
    post:
      operationId: createUser
      summary: Create user
      tags: [Users]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        "201":
          description: Created
components:
  schemas:
    User:
      type: object
      required: [name, email]
      properties:
        name:
          type: string
        email:
          type: string
          format: email
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
`
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "openapi30.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewOpenAPIParser(filepath.Join(tmpDir, "cache"))
	spec, err := p.ParseFile(specPath)
	if err != nil {
		t.Fatalf("OAS 3.0 parse failed: %v", err)
	}

	if spec.Title != "Test API 3.0" {
		t.Errorf("Title = %q, want %q", spec.Title, "Test API 3.0")
	}
	if spec.OpenAPIVersion != "3.0.3" {
		t.Errorf("OpenAPIVersion = %q, want %q", spec.OpenAPIVersion, "3.0.3")
	}
	if len(spec.Endpoints) != 2 {
		t.Errorf("Endpoints = %d, want 2", len(spec.Endpoints))
	}
	if len(spec.Schemas) != 1 {
		t.Errorf("Schemas = %d, want 1", len(spec.Schemas))
	}
	if len(spec.SecuritySchemes) != 1 {
		t.Errorf("SecuritySchemes = %d, want 1", len(spec.SecuritySchemes))
	}
	if spec.Schemas["User"] == nil {
		t.Error("User schema missing")
	} else if len(spec.Schemas["User"].Properties) != 2 {
		t.Errorf("User properties = %d, want 2", len(spec.Schemas["User"].Properties))
	}
}

func TestParseOAS31(t *testing.T) {
	specContent := `openapi: "3.1.0"
info:
  title: Test API 3.1
  version: "2.0.0"
paths:
  /items:
    get:
      operationId: listItems
      summary: List items
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Item'
components:
  schemas:
    Item:
      type: object
      properties:
        id:
          type: integer
        name:
          type:
            - "string"
            - "null"
        metadata:
          type: object
`
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "openapi31.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewOpenAPIParser(filepath.Join(tmpDir, "cache"))
	spec, err := p.ParseFile(specPath)
	if err != nil {
		t.Fatalf("OAS 3.1 parse failed: %v", err)
	}

	if spec.Title != "Test API 3.1" {
		t.Errorf("Title = %q, want %q", spec.Title, "Test API 3.1")
	}
	if spec.OpenAPIVersion != "3.1.0" {
		t.Errorf("OpenAPIVersion = %q, want %q", spec.OpenAPIVersion, "3.1.0")
	}
	if len(spec.Endpoints) != 1 {
		t.Errorf("Endpoints = %d, want 1", len(spec.Endpoints))
	}
	if len(spec.Schemas) != 1 {
		t.Errorf("Schemas = %d, want 1", len(spec.Schemas))
	}

	item := spec.Schemas["Item"]
	if item == nil {
		t.Fatal("Item schema missing")
	}
	// OAS 3.1 type array: name should have type "string" (first element)
	if name, ok := item.Properties["name"]; ok {
		if name.Type != "string" {
			t.Errorf("name.Type = %q, want %q", name.Type, "string")
		}
	} else {
		t.Error("name property missing from Item schema")
	}
}

func TestRefPreservation(t *testing.T) {
	specContent := `openapi: "3.0.3"
info:
  title: Ref Test
  version: "1.0.0"
paths:
  /users:
    get:
      operationId: listUsers
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
components:
  schemas:
    User:
      type: object
      properties:
        name:
          type: string
`
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "ref_test.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	p := NewOpenAPIParser(filepath.Join(tmpDir, "cache"))
	spec, err := p.ParseFile(specPath)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// The response schema should be a $ref, not inlined
	resp := spec.Endpoints[0].Responses["200"]
	if resp == nil {
		t.Fatal("200 response missing")
	}
	mt := resp.Content["application/json"]
	if mt == nil {
		t.Fatal("application/json content missing")
	}
	if mt.Schema == nil {
		t.Fatal("response schema missing")
	}
	if mt.Schema.Ref == "" {
		t.Error("expected $ref to be preserved, got empty ref")
	}
	if mt.Schema.Ref != "#/components/schemas/User" {
		t.Errorf("Ref = %q, want %q", mt.Schema.Ref, "#/components/schemas/User")
	}
}
