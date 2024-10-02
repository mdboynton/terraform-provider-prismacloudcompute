package system 

import (
	"fmt"
    "strings"
	"net/http"
    "sort"
    "strconv"

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
    //Templates
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
	var ans Vulnerabilities 
	if err := c.Request(http.MethodGet, VulnerabilitiesEndpoint, nil, nil, &ans); err != nil {
		return ans, fmt.Errorf("error getting vulnerabilities: %s", err)
	}
	return ans, nil
}

func GetComplianceHostVulnerabilities(c api.PrismaCloudComputeAPIClient) ([]Vulnerability, error) {
    // TODO: include custom compliance checks
    var complianceHostVulns []Vulnerability
    
    vulnerabilities, err := GetVulnerabilities(c)
    if err != nil {
		return complianceHostVulns, fmt.Errorf("error getting host compliance vulnerabilities: %s", err)
    }

    for _, vuln := range vulnerabilities.ComplianceVulnerabilities {
        if vuln.Type == "host_config" ||
            vuln.Type == "windows" ||
            vuln.Type == "linux" ||
            vuln.Type == "security_operations" ||
            vuln.Type == "daemon_config" ||
            vuln.Type == "daemon_config_files" ||
            strings.HasSuffix(vuln.Type, "_worker") ||
            strings.HasSuffix(vuln.Type, "_master") ||
            strings.HasSuffix(vuln.Type, "_federation") {
                complianceHostVulns = append(complianceHostVulns, vuln)
            }
    }
    
    sort.Slice(complianceHostVulns, func(i, j int) bool {
        val1 := strconv.Itoa(complianceHostVulns[i].Id)
        val2 := strconv.Itoa(complianceHostVulns[j].Id)
        return val1 < val2
    })

    return complianceHostVulns, nil
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
