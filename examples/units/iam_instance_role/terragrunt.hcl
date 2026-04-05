include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/iam_instance_role"
}

# Inputs for tests handled in go module. The below inputs are used in case you want to run terragrunt directly for testing.
inputs = {
  region    = "us-east-2"
  role_name = "test-instance"
}
