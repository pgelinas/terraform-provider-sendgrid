package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TemplateVersion is a Sendgrid transactional template version.
type TemplateVersion struct {
	ID                   string   `json:"id,omitempty"`
	TemplateID           string   `json:"template_id,omitempty"`   //nolint:tagliatelle
	UpdatedAt            string   `json:"updated_at,omitempty"`    //nolint:tagliatelle
	ThumbnailURL         string   `json:"thumbnail_url,omitempty"` //nolint:tagliatelle
	Warnings             []string `json:"warnings,omitempty"`
	Active               int      `json:"active,omitempty"`
	Name                 string   `json:"name,omitempty"`
	HTMLContent          string   `json:"html_content,omitempty"`           //nolint:tagliatelle
	PlainContent         string   `json:"plain_content,omitempty"`          //nolint:tagliatelle
	GeneratePlainContent bool     `json:"generate_plain_content,omitempty"` //nolint:tagliatelle
	Subject              string   `json:"subject,omitempty"`
	Editor               string   `json:"editor,omitempty"`
	TestData             string   `json:"test_data,omitempty"` //nolint:tagliatelle
}

func parseTemplateVersion(respBody string) (*TemplateVersion, error) {
	var body TemplateVersion

	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, fmt.Errorf("failed parsing template version: %w", err)
	}
	if body.ID == "" {
		return nil, fmt.Errorf("response is missing ID. %s", respBody)
	}
	return &body, nil
}

// CreateTemplateVersion creates a new version of a transactional template and returns it.
func (c *Client) CreateTemplateVersion(ctx context.Context, t TemplateVersion) (*TemplateVersion, error) {
	if t.TemplateID == "" {
		return nil, ErrTemplateIDRequired
	}

	if t.Name == "" {
		return nil, ErrTemplateVersionNameRequired
	}

	if t.Subject == "" {
		return nil, ErrTemplateVersionSubjectRequired
	}

	respBody, statusCode, err := c.Post(ctx, "POST", "/templates/"+t.TemplateID+"/versions", t)
	if err != nil {
		return nil, fmt.Errorf("failed creating template version: %w", err)
	}

	if statusCode >= 500 {
		return nil, fmt.Errorf("failed creating the template version. Status: %d, Body: %s", statusCode, respBody)
	}

	return parseTemplateVersion(respBody)
}

// ReadTemplateVersion retreives a version of a transactional template and returns it.
func (c *Client) ReadTemplateVersion(ctx context.Context, templateID, id string) (*TemplateVersion, error) {
	if templateID == "" {
		return nil, ErrTemplateIDRequired
	}

	if id == "" {
		return nil, ErrTemplateVersionIDRequired
	}

	respBody, _, err := c.Get(ctx, "GET", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed reading template version: %w", err)
	}

	return parseTemplateVersion(respBody)
}

// UpdateTemplateVersion edits a version of a transactional template and returns it.
func (c *Client) UpdateTemplateVersion(ctx context.Context, t TemplateVersion) (*TemplateVersion, error) {
	if t.ID == "" {
		return nil, ErrTemplateVersionIDRequired
	}

	if t.TemplateID == "" {
		return nil, ErrTemplateIDRequired
	}

	respBody, _, err := c.Post(ctx, "PATCH", "/templates/"+t.TemplateID+"/versions/"+t.ID, t)
	if err != nil {
		return nil, fmt.Errorf("failed updating template version: %w", err)
	}

	return parseTemplateVersion(respBody)
}

// DeleteTemplateVersion deletes a version of a transactional template.
func (c *Client) DeleteTemplateVersion(ctx context.Context, templateID, id string) (bool, error) {
	if templateID == "" {
		return false, ErrTemplateIDRequired
	}

	_, statusCode, err := c.Get(ctx, "DELETE", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return false, fmt.Errorf("failed deleting template version: %w", err)
	}
	if statusCode != http.StatusNoContent && statusCode != http.StatusNotFound {
		return false, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting template version: %d", statusCode),
		}
	}

	return true, nil
}
