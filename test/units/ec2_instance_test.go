package test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEC2InstanceCustomVPC deploys an EC2 instance into a subnet within a
// custom VPC, and validates that the instance is associated with that VPC.
func TestEC2InstanceCustomVPC(t *testing.T) {
	t.Parallel()

	// Step 1: deploy a VPC to use as the target
	vpcOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/private_vpc",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, vpcOptions)
	terraform.InitAndApply(t, vpcOptions)

	vpcID := terraform.Output(t, vpcOptions, "vpc_id")
	subnetIDs := terraform.OutputList(t, vpcOptions, "private_subnet_ids")
	require.NotEmpty(t, subnetIDs)

	// Step 2: deploy IAM instance role
	iamOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/iam_instance_role",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, iamOptions)
	terraform.InitAndApply(t, iamOptions)

	instanceProfileName := terraform.Output(t, iamOptions, "instance_profile_name")

	// Step 3: deploy shared cluster security group
	sgOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/instance_security_group",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, sgOptions)
	terraform.InitAndApply(t, sgOptions)

	clusterSGID := terraform.Output(t, sgOptions, "security_group_id")

	// Step 4: deploy EC2 into the first private subnet of the VPC
	ec2Options := &terraform.Options{
		TerraformDir:    "../../examples/units/ec2_instance",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, ec2Options)
	terraform.InitAndApply(t, ec2Options)

	instanceIDs := terraform.OutputList(t, ec2Options, "instance_ids")
	require.Len(t, instanceIDs, 1)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	require.NoError(t, err)

	ec2Client := ec2.NewFromConfig(cfg)

	out, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	require.NoError(t, err)
	require.Len(t, out.Reservations, 1)
	require.Len(t, out.Reservations[0].Instances, 1)

	instance := out.Reservations[0].Instances[0]

	// Instance is running
	assert.Equal(t, ec2types.InstanceStateNameRunning, instance.State.Name)

	// Instance is in the expected VPC
	assert.Equal(t, vpcID, *instance.VpcId)

	// Instance is in one of the expected subnets
	assert.Contains(t, subnetIDs, *instance.SubnetId)

	// No public IP assigned
	assert.Nil(t, instance.PublicIpAddress)

	// IAM instance profile is attached
	require.NotNil(t, instance.IamInstanceProfile)
	assert.Contains(t, *instance.IamInstanceProfile.Arn, instanceProfileName)

	// Cluster security group is attached
	sgIDs := make([]string, 0, len(instance.SecurityGroups))
	for _, sg := range instance.SecurityGroups {
		sgIDs = append(sgIDs, *sg.GroupId)
	}
	assert.Contains(t, sgIDs, clusterSGID)
}
