variable "region" {
  type        = string
  description = "AWS region"
}

variable "name" {
  type        = string
  description = "Name prefix for the security group"
}

variable "vpc_id" {
  type        = string
  description = "ID of the VPC"
}

variable "tags" {
  type        = map(string)
  description = "Additional tags to apply to all resources"
  default     = {}
}
