# only need to include if you want a shared remote state 
include "common" {
  path = find_in_parent_folders("common.hcl")
}

terraform {
  source = "../../../modules/github_oidc"
}

inputs = {
  ci_cd_role_name = "ci_cd_role"
  ci_cd_policy_name = "ci_cd_policy"
  github_repo = "Achan40/terramods"
}