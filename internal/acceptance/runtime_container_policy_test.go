package acceptance

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
)

const (
    resourceType string = "prismacloudcompute_container_runtime_policy"
)

func TestAccContainerRuntimePolicy_EmptyRules(t *testing.T) {
    resourceName := fmt.Sprintf("%s.accTestEmptyRules", resourceType)
    resourceConfig := `
        resource "prismacloudcompute_container_runtime_policy" "accTestEmptyRules" {
            automatic_runtime_learning = false
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
                    resource.TestCheckResourceAttr(resourceName, "automatic_runtime_learning", "false"),
                    resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
                ),
            },
        },
    })
}

func TestAccContainerRuntimePolicy_Base(t *testing.T) {
    resourceName1 := fmt.Sprintf("%s.accTestBase1", resourceType)
    resourceConfig1 := `
        resource "prismacloudcompute_container_runtime_policy" "accTestBase1" {
            automatic_runtime_learning = false
            rules = [
                {
                    order = 999
                    name = "baseTestRule"
                    collections = ["All"]
                    disabled = true
                    notes = "test notes"
                    anti_malware = {
                        malware_from_advanced_threat_protection = "block"
                        kubernetes_attacks = "prevent"
                        suspicious_cloud_provider_api_queries = "block"
                        wild_fire_analysis = "block"
                    }
                    processes = {
                        enabled = true
                        allow_only_learned_processes_from_parents = true 
                        allow_all_activity_in_attached_sessions = true
                        allowed_processes = ["ssh", "grep"]
                        processes_from_modified_binaries = "prevent"
                        crypto_miners = "prevent"
                        reverse_shell = "block"
                        lateral_movement_processes = "block"
                        processes_started_with_suid = "prevent"
                        denied_processes = {
                            effect = "block"
                            paths = ["netstat", "ipconfig"]
                        }
                        all_other_processes_effect = "prevent"
                    }
                    networking = {
                        ip_connectivity_enabled = true
                        allowed_listening_ports = ["100", "200-300"]
                        allowed_outbound_internet_ports = ["400", "500-600"]
                        allowed_outbound_ips = ["0.0.0.0", "1.1.1.1"]
                        network_activity_from_modified_binaries = "block"
                        port_scanning = "block"
                        raw_sockets = "disable"
                        denied_listening_ports = ["700", "800-900"]
                        denied_listening_ports_effect = "block"
                        denied_outbound_internet_ports = ["1000", "1100-1200"]
                        denied_outbound_internet_ports_effect = "block"
                        denied_outbound_ips = ["2.2.2.2", "3.3.3.3"]
                        denied_outbound_ips_effect = "block"
                        all_other_activity_effect = "block"
                        dns_enabled = true
                        allowed_dns_domains = ["allowed1.com", "allowed2.com"]
                        denied_dns_domains = ["denied1.com", "denied2.com"]
                        denied_dns_domains_effect = "prevent"
                        all_other_domains_effect = "prevent"
                    }
                    file_system = {
                        enabled = true
                        allowed_paths = ["/allowed/path/1", "/allowed/path/2"]
                        changes_to_binaries = "block"
                        detection_of_encrypted_binaries = "block"
                        changes_to_ssh_admin_account_config_files = "block"
                        suspicious_elf_headers = "block"
                        denied_paths = ["/denied/path/1", "/denied/path/2"]
                        denied_paths_effect = "block"
                        all_other_paths_effect = "block"
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
                    resource.TestCheckResourceAttr(resourceName1, "automatic_runtime_learning", "false"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.#", "1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.order", "999"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.name", "baseTestRule"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.collections.#", "1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.collections.0", "All"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.disabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.notes", "test notes"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.owner", os.Getenv(provider.UsernameEnvVar)),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.previous_name", ""),
                    // Anti-malware
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.anti_malware.malware_from_advanced_threat_protection", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.anti_malware.kubernetes_attacks", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.anti_malware.suspicious_cloud_provider_api_queries", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.anti_malware.wild_fire_analysis", "block"),
                    // Processes
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allow_only_learned_processes_from_parents", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allow_all_activity_in_attached_sessions", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allowed_processes.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allowed_processes.0", "ssh"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.allowed_processes.1", "grep"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.processes_from_modified_binaries", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.crypto_miners", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.reverse_shell", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.lateral_movement_processes", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.processes_started_with_suid", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.denied_processes.effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.denied_processes.paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.denied_processes.paths.0", "netstat"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.denied_processes.paths.1", "ipconfig"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.processes.all_other_processes_effect", "prevent"),
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
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.network_activity_from_modified_binaries", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.port_scanning", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.raw_sockets", "disable"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_listening_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_listening_ports.0", "700"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_listening_ports.1", "800-900"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_listening_ports_effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_internet_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_internet_ports.0", "1000"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_internet_ports.1", "1100-1200"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_internet_ports_effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_ips.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_ips.0", "2.2.2.2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_ips.1", "3.3.3.3"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_outbound_ips_effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.all_other_activity_effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.dns_enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.0", "allowed1.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.allowed_dns_domains.1", "allowed2.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_dns_domains.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_dns_domains.0", "denied1.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_dns_domains.1", "denied2.com"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.denied_dns_domains_effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.networking.all_other_domains_effect", "prevent"),
                    // File System
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.0", "/allowed/path/1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.allowed_paths.1", "/allowed/path/2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.changes_to_binaries", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.detection_of_encrypted_binaries", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.changes_to_ssh_admin_account_config_files", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.suspicious_elf_headers", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.denied_paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.denied_paths.0", "/denied/path/1"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.denied_paths.1", "/denied/path/2"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.denied_paths_effect", "block"),
                    resource.TestCheckResourceAttr(resourceName1, "rules.0.file_system.all_other_paths_effect", "block"),
                    // TODO: custom rules
                ),
            },
        },
    })
}

func TestAccContainerRuntimePolicy_Effects(t *testing.T) {
    testSteps, err := generateEffectsAttributesTestSteps("container", providerConfig)
    if err != nil {
        t.Error(err)
    }

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: testSteps,
    })
}


func TestAccContainerRuntimePolicy_ErrorOnSharedProcesses(t *testing.T) {
    resourceConfig := fmt.Sprintf(`
        resource "%s" "accTestSharedProcesses" {
            rules = [
                {
                    name = "sharedProcessesTestRule"
                    collections = ["All"]
                    processes = {
                        enabled = true
                        allowed_processes = ["ssh", "test", "grep"]
                        denied_processes = {
                            paths = ["netstat", "test", "ipconfig"]
                        }
                    }
                }
            ]
        }
    `, resourceType)

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: []resource.TestStep{
            {
                Config: providerConfig + resourceConfig,
                ExpectError: regexp.MustCompile("Invalid configuration at rules\\[0\\]\\.processes\nValues in allowed_processes and denied_processes.paths cannot be shared.\nThe following values must be removed from one of either argument: test"),
            },
        },
    })
}
