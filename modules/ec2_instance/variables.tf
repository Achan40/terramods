variable "region" {
  type        = string
  description = "AWS region"
}

variable "instance_name" {
  type        = string
  description = "Name tag for the EC2 instance"
}

variable "ami_id" {
  type        = string
  description = "AMI ID to use for the instance"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t3.micro"
}

variable "vpc_id" {
  type        = string
  description = "ID of the VPC to deploy the instance into."
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of subnet IDs to deploy instances into. Instances are distributed across subnets in round-robin order."
}

variable "instance_count" {
  type        = number
  description = "Number of EC2 instances to deploy."
  default     = 1
}

variable "eice_security_group_id" {
  type        = string
  description = "ID of the EC2 Instance Connect Endpoint security group to allow SSH ingress from"
}

variable "iam_instance_profile_name" {
  type        = string
  description = "Name of the IAM instance profile to attach to the instance"
}

variable "additional_security_group_ids" {
  type        = list(string)
  description = "Additional security group IDs to attach to the instance (e.g. a shared cluster SG)"
  default     = []
}

variable "tags" {
  type        = map(string)
  description = "Additional tags to apply to all resources"
  default     = {}
}
