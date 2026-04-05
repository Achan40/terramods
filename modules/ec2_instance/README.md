<!-- BEGIN_TF_DOCS -->
Deploy any number of EC2 machines to a vpc. Accessed via EICE or by tailscale once linked.

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
| <a name="input_ami_id"></a> [ami\_id](#input\_ami\_id) | AMI ID to use for the instance | `string` | n/a | yes |
| <a name="input_eice_security_group_id"></a> [eice\_security\_group\_id](#input\_eice\_security\_group\_id) | ID of the EC2 Instance Connect Endpoint security group to allow SSH ingress from | `string` | n/a | yes |
| <a name="input_iam_instance_profile_name"></a> [iam\_instance\_profile\_name](#input\_iam\_instance\_profile\_name) | Name of the IAM instance profile to attach to the instance | `string` | n/a | yes |
| <a name="input_instance_name"></a> [instance\_name](#input\_instance\_name) | Name tag for the EC2 instance | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | AWS region | `string` | n/a | yes |
| <a name="input_subnet_ids"></a> [subnet\_ids](#input\_subnet\_ids) | List of subnet IDs to deploy instances into. Instances are distributed across subnets in round-robin order. | `list(string)` | n/a | yes |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | ID of the VPC to deploy the instance into. | `string` | n/a | yes |
| <a name="input_additional_security_group_ids"></a> [additional\_security\_group\_ids](#input\_additional\_security\_group\_ids) | Additional security group IDs to attach to the instance (e.g. a shared cluster SG) | `list(string)` | `[]` | no |
| <a name="input_instance_count"></a> [instance\_count](#input\_instance\_count) | Number of EC2 instances to deploy. | `number` | `1` | no |
| <a name="input_instance_type"></a> [instance\_type](#input\_instance\_type) | EC2 instance type | `string` | `"t3.micro"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Additional tags to apply to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_instance_ids"></a> [instance\_ids](#output\_instance\_ids) | IDs of the EC2 instances |
| <a name="output_private_ips"></a> [private\_ips](#output\_private\_ips) | Private IP addresses of the EC2 instances |
| <a name="output_security_group_id"></a> [security\_group\_id](#output\_security\_group\_id) | ID of the shared instance security group |
<!-- END_TF_DOCS -->