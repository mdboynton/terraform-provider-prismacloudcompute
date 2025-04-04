package acceptance

import (
	"fmt"
    "os"
	"testing"
    "strings"
    "slices"
    "regexp"
    "errors"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/runtime"
)

var (
    providerConfig string
)

func testAccPreCheck(t *testing.T) {
	//if os.Getenv(configFileEnvVar) == "" {
	//	t.Fatalf("%s must be set for acceptance tests", configFileEnvVar)
	//}
    
    //config, err := p.GetConfigurationFromFile(credentialsFilePath)
    //if err != nil {
    //    // TODO: error message
    //    return
    //}
   
    //if config.ConsoleURL == nil {
    //    t.Fatalf("ConsoleURL is nil")
    //}

    //if config.Username == nil {
    //    t.Fatalf("Username is nil")
    //}

    //if config.Password == nil {
    //    t.Fatalf("Password is nil")
    //}

    //if config.RequestTimeout == nil {
    //    t.Fatalf("RequestTimeout is nil")
    //}

    consoleUrl := os.Getenv("PRISMACLOUDCOMPUTE_CONSOLE_URL") 
    username := os.Getenv("PRISMACLOUDCOMPUTE_USERNAME")
    password := os.Getenv("PRISMACLOUDCOMPUTE_PASSWORD")
    providerConfig = fmt.Sprintf(`
        provider "prismacloudcompute" {
            console_url = "%s"
            username = "%s"
            password = "%s"
            insecure = true 
            request_timeout = 60
        }
    `, consoleUrl, username, password)
}

func generateEffectsAttributesTestSteps(policyType string, providerConfig string) ([]resource.TestStep, error) {
    var (
        resourceConfig string
        attributeAssignment string
        testSteps []resource.TestStep = []resource.TestStep{}
    )

    if !slices.Contains([]string{"host", "container", "serverless", "app_embedded"}, policyType) {
        return testSteps, errors.New(
            fmt.Sprintf("Invalid policy type specified: %s", policyType),
        )
    }

    for attribute, effects := range policy.ValidEffects[policyType] {
        splitAttribute := strings.Split(attribute, ".")
        nestedAttribute := splitAttribute[0]
        
        if len(splitAttribute) == 2 {
            attributeAssignment = fmt.Sprintf("%s = \"invalid\"", splitAttribute[1])
        } else if len(splitAttribute) == 3 {
            attributeAssignment = fmt.Sprintf("%s = { %s = \"invalid\" }", splitAttribute[1], splitAttribute[2])
        }

        resourceConfig = fmt.Sprintf(`
            resource "prismacloudcompute_%s_runtime_policy" "accTestEffects" {
                rules = [
                    {
                        name = "effectsTestRule"
                        %s = {
                            %s
                        }
                    }
                ]
            }
        `, policyType, nestedAttribute, attributeAssignment)

        testSteps = append(testSteps, resource.TestStep{
            Config: providerConfig + resourceConfig,
            ExpectError: regexp.MustCompile(fmt.Sprintf("Invalid value \"%s\" specified for attribute\nrules\\[0\\]\\.%s\\.\nMust be one of the following: %s", "invalid", attribute, strings.Join(effects[:], ", "))),
        })
    }

    return testSteps, nil
}
