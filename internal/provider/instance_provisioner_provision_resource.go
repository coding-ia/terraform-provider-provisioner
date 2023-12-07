package provider

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ProvisionResource{}
var _ resource.ResourceWithConfigure = &ProvisionResource{}

type ProvisionResource struct {
	snsClient *sns.SNS
	snsTopic  string
}

type ProvisionerMessage struct {
	Action     string `json:"action"`
	ID         string `json:"id"`
	Hostname   string `json:"host"`
	PrivateIP  string `json:"private_ip"`
	InstanceId string `json:"instance_id"`
}

func (p *ProvisionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	data := request.ProviderData.(*InstanceProvisionerConfig)

	awsConfig := aws.NewConfig()
	if data.Region != "" {
		awsConfig = aws.NewConfig().WithRegion(data.Region)
	}
	sess := session.Must(session.NewSession())
	snsClient := sns.New(sess, awsConfig)

	p.snsClient = snsClient
	p.snsTopic = data.SNSTopic
}

func NewProvisionResource() resource.Resource {
	return &ProvisionResource{}
}

type ProvisionResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	InstanceID types.String `tfsdk:"instance_id"`
	PrivateIP  types.String `tfsdk:"private_ip"`
}

func (p *ProvisionResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_provision"
}

func (p *ProvisionResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Downloads a file from a website using the supplied URL.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Instance name",
				Required:            true,
			},
			"instance_id": schema.StringAttribute{
				MarkdownDescription: "Instance ID",
				Required:            true,
			},
			"private_ip": schema.StringAttribute{
				MarkdownDescription: "Instance Private IP",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
		},
	}
}

func (p *ProvisionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *ProvisionResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	// Need to give resource an ID, but there is none.  Just randomize one.
	newUUID, err := uuid.GenerateUUID()
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	data.Id = types.StringValue(newUUID)

	message := &ProvisionerMessage{
		ID:         newUUID,
		Action:     "create",
		InstanceId: data.InstanceID.ValueString(),
		Hostname:   data.Name.ValueString(),
		PrivateIP:  data.PrivateIP.ValueString(),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	_, err = p.snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageBody)),
		TopicArn: aws.String(p.snsTopic),
	})
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (p *ProvisionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *ProvisionResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (p *ProvisionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data *ProvisionResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	message := &ProvisionerMessage{
		ID:         data.Id.ValueString(),
		Action:     "update",
		InstanceId: data.InstanceID.ValueString(),
		Hostname:   data.Name.ValueString(),
		PrivateIP:  data.PrivateIP.ValueString(),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	_, err = p.snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageBody)),
		TopicArn: aws.String(p.snsTopic),
	})
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (p *ProvisionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *ProvisionResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	message := &ProvisionerMessage{
		ID:         data.Id.ValueString(),
		Action:     "delete",
		InstanceId: data.InstanceID.ValueString(),
		Hostname:   data.Name.ValueString(),
		PrivateIP:  data.PrivateIP.ValueString(),
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	_, err = p.snsClient.Publish(&sns.PublishInput{
		Message:  aws.String(string(messageBody)),
		TopicArn: aws.String(p.snsTopic),
	})
	if err != nil {
		response.Diagnostics.AddError("Create error", err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
