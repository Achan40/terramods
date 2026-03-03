# NAT Gateways sit in public subnets and allow instances in private subnets to access the internet for updates,
# etc. without exposing them to inbound traffic from the internet. 
# Each NAT Gateway is associated with an Elastic IP (EIP) for outbound traffic.
resource "aws_eip" "nat" {
  for_each = aws_subnet.public
  domain   = "vpc"
  tags     = { Name = "nat-eip-${each.key}" }
}

resource "aws_nat_gateway" "nat" {
  for_each      = aws_subnet.public
  allocation_id = aws_eip.nat[each.key].id
  subnet_id     = each.value.id
  tags          = { Name = "nat-${each.key}" }
  depends_on    = [aws_internet_gateway.igw]
}