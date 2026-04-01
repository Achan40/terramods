# ssm_gateway

Creates a private EC2 instance accessible via AWS Systems Manager Session Manager. No public IP, no open inbound ports, no bastion required.

Use the `security_group_id` output to grant this instance access to other resources in the VPC (e.g. allow inbound on port 9092 from this SG on your Kafka instances).

To access from a local machine:
1. First time only: `brew install --cask session-manager-plugin`
2. You can grab the instance ID directly from Terragrunt if you haven't already: `terragrunt output instance_id`
3. `aws ssm start-session --target <instance_id> --region <region>`

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
| <a name="input_name"></a> [name](#input\_name) | Name prefix for all resources | `string` | n/a | yes |
| <a name="input_subnet_id"></a> [subnet\_id](#input\_subnet\_id) | ID of the private subnet to deploy the gateway into | `string` | n/a | yes |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | ID of the VPC to deploy the gateway into | `string` | n/a | yes |
| <a name="input_instance_type"></a> [instance\_type](#input\_instance\_type) | EC2 instance type for the SSM gateway | `string` | `"t3.micro"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_instance_id"></a> [instance\_id](#output\_instance\_id) | ID of the SSM gateway EC2 instance |
| <a name="output_instance_private_ip"></a> [instance\_private\_ip](#output\_instance\_private\_ip) | Private IP of the SSM gateway EC2 instance |
| <a name="output_security_group_id"></a> [security\_group\_id](#output\_security\_group\_id) | ID of the SSM gateway security group |
<!-- END_TF_DOCS -->
