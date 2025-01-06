package system 

import (
	"fmt"
	"net/http"
    "sort"
    "strconv"

    "github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/api"
)

const VulnerabilitiesEndpoint = "api/v1/static/vulnerabilities"

type Vulnerabilities struct {
	ComplianceVulnerabilities []Vulnerability `json:"complianceVulnerabilities,omitempty"`
	CveVulnerabilities []Vulnerability `json:"cveVulnerabilities,omitempty"`
}

type Vulnerability struct {
    //Text
    Id int `json:"id"`
    Severity string `json:"severity"`
    //CVSS
    //Status
    //CVE
    //Cause
    //Description
    //Title
    //VecStr
    //Exploit
    //RiskFactors
    //Link
    Type string `json:"type"`
    //PackageName
    //PackageVersion
    //PackageType
    //LayerTime
    Templates *[]string `json:"templates"`
    Twistlock bool `json:"twistlock"`
    CRI bool `json:"cri"`
    //Published
    //FixDate
    //Discovered
    //FunctionLayer
    //WildfireMalware
    //Secret
}

func GetVulnerabilities(c api.PrismaCloudComputeAPIClient) (Vulnerabilities, error) {
	var vulns Vulnerabilities 
	if err := c.Request(http.MethodGet, VulnerabilitiesEndpoint, nil, nil, &vulns); err != nil {
		return vulns, fmt.Errorf("error getting vulnerabilities: %s", err)
	}

	return vulns, nil
}

func GetComplianceVulnerabilitiesMap(c api.PrismaCloudComputeAPIClient) (map[string][]Vulnerability, diag.Diagnostics) {
    var diags diag.Diagnostics

    vulnsMap := make(map[string][]Vulnerability)

    vulns, err := GetVulnerabilities(c)
    if err != nil {
        diags.AddError("API Request Error", fmt.Sprintf("Error occured while retrieving vulnerabilites from Prisma Cloud API: %s", err.Error()))
        return vulnsMap, diags
    }

    for _, vuln := range vulns.ComplianceVulnerabilities {
        if _, ok := vulnsMap[vuln.Type]; ok {
            vulnsMap[vuln.Type] = append(vulnsMap[vuln.Type], vuln)
        } else {
            vulnsMap[vuln.Type] = []Vulnerability{vuln}
        }
    }
    
    return vulnsMap, diags 
}

func GetComplianceVulnerabilitiesByPolicyType(c api.PrismaCloudComputeAPIClient, policyTypeFilter func (Vulnerability) bool) ([]Vulnerability, diag.Diagnostics) {
    var diags diag.Diagnostics

    vulnerabilities, err := GetVulnerabilities(c)
    if err != nil {
        diags.AddError("API Request Error", fmt.Sprintf("Error occured while retrieving vulnerabilites from Prisma Cloud API: %s", err.Error()))
        return []Vulnerability{}, diags
    }

    filteredComplianceVulnerabilities := []Vulnerability{}
    
    for _, complianceVulnerability := range vulnerabilities.ComplianceVulnerabilities {
        if policyTypeFilter(complianceVulnerability) {
            filteredComplianceVulnerabilities = append(filteredComplianceVulnerabilities, complianceVulnerability)
        }
    }
    
    sort.Slice(filteredComplianceVulnerabilities, func(i, j int) bool {
        val1 := strconv.Itoa(filteredComplianceVulnerabilities[i].Id)
        val2 := strconv.Itoa(filteredComplianceVulnerabilities[j].Id)
        return val1 < val2
    })

    return filteredComplianceVulnerabilities, diags
}

func GetHighOrCriticalVulnerabilities(complianceVulnerabilities []Vulnerability) []Vulnerability {
    var highOrCriticalVulns []Vulnerability 
    for _, vuln := range complianceVulnerabilities {
        if vuln.Severity == "high" || vuln.Severity == "critical" {
            highOrCriticalVulns = append(highOrCriticalVulns, vuln)
        }
    }
    return highOrCriticalVulns
}
