variable "ci_cd_role_name" {
  type = string
  description = "Name for the IAM role used in ci cd "
}

variable "ci_cd_policy_name" {
  type = string
  description = "Name for policy used in ci cd"
}

variable "github_repo" {
  type = string
  description = "Github repository name"
}