package test

import (
	"fmt"
	"os"
	"testing"

	sdkaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const awsRegion = "us-east-1"

func TestMain(m *testing.M) {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		fmt.Fprintln(os.Stderr, "ERROR: AWS_PROFILE is not set.")
		os.Exit(1)
	}
	if err := os.Setenv("AWS_PROFILE", profile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to set AWS_PROFILE: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func terraformOptions() *terraform.Options {
	return &terraform.Options{
		TerraformDir: "../",
	}
}

// TestVPCWithSubnets deploys a VPC with subnets, verifies the VPC CIDR
// and subnet count, then tears everything down.
func TestVPCWithSubnets(t *testing.T) {
	t.Parallel()

	opts := terraformOptions()
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	vpcID := terraform.Output(t, opts, "vpc_id")

	// Verify VPC CIDR via EC2 SDK (terratest Vpc struct does not expose CidrBlock)
	sess, err := aws.NewAuthenticatedSession(awsRegion)
	require.NoError(t, err)
	ec2Client := ec2.New(sess)
	out, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{sdkaws.String(vpcID)},
	})
	require.NoError(t, err)
	require.Len(t, out.Vpcs, 1)
	assert.Equal(t, "10.0.0.0/16", sdkaws.StringValue(out.Vpcs[0].CidrBlock))

	// Verify subnets
	subnets := aws.GetSubnetsForVpc(t, vpcID, awsRegion)
	assert.Equal(t, 3, len(subnets))
}

// TestVPCHasDNSSupport verifies the VPC has DNS support and hostnames enabled.
func TestVPCHasDNSSupport(t *testing.T) {
	t.Parallel()

	opts := terraformOptions()
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	vpcID := terraform.Output(t, opts, "vpc_id")

	sess, err := aws.NewAuthenticatedSession(awsRegion)
	require.NoError(t, err)
	ec2Client := ec2.New(sess)

	supportOut, err := ec2Client.DescribeVpcAttribute(&ec2.DescribeVpcAttributeInput{
		VpcId:     sdkaws.String(vpcID),
		Attribute: sdkaws.String("enableDnsSupport"),
	})
	require.NoError(t, err)
	assert.True(t, sdkaws.BoolValue(supportOut.EnableDnsSupport.Value))

	hostnamesOut, err := ec2Client.DescribeVpcAttribute(&ec2.DescribeVpcAttributeInput{
		VpcId:     sdkaws.String(vpcID),
		Attribute: sdkaws.String("enableDnsHostnames"),
	})
	require.NoError(t, err)
	assert.True(t, sdkaws.BoolValue(hostnamesOut.EnableDnsHostnames.Value))
}
