# Terraform Provider for Prisma Cloud Compute
[//]: # (You can find the Prisma Cloud Compute provider in the [Terraform Registry](https://registry.terraform.io/providers/PaloAltoNetworks/prismacloudcompute/latest).)

## Local Installation
1. Install requirements
    - Go v1.23.0+
    - GNU Make
2. Unzip provider files
3. Navigate to provider files directory in terminal
4. Run `make` (this will run the `install` step by default)
5. Populate provider data (see below) into `providers.tf`
6. Run `terraform init`

## Basic setup
```terraform
terraform {
  required_providers {
    prismacloudcompute = {
      source  = "registry.terraform.io/PaloAltoNetworks/prismacloudcompute"
      version = "1.0.0-alpha"
    }
  }
}

provider "prismacloudcompute" {
  console_url = "https://console.example.com" # Do not include trailing slash
  username = "username"
  password = "password"
  insecure = true

  # Alternatively, you can specify a file with the configuration data specified above 
  # config_file = "/path/to/config.json"
}
```

[//]: # (Complete documentation can be found in the [marketplace listing](https://registry.terraform.io/providers/PaloAltoNetworks/prismacloudcompute/latest/docs).)

## Contributing
Contributions are welcome!
Please read the [contributing guide](CONTRIBUTING.md) for more information.

## Support
Please read our [support document](SUPPORT.md) for details on how to get support for this project.
