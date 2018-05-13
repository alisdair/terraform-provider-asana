# terraform-provider-asana

This is a [HashiCorp Terraform][terraform] resource provider for
[Asana][asana]. Probably not a good idea.

[terraform]: https://www.terraform.io
[asana]: https://asana.com

## Usage

This is super early! But here's the basics:

1. [Create a personal access token in Asana][personal-access-token]
1. Store it in an environment variable `ASANA_PERSONAL_ACCESS_TOKEN`, preferably using [envchain][envchain]
1. Figure out your team or workspace ID (*FIXME* how?)
1. Write a Terraform configuration, something like this:

    ```hcl
    resource "asana_project" "asana" {
      name = "My Cool Project"
      workspace = "123456789"
      color = "dark-teal"
    }
    ```

1. Run `envchain asana terraform plan` and `envchain asana terraform apply` to create or update your project!
1. Run `envchain asana terraform destroy` to remove it once you're done.

[personal-acccess-token]: https://asana.com/guide/help/api/api#gl-access-tokens
[envchain]: https://github.com/sorah/envchain
