package models

import (
	"context"

	api "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/alertprofile"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type AlertProfileResourceModel struct {
	Name                                types.String `tfsdk:"name"`
	Policy                              types.Object `tfsdk:"policy"`
	VulnerabilityImmediateAlertsEnabled types.Bool   `tfsdk:"vuln_immediate_alerts_enabled"`
	Webhook                             types.Object `tfsdk:"webhook"`
}

func (m *AlertProfileResourceModel) ToCreateOrUpdateRequest(ctx context.Context, diagnostics *diag.Diagnostics) api.AlertProfile {
	var policy api.Policy
	diagnostics.Append(m.Policy.As(ctx, &policy, basetypes.ObjectAsOptions{})...)
	if diagnostics.HasError() {
		return api.AlertProfile{}
	}

	var webhook api.Webhook
	diagnostics.Append(m.Webhook.As(ctx, &webhook, basetypes.ObjectAsOptions{})...)
	if diagnostics.HasError() {
		return api.AlertProfile{}
	}

	return api.AlertProfile{
		Name:                                m.Name.ValueString(),
		Policy:                              policy,
		VulnerabilityImmediateAlertsEnabled: m.VulnerabilityImmediateAlertsEnabled.ValueBool(),
		Webhook:                             webhook,
	}
}

func (m *AlertProfileResourceModel) RefreshPropertyValues(ctx context.Context, diagnostics *diag.Diagnostics, data api.AlertProfile) {
	policy, diags := types.ObjectValueFrom(ctx, m.Policy.AttributeTypes(ctx), data.Policy)
	diagnostics.Append(diags...)

	webhook, diags := types.ObjectValueFrom(ctx, m.Webhook.AttributeTypes(ctx), data.Webhook)
	diagnostics.Append(diags...)

	m.Name = types.StringValue(data.Name)
	m.Policy = policy
	m.VulnerabilityImmediateAlertsEnabled = types.BoolValue(data.VulnerabilityImmediateAlertsEnabled)
	m.Webhook = webhook
}
