// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	devhub "terraform-provider-devhub/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &terradeskWorkspaceResource{}
	_ resource.ResourceWithConfigure   = &terradeskWorkspaceResource{}
	_ resource.ResourceWithImportState = &terradeskWorkspaceResource{}
)

func TerradeskWorkspaceResource() resource.Resource {
	return &terradeskWorkspaceResource{}
}

// TerradeskWorkspaceResourceModel describes the resource data model.
type terradeskWorkspaceResourceModel struct {
	Id                    types.String           `tfsdk:"id"`
	Name                  types.String           `tfsdk:"name"`
	Repository            types.String           `tfsdk:"repository"`
	InitArgs              types.String           `tfsdk:"init_args"`
	Path                  types.String           `tfsdk:"path"`
	RunPlansAutomatically types.Bool             `tfsdk:"run_plans_automatically"`
	RequiredApprovals     types.Int64            `tfsdk:"required_approvals"`
	DockerImage           types.String           `tfsdk:"docker_image"`
	CpuRequests           types.String           `tfsdk:"cpu_requests"`
	MemoryRequests        types.String           `tfsdk:"memory_requests"`
	AgentId               types.String           `tfsdk:"agent_id"`
	WorkloadIdentity      *workloadIdentityModel `tfsdk:"workload_identity"`
	EnvVars               []envVarModel          `tfsdk:"env_vars"`
	Secrets               []secretModel          `tfsdk:"secrets"`
}

type workloadIdentityModel struct {
	Enabled             types.Bool   `tfsdk:"enabled"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	Provider            types.String `tfsdk:"provider"`
}

type envVarModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type secretModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type terradeskWorkspaceResource struct {
	client *devhub.Client
}

func (r *terradeskWorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terradesk_workspace"
}

func (r *terradeskWorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "TerraDesk workspace resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Workspace id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name for users to use to identity the workspace.",
				Required:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "Which GitHub repository should be used in the format `owner/name`. Must have the GitHub integration enabled.",
				Required:            true,
			},
			"init_args": schema.StringAttribute{
				MarkdownDescription: "Args to pass to the init command.",
				Optional:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The file path of here the workspace is located in the provided GitHub repository. Defaults to the root of the repository.",
				Optional:            true,
			},
			"run_plans_automatically": schema.BoolAttribute{
				MarkdownDescription: "Whether to run plans automatically for PRs and pushes. Make sure to consider who can push to your GitHub repository if you have this setting on as it could grant sensitive access.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"required_approvals": schema.Int64Attribute{
				MarkdownDescription: "Specify how many reviews are required to apply plans.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"docker_image": schema.StringAttribute{
				MarkdownDescription: "The docker image to use for running commands, for example: hashicorp/terraform:1.10.",
				Required:            true,
			},
			"cpu_requests": schema.StringAttribute{
				MarkdownDescription: "How much cpu should be requested for the pod scheduled by the job, see kubernetes docs for allowable values.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("100m"),
			},
			"memory_requests": schema.StringAttribute{
				MarkdownDescription: "How much memory should be requested for the pod scheduled by the job, see kubernetes docs for allowable values.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("512M"),
			},
			"agent_id": schema.StringAttribute{
				MarkdownDescription: "The agent id for the database.",
				Optional:            true,
			},
			"workload_identity": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether to enable Google workload identity to authenticate with Google services.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"service_account_email": schema.StringAttribute{
						MarkdownDescription: "The service account email to use for workload identity.",
						Optional:            true,
					},
					"provider": schema.StringAttribute{
						MarkdownDescription: "The workload identity provider to use: `projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL}/providers/${PROVIDER}`",
						Optional:            true,
					},
				},
			},
			"env_vars": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				Default: listdefault.StaticValue(
					types.ListValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":    types.StringType,
								"name":  types.StringType,
								"value": types.StringType,
							},
						},
						[]attr.Value{},
					),
				),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Env var id.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name to use for the env var.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Env var value.",
							Required:            true,
						},
					},
				},
			},
			"secrets": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				Default: listdefault.StaticValue(
					types.ListValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":    types.StringType,
								"name":  types.StringType,
								"value": types.StringType,
							},
						},
						[]attr.Value{},
					),
				),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Secret id.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Secret name.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Secret value.",
							Required:            true,
							Sensitive:           true,
						},
					},
				},
			},
		},
	}
}

func (r *terradeskWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan terradeskWorkspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var envVars []devhub.EnvVar
	for _, envVar := range plan.EnvVars {
		envVars = append(envVars, devhub.EnvVar{
			Name:  envVar.Name.ValueString(),
			Value: envVar.Value.ValueString(),
		})
	}

	var secrets []devhub.Secret
	for _, secret := range plan.Secrets {
		secrets = append(secrets, devhub.Secret{
			Name:  secret.Name.ValueString(),
			Value: secret.Value.ValueString(),
		})
	}

	input := devhub.TerradeskWorkspace{
		Name:                  plan.Name.ValueString(),
		Repository:            plan.Repository.ValueString(),
		InitArgs:              plan.InitArgs.ValueString(),
		Path:                  plan.Path.ValueString(),
		RunPlansAutomatically: plan.RunPlansAutomatically.ValueBool(),
		RequiredApprovals:     int(plan.RequiredApprovals.ValueInt64()),
		DockerImage:           plan.DockerImage.ValueString(),
		CpuRequests:           plan.CpuRequests.ValueString(),
		MemoryRequests:        plan.MemoryRequests.ValueString(),
		AgentId:               plan.AgentId.ValueString(),
		EnvVars:               envVars,
		Secrets:               secrets,
	}

	if plan.WorkloadIdentity != nil {
		input.WorkloadIdentity = &devhub.WorkloadIdentity{
			Enabled:             plan.WorkloadIdentity.Enabled.ValueBool(),
			ServiceAccountEmail: plan.WorkloadIdentity.ServiceAccountEmail.ValueString(),
			Provider:            plan.WorkloadIdentity.Provider.ValueString(),
		}
	}

	workspace, err := r.client.CreateWorkspace(input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workspace",
			"Could not create workspace, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(workspace.Id)

	for index, envVar := range workspace.EnvVars {
		plan.EnvVars[index].Id = types.StringValue(envVar.Id)
	}

	for index, envVar := range workspace.Secrets {
		plan.Secrets[index].Id = types.StringValue(envVar.Id)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *terradeskWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state terradeskWorkspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, err := r.client.GetWorkspace(state.Id.ValueString())

	if err != nil && err.Error() == "not found" {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading workspace",
			"Could not read workspace "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(workspace.Name)
	state.Repository = types.StringValue(workspace.Repository)

	if workspace.InitArgs != "" {
		state.InitArgs = types.StringValue(workspace.InitArgs)
	}

	if workspace.Path != "" {
		state.Path = types.StringValue(workspace.Path)
	}

	state.RunPlansAutomatically = types.BoolValue(workspace.RunPlansAutomatically)
	state.RequiredApprovals = types.Int64Value(int64(workspace.RequiredApprovals))
	state.DockerImage = types.StringValue(workspace.DockerImage)
	state.CpuRequests = types.StringValue(workspace.CpuRequests)
	state.MemoryRequests = types.StringValue(workspace.MemoryRequests)

	if workspace.WorkloadIdentity == nil {
		state.WorkloadIdentity = nil
	} else {
		state.WorkloadIdentity = &workloadIdentityModel{
			Enabled:             types.BoolValue(workspace.WorkloadIdentity.Enabled),
			ServiceAccountEmail: types.StringValue(workspace.WorkloadIdentity.ServiceAccountEmail),
			Provider:            types.StringValue(workspace.WorkloadIdentity.Provider),
		}
	}

	if workspace.AgentId != "" {
		state.AgentId = types.StringValue(workspace.AgentId)
	} else {
		state.AgentId = types.StringNull()
	}

	if state.EnvVars == nil || len(state.EnvVars) != len(workspace.EnvVars) {
		state.EnvVars = make([]envVarModel, len(workspace.EnvVars))
	}

	for index, envVar := range workspace.EnvVars {
		state.EnvVars[index].Id = types.StringValue(envVar.Id)
		state.EnvVars[index].Name = types.StringValue(envVar.Name)
		state.EnvVars[index].Value = types.StringValue(envVar.Value)
	}

	if state.Secrets == nil || len(state.Secrets) != len(workspace.Secrets) {
		state.Secrets = make([]secretModel, len(workspace.Secrets))
	}

	for index, secret := range workspace.Secrets {
		state.Secrets[index].Id = types.StringValue(secret.Id)
		state.Secrets[index].Name = types.StringValue(secret.Name)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *terradeskWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan terradeskWorkspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var envVars []devhub.EnvVar
	for _, envVar := range plan.EnvVars {
		envVars = append(envVars, devhub.EnvVar{
			Id:    envVar.Id.ValueString(),
			Name:  envVar.Name.ValueString(),
			Value: envVar.Value.ValueString(),
		})
	}

	var secrets []devhub.Secret
	for _, secret := range plan.Secrets {
		secrets = append(secrets, devhub.Secret{
			Id:    secret.Id.ValueString(),
			Name:  secret.Name.ValueString(),
			Value: secret.Value.ValueString(),
		})
	}

	input := devhub.TerradeskWorkspace{
		Name:                  plan.Name.ValueString(),
		Repository:            plan.Repository.ValueString(),
		InitArgs:              plan.InitArgs.ValueString(),
		Path:                  plan.Path.ValueString(),
		RunPlansAutomatically: plan.RunPlansAutomatically.ValueBool(),
		RequiredApprovals:     int(plan.RequiredApprovals.ValueInt64()),
		DockerImage:           plan.DockerImage.ValueString(),
		CpuRequests:           plan.CpuRequests.ValueString(),
		MemoryRequests:        plan.MemoryRequests.ValueString(),
		AgentId:               plan.AgentId.ValueString(),
		EnvVars:               envVars,
		Secrets:               secrets,
	}

	if plan.WorkloadIdentity != nil {
		input.WorkloadIdentity = &devhub.WorkloadIdentity{
			Enabled:             plan.WorkloadIdentity.Enabled.ValueBool(),
			ServiceAccountEmail: plan.WorkloadIdentity.ServiceAccountEmail.ValueString(),
			Provider:            plan.WorkloadIdentity.Provider.ValueString(),
		}
	}

	_, err := r.client.UpdateWorkspace(plan.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workspace",
			"Could not update workspace, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *terradeskWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state terradeskWorkspaceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkspace(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workspace",
			"Could not delete workspace, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *terradeskWorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*devhub.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *terradeskWorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
