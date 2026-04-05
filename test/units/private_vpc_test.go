package test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPrivateVPCMultiAZ(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir:    "../../examples/units/private_vpc",
		TerraformBinary: "terragrunt",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	// Terraform Outputs
	region := terraform.Output(t, terraformOptions, "region")
	vpcCIDR := terraform.Output(t, terraformOptions, "vpc_cidr")
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
	assert.Equal(t, vpcCIDR, *vpc.CidrBlock)

	// Validate Subnets
	assert.Equal(t, 1, len(subnetIDs))

	subnetOut, err := ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: subnetIDs,
	})
	assert.NoError(t, err)
	assert.Len(t, subnetOut.Subnets, 1)

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

	// Validate EICE endpoints were created for each private subnet
	eiceOut, err := ec2Client.DescribeInstanceConnectEndpoints(ctx, &ec2.DescribeInstanceConnectEndpointsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	assert.NoError(t, err)

	endpointSubnetIDs := make([]string, 0, len(eiceOut.InstanceConnectEndpoints))
	for _, endpoint := range eiceOut.InstanceConnectEndpoints {
		assert.Equal(t, ec2types.Ec2InstanceConnectEndpointStateCreateComplete, endpoint.State)
		endpointSubnetIDs = append(endpointSubnetIDs, *endpoint.SubnetId)
	}

	assert.ElementsMatch(t, subnetIDs, endpointSubnetIDs)
}
