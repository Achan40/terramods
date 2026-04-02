include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/ec2_instance"
}

dependency "vpc" {
  config_path = "../private_vpc"
}

# Inputs for tests handled in go module. The below inputs are used in case you want to run terragrunt directly for testing.
inputs = {
  region        = "us-east-2"
  instance_name = "test-ec2-instance"
  ami_id        = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  # Required: sourced from private_vpc module outputs.
  eice_security_group_id = dependency.vpc.outputs.eice_security_group_id
  vpc_id                 = dependency.vpc.outputs.vpc_id
  subnet_id              = dependency.vpc.outputs.private_subnet_ids[0]
}
