package models

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomRuntimeRuleResourceModel struct {
    Id              types.Int64  `tfsdk:"id"`
    AttackTechniques types.List   `tfsdk:"attack_techniques"`
    Description     types.String `tfsdk:"description"`
    Message         types.String `tfsdk:"message"`
    MinVersion      types.String `tfsdk:"min_version"`
    Modified        types.Int64  `tfsdk:"modified"`
    Name            types.String `tfsdk:"name"`
    Owner           types.String `tfsdk:"owner"`
    Script          types.String `tfsdk:"script"`
    Type            types.String `tfsdk:"type"`
    VulnIDs         types.List   `tfsdk:"vuln_ids"`
}
