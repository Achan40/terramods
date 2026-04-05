resource "aws_security_group" "eice" {
  name        = "${var.vpc_name}-eice-sg"
  description = "Security group for EC2 Instance Connect Endpoints"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.vpc_name}-eice-sg" }
}

resource "aws_ec2_instance_connect_endpoint" "private" {
  for_each = aws_subnet.private

  subnet_id          = each.value.id
  security_group_ids = [aws_security_group.eice.id]

  tags = { Name = "${var.vpc_name}-eice-${each.key}" }
}
