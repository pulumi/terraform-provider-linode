// +build ignore

/**
 * Using this template:
 * - Copy resource_linode_template.go and resource_linode_template_test.go
 *   - Remove "// +build ignore"
 *   - Replace "Template" with Linode Resource Name
 *   - Replace "template" with Linode resource name
 * - Add Resource to provider.go
 */
package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeTemplateCreate,
		Read:   resourceLinodeTemplateRead,
		Update: resourceLinodeTemplateUpdate,
		Delete: resourceLinodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode Template.",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeInt,
				Description: "The status of the template, indicating the current readiness state.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Template ID %s as int: %s", d.Id(), err)
	}

	template, err := client.GetTemplate(context.Background(), int(id))

	if err != nil {
		return fmt.Errorf("Error finding the specified Linode Template: %s", err)
	}

	d.Set("label", template.Label)
	d.Set("status", template.Status)

	return nil
}

func resourceLinodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(*ProviderMeta).Client
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Template")
	}

	createOpts := linodego.TemplateCreateOptions{
		Label: d.Get("label").(string),
	}
	template, err := client.CreateTemplate(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode Template: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", template.ID))
	d.Set("label", template.Label)

	return resourceLinodeTemplateRead(d, meta)
}

func resourceLinodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Template id %s as int: %s", d.Id(), err)
	}

	template, err := client.GetTemplate(int(id))
	if err != nil {
		return fmt.Errorf("Error fetching data about the current Linode Template: %s", err)
	}

	if d.HasChange("label") {
		if template, err = client.RenameTemplate(context.Background(), template.ID, d.Get("label").(string)); err != nil {
			return err
		}
		d.Set("label", template.Label)
	}

	return resourceLinodeTemplateRead(d, meta)
}

func resourceLinodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Template id %s as int", d.Id())
	}
	err = client.DeleteTemplate(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode Template %d: %s", id, err)
	}
	return nil
}
