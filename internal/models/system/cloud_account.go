package models

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type GcpCloudAccountResourceModel struct {
    AccountName types.String `tfsdk:"account_name"`
    Description types.String `tfsdk:"description"`
    CredentialLevel types.String `tfsdk:"credential_level"`
    ServiceAccountKey types.String `tfsdk:"service_account_key"`
    ApiToken types.String `tfsdk:"api_token"`
    AgentlessScanning AgentlessScanningResourceModel `tfsdk:"agentless_scanning"`
    ServerlessScanning ServerlessScanningResourceModel `tfsdk:"serverless_scanning"`
    EnableCloudDiscovery types.Bool `tfsdk:"enable_cloud_discovery"`
}

type AgentlessScanningResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    IsHubAccount types.Bool `tfsdk:"is_hub_account"`
    HubAccountId types.String `tfsdk:"hub_account_id"`
    ConsoleURL types.String `tfsdk:"console_url"`
    Port types.Int32 `tfsdk:"port"`
    ProxyAddress types.String `tfsdk:"proxy_address"`
    ProxyCertificate types.String `tfsdk:"proxy_certificate"`
    AutoScale types.Bool `tfsdk:"auto_scale"`
    CustomLabels []CustomLabelResourceModel `tfsdk:"custom_labels"`
    MaxNumberOfScanners types.Int32 `tfsdk:"max_number_of_scanners"`
    EnforcePermissionsCheck types.Bool `tfsdk:"enforce_permissions_check"`
    Regions types.Set `tfsdk:"regions"`
    ScanScope types.String `tfsdk:"scan_scope"`
    ScanNonRunningHosts types.Bool `tfsdk:"scan_non_running_hosts"`
    ScopeByLabels types.String `tfsdk:"scope_by_labels"`
    Subnet types.String `tfsdk:"subnet"`
}

type CustomLabelResourceModel struct {
    Key types.String `tfsdk:"key"`
    Value types.String `tfsdk:"value"`
}

type ServerlessScanningResourceModel struct {
    Enabled types.Bool `tfsdk:"enabled"`
    Limit types.Int32 `tfsdk:"limit"`
}
