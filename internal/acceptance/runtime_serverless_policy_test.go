package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
)

// TODO: Defaults test
// TODO: Rule ordering/sorting test

func TestAccServerlessRuntimePolicy_EmptyRules(t *testing.T) {
    resourceName := "prismacloudcompute_serverless_runtime_policy.accTestEmptyRules"
    resourceConfig := `
        resource "prismacloudcompute_serverless_runtime_policy" "accTestEmptyRules" {
            rules = []
        }
    `

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: []resource.TestStep{
            {
                Config: providerConfig + resourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
                ),
            },
        },
    })
}

func TestAccServerlessRuntimePolicy_Base(t *testing.T) {
    resourceName1 := "prismacloudcompute_serverless_runtime_policy.accTestBase1"
    resourceConfig1 := `
        resource "prismacloudcompute_serverless_runtime_policy" "accTestBase1" {
            rules = [
                {
                    order = 999
                    name = "baseTestRule1"
                    collections = ["All"]
                    disabled = true
                    notes = "test notes"
                    advanced_threat_protection = true
                    networking = {
                        ip_connectivity_enabled = true
                        allowed_listening_ports = ["100", "200-300"]
                        allowed_outbound_internet_ports = ["400", "500-600"]
                        allowed_outbound_ips = ["0.0.0.0", "1.1.1.1"]
                        denied_ips_ports_effect = "prevent"
                        dns_enabled = true
                        allowed_dns_domains = ["allowed1.com", "allowed2.com"]
                        denied_dns_domains_effect = "prevent"
                    }
                    processes = {
                        enabled = true
                        allowed_processes = []
                        denied_processes_effect = "prevent"
                        crypto_miners = true
                        block_all_processes_except_main = true 
                    }
                    file_system = {
                        enabled = true
                        allowed_paths = ["/allowed/path/1", "/allowed/path/2"]
                        denied_paths_effect = "prevent"
                    }
                    //custom_rules = [
                    //    {
                    //        name = "test123"
                    //        //id = 61
                    //        //log_as = "incident"
                    //        //effect = "alert"
                    //        effect = "allow"
                    //        //log_as = "incident"
                    //    },
                    //]
                }
            ]
        }
    `

    // Changes:
    //      processes: 
    //          allowed_processes non-empty
    //          block_all_processes_except_main = false
    //      file_system:
    //          allowed_paths empty
    //          denied_paths non-empty
    resourceName2 := "prismacloudcompute_serverless_runtime_policy.accTestBase2"
    resourceConfig2 := `
        resource "prismacloudcompute_serverless_runtime_policy" "accTestBase2" {
            rules = [
                {
                    order = 999
                    name = "baseTestRule2"
                    collections = ["All"]
                    disabled = true
                    notes = "test notes"
                    networking = {
                        ip_connectivity_enabled = true
                        allowed_listening_ports = ["100", "200-300"]
                        allowed_outbound_internet_ports = ["400", "500-600"]
                        allowed_outbound_ips = ["0.0.0.0", "1.1.1.1"]
                        denied_ips_ports_effect = "prevent"
                        dns_enabled = true
                        allowed_dns_domains = ["allowed1.com", "allowed2.com"]
                        denied_dns_domains_effect = "prevent"
                    }
                    processes = {
                        enabled = true
                        allowed_processes = ["bash", "ssh"]
                        denied_processes_effect = "prevent"
                        crypto_miners = true
                        block_all_processes_except_main = false
                    }
                    file_system = {
                        enabled = true
                        denied_paths = ["/denied/path/1", "/denied/path/2"]
                        denied_paths_effect = "prevent"
                    }
                    //custom_rules = [
                    //    {
                    //        name = "test123"
                    //        //id = 61
                    //        //log_as = "incident"
                    //        //effect = "alert"
                    //        effect = "allow"
                    //        //log_as = "incident"
                    //    },
                    //]
                }
            ]
        }
    `

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: []resource.TestStep{
            {
                Config: providerConfig + resourceConfig1,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName1, "rules.#", "1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.order", "999"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.name", "baseTestRule1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.collections.#", "1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.collections.0", "All"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.disabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.advanced_threat_protection", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.notes", "test notes"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.owner", os.Getenv(provider.UsernameEnvVar)),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.previous_name", ""),
                    // Networking 
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.ip_connectivity_enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_listening_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_listening_ports.0", "100"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_listening_ports.1", "200-300"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_internet_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_internet_ports.0", "400"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_internet_ports.1", "500-600"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_ips.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_ips.0", "0.0.0.0"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_outbound_ips.1", "1.1.1.1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_ips_ports_effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.dns_enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.0", "allowed1.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.1", "allowed2.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_dns_domains_effect", "prevent"),
                    // Processes
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allowed_processes.#", "0"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.denied_processes_effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.crypto_miners", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.block_all_processes_except_main", "true"),
                    // File System
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.0", "/allowed/path/1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.1", "/allowed/path/2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.denied_paths_effect", "prevent"),
                    // TODO: custom rules
                ),
            },
            {
                Config: providerConfig + resourceConfig2,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName2, "rules.#", "1"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.order", "999"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.name", "baseTestRule2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.collections.#", "1"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.collections.0", "All"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.disabled", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.advanced_threat_protection", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.notes", "test notes"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.owner", os.Getenv(provider.UsernameEnvVar)),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.previous_name", ""),
                    // Networking 
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.ip_connectivity_enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_listening_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_listening_ports.0", "100"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_listening_ports.1", "200-300"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_internet_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_internet_ports.0", "400"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_internet_ports.1", "500-600"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_ips.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_ips.0", "0.0.0.0"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_outbound_ips.1", "1.1.1.1"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.denied_ips_ports_effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.dns_enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_dns_domains.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_dns_domains.0", "allowed1.com"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.allowed_dns_domains.1", "allowed2.com"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.networking.denied_dns_domains_effect", "prevent"),
                    // Processes
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.allowed_processes.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.allowed_processes.0", "bash"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.allowed_processes.1", "ssh"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.denied_processes_effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.crypto_miners", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.processes.block_all_processes_except_main", "false"),
                    // File System
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.file_system.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.file_system.denied_paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.file_system.denied_paths.0", "/denied/path/1"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.file_system.denied_paths.1", "/denied/path/2"),
                    resource.TestCheckResourceAttr(resourceName2, "rules.0.file_system.denied_paths_effect", "prevent"),
                    // TODO: custom rules
                ),
            },
        },
    })
}

func TestAccServerlessRuntimePolicy_Effects(t *testing.T) {
    testSteps, err := generateEffectsAttributesTestSteps("serverless", providerConfig)
    if err != nil {
        t.Error(err)
    }

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: testSteps,
    })
}
