variable "name" {
  type        = string
  description = "Name prefix for all resources"
}

variable "vpc_id" {
  type        = string
  description = "ID of the VPC to deploy the gateway into"
}

variable "subnet_id" {
  type        = string
  description = "ID of the private subnet to deploy the gateway into"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type for the SSM gateway"
  default     = "t3.micro"
}
