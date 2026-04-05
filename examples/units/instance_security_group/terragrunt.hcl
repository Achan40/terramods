include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/instance_security_group"
}

dependency "vpc" {
  config_path = "../private_vpc"
}

# Inputs for tests handled in go module. The below inputs are used in case you want to run terragrunt directly for testing.
inputs = {
  region = "us-east-2"
  name   = "test-cluster"
  vpc_id = dependency.vpc.outputs.vpc_id
}
