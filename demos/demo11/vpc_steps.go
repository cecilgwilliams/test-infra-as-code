package bddtest

// VPC BDD step definitions using godog + Terratest.
//
// These steps implement the Gherkin scenarios in features/vpc.feature.
// They use Terratest to deploy real AWS resources and make assertions.

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cucumber/godog"
	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// VPCTestContext holds shared state between steps for VPC scenarios.
type VPCTestContext struct {
	TerraformOptions *terraform.Options
	VPCID            string
	Region           string
}

// NewVPCTestContext creates a new test context pointing at the VPC Terraform module.
func NewVPCTestContext() *VPCTestContext {
	return &VPCTestContext{
		TerraformOptions: &terraform.Options{
			TerraformDir: "../demo9",
		},
		Region: awsRegion,
	}
}

func (v *VPCTestContext) theTerraformModuleHasBeenApplied(ctx context.Context) error {
	t := GetTestingT(ctx)
	terraform.InitAndApply(t, v.TerraformOptions)
	v.VPCID = terraform.Output(t, v.TerraformOptions, "vpc_id")
	return nil
}

func (v *VPCTestContext) theVPCCIDRBlockShouldBe(ctx context.Context, expectedCIDR string) error {
	t := GetTestingT(ctx)
	sess, err := terraaws.NewAuthenticatedSession(v.Region)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %s", err)
	}
	ec2Client := ec2.New(sess)
	out, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(v.VPCID)},
	})
	if err != nil {
		return fmt.Errorf("failed to describe VPC: %s", err)
	}
	assert.True(t, len(out.Vpcs) != 0)
	actualCIDR := aws.StringValue(out.Vpcs[0].CidrBlock)
	assert.Equal(t, expectedCIDR, actualCIDR)
	return nil
}

func (v *VPCTestContext) theVPCShouldHaveSubnets(ctx context.Context, expectedCount int) error {
	t := GetTestingT(ctx)
	subnets := terraaws.GetSubnetsForVpc(t, v.VPCID, v.Region)
	assert.Equal(t, expectedCount, len(subnets))
	return nil
}

func (v *VPCTestContext) theVPCShouldHaveDNSSupportEnabled(ctx context.Context) error {
	t := GetTestingT(ctx)
	enabled, err := v.getVpcAttribute("enableDnsSupport")
	if err != nil {
		return err
	}
	assert.True(t, enabled)
	return nil
}

func (v *VPCTestContext) theVPCShouldHaveDNSHostnamesEnabled(ctx context.Context) error {
	t := GetTestingT(ctx)
	enabled, err := v.getVpcAttribute("enableDnsHostnames")
	if err != nil {
		return err
	}
	assert.True(t, enabled)
	return nil
}

func (v *VPCTestContext) getVpcAttribute(attribute string) (bool, error) {
	sess, err := terraaws.NewAuthenticatedSession(v.Region)
	if err != nil {
		return false, fmt.Errorf("failed to create AWS session: %s", err)
	}
	ec2Client := ec2.New(sess)
	out, err := ec2Client.DescribeVpcAttribute(&ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String(v.VPCID),
		Attribute: aws.String(attribute),
	})
	if err != nil {
		return false, fmt.Errorf("failed to describe VPC attribute %s: %s", attribute, err)
	}
	switch attribute {
	case "enableDnsSupport":
		return aws.BoolValue(out.EnableDnsSupport.Value), nil
	case "enableDnsHostnames":
		return aws.BoolValue(out.EnableDnsHostnames.Value), nil
	}
	return false, fmt.Errorf("unknown attribute: %s", attribute)
}

// Cleanup destroys the Terraform resources.
func (v *VPCTestContext) Cleanup(ctx context.Context) {
	t := GetTestingT(ctx)
	terraform.Destroy(t, v.TerraformOptions)
}

// InitializeVPCScenario registers all VPC step definitions with godog.
func InitializeVPCScenario(ctx *godog.ScenarioContext) {
	v := NewVPCTestContext()

	ctx.Given(`^the Terraform module has been applied$`, v.theTerraformModuleHasBeenApplied)
	ctx.Then(`^the VPC CIDR block should be "([^"]*)"$`, v.theVPCCIDRBlockShouldBe)
	ctx.Then(`^the VPC should have (\d+) subnets$`, v.theVPCShouldHaveSubnets)
	ctx.Then(`^the VPC should have DNS support enabled$`, v.theVPCShouldHaveDNSSupportEnabled)
	ctx.Then(`^the VPC should have DNS hostnames enabled$`, v.theVPCShouldHaveDNSHostnamesEnabled)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		v.Cleanup(ctx)
		return ctx, nil
	})
}
