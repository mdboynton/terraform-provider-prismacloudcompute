package acceptance

import (
    "fmt"
	//"context"
	"os"
	"regexp"
	"testing"
    "strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
	"github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/resources/policy/runtime"
)


func TestAccHostRuntimePolicy_EmptyRules(t *testing.T) {
    resourceConfig := `
        resource "prismacloudcompute_host_runtime_policy" "accTestEmptyRules" {
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
                    resource.TestCheckResourceAttr("prismacloudcompute_host_runtime_policy.accTestEmptyRules", "rules.#", "0"),
                ),
            },
        },
    })
}

func TestAccHostRuntimePolicy_Base(t *testing.T) {
    resourceName := "prismacloudcompute_host_runtime_policy.accTestBase"
    resourceConfig := `
        resource "prismacloudcompute_host_runtime_policy" "accTestBase" {
            rules = [
                {
                    order = 999
                    name = "baseTestRule"
                    collections = ["All"]
                    disabled = true
                    notes = "test notes"
                    anti_malware = {
                        allowed_processes = ["/path/to/process", "/other/path"]
                        denied_processes = {
                            effect = "prevent"
                            paths = ["/another/path", "/yet/another"]
                        }
                        crypto_miners = "alert"
                        non_packaged_binaries_service = "disable"
                        non_packaged_binaries_user = "prevent"
                        processes_temporary_storage = "prevent"
                        web_shell = "prevent"
                        reverse_shell = "disable"
                        execution_flow_hijacking = "disable"
                        encrypted_binaries = "disable"
                        suspicious_elf_headers = "disable"
                        malware_from_custom_feed = "disable"
                        malware_from_advanced_threat_protection = "disable"
                    }
                    networking = {
                        allowed_outbound_ips = ["0.0.0.0", "3.3.3.3"]
                        suspicious_ips_custom_feed = "disable"
                        denied_listening_ports = ["555", "666"]
                        denied_outbound_ips = ["1.1.1.1", "2.2.2.2"]
                        denied_outbound_ports = ["222", "333"]
                        denied_ips_ports_effect = "disable"
                        suspicious_ips_advanced_threat_protection_effect = "disable"
                        allowed_dns_domains = ["test.com"]
                        denied_dns_domains = ["bad.com"]
                        denied_dns_domains_effect = "disable"
                        suspicious_domains_advanced_threat_protection_effect = "disable"
                    }
                    log_inspection_rules = [
                        {
                            path = "/path/to/file"
                            regex = ["test1", "test2"]
                        },
                        {
                            path = "/path/to/another/file"
                            regex = ["test3", "test4"]
                        },
                    ]
                    file_integrity_rules = [
                        {
                            file_path = "/path/to/file"
                            allowed_processes = ["bash"]
                            monitor_write_ops = true
                        },
                    ]
                    activities = {
                        host_activity_monitoring = {
                            enabled = true
                            docker_commands = {
                                enabled = true
                                include_read_only_events = true
                            }
                            sshd_sessions = true
                            sudo_commands = true
                            log_background_apps = true
                        }
                        track_ssh_events = true
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
                Config: providerConfig + resourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.order", "999"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.name", "baseTestRule"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.collections.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.collections.0", "All"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.disabled", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.notes", "test notes"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.owner", os.Getenv(provider.UsernameEnvVar)),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.previous_name", ""),
                    // Activities
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.docker_commands.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.docker_commands.include_read_only_events", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.enabled", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.log_background_apps", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.sshd_sessions", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.host_activity_monitoring.sudo_commands", "true"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.activities.track_ssh_events", "true"),
                    // Anti-malware
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.allowed_processes.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.allowed_processes.0", "/path/to/process"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.allowed_processes.1", "/other/path"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.crypto_miners", "alert"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.denied_processes.effect", "prevent"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.denied_processes.paths.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.denied_processes.paths.0", "/another/path"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.denied_processes.paths.1", "/yet/another"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.encrypted_binaries", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.execution_flow_hijacking", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.malware_from_advanced_threat_protection", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.malware_from_custom_feed", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.non_packaged_binaries_service", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.non_packaged_binaries_user", "prevent"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.processes_temporary_storage", "prevent"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.reverse_shell", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.suppress_compiler_generated_binaries", "false"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.suspicious_elf_headers", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.web_shell", "prevent"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.anti_malware.wild_fire_analysis", "alert"),
                    // File Integrity Rules 
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.allowed_processes.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.allowed_processes.0", "bash"),
                    resource.TestCheckNoResourceAttr(resourceName, "rules.0.file_integrity_rules.0.excluded_file_patterns"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.file_path", "/path/to/file"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.monitor_metadata_changes", "false"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.monitor_read_ops", "false"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.monitor_subdirectories", "false"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.file_integrity_rules.0.monitor_write_ops", "true"),
                    // Log Inspection Rules
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.0.path", "/path/to/file"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.0.regex.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.0.regex.0", "test1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.0.regex.1", "test2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.1.path", "/path/to/another/file"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.1.regex.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.1.regex.0", "test3"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.log_inspection_rules.1.regex.1", "test4"),
                    // Networking
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.allowed_dns_domains.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.allowed_dns_domains.0", "test.com"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.allowed_outbound_ips.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.allowed_outbound_ips.0", "0.0.0.0"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.allowed_outbound_ips.1", "3.3.3.3"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_dns_domains.#", "1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_dns_domains.0", "bad.com"),

                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_dns_domains_effect", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_ips_ports_effect", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_listening_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_listening_ports.0", "555"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_listening_ports.1", "666"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_outbound_ips.0", "1.1.1.1"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_outbound_ips.1", "2.2.2.2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_outbound_ports.#", "2"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_outbound_ports.0", "222"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.denied_outbound_ports.1", "333"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.suspicious_domains_advanced_threat_protection_effect", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.suspicious_ips_advanced_threat_protection_effect", "disable"),
                    resource.TestCheckResourceAttr(resourceName, "rules.0.networking.suspicious_ips_custom_feed", "disable"),
                    // TODO: custom rules
                ),
            },
        },
    })
}

func TestAccHostRuntimePolicy_Effects(t *testing.T) {
    var (
        resourceConfig string
        attributeAssignment string
        testSteps []resource.TestStep = []resource.TestStep{}
    )

    for attribute, effects := range policy.ValidEffects["host"] {
        splitAttribute := strings.Split(attribute, ".")
        nestedAttribute := splitAttribute[0]
        
        if len(splitAttribute) == 2 {
            attributeAssignment = fmt.Sprintf("%s = \"invalid\"", splitAttribute[1])
        } else if len(splitAttribute) == 3 {
            attributeAssignment = fmt.Sprintf("%s = { %s = \"invalid\" }", splitAttribute[1], splitAttribute[2])
        }

        resourceConfig = fmt.Sprintf(`
            resource "prismacloudcompute_host_runtime_policy" "accTestBase" {
                rules = [
                    {
                        name = "effectsTestRule"
                        %s = {
                            %s
                        }
                    }
                ]
            }
        `, nestedAttribute, attributeAssignment)

        testSteps = append(testSteps, resource.TestStep{
            Config: providerConfig + resourceConfig,
            ExpectError: regexp.MustCompile(fmt.Sprintf("Invalid value \"%s\" specified for attribute\nrules\\[0\\]\\.%s\\.\nMust be one of the following: %s", "invalid", attribute, strings.Join(effects[:], ", "))),
        })
    }

    resource.UnitTest(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: protoV6ProviderFactories(),
        Steps: testSteps,
    })
}
