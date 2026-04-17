include "common" {
  path = find_in_parent_folders("common.hcl")
}

locals {
  common = read_terragrunt_config(find_in_parent_folders("common.hcl"))
  region = local.common.locals.region
}

terraform {
  source = "../../../modules/ec2_instance"
}

dependency "vpc" {
  config_path = "../private_vpc"
}

dependency "iam" {
  config_path = "../iam_instance_role"
}

dependency "sg" {
  config_path = "../instance_security_group"
}

inputs = {
  region        = local.region
  instance_name = "test-ec2-instance"
  ami_id        = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  instance_count = 2

  # Required: sourced from private_vpc module outputs.
  eice_security_group_id = dependency.vpc.outputs.eice_security_group_id
  vpc_id                 = dependency.vpc.outputs.vpc_id
  subnet_ids             = dependency.vpc.outputs.private_subnet_ids

  # IAM instance profile sourced from iam_instance_role module.
  iam_instance_profile_name = dependency.iam.outputs.instance_profile_name

  # Shared cluster SG sourced from instance_security_group module.
  additional_security_group_ids = [dependency.sg.outputs.security_group_id]

  # If supplied connects to tailscale as a subnet router, allowing access to other AWS services on the same VPC without them having to connect to tailscale
  tailscale_auth_key = get_env("TAILSCALE_AUTH_KEY", "") 
}
