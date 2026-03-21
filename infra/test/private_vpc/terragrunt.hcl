include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/private_vpc"
}

# Inputs for tests handled in go module
inputs = {

}