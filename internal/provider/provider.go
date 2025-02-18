package provider

import (
	"context"
	"os"
	devhub "terraform-provider-devhub/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &devhubProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &devhubProvider{
			version: version,
		}
	}
}

type devhubProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

type devhubProvider struct {
	version string
}

func (p *devhubProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devhub"
	resp.Version = p.version
}

func (p *devhubProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Alternatively, can be configured using the `DEVHUB_API_KEY` environment variable.",
			},
		},
	}
}

func (p *devhubProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config devhubProviderModel

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Devhub API Host",
			"The provider cannot create the Devhub API client as there is an unknown configuration value for the Devhub API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DEVHUB_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("DEVHUB_HOST")
	api_key := os.Getenv("DEVHUB_API_KEY")

	if host == "" {
		host = config.Host.ValueString()
	}

	if api_key == "" {
		api_key = config.ApiKey.ValueString()
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Devhub API Key",
			"The provider cannot create the Devhub API client as there is a missing or empty value for the Devhub API key. "+
				"Set the api key value in the configuration or use the DEVHUB_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := devhub.NewClient(&host, &api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Devhub API Client",
			"An unexpected error occurred when creating the Devhub API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Devhub Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *devhubProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *devhubProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		DatabaseResource,
		TerradeskWorkspaceResource,
	}
}
