terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.10.0"
    }
  }
  backend "s3" {}
}

provider "aws" {
  region = var.region
}

# Shared security group for instance-to-instance communication.
# The self-referencing ingress rule allows all traffic between any instances that share this SG.
# At the moment, this may be a bit too permissive, but should be fine for a development environment. We can always add more specific rules later if needed.
resource "aws_security_group" "instance" {
  name        = "${var.name}-sg"
  description = "Shared security group for instance-to-instance communication"
  vpc_id      = var.vpc_id

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  tags = merge(var.tags, { Name = "${var.name}-sg" })
}
