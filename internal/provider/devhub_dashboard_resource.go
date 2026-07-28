// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	devhub "terraform-provider-devhub/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &dashboardResource{}
	_ resource.ResourceWithConfigure   = &dashboardResource{}
	_ resource.ResourceWithImportState = &dashboardResource{}
)

func DashboardResource() resource.Resource {
	return &dashboardResource{}
}

// DashboardResourceModel describes the resource data model.
type dashboardResourceModel struct {
	Id               types.String          `tfsdk:"id"`
	Name             types.String          `tfsdk:"name"`
	RestrictedAccess types.Bool            `tfsdk:"restricted_access"`
	Panels           []dashboardPanelModel `tfsdk:"panels"`
}

type dashboardPanelModel struct {
	Id           types.String                     `tfsdk:"id"`
	Title        types.String                     `tfsdk:"title"`
	Inputs       []dashboardPanelInputModel       `tfsdk:"inputs"`
	Actions      []dashboardPanelActionModel      `tfsdk:"actions"`
	QueryDetails *dashboardPanelQueryDetailsModel `tfsdk:"query_details"`
}

type dashboardPanelInputModel struct {
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`
}

type dashboardPanelActionModel struct {
	Id            types.String                            `tfsdk:"id"`
	Label         types.String                            `tfsdk:"label"`
	WorkflowId    types.String                            `tfsdk:"workflow_id"`
	InputMappings []dashboardPanelActionInputMappingModel `tfsdk:"input_mappings"`
}

type dashboardPanelActionInputMappingModel struct {
	WorkflowInputKey types.String `tfsdk:"workflow_input_key"`
	Column           types.String `tfsdk:"column"`
	Prompt           types.String `tfsdk:"prompt"`
	PromptUser       types.Bool   `tfsdk:"prompt_user"`
}

type dashboardPanelQueryDetailsModel struct {
	Query        types.String `tfsdk:"query"`
	CredentialId types.String `tfsdk:"credential_id"`
}

type dashboardResource struct {
	client *devhub.Client
}

func (r *dashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (r *dashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dashboard resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Dashboard id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dashboard.",
				Required:            true,
			},
			"restricted_access": schema.BoolAttribute{
				MarkdownDescription: "Whether the dashboard is restricted to certain users.",
				Optional:            true,
				Computed:            true,
			},
			"panels": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Panel id.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The title of the panel.",
							Required:            true,
						},
						"inputs": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										MarkdownDescription: "The key for this input.",
										Required:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "A description of what this input is for.",
										Optional:            true,
									},
								},
							},
						},
						"actions": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Row actions that trigger a workflow when the user clicks the button on a row.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Action id.",
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"label": schema.StringAttribute{
										MarkdownDescription: "Button label shown to the user.",
										Required:            true,
									},
									"workflow_id": schema.StringAttribute{
										MarkdownDescription: "ID of the workflow to trigger when the button is clicked.",
										Required:            true,
									},
									"input_mappings": schema.ListNestedAttribute{
										Optional:            true,
										MarkdownDescription: "One entry per workflow input: either map a row column value or prompt the user at run time.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"workflow_input_key": schema.StringAttribute{
													MarkdownDescription: "The workflow input key this mapping fills.",
													Required:            true,
												},
												"column": schema.StringAttribute{
													MarkdownDescription: "Row column to use when prompt_user is false.",
													Optional:            true,
												},
												"prompt": schema.StringAttribute{
													MarkdownDescription: "Prompt text shown to the user when prompt_user is true.",
													Optional:            true,
												},
												"prompt_user": schema.BoolAttribute{
													MarkdownDescription: "When true, prompt the user for this value at run time instead of pulling it from a row column. Defaults to false.",
													Optional:            true,
													Computed:            true,
													Default:             booldefault.StaticBool(false),
												},
											},
										},
									},
								},
							},
						},
						"query_details": schema.SingleNestedAttribute{
							Optional: true,
							Validators: []validator.Object{
								objectvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("query_details"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"query": schema.StringAttribute{
									MarkdownDescription: "The SQL query to execute.",
									Required:            true,
								},
								"credential_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the database credential to use.",
									Required:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// The API omits unset optional strings as null, which decodes to "". Reading that
// back as an empty string would conflict with the null in configuration and leave a
// perpetual diff, so map it back to null.
func optionalString(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func planPanelsToDevhub(plan []dashboardPanelModel, includeIds bool) []devhub.DashboardPanel {
	var panels []devhub.DashboardPanel
	for _, statePanel := range plan {
		var inputs []devhub.DashboardPanelInput
		for _, input := range statePanel.Inputs {
			inputs = append(inputs, devhub.DashboardPanelInput{
				Key:         input.Key.ValueString(),
				Description: input.Description.ValueString(),
			})
		}

		var actions []devhub.DashboardPanelAction
		for _, action := range statePanel.Actions {
			var mappings []devhub.DashboardPanelActionInputMapping
			for _, mapping := range action.InputMappings {
				mappings = append(mappings, devhub.DashboardPanelActionInputMapping{
					WorkflowInputKey: mapping.WorkflowInputKey.ValueString(),
					Column:           mapping.Column.ValueString(),
					Prompt:           mapping.Prompt.ValueString(),
					PromptUser:       mapping.PromptUser.ValueBool(),
				})
			}

			actionInput := devhub.DashboardPanelAction{
				Label:         action.Label.ValueString(),
				WorkflowId:    action.WorkflowId.ValueString(),
				InputMappings: mappings,
			}
			if includeIds {
				actionInput.Id = action.Id.ValueString()
			}

			actions = append(actions, actionInput)
		}

		panel := devhub.DashboardPanel{
			Title:   statePanel.Title.ValueString(),
			Inputs:  inputs,
			Actions: actions,
		}
		if includeIds {
			panel.Id = statePanel.Id.ValueString()
		}

		if statePanel.QueryDetails != nil {
			panel.Details = &devhub.DashboardPanelDetails{
				Type:         "query",
				Query:        statePanel.QueryDetails.Query.ValueString(),
				CredentialId: statePanel.QueryDetails.CredentialId.ValueString(),
			}
		}

		panels = append(panels, panel)
	}
	return panels
}

func (r *dashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dashboardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	panels := planPanelsToDevhub(plan.Panels, false)

	input := devhub.Dashboard{
		Name:             plan.Name.ValueString(),
		RestrictedAccess: plan.RestrictedAccess.ValueBool(),
		Panels:           panels,
	}

	createdDashboard, err := r.client.CreateDashboard(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dashboard",
			"Could not create dashboard, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	plan.Id = types.StringValue(createdDashboard.Id)
	plan.RestrictedAccess = types.BoolValue(createdDashboard.RestrictedAccess)

	for index, panel := range createdDashboard.Panels {
		plan.Panels[index].Id = types.StringValue(panel.Id)
		for actionIndex, action := range panel.Actions {
			plan.Panels[index].Actions[actionIndex].Id = types.StringValue(action.Id)
			for mappingIndex, mapping := range action.InputMappings {
				plan.Panels[index].Actions[actionIndex].InputMappings[mappingIndex].PromptUser = types.BoolValue(mapping.PromptUser)
			}
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dashboardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboard, err := r.client.GetDashboard(state.Id.ValueString())

	if err != nil && err.Error() == "not found" {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading dashboard",
			"Could not read dashboard ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(dashboard.Id)
	state.Name = types.StringValue(dashboard.Name)
	state.RestrictedAccess = types.BoolValue(dashboard.RestrictedAccess)

	var statePanels []dashboardPanelModel
	for _, panel := range dashboard.Panels {

		var inputs []dashboardPanelInputModel
		for _, input := range panel.Inputs {
			inputs = append(inputs, dashboardPanelInputModel{
				Key:         types.StringValue(input.Key),
				Description: optionalString(input.Description),
			})
		}

		var actions []dashboardPanelActionModel
		for _, action := range panel.Actions {
			var mappings []dashboardPanelActionInputMappingModel
			for _, mapping := range action.InputMappings {
				mappings = append(mappings, dashboardPanelActionInputMappingModel{
					WorkflowInputKey: types.StringValue(mapping.WorkflowInputKey),
					Column:           optionalString(mapping.Column),
					Prompt:           optionalString(mapping.Prompt),
					PromptUser:       types.BoolValue(mapping.PromptUser),
				})
			}

			actions = append(actions, dashboardPanelActionModel{
				Id:            types.StringValue(action.Id),
				Label:         types.StringValue(action.Label),
				WorkflowId:    types.StringValue(action.WorkflowId),
				InputMappings: mappings,
			})
		}

		panelModel := dashboardPanelModel{
			Id:      types.StringValue(panel.Id),
			Title:   types.StringValue(panel.Title),
			Inputs:  inputs,
			Actions: actions,
		}

		if panel.Details.Type == "query" {
			panelModel.QueryDetails = &dashboardPanelQueryDetailsModel{
				Query:        types.StringValue(panel.Details.Query),
				CredentialId: types.StringValue(panel.Details.CredentialId),
			}
		}

		statePanels = append(statePanels, panelModel)
	}

	state.Panels = statePanels

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dashboardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	panels := planPanelsToDevhub(plan.Panels, true)

	dashboard := devhub.Dashboard{
		Id:               plan.Id.ValueString(),
		Name:             plan.Name.ValueString(),
		RestrictedAccess: plan.RestrictedAccess.ValueBool(),
		Panels:           panels,
	}

	updatedDashboard, err := r.client.UpdateDashboard(plan.Id.ValueString(), dashboard)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating dashboard",
			"Could not update dashboard, unexpected error: "+err.Error(),
		)
		return
	}

	plan.RestrictedAccess = types.BoolValue(updatedDashboard.RestrictedAccess)

	for index, panel := range updatedDashboard.Panels {
		plan.Panels[index].Id = types.StringValue(panel.Id)
		for actionIndex, action := range panel.Actions {
			plan.Panels[index].Actions[actionIndex].Id = types.StringValue(action.Id)
			for mappingIndex, mapping := range action.InputMappings {
				plan.Panels[index].Actions[actionIndex].InputMappings[mappingIndex].PromptUser = types.BoolValue(mapping.PromptUser)
			}
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *dashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dashboardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDashboard(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting dashboard",
			"Could not delete dashboard, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *dashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*devhub.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *devhub.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *dashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
