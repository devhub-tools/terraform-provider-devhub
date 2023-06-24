// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// CreateDatabaseCreateDatabaseCreateDatabaseResult includes the requested fields of the GraphQL type CreateDatabaseResult.
// The GraphQL type's documentation follows.
//
// The result of the :create_database mutation
type CreateDatabaseCreateDatabaseCreateDatabaseResult struct {
	// The successful result of the mutation
	Result CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase `json:"result"`
	// Any errors generated, if the mutation failed
	Errors []CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError `json:"errors"`
}

// GetResult returns CreateDatabaseCreateDatabaseCreateDatabaseResult.Result, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResult) GetResult() CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase {
	return v.Result
}

// GetErrors returns CreateDatabaseCreateDatabaseCreateDatabaseResult.Errors, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResult) GetErrors() []CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError {
	return v.Errors
}

// CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError includes the requested fields of the GraphQL type MutationError.
// The GraphQL type's documentation follows.
//
// An error generated by a failed mutation
type CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError struct {
	// An error code for the given error
	Code string `json:"code"`
	// The field or fields that produced the error
	Fields []string `json:"fields"`
	// The human readable error message
	Message string `json:"message"`
}

// GetCode returns CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError.Code, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError) GetCode() string {
	return v.Code
}

// GetFields returns CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError.Fields, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError) GetFields() []string {
	return v.Fields
}

// GetMessage returns CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError.Message, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultErrorsMutationError) GetMessage() string {
	return v.Message
}

// CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase includes the requested fields of the GraphQL type Database.
type CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	RestrictAccess bool   `json:"restrictAccess"`
}

// GetId returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Id, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetId() string { return v.Id }

// GetName returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Name, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetName() string {
	return v.Name
}

// GetAdapter returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Adapter, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetAdapter() string {
	return v.Adapter
}

// GetHostname returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Hostname, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetHostname() string {
	return v.Hostname
}

// GetDatabase returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Database, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetDatabase() string {
	return v.Database
}

// GetSsl returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.Ssl, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetSsl() bool { return v.Ssl }

// GetRestrictAccess returns CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase.RestrictAccess, and is useful for accessing the field via an interface.
func (v *CreateDatabaseCreateDatabaseCreateDatabaseResultResultDatabase) GetRestrictAccess() bool {
	return v.RestrictAccess
}

type CreateDatabaseInput struct {
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	RestrictAccess bool   `json:"restrictAccess"`
	Cacertfile     string `json:"cacertfile"`
	Keyfile        string `json:"keyfile"`
	Certfile       string `json:"certfile"`
	AgentId        string `json:"agentId"`
}

// GetName returns CreateDatabaseInput.Name, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetName() string { return v.Name }

// GetAdapter returns CreateDatabaseInput.Adapter, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetAdapter() string { return v.Adapter }

// GetHostname returns CreateDatabaseInput.Hostname, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetHostname() string { return v.Hostname }

// GetDatabase returns CreateDatabaseInput.Database, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetDatabase() string { return v.Database }

// GetSsl returns CreateDatabaseInput.Ssl, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetSsl() bool { return v.Ssl }

// GetRestrictAccess returns CreateDatabaseInput.RestrictAccess, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetRestrictAccess() bool { return v.RestrictAccess }

// GetCacertfile returns CreateDatabaseInput.Cacertfile, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetCacertfile() string { return v.Cacertfile }

// GetKeyfile returns CreateDatabaseInput.Keyfile, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetKeyfile() string { return v.Keyfile }

// GetCertfile returns CreateDatabaseInput.Certfile, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetCertfile() string { return v.Certfile }

// GetAgentId returns CreateDatabaseInput.AgentId, and is useful for accessing the field via an interface.
func (v *CreateDatabaseInput) GetAgentId() string { return v.AgentId }

// CreateDatabaseResponse is returned by CreateDatabase on success.
type CreateDatabaseResponse struct {
	CreateDatabase CreateDatabaseCreateDatabaseCreateDatabaseResult `json:"createDatabase"`
}

// GetCreateDatabase returns CreateDatabaseResponse.CreateDatabase, and is useful for accessing the field via an interface.
func (v *CreateDatabaseResponse) GetCreateDatabase() CreateDatabaseCreateDatabaseCreateDatabaseResult {
	return v.CreateDatabase
}

// GetDatabaseDatabase includes the requested fields of the GraphQL type Database.
type GetDatabaseDatabase struct {
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	RestrictAccess bool   `json:"restrictAccess"`
}

// GetName returns GetDatabaseDatabase.Name, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetName() string { return v.Name }

// GetAdapter returns GetDatabaseDatabase.Adapter, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetAdapter() string { return v.Adapter }

// GetHostname returns GetDatabaseDatabase.Hostname, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetHostname() string { return v.Hostname }

// GetDatabase returns GetDatabaseDatabase.Database, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetDatabase() string { return v.Database }

// GetSsl returns GetDatabaseDatabase.Ssl, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetSsl() bool { return v.Ssl }

// GetRestrictAccess returns GetDatabaseDatabase.RestrictAccess, and is useful for accessing the field via an interface.
func (v *GetDatabaseDatabase) GetRestrictAccess() bool { return v.RestrictAccess }

// GetDatabaseResponse is returned by GetDatabase on success.
type GetDatabaseResponse struct {
	Database GetDatabaseDatabase `json:"database"`
}

// GetDatabase returns GetDatabaseResponse.Database, and is useful for accessing the field via an interface.
func (v *GetDatabaseResponse) GetDatabase() GetDatabaseDatabase { return v.Database }

type UpdateDatabaseInput struct {
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	RestrictAccess bool   `json:"restrictAccess"`
	NewCacertfile  string `json:"newCacertfile"`
	NewKeyfile     string `json:"newKeyfile"`
	NewCertfile    string `json:"newCertfile"`
	AgentId        string `json:"agentId"`
}

// GetName returns UpdateDatabaseInput.Name, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetName() string { return v.Name }

// GetAdapter returns UpdateDatabaseInput.Adapter, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetAdapter() string { return v.Adapter }

// GetHostname returns UpdateDatabaseInput.Hostname, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetHostname() string { return v.Hostname }

// GetDatabase returns UpdateDatabaseInput.Database, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetDatabase() string { return v.Database }

// GetSsl returns UpdateDatabaseInput.Ssl, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetSsl() bool { return v.Ssl }

// GetRestrictAccess returns UpdateDatabaseInput.RestrictAccess, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetRestrictAccess() bool { return v.RestrictAccess }

// GetNewCacertfile returns UpdateDatabaseInput.NewCacertfile, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetNewCacertfile() string { return v.NewCacertfile }

// GetNewKeyfile returns UpdateDatabaseInput.NewKeyfile, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetNewKeyfile() string { return v.NewKeyfile }

// GetNewCertfile returns UpdateDatabaseInput.NewCertfile, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetNewCertfile() string { return v.NewCertfile }

// GetAgentId returns UpdateDatabaseInput.AgentId, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseInput) GetAgentId() string { return v.AgentId }

// UpdateDatabaseResponse is returned by UpdateDatabase on success.
type UpdateDatabaseResponse struct {
	UpdateDatabase UpdateDatabaseUpdateDatabaseUpdateDatabaseResult `json:"updateDatabase"`
}

// GetUpdateDatabase returns UpdateDatabaseResponse.UpdateDatabase, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseResponse) GetUpdateDatabase() UpdateDatabaseUpdateDatabaseUpdateDatabaseResult {
	return v.UpdateDatabase
}

// UpdateDatabaseUpdateDatabaseUpdateDatabaseResult includes the requested fields of the GraphQL type UpdateDatabaseResult.
// The GraphQL type's documentation follows.
//
// The result of the :update_database mutation
type UpdateDatabaseUpdateDatabaseUpdateDatabaseResult struct {
	// The successful result of the mutation
	Result UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase `json:"result"`
	// Any errors generated, if the mutation failed
	Errors []UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError `json:"errors"`
}

// GetResult returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResult.Result, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResult) GetResult() UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase {
	return v.Result
}

// GetErrors returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResult.Errors, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResult) GetErrors() []UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError {
	return v.Errors
}

// UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError includes the requested fields of the GraphQL type MutationError.
// The GraphQL type's documentation follows.
//
// An error generated by a failed mutation
type UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError struct {
	// An error code for the given error
	Code string `json:"code"`
	// The field or fields that produced the error
	Fields []string `json:"fields"`
	// The human readable error message
	Message string `json:"message"`
}

// GetCode returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError.Code, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError) GetCode() string {
	return v.Code
}

// GetFields returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError.Fields, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError) GetFields() []string {
	return v.Fields
}

// GetMessage returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError.Message, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultErrorsMutationError) GetMessage() string {
	return v.Message
}

// UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase includes the requested fields of the GraphQL type Database.
type UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase struct {
	Name           string `json:"name"`
	Adapter        string `json:"adapter"`
	Hostname       string `json:"hostname"`
	Database       string `json:"database"`
	Ssl            bool   `json:"ssl"`
	RestrictAccess bool   `json:"restrictAccess"`
}

// GetName returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.Name, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetName() string {
	return v.Name
}

// GetAdapter returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.Adapter, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetAdapter() string {
	return v.Adapter
}

// GetHostname returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.Hostname, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetHostname() string {
	return v.Hostname
}

// GetDatabase returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.Database, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetDatabase() string {
	return v.Database
}

// GetSsl returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.Ssl, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetSsl() bool { return v.Ssl }

// GetRestrictAccess returns UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase.RestrictAccess, and is useful for accessing the field via an interface.
func (v *UpdateDatabaseUpdateDatabaseUpdateDatabaseResultResultDatabase) GetRestrictAccess() bool {
	return v.RestrictAccess
}

// __CreateDatabaseInput is used internally by genqlient
type __CreateDatabaseInput struct {
	Input CreateDatabaseInput `json:"input"`
}

// GetInput returns __CreateDatabaseInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateDatabaseInput) GetInput() CreateDatabaseInput { return v.Input }

// __GetDatabaseInput is used internally by genqlient
type __GetDatabaseInput struct {
	Id string `json:"id"`
}

// GetId returns __GetDatabaseInput.Id, and is useful for accessing the field via an interface.
func (v *__GetDatabaseInput) GetId() string { return v.Id }

// __UpdateDatabaseInput is used internally by genqlient
type __UpdateDatabaseInput struct {
	Id    string              `json:"id"`
	Input UpdateDatabaseInput `json:"input"`
}

// GetId returns __UpdateDatabaseInput.Id, and is useful for accessing the field via an interface.
func (v *__UpdateDatabaseInput) GetId() string { return v.Id }

// GetInput returns __UpdateDatabaseInput.Input, and is useful for accessing the field via an interface.
func (v *__UpdateDatabaseInput) GetInput() UpdateDatabaseInput { return v.Input }

// The query or mutation executed by CreateDatabase.
const CreateDatabase_Operation = `
mutation CreateDatabase ($input: CreateDatabaseInput!) {
	createDatabase(input: $input) {
		result {
			id
			name
			adapter
			hostname
			database
			ssl
			restrictAccess
		}
		errors {
			code
			fields
			message
		}
	}
}
`

func CreateDatabase(
	ctx context.Context,
	client graphql.Client,
	input CreateDatabaseInput,
) (*CreateDatabaseResponse, error) {
	req := &graphql.Request{
		OpName: "CreateDatabase",
		Query:  CreateDatabase_Operation,
		Variables: &__CreateDatabaseInput{
			Input: input,
		},
	}
	var err error

	var data CreateDatabaseResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by GetDatabase.
const GetDatabase_Operation = `
query GetDatabase ($id: ID!) {
	database(id: $id) {
		name
		adapter
		hostname
		database
		ssl
		restrictAccess
	}
}
`

func GetDatabase(
	ctx context.Context,
	client graphql.Client,
	id string,
) (*GetDatabaseResponse, error) {
	req := &graphql.Request{
		OpName: "GetDatabase",
		Query:  GetDatabase_Operation,
		Variables: &__GetDatabaseInput{
			Id: id,
		},
	}
	var err error

	var data GetDatabaseResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

// The query or mutation executed by UpdateDatabase.
const UpdateDatabase_Operation = `
mutation UpdateDatabase ($id: ID!, $input: UpdateDatabaseInput!) {
	updateDatabase(id: $id, input: $input) {
		result {
			name
			adapter
			hostname
			database
			ssl
			restrictAccess
		}
		errors {
			code
			fields
			message
		}
	}
}
`

func UpdateDatabase(
	ctx context.Context,
	client graphql.Client,
	id string,
	input UpdateDatabaseInput,
) (*UpdateDatabaseResponse, error) {
	req := &graphql.Request{
		OpName: "UpdateDatabase",
		Query:  UpdateDatabase_Operation,
		Variables: &__UpdateDatabaseInput{
			Id:    id,
			Input: input,
		},
	}
	var err error

	var data UpdateDatabaseResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
