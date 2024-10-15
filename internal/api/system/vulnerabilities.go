package system 

import (
	"fmt"
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
	var vulns Vulnerabilities 
	if err := c.Request(http.MethodGet, VulnerabilitiesEndpoint, nil, nil, &vulns); err != nil {
		return vulns, fmt.Errorf("error getting vulnerabilities: %s", err)
	}

	return vulns, nil
}

func GetComplianceVulnerabilitiesMap(c api.PrismaCloudComputeAPIClient) (map[string][]Vulnerability, error) {
    vulnsMap := make(map[string][]Vulnerability)

    vulns, err := GetVulnerabilities(c)
    if err != nil {
        return vulnsMap, err
    }

    for _, vuln := range vulns.ComplianceVulnerabilities {
        if _, ok := vulnsMap[vuln.Type]; ok {
            vulnsMap[vuln.Type] = append(vulnsMap[vuln.Type], vuln)
        } else {
            vulnsMap[vuln.Type] = []Vulnerability{vuln}
        }
    }
    
    return vulnsMap, nil
}

func GetComplianceVulnerabilities(c api.PrismaCloudComputeAPIClient, policyType string) ([]Vulnerability, error) {
    var complianceVulns []Vulnerability
    var vulnTypes []string

    vulnsMap, err := GetComplianceVulnerabilitiesMap(c)
    if err != nil {
		return complianceVulns, fmt.Errorf("error getting vm image compliance vulnerabilities: %s", err)
    }

    switch policyType {
        case "hostCompliance":
            vulnTypes = getHostComplianceVulnTypes()
        case "containerCompliance":
            vulnTypes = getContainerComplianceVulnTypes()
        case "ciImagesCompliance":
            //vulnTypes = []string{"image"}
            break
        case "vmCompliance":
            vulnTypes = getVmImageComplianceVulnTypes()
        default:
            return complianceVulns, fmt.Errorf("invalid compliance policy type supplied: \"%s\"", policyType)
    }


    if policyType == "ciImagesCompliance" {
        for _, vuln := range vulnsMap["image"] {
            if vuln.Id == 406 || vuln.Id == 408 || vuln.Id == 41 || vuln.Id == 422 || vuln.Id == 424 ||
                vuln.Id == 425 || vuln.Id == 426 || vuln.Id == 448 || vuln.Id == 5041 || vuln.Id == 5048 {
                complianceVulns = append(complianceVulns, vuln)    
            }
        }
    } else {
        for _, vulnType := range vulnTypes {
            complianceVulns = append(complianceVulns, vulnsMap[vulnType]...)
        }
    }

    sort.Slice(complianceVulns, func(i, j int) bool {
        val1 := strconv.Itoa(complianceVulns[i].Id)
        val2 := strconv.Itoa(complianceVulns[j].Id)
        return val1 < val2
    })

    return complianceVulns, nil
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

func getHostComplianceVulnTypes() []string {
    return []string{
        "host_config", 
        "daemon_config", 
        "daemon_config_files", 
        "security_operations", 
        "linux", 
        "windows", 
        "k8s_worker",
        "eks_worker",
        "aks_worker",
        "openshift_worker",
        "k8s_master",
        "openshift_master",
        "k8s_federation",
    }
}

func getContainerComplianceVulnTypes() []string {
    return []string{
        "container", 
        "istio", 
        "image",
    }
}

func getVmImageComplianceVulnTypes() []string {
    return []string{
        "host", 
        "host_config", 
        "daemon_config", 
        "daemon_config_files", 
        "security_operations", 
        "linux",
    }
}
