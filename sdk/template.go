package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Template is a Sendgrid transactional template.
type Template struct {
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Generation string            `json:"generation,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"`
	Versions   []TemplateVersion `json:"versions,omitempty"`
	Warnings   []string          `json:"warnings,omitempty"`
}

type Templates struct {
	Result []Template `json:"result"`
}

func parseTemplate(respBody string) (*Template, error) {
	var body Template

	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}
	if body.ID == "" {
		return nil, fmt.Errorf("response is missing ID. %s", respBody)
	}
	return &body, nil
}

func parseTemplates(respBody string) ([]Template, error) {
	var body Templates

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return body.Result, nil
}

// CreateTemplate creates a transactional template and returns it.
func (c *Client) CreateTemplate(ctx context.Context, name, generation string) (*Template, error) {
	if name == "" {
		return nil, ErrTemplateNameRequired
	}

	if generation == "" {
		generation = "dynamic"
	}

	respBody, statusCode, err := c.Post(ctx, http.MethodPost, "/templates", &Template{
		Name:       name,
		Generation: generation,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating template: %w", err)
	}
	if statusCode != http.StatusCreated {
		return nil, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingTemplate, statusCode, respBody),
		}
	}

	return parseTemplate(respBody)
}

// ReadTemplate retreives a transactional template and returns it.
func (c *Client) ReadTemplate(ctx context.Context, id string) (*Template, error) {
	if id == "" {
		return nil, ErrTemplateIDRequired
	}

	respBody, statusCode, err := c.Get(ctx, http.MethodGet, "/templates/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed reading template: %w", err)
	}
	if statusCode != http.StatusOK {
		return nil, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedGettingTemplate, statusCode, respBody),
		}
	}

	return parseTemplate(respBody)
}

func (c *Client) ReadTemplates(ctx context.Context, generation string) ([]Template, error) {
	respBody, _, err := c.Get(ctx, "GET", "/templates?page_size=200&generations="+generation)
	if err != nil {
		return nil, fmt.Errorf("failed reading template: %w", err)
	}

	return parseTemplates(respBody)
}

// UpdateTemplate edits a transactional template and returns it.
// We can't change the "generation" of a transactional template.
func (c *Client) UpdateTemplate(ctx context.Context, id, name string) (*Template, error) {
	if id == "" {
		return nil, ErrTemplateIDRequired
	}

	if name == "" {
		return nil, ErrTemplateNameRequired
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/templates/"+id, &Template{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed updating template: %w", err)
	}
	if statusCode != http.StatusOK {
		return nil, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedUpdatingTemplate, statusCode, respBody),
		}
	}

	return parseTemplate(respBody)
}

// DeleteTemplate deletes a transactional template.
func (c *Client) DeleteTemplate(ctx context.Context, id string) (bool, *RequestError) {
	if id == "" {
		return false, &RequestError{
			Err: ErrTemplateIDRequired,
		}
	}

	_, statusCode, err := c.Get(ctx, http.MethodDelete, "/templates/"+id)
	if err != nil {
		return false, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting template: %w", err),
		}
	}
	if statusCode != http.StatusNoContent && statusCode != http.StatusNotFound {
		return false, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting template: %d", statusCode),
		}
	}
	return true, nil
}
