variable "region" {
  type        = string
  description = "AWS region"
}

variable "role_name" {
  type        = string
  description = "Name prefix for the IAM role and instance profile"
}

variable "tags" {
  type        = map(string)
  description = "Additional tags to apply to all resources"
  default     = {}
}
