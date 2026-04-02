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

variable "subnet_id" {
  type        = string
  description = "ID of the subnet to deploy the instance into."
}

variable "eice_security_group_id" {
  type        = string
  description = "ID of the EC2 Instance Connect Endpoint security group to allow SSH ingress from"
}

variable "tags" {
  type        = map(string)
  description = "Additional tags to apply to all resources"
  default     = {}
}
