include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/ssm_gateway"
}

dependency "private_vpc" {
  config_path = "../private_vpc"
}

# Inputs for tests handled in go module. The below inputs are used in case you want to run terragrunt directly for testing.
inputs = {
  name          = "test-ssm-gateway"
  vpc_id        = dependency.private_vpc.outputs.vpc_id
  subnet_id     = dependency.private_vpc.outputs.private_subnet_ids[0]
  instance_type = "t3.micro"
}
