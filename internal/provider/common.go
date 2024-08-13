package provider

const (
	policyTypeAdmission               = "admission"
	policyTypeComplianceCiImage       = "ciImagesCompliance"
	policyTypeComplianceCoderepo      = "codeRepoCompliance"
	policyTypeComplianceCiCoderepo    = "ciCodeRepoCompliance"
	policyTypeComplianceContainer     = "containerCompliance"
	policyTypeComplianceHost          = "hostCompliance"
    policyTypeComplianceTrustedImages = "trust"
	policyTypeRuntimeContainer        = "containerRuntime"
	policyTypeRuntimeHost             = "hostRuntime"
	policyTypeVulnerabilityCiCoderepo = "ciCodeRepoVulnerability"
	policyTypeVulnerabilityCiImage    = "ciImagesVulnerability"
	policyTypeVulnerabilityCoderepo   = "codeRepoVulnerability"
	policyTypeVulnerabilityHost       = "hostVulnerability"
	policyTypeVulnerabilityImage      = "containerVulnerability"
)
