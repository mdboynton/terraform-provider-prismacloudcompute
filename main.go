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

    opts := providerserver.ServeOpts {
        Address:    "registry.terraform.io/PaloAltoNetworks/prismacloudcompute",
        Debug:      debug,
    }

    err := providerserver.Serve(context.Background(), p.New(version), opts)
        
    if err != nil {
        log.Fatal(err)
    }
}
