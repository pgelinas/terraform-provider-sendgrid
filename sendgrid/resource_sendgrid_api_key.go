/*
Provide a resource to manage an API key.
Example Usage
```hcl
resource "sendgrid_api_key" "api_key" {
	name   = "my-api-key"
	scopes = [
		"mail.send",
		"sender_verification_eligible",
	]
}
```
Import
An API key can be imported, e.g.
```hcl
$ terraform import sendgrid_api_key.api_key apiKeyID
```
*/
package sendgrid

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sendgrid "github.com/taharah/terraform-provider-sendgrid/sdk"
)

func resourceSendgridAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridAPIKeyCreate,
		ReadContext:   resourceSendgridAPIKeyRead,
		UpdateContext: resourceSendgridAPIKeyUpdate,
		DeleteContext: resourceSendgridAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name you will use to describe this API Key.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, maxStringLength),
			},
			"scopes": {
				Type:        schema.TypeSet,
				Description: "The individual permissions that you are giving to this API Key.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return scopeInScopes([]string{"2fa_required", "sender_verification_eligible", "sender_verification_legacy"}, old)
					},
				},
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The API key created by the API.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceSendgridAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)
	req := &sendgrid.APIKey{}

	if v, ok := d.GetOk("name"); ok {
		req.Name = v.(string)
		log.Printf("[DEBUG] create api_key name: %v", req.Name)
	}

	if v, ok := d.GetOk("scopes"); ok {
		var scopes []string
		vl := v.(*schema.Set).List()

		for _, l := range vl {
			scopes = append(scopes, l.(string))
		}

		req.Scopes = scopes
		log.Printf("[DEBUG] create api_key scopes: %v", req.Scopes)
	}

	log.Printf("[DEBUG] creating API Key: %s", req.Name)
	apiKeyStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateAPIKey(req)
	})
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] created API Key: %s", req.Name)

	apiKey := apiKeyStruct.(*sendgrid.APIKey)
	d.SetId(apiKey.ID)
	d.Set("api_key", apiKey.APIKey)

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	apiKey, err := c.ReadAPIKey(d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	d.Set("name", apiKey.Name)
	d.Set("scopes", apiKey.Scopes)
	return nil
}

func resourceSendgridAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)
	req := &sendgrid.APIKey{}

	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
		log.Printf("[DEBUG] update api_key name: %v", req.Name)
	}

	if d.HasChange("scopes") {
		var scopes []string
		vl := d.Get("scopes").(*schema.Set).List()

		for _, l := range vl {
			scopes = append(scopes, l.(string))
		}

		req.Scopes = scopes
		log.Printf("[DEBUG] update api_key scopes: %v", req.Scopes)
	}

	log.Printf("[DEBUG] updating API Key: %s", req.Name)
	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.UpdateAPIKey(d.Id(), req)
	})
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] updated API Key: %s", req.Name)

	return resourceSendgridAPIKeyRead(ctx, d, m)
}

func resourceSendgridAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	if _, err := c.DeleteAPIKey(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func scopeInScopes(scopes []string, scope string) bool {
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}
	return false
}
