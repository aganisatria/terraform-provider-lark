# Terraform Provider Lark

Terraform Provider Lark is a custom provider that enables you to manage and integrate resources from the Lark API using Terraform. With this provider, you can automate the configuration and management of resources such as group chats, user groups, and user roles in Lark.

## Features

### Resource

| Resource | Description |
|---|---|
| lark_group_chat | Create, update, and delete group chats in Lark |
| lark_group_chat_member | Manage members for group chats in Lark |
| lark_user_group | Create, update, and delete user groups in Lark |
| lark_user_group_member | Manage members for user groups in Lark |
| lark_role | Create, update, and delete roles in Lark |
| lark_role_member | Manage members for roles in Lark |

### Data Source

| Data Source | Description |
|---|---|
| lark_user_by_email | Retrieve user data based on email |
| lark_user_by_id | Retrieve user data based on user ID, open ID, or union ID |

## How to use
See [docs](docs/index.md)

## How to contribute
Open an issue or a PR.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Installation

1. golangcli-lint local => brew install golangci-lint
2. goconvey => go install github.com/smartystreets/goconvey
3. tfenv => brew install tfenv
4. terraform => tfenv install latest

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

Created by [@aganisatria](https://github.com/aganisatria)
