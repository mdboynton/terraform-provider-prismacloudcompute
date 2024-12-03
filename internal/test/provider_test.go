package test 

import (
	//"fmt"
	//"os"
	//"testing"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	PrismacloudcomputeJsonConfigFileEnvVar = "PRISMACLOUDCOMPUTE_CONFIG_FILE"

    providerConfig = `
provider "prismacloudcompute" {
    console_url = ""
    username = ""
    password = ""
    insecure = true
}
`
)

var (
	//testAccProviders                   map[string]*schema.Provider
	//testAccProvider                    *schema.Provider
	//sessionTimeoutOrig, sessionTimeout int

    testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
)

func init() {
    testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error) {
        "prismacloudcompute": providerserver.NewProtocol6WithError(provider.New("test")()), 
    }
}
