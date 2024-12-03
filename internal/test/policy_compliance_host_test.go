package test 

import (
    //"context"
    "testing"
    "fmt" 

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
    FullRule1Name string = "testrule1"
    FullRule2Name string = "testrule2"
    FullRule3Name string = "testrule3"
    FullRule2Order string = "2"
    FullRule3Order string = "3"
    FullRule2BlockMessage string = "testrule2 block message"
    FullRule3BlockMessage string = "testrule3 block message"
    FullRule2Collections string = "[\"Security Suite Hosts\"]"
    FullRule2Effect string = "alert"
    FullRule3Effect string = "block"
    FullRule2Notes string = "testrule2 notes"
    FullRule3Notes string = "testrule3 notes"
)

func TestAccHostCompliancePolicy_EmptyRules(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: providerConfig + `
resource "prismacloudcompute_host_compliance_policy" "test" {
    rules = [] 
}`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.#", "0"),
                ),
            },
        },
    })
}

func TestAccHostCompliancePolicy_Full(t *testing.T) {
    t.Parallel()

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                // testrule1: minimum configuration
                // testrule2: full configuration with "alert"
                // testrule2: full configuration with "block"
                Config: fmt.Sprintf(`
%s

resource "prismacloudcompute_host_compliance_policy" "test" {
    rules = [
        {
            name = "%s"
            order = %s
            block_message = "%s"
            disabled = true
            effect = "%s"
            notes = "%s"
            verbose = true
        },
        {
            name = "%s"
            order = %s
            block_message = "%s"
            collections = %s
            effect = "%s"
            notes = "%s"
            verbose = true
        },
        {
            name = "%s"
        }
    ] 
}`, providerConfig, FullRule3Name, FullRule3Order, FullRule3BlockMessage, FullRule3Effect, FullRule3Notes, FullRule2Name, FullRule2Order, FullRule2BlockMessage, FullRule2Collections, FullRule2Effect, FullRule2Notes, FullRule1Name),
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of rules
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.#", "3"),
                    // Verify first rule
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.name", FullRule3Name),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.order", FullRule3Order),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.block_message", FullRule3BlockMessage),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.collections.#", "1"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.collections.0", "All"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.disabled", "true"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.effect", FullRule3Effect),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.notes", FullRule3Notes),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.0.verbose", "true"),
                    // Verify second rule
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.name", FullRule2Name),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.order", FullRule2Order),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.block_message", FullRule2BlockMessage),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.collections.#", "1"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.collections.0", "Security Suite Hosts"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.disabled", "false"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.effect", FullRule2Effect),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.notes", FullRule2Notes),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.1.verbose", "true"),
                    // Verify third rule
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.name", FullRule1Name),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.order", "3"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.block_message", ""),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.collections.#", "1"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.collections.0", "All"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.disabled", "false"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.effect", "alert"),
                    resource.TestCheckNoResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.notes"),
                    resource.TestCheckResourceAttr("prismacloudcompute_host_compliance_policy.test", "rules.2.verbose", "false"),
                ),
            },
//            {
//                Config: providerConfig + `
//removed {
//    from = prismacloudcompute_host_compliance_policy.test
//
//    lifecycle {
//        destroy = false
//    }
//}`,
//            },
        },
    })
}
