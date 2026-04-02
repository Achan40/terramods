data "aws_iam_policy_document" "ec2_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "instance" {
  name               = "${var.instance_name}-role"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume_role.json
  tags               = merge(var.tags, { Name = "${var.instance_name}-role" })
}

resource "aws_iam_instance_profile" "instance" {
  name = "${var.instance_name}-profile"
  role = aws_iam_role.instance.name
}
