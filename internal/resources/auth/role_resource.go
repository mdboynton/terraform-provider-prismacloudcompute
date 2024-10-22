package auth

import (
	"context"
    "fmt"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api/auth"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/util"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
    //"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.GetSchema()
}

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan RoleResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    role, diags := schemaToRole(ctx, &plan, getSchemaToRolePermissionsMap())
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Create new role 
    err := auth.CreateRole(*r.client, role)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error creating Role resource", 
            fmt.Sprintf("Failed to create role resource: %s", err.Error()),
        )
        return
	}

    response, err := auth.GetRole(*r.client, role.Name)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error retrieving created Role resource", 
            fmt.Sprintf("Failed to retrieve created role: %s", err.Error()),
        )
        return
	}

    createdRole, diags := roleToSchema(ctx, *response, getRoleToSchemaPermissionsMap())
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set state to plan data
    diags = resp.State.Set(ctx, createdRole)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state RoleResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get role value from Prisma Cloud
    role, err := auth.GetRole(*r.client, state.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading Role resource", 
            fmt.Sprintf("Failed to retrieve role with name \"%s\": %s", state.Name.ValueString(), err.Error()),
        )
        return
    }
  
    // Convert role value to schema
    state, diags = roleToSchema(ctx, *role, getRoleToSchemaPermissionsMap()) 
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan RoleResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    role, diags := schemaToRole(ctx, &plan, getSchemaToRolePermissionsMap())
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    util.DLog(ctx, fmt.Sprintf("schemaToRole in Update() returned with value:\n\n %+v", role))

    // Update existing role 
	err := auth.UpdateRole(*r.client, role)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error updating Role resource", 
            fmt.Sprintf("Failed to update role: %s", err.Error()),
        )
        return
	}

    // Fetch updated role from Prisma Cloud
    updatedRole, err := auth.GetRole(*r.client, plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating Role resource", 
            fmt.Sprintf("Failed to retieve role with name \"%s\": %s", plan.Name.ValueString(), err.Error()),
        )
        return
    }

    util.DLog(ctx, fmt.Sprintf("updatedRole: \n\n %+v", updatedRole))

    // Convert updated user to schema
    plan, diags = roleToSchema(ctx, *updatedRole, getRoleToSchemaPermissionsMap())
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    util.DLog(ctx, fmt.Sprintf("setting state from Update() with value:\n\n %+v", plan))
    util.DLog(ctx, fmt.Sprintf("%+v", *plan.Permissions))

    // Set updated state
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
	var state RoleResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    
    // Delete existing role 
    role := state.Name.ValueString()
    err := auth.DeleteRole(*r.client, role)
	if err != nil {
		resp.Diagnostics.AddError(
            "Error deleting Role resource", 
            fmt.Sprintf("Failed to delete role: %s", err.Error()),
        )
        return
	}
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RoleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    util.DLog(ctx, "entering ModifyPlan")
    
    //var plan *RoleResourceModel
    //diags := req.Plan.Get(ctx, &plan)
    //resp.Diagnostics.Append(diags...)
    //if resp.Diagnostics.HasError() {
    //    return
    //}

    //if plan == nil {
    //    return
    //}

    ////if plan.Permissions == nil {
    ////    resp.Plan.SetAttribute(ctx, path.Root("permissions"), &[]RolePermissionResourceModel{})
    ////}
    
    //userPermissionExists := false
    //for _, permission := range *plan.Permissions {
    //    if permission.Name.ValueString() == "user" {
    //        userPermissionExists = true
    //        break
    //    }
    //}

    //if !userPermissionExists {
    //    *plan.Permissions = append(*plan.Permissions, RolePermissionResourceModel{
    //        Name: types.StringValue("user"),
    //        ReadWrite: types.BoolValue(true),
    //    })
    //    diags.Append(resp.Plan.SetAttribute(ctx, path.Root("permissions"), &plan.Permissions)...)
    //    util.DLog(ctx, "no user permission")
    //}

    util.DLog(ctx, "exiting ModifyPlan")
}

func getSchemaToRolePermissionsMap() map[string]string {
    return map[string]string{
        "cloud_radar": "radarsCloud",
        "hosts_radar": "radarsHosts",
        "containers_radar": "radarsContainers",
        "serverless_radar": "radarsServerless",
        "containers_policies": "policyContainers",
        "hosts_policies": "policyHosts",
        "serverless_policies": "policyServerless",
        "custom_compliance_policies": "policyComplianceCustomRules",
        "container_runtime_policy": "policyRuntimeContainer",
        "host_runtime_policy": "policyRuntimeHosts",
        "serverless_app_embedded_runtime_policy": "policyRuntimeServerless",
        "cnnf_policies": "policyCNNF",
        "waas_policies": "policyWAAS",
        "secrets_policies": "policyAccessSecrets",
        "kubernetes_admission_policies": "policyAccessKubernetes",
        "custom_rules": "policyCustomRules",
        "vulnerabilities_dashboard": "monitorVuln",
        "compliance_dashboard": "monitorCompliance",
        "runtime_dashboards": "monitorRuntimeIncidents",
        "containers_results": "monitorImages",
        "hosts_results": "monitorHosts",
        "serverless_app_embedded_results": "monitorServerless",
        "ci_results": "monitorCI",
        "container_runtime_results": "monitorRuntimeContainers",
        "host_runtime_results": "monitorRuntimeHosts",
        "serverless_app_embedded_runtime_results": "monitorRuntimeServerless",
        "image_analysis_sandbox": "sandbox",
        "waas_events": "monitorWAAS",
        "cnns_runtime_results": "monitorCNNF",
        "docker_runtime_results": "monitorAccessDocker",
        "kubernetes_admission_runtime_results": "monitorAccessKubernetes",
        "client_browser_data_updates": "uIEventSubscriber",
        "cloud_account_policy": "policyCloud",
        "cloud_discovery_results": "monitorCloud",
        "logs": "systemLogs",
        "defenders_management": "manageDefenders",
        "alerts": "manageAlerts",
        "collections_and_tags": "collections",
        "credentials_store": "manageCreds",
        "authentication_configuration": "authConfiguration",
        "user_management": "userManagement",
        "system": "systemOperations",
        "system_privileged": "privilegedOperations",
        "utilities": "downloads",
        "access_console_UI": "accessUI",
        "user": "user",
    }
}

func getRoleToSchemaPermissionsMap() map[string]string {
    return map[string]string{
        "radarsCloud": "cloud_radar",
        "radarsHosts": "hosts_radar",
        "radarsContainers": "containers_radar",
        "radarsServerless": "serverless_radar",
        "policyContainers": "containers_policies",
        "policyHosts": "hosts_policies",
        "policyServerless": "serverless_policies",
        "policyComplianceCustomRules": "custom_compliance_policies",
        "policyRuntimeContainer": "container_runtime_policy",
        "policyRuntimeHosts": "host_runtime_policy",
        "policyRuntimeServerless": "serverless_app_embedded_runtime_policy",
        "policyCNNF": "cnnf_policies",
        "policyWAAS": "waas_policies",
        "policyAccessSecrets": "secrets_policies",
        "policyAccessKubernetes": "kubernetes_admission_policies",
        "policyCustomRules": "custom_rules",
        "monitorVuln": "vulnerabilities_dashboard",
        "monitorCompliance": "compliance_dashboard",
        "monitorRuntimeIncidents": "runtime_dashboards",
        "monitorImages": "containers_results",
        "monitorHosts": "hosts_results",
        "monitorServerless": "serverless_app_embedded_results",
        "monitorCI": "ci_results",
        "monitorRuntimeContainers": "container_runtime_results",
        "monitorRuntimeHosts": "host_runtime_results",
        "monitorRuntimeServerless": "serverless_app_embedded_runtime_results",
        "sandbox": "image_analysis_sandbox",
        "monitorWAAS": "waas_events",
        "monitorCNNF": "cnns_runtime_results",
        "monitorAccessDocker": "docker_runtime_results",
        "monitorAccessKubernetes": "kubernetes_admission_runtime_results",
        "uIEventSubscriber": "client_browser_data_updates",
        "policyCloud": "cloud_account_policy",
        "monitorCloud": "cloud_discovery_results",
        "systemLogs": "logs",
        "manageDefenders": "defenders_management",
        "manageAlerts": "alerts",
        "collections": "collections_and_tags",
        "manageCreds": "credentials_store",
        "authConfiguration": "authentication_configuration",
        "userManagement": "user_management",
        "systemOperations": "system",
        "privilegedOperations": "system_privileged",
        "downloads": "utilities",
        "accessUI": "access_console_UI",
        "user": "user",
    }
}

func schemaToRole(ctx context.Context, plan *RoleResourceModel, permissionsMap map[string]string) (auth.Role, diag.Diagnostics) {
    util.DLog(ctx, "entering schemaToRole")

    var diags diag.Diagnostics

    role := auth.Role{
        Name: plan.Name.ValueString(),
        Description: plan.Description.ValueString(),
        System: plan.System.ValueBool(),
    }

    permissions := []auth.RolePermission{
        auth.RolePermission{
            Name: "user",
            ReadWrite: true,
        },
    }

    for _, planPermission := range *plan.Permissions {
        planPermissionName := planPermission.Name.ValueString()
        if permissionName, ok := permissionsMap[planPermissionName]; ok {
            permission := auth.RolePermission{
                Name: permissionName,
                ReadWrite: planPermission.ReadWrite.ValueBool(),
            }
            permissions = append(permissions, permission)    
        } else {
            diags.AddError(
                "Invalid Role Permission Attribute",
                fmt.Sprintf("Unknown permission \"%s\" specified for role \"%s\".", planPermissionName, plan.Name.ValueString()),
            )
        }
    }
    role.Permissions = permissions 

    util.DLog(ctx, "exiting schemaToRole")
	
    return role, diags 
}

func roleToSchema(ctx context.Context, role auth.Role, permissionsMap map[string]string) (RoleResourceModel, diag.Diagnostics) {
    util.DLog(ctx, "entering roleToSchema")

    var diags diag.Diagnostics

    schema := RoleResourceModel{
        Name: types.StringValue(role.Name),
        Description: types.StringValue(role.Description),
        System: types.BoolValue(role.System),
    }

    permissions := []RolePermissionResourceModel{}

    for _, schemaPermission := range role.Permissions {
        schemaPermissionName := schemaPermission.Name

        if schemaPermissionName == "user" {
            continue
        }

        if permissionName, ok := permissionsMap[schemaPermissionName]; ok {
            permission := RolePermissionResourceModel {
                Name: types.StringValue(permissionName),
                ReadWrite: types.BoolValue(schemaPermission.ReadWrite),
            }
            permissions = append(permissions, permission)
        } else {
            diags.AddError(
                "Invalid Role Permission Attribute",
                fmt.Sprintf("API response returned unknown permission name for role \"%s\": %s", role.Name, schemaPermissionName),
            )
        }
    }
    schema.Permissions = &permissions
    
    util.DLog(ctx, "exiting roleToSchema")

    return schema, diags
}
