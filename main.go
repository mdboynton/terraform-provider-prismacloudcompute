package main

import (
    "context"
    "flag"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"

	p "github.com/PaloAltoNetworks/terraform-provider-prismacloudcompute/internal/provider"
)

var (
    version string = "dev"
)

func main() {
    var debug bool

    flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

    err := providerserver.Serve(context.Background(), p.New(version), providerserver.ServeOpts {
        Address:    "registry.terraform.io/PaloAltoNetworks/terraform-provider-prismacloudcompute",
        Debug:      debug,
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
