include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/private_vpc"
}

# Inputs for tests handled in go module. The below inputs are used in case you want to run terragrunt directly for testing.
inputs = {
  vpc_name           = "test-private-vpc"
  region             = "us-east-2"
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-east-2b"]
}