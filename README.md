# Dynatrace Terraform Provider

- Website: [https://www.terraform.io](https://www.terraform.io)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

## Requirements

- [Terraform] 0.12.1+
- [Go] 1.13 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/dynatrace-ace/terraform-provider-dynatrace`

```sh
mkdir -p $GOPATH/src/github.com/dynatrace-ace; cd $GOPATH/src/github.com/dynatrace-ace
git clone https://github.com/dynatrace-ace/terraform-provider-dynatrace.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/dynatrace-ace/terraform-provider-dynatrace
make build
```

## Using the provider

If you want to run Terraform with the dynatrace provider plugin on your system, complete the following steps:

1. [Download and install Terraform for your system].

1. [Download the dynatrace provider plugin for Terraform].

1. Unzip the release archive to extract the plugin binary (`terraform-provider-dynatrace_vX.Y.Z`).

1. Move the binary into the Terraform [plugins directory] for the platform.
    - Linux/Unix/macOS: `~/.terraform.d/plugins`
    - Windows: `%APPDATA%\terraform.d\plugins`

1. Export dynatrace environment URL and API token as environment variables. Generate the [API token] with the following scope:

    - Read Configuration
    - Write Configuration

    ```sh
    Managed

    export DYNATRACE_ENV_URL="https://{your-domain}/e/{your-environment-id}"
    export DYNATRACE_API_TOKEN="{your-dynatrace-api-token}"
    ```

    ```sh
    SaaS

    export DYNATRACE_ENV_URL=" https://{your-environment-id}.live.dynatrace.com"
    export DYNATRACE_API_TOKEN="{your-dynatrace-api-token}"
    ```

1. Add the plug-in provider to the Terraform configuration file.

    ```go
    provider "dynatrace" {}
    ```

1. Alternatively you may also define your dynatrace environment URL and API token directly on the Terraform configuration file.

    ```sh
    Managed

    provider "dynatrace" {
        dt_env_url    = "https://{your-domain}/e/{your-environment-id}"
        dt_api_token  = "{your-dynatrace-api-token}"
    }
    ```

    ```sh
    SaaS

    provider "dynatrace" {
        dt_env_url    = "https://{your-environment-id}.live.dynatrace.com"
        dt_api_token  = "{your-dynatrace-api-token}"
    }
    ```

## Developing the Provider

To contribute to the provider, [Go](http://www.golang.org) is required to be installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
make build

$GOPATH/bin/terraform-provider-dynatrace
```

To directly add provider to the plugins directory in Linux, `~/.terraform/plugins` run make install.

```sh
make install

~/.terraform/plugins/terraform-provider-dynatrace
```

## Known limitations

For management zones, dynamic json fields like `value` in the ComparisonBasic object only accept map[string]interface{} data types. Dynamic fields like `DynamicKey` in the ConditionKey Object are currently not supported.

[Terraform]: (https://www.terraform.io/downloads.html)
[Go]: (https://golang.org/doc/install)
[Download and install Terraform for your system]: (https://www.terraform.io/intro/getting-started/install.html)
[Download the dynatrace provider plugin for Terraform]: (https://github.com/dynatrace-ace/terraform-provider-dynatrace/releases)
[plugins directory]: (https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
[API token]: https://www.dynatrace.com/support/help/dynatrace-api/basics/dynatrace-api-authentication/
