# only need to include if you want a shared remote state 
include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/oidc"
}

inputs = {
  github_repo = "Achan40/terramods"
}