# github_oidc

Create the github odic provider so that CI/CD works correctly.

<!-- BEGIN_TF_DOCS -->


## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | 6.10.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.10.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_ci_cd_policy_name"></a> [ci\_cd\_policy\_name](#input\_ci\_cd\_policy\_name) | Name for policy used in ci cd | `string` | n/a | yes |
| <a name="input_ci_cd_role_name"></a> [ci\_cd\_role\_name](#input\_ci\_cd\_role\_name) | Name for the IAM role used in ci cd | `string` | n/a | yes |
| <a name="input_github_repo"></a> [github\_repo](#input\_github\_repo) | Github repository name | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_oidc_provider_arn"></a> [oidc\_provider\_arn](#output\_oidc\_provider\_arn) | ARN of the oidc provider |
| <a name="output_oidc_role_arn"></a> [oidc\_role\_arn](#output\_oidc\_role\_arn) | ARN of the CI/CD role for github |
<!-- END_TF_DOCS -->