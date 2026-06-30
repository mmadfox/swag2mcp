package types

import "github.com/mmadfox/swag2mcp/internal/spec"

type Spec struct {
	ID             string `json:"id"`
	Domain         string `json:"domain"`
	LLMTitle       string `json:"llmtitle"`
	LLMInstruction string `json:"llminstruction"`
	Stats          struct {
		Collections int `json:"collections"`
		Tags        int `json:"tags"`
		Methods     int `json:"methods"`
	}
}

type Collection struct {
	ID             string `json:"id"`
	SpecID         string `json:"specId"`
	LLMTitle       string `json:"llmtitle"`
	LLMInstruction string `json:"llminstruction"`
	Title          string `json:"title"`
	Stats          struct {
		Tags    int `json:"tags"`
		Methods int `json:"methods"`
	}
}

type Tag struct {
	ID           string `json:"id"`
	CollectionID string `json:"collectionId"`
	SpecID       string `json:"specId"`
	Name         string `json:"name"`
	Stats        struct {
		Methods int `json:"methods"`
	}
}

type Endpoint struct {
	ID           string          `json:"id"`
	TagID        string          `json:"tagId"`
	CollectionID string          `json:"collectionId"`
	SpecID       string          `json:"specId"`
	Name         string          `json:"method"`
	Path         string          `json:"path"`
	Tag          string          `json:"tag"`
	Operation    *spec.Operation `json:"-"`
}

func (e *Endpoint) SummaryOrFallback() string {
	if e.Operation.Summary != "" {
		return e.Operation.Summary
	}
	if e.Operation.Description != "" {
		return e.Operation.Description
	}
	return e.Name + " " + e.Path
}
