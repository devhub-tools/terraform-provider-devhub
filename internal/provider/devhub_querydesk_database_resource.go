// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	devhub "terraform-provider-devhub/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &databaseResource{}
	_ resource.ResourceWithConfigure   = &databaseResource{}
	_ resource.ResourceWithImportState = &databaseResource{}
)

func DatabaseResource() resource.Resource {
	return &databaseResource{}
}

// DatabaseResourceModel describes the resource data model.
type databaseResourceModel struct {
	Id             types.String              `tfsdk:"id"`
	Name           types.String              `tfsdk:"name"`
	Adapter        types.String              `tfsdk:"adapter"`
	Hostname       types.String              `tfsdk:"hostname"`
	Port           types.Int64               `tfsdk:"port"`
	Database       types.String              `tfsdk:"database"`
	Ssl            types.Bool                `tfsdk:"ssl"`
	Cacertfile     types.String              `tfsdk:"cacertfile"`
	Keyfile        types.String              `tfsdk:"keyfile"`
	Certfile       types.String              `tfsdk:"certfile"`
	RestrictAccess types.Bool                `tfsdk:"restrict_access"`
	Group          types.String              `tfsdk:"group"`
	SlackChannel   types.String              `tfsdk:"slack_channel"`
	AgentId        types.String              `tfsdk:"agent_id"`
	Credentials    []databaseCredentialModel `tfsdk:"credentials"`
	CredentialIds  types.Map                 `tfsdk:"credential_ids"`
}

type databaseCredentialModel struct {
	Id                types.String `tfsdk:"id"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Hostname          types.String `tfsdk:"hostname"`
	ReviewsRequired   types.Int64  `tfsdk:"reviews_required"`
	DefaultCredential types.Bool   `tfsdk:"default_credential"`
}

type databaseResource struct {
	client *devhub.Client
}

func (r *databaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_querydesk_database"
}

func (r *databaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Database resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Database id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name for users to use to identity the database.",
				Required:            true,
			},
			"adapter": schema.StringAttribute{
				MarkdownDescription: "The adapter to use to establish the connection. Currently only `POSTGRES`, `MYSQL` and `CLICKHOUSE` are supported.",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port to connect to the database on, if not specified the default port for the database type will be used.",
				Optional:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "The name of the database to connect to.",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname for connecting to the database, either an ip or url.",
				Required:            true,
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` to turn on ssl connections for this database.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cacertfile": schema.StringAttribute{
				MarkdownDescription: "The server ca cert to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"keyfile": schema.StringAttribute{
				MarkdownDescription: "The client key to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"certfile": schema.StringAttribute{
				MarkdownDescription: "The client cert to use with ssl connections, `ssl` must be set to `true`.",
				Optional:            true,
				Sensitive:           true,
			},
			"restrict_access": schema.BoolAttribute{
				MarkdownDescription: "Whether access to this databases should be explicitly granted to users or if any authenticated user can access it.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The group this database belongs to, used for UI grouping.",
				Optional:            true,
			},
			"slack_channel": schema.StringAttribute{
				MarkdownDescription: "The slack channel to send query request notifications to.",
				Optional:            true,
			},
			"agent_id": schema.StringAttribute{
				MarkdownDescription: "The agent id for the database.",
				Optional:            true,
			},
			"credential_ids": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "A map of credential IDs by username.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Credential id.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "The username to use for connecting to the database.",
							Required:            true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "The username to use for connecting to the database.",
							Required:            true,
							Sensitive:           true,
						},
						"hostname": schema.StringAttribute{
							MarkdownDescription: "The hostname to use for connecting to the database when using this credential (overrides the default hostname).",
							Optional:            true,
						},
						"reviews_required": schema.Int64Attribute{
							MarkdownDescription: "The number of reviews required before a query can be executed.",
							Required:            true,
						},
						"default_credential": schema.BoolAttribute{
							MarkdownDescription: "Whether this is the default credential for the database.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan databaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var credentials []devhub.DatabaseCredential
	for _, credential := range plan.Credentials {
		credentials = append(credentials, devhub.DatabaseCredential{
			Username:          credential.Username.ValueString(),
			Password:          credential.Password.ValueString(),
			Hostname:          credential.Hostname.ValueString(),
			ReviewsRequired:   int(credential.ReviewsRequired.ValueInt64()),
			DefaultCredential: credential.DefaultCredential.ValueBool(),
		})
	}

	var port *int64
	if !plan.Port.IsNull() {
		portValue := plan.Port.ValueInt64()
		port = &portValue
	}

	input := devhub.Database{
		Name:           plan.Name.ValueString(),
		Adapter:        strings.ToLower(plan.Adapter.ValueString()),
		Hostname:       plan.Hostname.ValueString(),
		Port:           port,
		Database:       plan.Database.ValueString(),
		Ssl:            plan.Ssl.ValueBool(),
		Cacertfile:     plan.Cacertfile.ValueString(),
		Keyfile:        plan.Keyfile.ValueString(),
		Certfile:       plan.Certfile.ValueString(),
		RestrictAccess: plan.RestrictAccess.ValueBool(),
		Group:          plan.Group.ValueString(),
		SlackChannel:   plan.SlackChannel.ValueString(),
		AgentId:        plan.AgentId.ValueString(),
		Credentials:    credentials,
	}

	database, err := r.client.CreateDatabase(input)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(database.Id)

	credentialIds := make(map[string]attr.Value)

	for index, credential := range database.Credentials {
		plan.Credentials[index].Id = types.StringValue(credential.Id)
		plan.Credentials[index].DefaultCredential = types.BoolValue(credential.DefaultCredential)

		credentialIds[credential.Username] = types.StringValue(credential.Id)
	}

	plan.CredentialIds = types.MapValueMust(types.StringType, credentialIds)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	database, err := r.client.GetDatabase(state.Id.ValueString())

	if err != nil && err.Error() == "not found" {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database",
			"Could not read Database "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(database.Name)
	state.Adapter = types.StringValue(strings.ToUpper(database.Adapter))
	state.Hostname = types.StringValue(database.Hostname)
	state.Database = types.StringValue(database.Database)
	state.Ssl = types.BoolValue(database.Ssl)
	state.RestrictAccess = types.BoolValue(database.RestrictAccess)

	state.Port = types.Int64Null()
	state.Group = types.StringNull()
	state.SlackChannel = types.StringNull()
	state.AgentId = types.StringNull()

	if database.Port != nil {
		state.Port = types.Int64Value(*database.Port)
	}

	if database.Group != "" {
		state.Group = types.StringValue(database.Group)
	}

	if database.SlackChannel != "" {
		state.SlackChannel = types.StringValue(database.SlackChannel)
	}

	if database.AgentId != "" {
		state.AgentId = types.StringValue(database.AgentId)
	}

	if state.Credentials == nil || len(state.Credentials) != len(database.Credentials) {
		state.Credentials = make([]databaseCredentialModel, len(database.Credentials))
	}

	credentialIds := make(map[string]attr.Value)

	for index, credential := range database.Credentials {
		state.Credentials[index].Id = types.StringValue(credential.Id)
		state.Credentials[index].Username = types.StringValue(credential.Username)
		state.Credentials[index].ReviewsRequired = types.Int64Value(int64(credential.ReviewsRequired))
		state.Credentials[index].DefaultCredential = types.BoolValue(credential.DefaultCredential)

		state.Credentials[index].Hostname = types.StringNull()

		if credential.Hostname != "" {
			state.Credentials[index].Hostname = types.StringValue(credential.Hostname)
		}

		credentialIds[credential.Username] = types.StringValue(credential.Id)
	}

	state.CredentialIds = types.MapValueMust(types.StringType, credentialIds)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan databaseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var credentials []devhub.DatabaseCredential
	for _, credential := range plan.Credentials {
		credentials = append(credentials, devhub.DatabaseCredential{
			Id:                credential.Id.ValueString(),
			Username:          credential.Username.ValueString(),
			Password:          credential.Password.ValueString(),
			Hostname:          credential.Hostname.ValueString(),
			ReviewsRequired:   int(credential.ReviewsRequired.ValueInt64()),
			DefaultCredential: credential.DefaultCredential.ValueBool(),
		})
	}

	var port *int64
	if !plan.Port.IsNull() {
		portValue := plan.Port.ValueInt64()
		port = &portValue
	}

	input := devhub.Database{
		Name:           plan.Name.ValueString(),
		Adapter:        strings.ToLower(plan.Adapter.ValueString()),
		Hostname:       plan.Hostname.ValueString(),
		Port:           port,
		Database:       plan.Database.ValueString(),
		Ssl:            plan.Ssl.ValueBool(),
		Cacertfile:     plan.Cacertfile.ValueString(),
		Keyfile:        plan.Keyfile.ValueString(),
		Certfile:       plan.Certfile.ValueString(),
		RestrictAccess: plan.RestrictAccess.ValueBool(),
		Group:          plan.Group.ValueString(),
		SlackChannel:   plan.SlackChannel.ValueString(),
		AgentId:        plan.AgentId.ValueString(),
		Credentials:    credentials,
	}

	// Update existing order
	database, err := r.client.UpdateDatabase(plan.Id.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Database",
			"Could not update database, unexpected error: "+err.Error(),
		)
		return
	}

	credentialIds := make(map[string]attr.Value)

	for index, credential := range database.Credentials {
		plan.Credentials[index].Id = types.StringValue(credential.Id)

		credentialIds[credential.Username] = types.StringValue(credential.Id)
	}

	plan.CredentialIds = types.MapValueMust(types.StringType, credentialIds)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state databaseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDatabase(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting database",
			"Could not delete database, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *databaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
