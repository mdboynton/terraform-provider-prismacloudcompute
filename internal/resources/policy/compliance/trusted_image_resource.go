package policy

import (
	"context"
	"fmt"
	"time"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	collectionAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/collection"
	policyAPI "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/policy"
	models "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/models/policy"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

	//"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *TrustedImagesPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_images_compliance_policy"
}

func (r *TrustedImagesPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.GetSchema()
}

func (r *TrustedImagesPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.PrismaCloudComputeAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TrustedImagesPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.TrustedImagesPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	data, diags := schemaToTrustedImagesPolicy(ctx, *r.client, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new trusted images compliance policy
	err := policyAPI.UpsertTrustedImagesPolicy(*r.client, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trusted Images Policy resource",
			fmt.Sprintf("Failed to create trusted images policy: %s", err.Error()),
		)
		return
	}

	// Retrieve newly created trusted images compliance policy
	response, err := policyAPI.GetTrustedImagesPolicy(*r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created Trusted Images Policy resource",
			"Failed to retrieve created trusted images policy: "+err.Error(),
		)
		return
	}

	createdPolicy, diags := trustedImagesPolicyToSchema(ctx, response)
	if diags.HasError() {
		return
	}

	// Set state to policy data
	diags = resp.State.Set(ctx, createdPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TrustedImagesPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.TrustedImagesPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get policy value from Prisma Cloud
	data, err := policyAPI.GetTrustedImagesPolicy(*r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Trusted Images Policy resource",
			"Failed to read trusted images policy: "+err.Error(),
		)
		return
	}

	// Overwrite state values with Prisma Cloud data
	policySchema, diags := trustedImagesPolicyToSchema(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &policySchema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TrustedImagesPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.TrustedImagesPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan models.TrustedImagesPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	planPolicy, diags := schemaToTrustedImagesPolicy(ctx, *r.client, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing policy
	err := policyAPI.UpsertTrustedImagesPolicy(*r.client, planPolicy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trusted Images Policy resource",
			"Failed to update trusted images policy: "+err.Error(),
		)
		return
	}

	// Get updated policy value from Prisma Cloud
	updatedPolicy, err := policyAPI.GetTrustedImagesPolicy(*r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Trusted Images Policy resource",
			"Failed to read trusted images policy: "+err.Error(),
		)
		return
	}

	// Convert updated policy into schema
	policySchema, diags := trustedImagesPolicyToSchema(ctx, updatedPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set updated state
	diags = resp.State.Set(ctx, &policySchema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TrustedImagesPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.TrustedImagesPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Clear policy rules
	state.Policy.Rules = &[]models.TrustedImagesPolicyRuleResourceModel{}

	// Clear groups
	state.Groups = &[]models.TrustGroupResourceModel{}

	// Disable policy
	state.Policy.Enabled = types.BoolValue(false)

	// Generate API request body from plan
	updatedPlan, diags := schemaToTrustedImagesPolicy(ctx, *r.client, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := policyAPI.UpsertTrustedImagesPolicy(*r.client, updatedPlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Trusted Images Policy resource",
			"Failed to delete trusted images policy: "+err.Error(),
		)
		return
	}
}

func (r *TrustedImagesPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	util.LogDebug(ctx, "executing ImportState")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TrustedImagesPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	util.LogDebug(ctx, "entering ModifyPlan")

	var plan *models.TrustedImagesPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: put this into ValidateConfig() (and do the same for other such checks in ModifyPlan funcs)
	if plan != nil && len(*plan.Groups) == 0 && len(*plan.Policy.Rules) > 0 {
		resp.Diagnostics.AddError(
			"Resource Configuration Error",
			"Trusted image policy cannot contain zero trust groups and a non-zero amount of rules",
		)
		return
	}

	util.LogDebug(ctx, "exiting ModifyPlan")
}

func schemaToTrustedImagesPolicy(ctx context.Context, client api.PrismaCloudComputeAPIClient, plan *models.TrustedImagesPolicyResourceModel) (policyAPI.TrustedImages, diag.Diagnostics) {
	util.LogDebug(ctx, "entering schemaToTrustedImagesPolicy")
	var diags diag.Diagnostics

	policy := policyAPI.TrustedImages{}

	if len(*plan.Groups) == 0 && len(*plan.Policy.Rules) > 0 {
		diags.AddError(
			"Resource Configuration Error",
			"Trusted image policy cannot contain zero trust groups and a non-zero amount of rules",
		)
		return policy, diags
	}

	groups, diags := schemaToTrustedImagesPolicyGroups(ctx, *plan.Groups)
	if diags.HasError() {
		return policy, diags
	}

	policy.Groups = groups

	rules, diags := schemaToTrustedImagesPolicyRules(ctx, client, plan.Policy)
	if diags.HasError() {
		return policy, diags
	}

	policy.Policy = rules

	// TODO: may want to raise an exception here if groups is empty but rules isnt

	util.LogDebug(ctx, "exiting schemaToTrustedImagesPolicy")
	return policy, diags
}

func schemaToTrustedImagesPolicyGroups(ctx context.Context, schemaGroups []models.TrustGroupResourceModel) ([]policyAPI.TrustGroup, diag.Diagnostics) {
	util.LogDebug(ctx, "entering schemaToTrustedImagesPolicyGroups")

	var diags diag.Diagnostics
	groups := make([]policyAPI.TrustGroup, 0, len(schemaGroups))

	for _, schemaGroup := range schemaGroups {
		images := make([]string, 0, len(schemaGroup.Images.Elements()))
		diags := schemaGroup.Images.ElementsAs(ctx, &images, false)
		if diags.HasError() {
			return groups, diags
		}

		groups = append(groups, policyAPI.TrustGroup{
			Id:           schemaGroup.Id.ValueString(),
			Images:       images,
			Modified:     time.Now().Format("2006-01-02T15:04:05.000Z"),
			Name:         "",
			Owner:        schemaGroup.Owner.ValueString(),
			PreviousName: schemaGroup.PreviousName.ValueString(),
		})
	}

	util.LogDebug(ctx, "exiting schemaToTrustedImagesPolicyGroups")

	return groups, diags
}

func schemaToTrustedImagesPolicyRules(ctx context.Context, client api.PrismaCloudComputeAPIClient, schemaRules models.TrustedImagesPolicyRulesResourceModel) (policyAPI.TrustedImagesPolicy, diag.Diagnostics) {
	util.LogDebug(ctx, "entering schemaToTrustedImagesPolicyRules")

	var diags diag.Diagnostics
	rules := policyAPI.TrustedImagesPolicy{
		Id:      schemaRules.Id.ValueString(),
		Enabled: schemaRules.Enabled.ValueBool(),
	}

	policyRules := []policyAPI.TrustedImagesPolicyRule{}
	if schemaRules.Rules != nil {
		for _, schemaRule := range *schemaRules.Rules {
			//action := []string{}
			//diags = schemaRule.Action.ElementsAs(ctx, &action, false)
			//if diags.HasError() {
			//    return rules, diags
			//}

			allowedGroups := []string{}
			diags = schemaRule.AllowedGroups.ElementsAs(ctx, &allowedGroups, false)
			if diags.HasError() {
				return rules, diags
			}

			deniedGroups := []string{}
			diags = schemaRule.DeniedGroups.ElementsAs(ctx, &deniedGroups, false)
			if diags.HasError() {
				return rules, diags
			}

			collectionNames := []string{}
			diags = schemaRule.Collections.ElementsAs(ctx, &collectionNames, false)
			if diags.HasError() {
				return rules, diags
			}

			collections, err := collectionAPI.GetCollections(client, collectionNames)
			if err != nil {
				diags.AddError(
					"Value Conversion Error",
					fmt.Sprintf("Error converting trusted image policy rules to schema: %s", err.Error()),
				)
				return rules, diags
			}

			modified := time.Now().Format("2006-01-02T15:04:05.000Z")

			policyRules = append(policyRules, policyAPI.TrustedImagesPolicyRule{
				//Action: action,
				AllowedGroups: allowedGroups,
				Collections:   collections,
				DeniedGroups:  deniedGroups,
				Effect:        schemaRule.Effect.ValueString(),
				Modified:      &modified,
				Name:          schemaRule.Name.ValueString(),
				Notes:         schemaRule.Notes.ValueString(),
				Owner:         schemaRule.Owner.ValueString(),
				PreviousName:  schemaRule.PreviousName.ValueString(),
			})
		}
	}

	rules.Rules = policyRules

	util.LogDebug(ctx, "exiting schemaToTrustedImagesPolicyRules")

	return rules, diags
}

func trustedImagesPolicyToSchema(ctx context.Context, policy policyAPI.TrustedImages) (models.TrustedImagesPolicyResourceModel, diag.Diagnostics) {
	util.LogDebug(ctx, "entering trustedImagesPolicyToSchema")

	schemaPolicy := models.TrustedImagesPolicyResourceModel{}

	schemaGroups, diags := trustedImagesPolicyGroupsToSchema(ctx, policy.Groups)
	if diags.HasError() {
		return schemaPolicy, diags
	}
	schemaPolicy.Groups = &schemaGroups

	schemaRules, diags := trustedImagesPolicyRulesToSchema(ctx, policy.Policy)
	if diags.HasError() {
		return schemaPolicy, diags
	}
	schemaPolicy.Policy = schemaRules

	util.LogDebug(ctx, "exiting trustedImagesPolicyToSchema")

	return schemaPolicy, diags
}

func trustedImagesPolicyGroupsToSchema(ctx context.Context, groups []policyAPI.TrustGroup) ([]models.TrustGroupResourceModel, diag.Diagnostics) {
	util.LogDebug(ctx, "entering trustedImagesPolicyGroupsToSchema")

	var diags diag.Diagnostics

	schemaGroups := []models.TrustGroupResourceModel{}

	for _, group := range groups {
		images, diags := types.SetValueFrom(ctx, types.StringType, group.Images)
		if diags.HasError() {
			return schemaGroups, diags
		}

		schemaGroup := models.TrustGroupResourceModel{
			Id:     types.StringValue(group.Id),
			Images: images,
			//Modified: types.StringValue(group.Modified),
			Modified:     types.StringValue(""),
			Name:         types.StringValue(group.Name),
			Owner:        types.StringValue(group.Owner),
			PreviousName: types.StringValue(group.PreviousName),
		}

		schemaGroups = append(schemaGroups, schemaGroup)
	}

	util.LogDebug(ctx, "exiting trustedImagesPolicyGroupsToSchema")

	return schemaGroups, diags
}

func trustedImagesPolicyRulesToSchema(ctx context.Context, rules policyAPI.TrustedImagesPolicy) (models.TrustedImagesPolicyRulesResourceModel, diag.Diagnostics) {
	util.LogDebug(ctx, "entering trustedImagesPolicyRulesToSchema")

	var diags diag.Diagnostics

	schemaPolicy := models.TrustedImagesPolicyRulesResourceModel{
		Id:      types.StringValue(rules.Id),
		Enabled: types.BoolValue(rules.Enabled),
	}

	if len(rules.Rules) == 0 {
		emptyRules := []models.TrustedImagesPolicyRuleResourceModel{}
		schemaPolicy.Rules = &emptyRules
		return schemaPolicy, diags
	}

	schemaRules := []models.TrustedImagesPolicyRuleResourceModel{}
	for _, rule := range rules.Rules {
		//action, diags := types.SetValueFrom(ctx, types.StringType, rule.Action)
		//if diags.HasError() {
		//    return schemaPolicy, diags
		//}

		var allowedGroups basetypes.SetValue
		if len(rule.AllowedGroups) == 0 {
			allowedGroups, diags = basetypes.NewSetValue(types.StringType, []attr.Value{})
		} else {
			allowedGroups, diags = types.SetValueFrom(ctx, types.StringType, rule.AllowedGroups)
		}

		if diags.HasError() {
			return schemaPolicy, diags
		}

		var deniedGroups basetypes.SetValue
		if len(rule.DeniedGroups) == 0 {
			deniedGroups, diags = basetypes.NewSetValue(types.StringType, []attr.Value{})
		} else {
			deniedGroups, diags = types.SetValueFrom(ctx, types.StringType, rule.DeniedGroups)
		}

		if diags.HasError() {
			return schemaPolicy, diags
		}

		collectionNames := []string{}
		for _, collection := range rule.Collections {
			collectionNames = append(collectionNames, collection.Name)
		}

		collections, diags := types.SetValueFrom(ctx, types.StringType, collectionNames)
		if diags.HasError() {
			return schemaPolicy, diags
		}

		schemaRule := models.TrustedImagesPolicyRuleResourceModel{
			//Action: action,
			AllowedGroups: allowedGroups,
			Collections:   collections,
			DeniedGroups:  deniedGroups,
			Effect:        types.StringValue(rule.Effect),
			Modified:      types.StringValue(""),
			Name:          types.StringValue(rule.Name),
			Notes:         types.StringValue(rule.Notes),
			Owner:         types.StringValue(rule.Owner),
			PreviousName:  types.StringValue(rule.PreviousName),
		}

		schemaRules = append(schemaRules, schemaRule)
	}

	schemaPolicy.Rules = &schemaRules

	util.LogDebug(ctx, "exiting trustedImagesPolicyRulesToSchema")

	return schemaPolicy, diags
}
