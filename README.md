## Local Installation
1. Install requirements
    - Go v1.23.0+
    - GNU Make
2. Clone this branch (or download the source as a zip and extract the contents)
3. Navigate to the repo directory in your terminal
4. Run `make` (this will run the `install` step by default)
    - This assumes that Terraform will look in the default location for providers (`~/.terraform.d/plugins`)
5. After the provider has been built and moved to the correct location, populate provider data (see below) into `providers.tf`
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

## Contributing
Contributions are welcome!
Please read the [contributing guide](CONTRIBUTING.md) for more information.

## Support
Please read our [support document](SUPPORT.md) for details on how to get support for this project.
