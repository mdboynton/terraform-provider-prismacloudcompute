## Local Installation

### MacOS/Linux
1. Install requirements
    - Go v1.23.0+
    - GNU Make
2. Clone this branch (or download the source as a zip and extract the contents)
3. Navigate to the repo directory in your terminal
4. Run `make` (this will run the `install` step by default)
    - This assumes that Terraform will look in the default location for providers (`~/.terraform.d/plugins`)
5. After the provider has been built and moved to the correct location, populate provider data (see "Basic Setup" section below) into `providers.tf`
6. Run `terraform init`

### Windows [WIP]
1. Install requirements
    - Go v1.23.0+
2. Clone this branch (or download the source as a zip and extract the contents)
3. Navigate to the repo directory in cmd or PowerShell 
4. Run `go build -o terraform-provider-prismacloudcompute` to compile the provider
5. Move the resulting binary to the default provider location (`%APPDATA%\terraform.d\plugins`)
6. Populate provider data (see "Basic Setup" section below) into `providers.tf`
7. Run `terraform init`

## Basic setup

Be sure to include the API port (8083 by default) and trailing slash in console URL. 

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
  console_url = "https://192.168.0.01:8083/"
  username = "username"
  password = "password"
  insecure = true
  request_timeout = 60
}
```

Alternatively, you can specify a file in the provider configuration that contains the configuration data:

```terraform
provider "prismacloudcompute" {
  config_file = "/path/to/config.json"

  # Example config:
  #
  # {
  #     "console_url": "https://192.168.0.1:8083/",
  #     "username": "admin",
  #     "password": "password",
  #     "insecure": true,
  #     "request_timeout": 60
  # }
}
```

## Contributing
Contributions are welcome!
Please read the [contributing guide](CONTRIBUTING.md) for more information.

## Support
Please read our [support document](SUPPORT.md) for details on how to get support for this project.
