package test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceSecurityGroup(t *testing.T) {
	t.Parallel()

	// Step 1: deploy a VPC to host the security group
	vpcOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/private_vpc",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, vpcOptions)
	terraform.InitAndApply(t, vpcOptions)

	vpcID := terraform.Output(t, vpcOptions, "vpc_id")

	// Step 2: deploy the cluster security group
	sgOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/instance_security_group",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, sgOptions)
	terraform.InitAndApply(t, sgOptions)

	sgID := terraform.Output(t, sgOptions, "security_group_id")
	assert.NotEmpty(t, sgID)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	require.NoError(t, err)

	ec2Client := ec2.NewFromConfig(cfg)

	out, err := ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{sgID},
	})
	require.NoError(t, err)
	require.Len(t, out.SecurityGroups, 1)

	sg := out.SecurityGroups[0]

	// Security group belongs to the correct VPC
	assert.Equal(t, vpcID, *sg.VpcId)

	// Has a self-referencing ingress rule allowing all traffic
	hasSelfRef := false
	for _, perm := range sg.IpPermissions {
		if aws.ToInt32(perm.FromPort) == 0 && aws.ToInt32(perm.ToPort) == 0 && aws.ToString(perm.IpProtocol) == "-1" {
			for _, pair := range perm.UserIdGroupPairs {
				if aws.ToString(pair.GroupId) == sgID {
					hasSelfRef = true
				}
			}
		}
	}
	assert.True(t, hasSelfRef, "expected a self-referencing ingress rule allowing all traffic")
}
