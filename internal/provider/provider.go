package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &InstanceProvisioner{}

type InstanceProvisioner struct {
	version string
}

type InstanceProvisionerModel struct {
	SNSTopic types.String `tfsdk:"sns_topic"`
	Region   types.String `tfsdk:"region"`
}

type InstanceProvisionerConfig struct {
	SNSTopic string
	Region   string
}

func (i *InstanceProvisioner) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "provisioner"
	response.Version = i.version
}

func (i *InstanceProvisioner) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "The provisioner Terraform provider allows you to send events to a SNS topic when provisioning instances.",
		Attributes: map[string]schema.Attribute{
			"sns_topic": schema.StringAttribute{
				MarkdownDescription: "SNS topic to send provisioning events to.",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "AWS region.",
				Optional:            true,
			},
		},
	}
}

func (i *InstanceProvisioner) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var data InstanceProvisionerModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	config := &InstanceProvisionerConfig{
		SNSTopic: data.SNSTopic.ValueString(),
		Region:   data.Region.ValueString(),
	}

	response.ResourceData = config
}

func (i *InstanceProvisioner) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (i *InstanceProvisioner) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProvisionResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InstanceProvisioner{
			version: version,
		}
	}
}
