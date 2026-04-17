terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.10.0"
    }
  }
  backend "s3" {
  }
}

provider "aws" {
  region = var.region
}

resource "aws_security_group" "instance" {
  name        = "${var.instance_name}-sg"
  description = "Security group for ${var.instance_name}"
  vpc_id      = var.vpc_id

  # only allow SSH from the EICE security group
  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [var.eice_security_group_id]
  }

  # allows all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, { Name = "${var.instance_name}-sg" })
}

resource "aws_instance" "main" {
  count = var.instance_count

  ami                    = var.ami_id
  instance_type          = var.instance_type
  subnet_id              = var.subnet_ids[count.index % length(var.subnet_ids)]
  vpc_security_group_ids = concat([aws_security_group.instance.id], var.additional_security_group_ids)
  iam_instance_profile   = var.iam_instance_profile_name

  source_dest_check = var.tailscale_auth_key != null ? false : true

  user_data = <<-EOF
    #!/bin/bash
    apt-get update -y
    apt-get install -y ec2-instance-connect
    %{if var.tailscale_auth_key != null}
    curl -fsSL https://tailscale.com/install.sh | sh
    echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
    sysctl -p
    IFACE=$(ip route | grep default | awk '{print $5}' | head -1)
    SUBNET=$(ip route | grep "$IFACE" | grep -v default | awk '{print $1}' | head -1)
    tailscale up --authkey=${var.tailscale_auth_key} --advertise-routes=$SUBNET
    %{endif}
  EOF

  tags = merge(var.tags, { Name = "${var.instance_name}-${count.index + 1}" })
}
