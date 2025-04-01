package acceptance

import (
	"fmt"
    "os"
	"testing"
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
