package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPrivateVPCMultiAZ(t *testing.T) {
	t.Parallel()

	region := "us-east-2"

	uniqueID := fmt.Sprintf("%d", time.Now().Unix())

	terraformOptions := &terraform.Options{
		TerraformDir:    "../../infra/test/private_vpc",
		TerraformBinary: "terragrunt",
		Vars: map[string]interface{}{
			"vpc_name": "test-private-vpc-" + uniqueID,
			"region":   region,
			"vpc_cidr": "10.0.0.0/16",
			"availability_zones": []string{
				"us-east-2a",
				"us-east-2b",
			},
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Terraform Outputs
	vpcID := terraform.Output(t, terraformOptions, "vpc_id")
	subnetIDs := terraform.OutputList(t, terraformOptions, "private_subnet_ids")

	// Load AWS config
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	assert.NoError(t, err)

	ec2Client := ec2.NewFromConfig(cfg)

	// Validate VPC
	vpcOut, err := ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	})
	assert.NoError(t, err)
	assert.Len(t, vpcOut.Vpcs, 1)

	vpc := vpcOut.Vpcs[0]
	assert.Equal(t, "10.0.0.0/16", *vpc.CidrBlock)

	// Validate Subnets
	assert.Equal(t, 2, len(subnetIDs))

	subnetOut, err := ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: subnetIDs,
	})
	assert.NoError(t, err)
	assert.Len(t, subnetOut.Subnets, 2)

	for _, subnet := range subnetOut.Subnets {

		// Belongs to correct VPC
		assert.Equal(t, vpcID, *subnet.VpcId)

		// Private subnet check
		assert.NotNil(t, subnet.MapPublicIpOnLaunch)
		assert.False(t, *subnet.MapPublicIpOnLaunch)

		// Has CIDR
		assert.NotNil(t, subnet.CidrBlock)
		assert.NotEmpty(t, *subnet.CidrBlock)
	}
}
