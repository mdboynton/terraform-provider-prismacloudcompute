package models

import (
	"context"

	api "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/rule"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomRuleResourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	AttackTechniques types.List   `tfsdk:"attack_techniques"`
	Description      types.String `tfsdk:"description"`
	Message          types.String `tfsdk:"message"`
	MinVersion       types.String `tfsdk:"min_version"`
	Modified         types.Int64  `tfsdk:"modified"`
	Name             types.String `tfsdk:"name"`
	Owner            types.String `tfsdk:"owner"`
	Script           types.String `tfsdk:"script"`
	Type             types.String `tfsdk:"type"`
	VulnIDs          types.List   `tfsdk:"vuln_ids"`
}

func (m *CustomRuleResourceModel) ToTerraform(ctx context.Context, diagnostics *diag.Diagnostics) api.CustomRule {
	var attackTechniques []string
	diagnostics.Append(m.AttackTechniques.ElementsAs(ctx, &attackTechniques, false)...)
	if diagnostics.HasError() {
		return api.CustomRule{}
	}

	var vulnIDs []string
	if m.VulnIDs.IsNull() || m.VulnIDs.IsUnknown() {
		vulnIDs = []string{}
	} else {
		diagnostics.Append(m.VulnIDs.ElementsAs(ctx, &vulnIDs, false)...)
		if diagnostics.HasError() {
			return api.CustomRule{}
		}
	}

	return api.CustomRule{
		Id:               int(m.Id.ValueInt64()),
		AttackTechniques: attackTechniques,
		Description:      m.Description.ValueString(),
		Message:          m.Message.ValueString(),
		MinVersion:       m.MinVersion.ValueString(),
		Modified:         int(m.Modified.ValueInt64()),
		Name:             m.Name.ValueString(),
		Owner:            m.Owner.ValueString(),
		Script:           m.Script.ValueString(),
		Type:             m.Type.ValueString(),
		VulnIDs:          vulnIDs,
	}
}

func (m *CustomRuleResourceModel) UpdateValues(ctx context.Context, diagnostics *diag.Diagnostics, data api.CustomRule) {
	attackTechniques, diags := types.ListValueFrom(ctx, types.StringType, data.AttackTechniques)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}

	vulnIDs, diags := types.ListValueFrom(ctx, types.StringType, data.VulnIDs)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}

	m.Id = types.Int64Value(int64(data.Id))
	m.AttackTechniques = attackTechniques
	m.Description = types.StringValue(data.Description)
	m.Message = types.StringValue(data.Message)
	m.MinVersion = types.StringValue(data.MinVersion)
	m.Modified = types.Int64Value(int64(data.Modified))
	m.Name = types.StringValue(data.Name)
	m.Owner = types.StringValue(data.Owner)
	m.Script = types.StringValue(data.Script)
	m.Type = types.StringValue(data.Type)
	m.VulnIDs = vulnIDs
}
