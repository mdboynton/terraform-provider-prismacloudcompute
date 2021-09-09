package prismacloudcompute

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	pcc "github.com/paloaltonetworks/prisma-cloud-compute-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"console_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Prisma Cloud Compute Console URL",
				DefaultFunc: schema.EnvDefaultFunc("PRISMACLOUDCOMPUTE_CONSOLE_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prisma Cloud Compute username",
				DefaultFunc: schema.EnvDefaultFunc("PRISMACLOUDCOMPUTE_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prisma Cloud Compute password",
				DefaultFunc: schema.EnvDefaultFunc("PRISMACLOUDCOMPUTE_PASSWORD", nil),
				Sensitive:   true,
			},
			"skip_cert_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not to skip certificate verification",
				Default:     true,
			},
			"config_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Configuration file in JSON format. See examples/creds.json",
				DefaultFunc: schema.EnvDefaultFunc("PRISMACLOUDCOMPUTE_CONFIG_FILE", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"prismacloudcompute_collection":                    resourceCollection(),
			"prismacloudcompute_ci_image_compliance_policy":    resourcePoliciesComplianceCiImages(),
			"prismacloudcompute_container_compliance_policy":   resourcePoliciesComplianceContainer(),
			"prismacloudcompute_host_compliance_policy":        resourcePoliciesComplianceHost(),
			"prismacloudcompute_container_runtime_policy":      resourcePoliciesRuntimeContainer(),
			"prismacloudcompute_host_runtime_policy":           resourcePoliciesRuntimeHost(),
			"prismacloudcompute_ci_image_vulnerability_policy": resourcePoliciesVulnerabilityCiImages(),
			"prismacloudcompute_host_vulnerability_policy":     resourcePoliciesVulnerabilityHost(),
			"prismacloudcompute_image_vulnerability_policy":    resourcePoliciesVulnerabilityImages(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"prismacloudcompute_collections":                   dataSourceCollections(),
			"prismacloudcompute_ci_image_compliance_policy":    dataSourcePoliciesComplianceCiImages(),
			"prismacloudcompute_container_compliance_policy":   dataSourcePoliciesComplianceContainer(),
			"prismacloudcompute_host_compliance_policy":        dataSourcePoliciesComplianceHost(),
			"prismacloudcompute_container_runtime_policy":      dataSourcePoliciesRuntimeContainer(),
			"prismacloudcompute_host_runtime_policy":           dataSourcePoliciesRuntimeHost(),
			"prismacloudcompute_ci_image_vulnerability_policy": dataSourcePoliciesVulnerabilityCiImages(),
			"prismacloudcompute_host_vulnerability_policy":     dataSourcePoliciesVulnerabilityHost(),
			"prismacloudcompute_image_vulnerability_policy":    dataSourcePoliciesVulnerabilityImages(),
		},

		ConfigureFunc: configure,
	}
}

func configure(d *schema.ResourceData) (interface{}, error) {
	var config pcc.Credentials
	if d.Get("config_file") != nil {
		configFile, err := os.Open(d.Get("config_file").(string))
		if err != nil {
			fmt.Printf("error opening config file: %v", err)
		}
		defer configFile.Close()

		fileContent, err := ioutil.ReadAll(configFile)
		if err != nil {
			fmt.Printf("error reading config file: %v", err)
			return nil, err
		}
		if err := json.Unmarshal(fileContent, &config); err != nil {
			fmt.Printf("error unmarshalling config file: %v", err)
			return nil, err
		}
	}

	// if d.Get("console_url") != nil {
	// 	config.ConsoleURL = d.Get("console_url").(string)
	// }
	// if d.Get("username") != nil {
	// 	config.Username = d.Get("username").(string)
	// }
	// if d.Get("password") != nil {
	// 	config.ConsoleURL = d.Get("password").(string)
	// }
	// if d.Get("skip_cert_verification") != nil {
	// 	config.SkipCertVerification = d.Get("skip_cert_verification").(bool)
	// }

	client, err := pcc.APIClient(config.ConsoleURL, config.Username, config.Password, config.SkipCertVerification)
	if err != nil {
		fmt.Printf("failed creating API client")
		return nil, err
	}

	return client, err
}
