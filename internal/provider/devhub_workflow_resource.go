// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	devhub "terraform-provider-devhub/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &workflowResource{}
	_ resource.ResourceWithConfigure   = &workflowResource{}
	_ resource.ResourceWithImportState = &workflowResource{}
)

func WorkflowResource() resource.Resource {
	return &workflowResource{}
}

// WorkflowResourceModel describes the resource data model.
type workflowResourceModel struct {
	Id                     types.String         `tfsdk:"id"`
	Name                   types.String         `tfsdk:"name"`
	TriggerLinearLabelName types.String         `tfsdk:"trigger_linear_label_name"`
	Inputs                 []workflowInputModel `tfsdk:"inputs"`
	Steps                  []workflowStepModel  `tfsdk:"steps"`
}

type workflowInputModel struct {
	Key         types.String `tfsdk:"key"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

type workflowStepModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// Action type fields
	ApiAction        *workflowApiActionModel        `tfsdk:"api_action"`
	ApprovalAction   *workflowApprovalActionModel   `tfsdk:"approval_action"`
	QueryAction      *workflowQueryActionModel      `tfsdk:"query_action"`
	SlackAction      *workflowSlackActionModel      `tfsdk:"slack_action"`
	SlackReplyAction *workflowSlackReplyActionModel `tfsdk:"slack_reply_action"`
}

type permissionModel struct {
	Id                 types.String `tfsdk:"id"`
	Permission         types.String `tfsdk:"permission"`
	RoleId             types.String `tfsdk:"role_id"`
	OrganizationUserId types.String `tfsdk:"organization_user_id"`
}

type workflowApiActionModel struct {
	Endpoint           types.String                    `tfsdk:"endpoint"`
	Method             types.String                    `tfsdk:"method"`
	Headers            []workflowApiActionHeadersModel `tfsdk:"headers"`
	Body               types.String                    `tfsdk:"body"`
	ExpectedStatusCode types.Int64                     `tfsdk:"expected_status_code"`
	IncludeDevhubJwt   types.Bool                      `tfsdk:"include_devhub_jwt"`
}

type workflowApiActionHeadersModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type workflowApprovalActionModel struct {
	ReviewsRequired types.Int64        `tfsdk:"reviews_required"`
	Permissions     []*permissionModel `tfsdk:"permissions"`
}

type workflowQueryActionModel struct {
	Timeout      types.Int64  `tfsdk:"timeout"`
	Query        types.String `tfsdk:"query"`
	CredentialId types.String `tfsdk:"credential_id"`
}

type workflowSlackActionModel struct {
	SlackChannel types.String `tfsdk:"slack_channel"`
	Message      types.String `tfsdk:"message"`
	LinkText     types.String `tfsdk:"link_text"`
}

type workflowSlackReplyActionModel struct {
	ReplyToStepName types.String `tfsdk:"reply_to_step_name"`
	Message         types.String `tfsdk:"message"`
}

type workflowResource struct {
	client *devhub.Client
}

func (r *workflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *workflowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workflow resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Workflow id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workflow.",
				Required:            true,
			},
			"trigger_linear_label_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Linear label that should trigger the workflow.",
				Optional:            true,
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
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of this input (string, float, integer, boolean).",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"string",
									"float",
									"integer",
									"boolean",
								),
							},
						},
					},
				},
			},
			"steps": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Step id.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the step.",
							Optional:            true,
						},
						"api_action": schema.SingleNestedAttribute{
							Optional: true,
							Validators: []validator.Object{
								objectvalidator.ExactlyOneOf(
									path.MatchRelative().AtParent().AtName("api_action"),
									path.MatchRelative().AtParent().AtName("approval_action"),
									path.MatchRelative().AtParent().AtName("query_action"),
									path.MatchRelative().AtParent().AtName("slack_action"),
									path.MatchRelative().AtParent().AtName("slack_reply_action"),
								),
							},
							Attributes: map[string]schema.Attribute{
								"endpoint": schema.StringAttribute{
									MarkdownDescription: "The endpoint for the API request.",
									Required:            true,
								},
								"method": schema.StringAttribute{
									MarkdownDescription: "The HTTP method for the API request.",
									Required:            true,
								},
								"headers": schema.ListNestedAttribute{
									MarkdownDescription: "Headers for the API request.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												MarkdownDescription: "The key for the header.",
												Required:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "The value for the header.",
												Required:            true,
											},
										},
									},
								},
								"body": schema.StringAttribute{
									MarkdownDescription: "The request body for the API request.",
									Optional:            true,
								},
								"expected_status_code": schema.Int64Attribute{
									MarkdownDescription: "The expected status code for the API request.",
									Required:            true,
								},
								"include_devhub_jwt": schema.BoolAttribute{
									MarkdownDescription: "Whether to include the Devhub JWT in the API request.",
									Required:            true,
								},
							},
						},
						"approval_action": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"reviews_required": schema.Int64Attribute{
									MarkdownDescription: "Number of required approvals.",
									Required:            true,
								},
								"permissions": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												MarkdownDescription: "Permission ID.",
												Computed:            true,
											},
											"permission": schema.StringAttribute{
												MarkdownDescription: "The permission granted to the role or user.",
												Required:            true,
											},
											"role_id": schema.StringAttribute{
												MarkdownDescription: "The id of the role.",
												Optional:            true,
												Validators: []validator.String{
													stringvalidator.ExactlyOneOf(
														path.MatchRelative().AtParent().AtName("role_id"),
														path.MatchRelative().AtParent().AtName("organization_user_id"),
													),
												},
											},
											"organization_user_id": schema.StringAttribute{
												MarkdownDescription: "The id of the organization user.",
												Optional:            true,
											},
										},
									},
								},
							},
						},
						"query_action": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"timeout": schema.Int64Attribute{
									MarkdownDescription: "The timeout for the query.",
									Required:            true,
								},
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
						"slack_action": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"slack_channel": schema.StringAttribute{
									MarkdownDescription: "The Slack channel to post to.",
									Required:            true,
								},
								"message": schema.StringAttribute{
									MarkdownDescription: "The message to post.",
									Required:            true,
								},
								"link_text": schema.StringAttribute{
									MarkdownDescription: "The text to display in the link.",
									Required:            true,
								},
							},
						},
						"slack_reply_action": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"reply_to_step_name": schema.StringAttribute{
									MarkdownDescription: "The name of the step to reply to.",
									Required:            true,
								},
								"message": schema.StringAttribute{
									MarkdownDescription: "The message to reply with.",
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

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workflowResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var inputs []devhub.WorkflowInput
	for _, planInput := range plan.Inputs {
		input := devhub.WorkflowInput{
			Key:  planInput.Key.ValueString(),
			Type: planInput.Type.ValueString(),
		}

		if planInput.Description.ValueString() != "" {
			input.Description = planInput.Description.ValueString()
		}

		inputs = append(inputs, input)
	}

	var steps []devhub.WorkflowStep
	for _, step := range plan.Steps {
		workflowStep := devhub.WorkflowStep{}

		if step.Name.ValueString() != "" {
			workflowStep.Name = step.Name.ValueString()
		}

		workflowStep.Permissions = make([]*devhub.Permission, 0)

		if step.ApiAction != nil {
			var headers []devhub.WorkflowStepActionApiHeader
			for _, header := range step.ApiAction.Headers {
				headers = append(headers, devhub.WorkflowStepActionApiHeader{
					Key:   header.Key.ValueString(),
					Value: header.Value.ValueString(),
				})
			}

			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:               "api",
				Endpoint:           step.ApiAction.Endpoint.ValueString(),
				Method:             step.ApiAction.Method.ValueString(),
				Headers:            headers,
				Body:               step.ApiAction.Body.ValueString(),
				ExpectedStatusCode: step.ApiAction.ExpectedStatusCode.ValueInt64(),
				IncludeDevhubJwt:   step.ApiAction.IncludeDevhubJwt.ValueBool(),
			}
		}

		if step.ApprovalAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:            "approval",
				ReviewsRequired: int(step.ApprovalAction.ReviewsRequired.ValueInt64()),
			}

			for _, permission := range step.ApprovalAction.Permissions {
				workflowStep.Permissions = append(workflowStep.Permissions, &devhub.Permission{
					Permission:         permission.Permission.ValueString(),
					RoleId:             permission.RoleId.ValueString(),
					OrganizationUserId: permission.OrganizationUserId.ValueString(),
				})
			}
		}

		if step.QueryAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:         "query",
				Timeout:      int(step.QueryAction.Timeout.ValueInt64()),
				Query:        step.QueryAction.Query.ValueString(),
				CredentialId: step.QueryAction.CredentialId.ValueString(),
			}
		}

		if step.SlackAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:         "slack",
				SlackChannel: step.SlackAction.SlackChannel.ValueString(),
				Message:      step.SlackAction.Message.ValueString(),
				LinkText:     step.SlackAction.LinkText.ValueString(),
			}
		}

		if step.SlackReplyAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:            "slack_reply",
				ReplyToStepName: step.SlackReplyAction.ReplyToStepName.ValueString(),
				Message:         step.SlackReplyAction.Message.ValueString(),
			}
		}

		steps = append(steps, workflowStep)
	}

	input := devhub.Workflow{
		Name:               plan.Name.ValueString(),
		TriggerLinearLabel: devhub.TriggerLinearLabel{Name: plan.TriggerLinearLabelName.ValueString()},
		Inputs:             inputs,
		Steps:              steps,
	}

	createdWorkflow, err := r.client.CreateWorkflow(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workflow",
			"Could not create workflow, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	plan.Id = types.StringValue(createdWorkflow.Id)

	for index, step := range createdWorkflow.Steps {
		plan.Steps[index].Id = types.StringValue(step.Id)

		if step.Action.Type == "approval" {
			for permissionIndex, permission := range step.Permissions {
				plan.Steps[index].ApprovalAction.Permissions[permissionIndex].Id = types.StringValue(permission.Id)
			}
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflow, err := r.client.GetWorkflow(state.Id.ValueString())

	if err != nil && err.Error() == "not found" {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workflow",
			"Could not read workflow ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(workflow.Id)
	state.Name = types.StringValue(workflow.Name)

	if workflow.TriggerLinearLabel.Name == "" {
		state.TriggerLinearLabelName = types.StringNull()
	} else {
		state.TriggerLinearLabelName = types.StringValue(workflow.TriggerLinearLabel.Name)
	}

	var stateInputs []workflowInputModel
	for _, stateInput := range workflow.Inputs {
		input := workflowInputModel{
			Key:  types.StringValue(stateInput.Key),
			Type: types.StringValue(stateInput.Type),
		}

		if stateInput.Description != "" {
			input.Description = types.StringValue(stateInput.Description)
		}

		stateInputs = append(stateInputs, input)
	}

	state.Inputs = stateInputs

	var stateSteps []workflowStepModel
	for _, step := range workflow.Steps {
		stepModel := workflowStepModel{
			Id: types.StringValue(step.Id),
		}

		if step.Name != "" {
			stepModel.Name = types.StringValue(step.Name)
		}

		switch step.Action.Type {
		case "api":
			var headers []workflowApiActionHeadersModel
			for _, header := range step.Action.Headers {
				headers = append(headers, workflowApiActionHeadersModel{
					Key:   types.StringValue(header.Key),
					Value: types.StringValue(header.Value),
				})
			}

			stepModel.ApiAction = &workflowApiActionModel{
				Endpoint:           types.StringValue(step.Action.Endpoint),
				Method:             types.StringValue(step.Action.Method),
				Headers:            headers,
				ExpectedStatusCode: types.Int64Value(step.Action.ExpectedStatusCode),
				IncludeDevhubJwt:   types.BoolValue(step.Action.IncludeDevhubJwt),
			}

			if step.Action.Body != "" {
				stepModel.ApiAction.Body = types.StringValue(step.Action.Body)
			}
		case "approval":
			permissions := make([]*permissionModel, 0)

			for _, permission := range step.Permissions {
				permissionModel := permissionModel{
					Id:         types.StringValue(permission.Id),
					Permission: types.StringValue(permission.Permission),
				}

				if permission.RoleId != "" {
					permissionModel.RoleId = types.StringValue(permission.RoleId)
				}

				if permission.OrganizationUserId != "" {
					permissionModel.OrganizationUserId = types.StringValue(permission.OrganizationUserId)
				}

				permissions = append(permissions, &permissionModel)
			}

			stepModel.ApprovalAction = &workflowApprovalActionModel{
				ReviewsRequired: types.Int64Value(int64(step.Action.ReviewsRequired)),
				Permissions:     permissions,
			}

		case "query":
			stepModel.QueryAction = &workflowQueryActionModel{
				Timeout:      types.Int64Value(int64(step.Action.Timeout)),
				Query:        types.StringValue(step.Action.Query),
				CredentialId: types.StringValue(step.Action.CredentialId),
			}
		case "slack":
			stepModel.SlackAction = &workflowSlackActionModel{
				SlackChannel: types.StringValue(step.Action.SlackChannel),
				Message:      types.StringValue(step.Action.Message),
				LinkText:     types.StringValue(step.Action.LinkText),
			}
		case "slack_reply":
			stepModel.SlackReplyAction = &workflowSlackReplyActionModel{
				ReplyToStepName: types.StringValue(step.Action.ReplyToStepName),
				Message:         types.StringValue(step.Action.Message),
			}
		}

		stateSteps = append(stateSteps, stepModel)
	}

	state.Steps = stateSteps

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workflowResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var inputs []devhub.WorkflowInput
	for _, planInput := range plan.Inputs {
		input := devhub.WorkflowInput{
			Key:  planInput.Key.ValueString(),
			Type: planInput.Type.ValueString(),
		}

		if planInput.Description.ValueString() != "" {
			input.Description = planInput.Description.ValueString()
		}

		inputs = append(inputs, input)
	}

	var steps []devhub.WorkflowStep
	for _, step := range plan.Steps {
		workflowStep := devhub.WorkflowStep{
			Id: step.Id.ValueString(),
		}

		if step.Name.ValueString() != "" {
			workflowStep.Name = step.Name.ValueString()
		}

		workflowStep.Permissions = make([]*devhub.Permission, 0)

		if step.ApiAction != nil {
			var headers []devhub.WorkflowStepActionApiHeader
			for _, header := range step.ApiAction.Headers {
				headers = append(headers, devhub.WorkflowStepActionApiHeader{
					Key:   header.Key.ValueString(),
					Value: header.Value.ValueString(),
				})
			}

			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:               "api",
				Endpoint:           step.ApiAction.Endpoint.ValueString(),
				Method:             step.ApiAction.Method.ValueString(),
				Headers:            headers,
				Body:               step.ApiAction.Body.ValueString(),
				ExpectedStatusCode: step.ApiAction.ExpectedStatusCode.ValueInt64(),
				IncludeDevhubJwt:   step.ApiAction.IncludeDevhubJwt.ValueBool(),
			}
		} else if step.ApprovalAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:            "approval",
				ReviewsRequired: int(step.ApprovalAction.ReviewsRequired.ValueInt64()),
			}

			for _, permission := range step.ApprovalAction.Permissions {
				permissionModel := devhub.Permission{
					Id:         permission.Id.ValueString(),
					Permission: permission.Permission.ValueString(),
				}

				if permission.RoleId.ValueString() != "" {
					permissionModel.RoleId = permission.RoleId.ValueString()
				}

				if permission.OrganizationUserId.ValueString() != "" {
					permissionModel.OrganizationUserId = permission.OrganizationUserId.ValueString()
				}

				workflowStep.Permissions = append(workflowStep.Permissions, &permissionModel)
			}
		} else if step.QueryAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:         "query",
				Timeout:      int(step.QueryAction.Timeout.ValueInt64()),
				Query:        step.QueryAction.Query.ValueString(),
				CredentialId: step.QueryAction.CredentialId.ValueString(),
			}
		} else if step.SlackAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:         "slack",
				SlackChannel: step.SlackAction.SlackChannel.ValueString(),
				Message:      step.SlackAction.Message.ValueString(),
				LinkText:     step.SlackAction.LinkText.ValueString(),
			}
		} else if step.SlackReplyAction != nil {
			workflowStep.Action = &devhub.WorkflowStepAction{
				Type:            "slack_reply",
				ReplyToStepName: step.SlackReplyAction.ReplyToStepName.ValueString(),
				Message:         step.SlackReplyAction.Message.ValueString(),
			}
		}

		steps = append(steps, workflowStep)
	}

	workflow := devhub.Workflow{
		Id:                 plan.Id.ValueString(),
		Name:               plan.Name.ValueString(),
		TriggerLinearLabel: devhub.TriggerLinearLabel{Name: plan.TriggerLinearLabelName.ValueString()},
		Inputs:             inputs,
		Steps:              steps,
	}

	updatedWorkflow, err := r.client.UpdateWorkflow(plan.Id.ValueString(), workflow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workflow",
			"Could not update workflow, unexpected error: "+err.Error(),
		)
		return
	}

	for index, step := range updatedWorkflow.Steps {
		plan.Steps[index].Id = types.StringValue(step.Id)

		if step.Action.Type == "approval" {
			for permissionIndex, permission := range step.Permissions {
				plan.Steps[index].ApprovalAction.Permissions[permissionIndex].Id = types.StringValue(permission.Id)
			}
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workflowResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkflow(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting workflow",
			"Could not delete workflow, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *workflowResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
