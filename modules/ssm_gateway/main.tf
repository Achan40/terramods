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

# look up the latest amazon linux 2023 ami - ssm agent is pre-installed
data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "state"
    values = ["available"]
  }
}

# iam role so the ec2 instance can communicate with ssm
resource "aws_iam_role" "ssm_gateway" {
  name = "${var.name}-ssm-gateway-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# attach the aws managed policy that enables ssm session manager
resource "aws_iam_role_policy_attachment" "ssm_core" {
  role       = aws_iam_role.ssm_gateway.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# instance profile wraps the role so it can be attached to ec2
resource "aws_iam_instance_profile" "ssm_gateway" {
  name = "${var.name}-ssm-gateway-profile"
  role = aws_iam_role.ssm_gateway.name
}

# security group - no inbound needed, only outbound 443 for ssm
resource "aws_security_group" "ssm_gateway" {
  name        = "${var.name}-ssm-gateway-sg"
  description = "Security group for SSM gateway instance"
  vpc_id      = var.vpc_id

  egress {
    description = "Allow HTTPS outbound for SSM endpoints"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ec2 instance in a private subnet - no public ip needed
resource "aws_instance" "ssm_gateway" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = var.instance_type
  subnet_id              = var.subnet_id
  iam_instance_profile   = aws_iam_instance_profile.ssm_gateway.name
  vpc_security_group_ids = [aws_security_group.ssm_gateway.id]

  associate_public_ip_address = false

  metadata_options {
    http_tokens = "required"
  }

  tags = {
    Name = "${var.name}-ssm-gateway"
  }
}