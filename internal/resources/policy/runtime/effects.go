package policy

var ValidEffects = map[string]map[string][]string{
    "host": {
        "anti_malware.crypto_miners": []string{"disable", "alert", "prevent"},
        "anti_malware.denied_processes.effect": []string{"alert", "prevent"},
        "anti_malware.encrypted_binaries": []string{"disable", "alert"},
        "anti_malware.execution_flow_hijacking": []string{"disable", "alert"},
        "anti_malware.malware_from_advanced_threat_protection": []string{"disable", "alert"},
        "anti_malware.malware_from_custom_feed": []string{"disable", "alert"},
        "anti_malware.non_packaged_binaries_service": []string{"disable", "alert", "prevent"},
        "anti_malware.non_packaged_binaries_user": []string{"disable", "alert", "prevent"},
        "anti_malware.processes_temporary_storage": []string{"disable", "alert", "prevent"},
        "anti_malware.reverse_shell": []string{"disable", "alert"},
        "anti_malware.suspicious_elf_headers": []string{"disable", "alert"},
        "anti_malware.web_shell": []string{"disable", "alert", "prevent"},
        "anti_malware.wild_fire_analysis": []string{"disable", "alert"},
    },
}
