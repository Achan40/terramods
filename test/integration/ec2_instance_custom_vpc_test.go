package test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	uniqueID := fmt.Sprintf("%d", time.Now().Unix())

	// Step 1: deploy a VPC to use as the target
	vpcOptions := &terraform.Options{
		TerraformDir:    "../../infra/test/private_vpc",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"vpc_name":           "test-ec2-vpc-" + uniqueID,
			"region":             "us-east-2",
			"vpc_cidr":           "10.1.0.0/16",
			"availability_zones": []string{"us-east-2a", "us-east-2b"},
		},
	}

	defer terraform.Destroy(t, vpcOptions)
	terraform.InitAndApply(t, vpcOptions)

	vpcID := terraform.Output(t, vpcOptions, "vpc_id")
	subnetIDs := terraform.OutputList(t, vpcOptions, "private_subnet_ids")
	eiceSGID := terraform.Output(t, vpcOptions, "eice_security_group_id")
	require.NotEmpty(t, subnetIDs)

	// Step 2: deploy EC2 into the first private subnet of the VPC
	ec2Options := &terraform.Options{
		TerraformDir:    "../../infra/test/ec2_instance",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"region":                  "us-east-2",
			"instance_name":           "test-ec2-custom-vpc-" + uniqueID,
			"ami_id":                  "ami-0900fe555666598a2",
			"instance_type":           "t3.micro",
			"vpc_id":                  vpcID,
			"subnet_id":               subnetIDs[0],
			"eice_security_group_id":  eiceSGID,
		},
	}

	defer terraform.Destroy(t, ec2Options)
	terraform.InitAndApply(t, ec2Options)

	instanceID := terraform.Output(t, ec2Options, "instance_id")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	require.NoError(t, err)

	ec2Client := ec2.NewFromConfig(cfg)

	out, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	require.NoError(t, err)
	require.Len(t, out.Reservations, 1)
	require.Len(t, out.Reservations[0].Instances, 1)

	instance := out.Reservations[0].Instances[0]

	// Instance is running
	assert.Equal(t, ec2types.InstanceStateNameRunning, instance.State.Name)

	// Instance is in the expected VPC
	assert.Equal(t, vpcID, *instance.VpcId)

	// Instance is in the expected subnet
	assert.Equal(t, subnetIDs[0], *instance.SubnetId)

	// No public IP assigned
	assert.Nil(t, instance.PublicIpAddress)
}
