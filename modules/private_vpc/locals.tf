# Calculations for subnet CIDR blocks based on the VPC CIDR and availability zones
locals {
  public_subnet_cidrs = {
    for idx, az in var.availability_zones :
    az => cidrsubnet(var.vpc_cidr, 8, idx)
  }
  private_subnet_cidrs = {
    for idx, az in var.availability_zones :
    az => cidrsubnet(var.vpc_cidr, 8, idx + 10)
  }
}