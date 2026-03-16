# features/vpc.feature
# BDD feature spec for VPC networking requirements.
# Run: go test -v -run TestBDDVPC

Feature: Production VPC
  As a platform engineer
  I need a properly configured VPC
  So that our workloads are network-isolated

  Background:
    Given the Terraform module has been applied

  Scenario: VPC has correct CIDR
    Then the VPC CIDR block should be "10.0.0.0/16"

  Scenario: VPC has three subnets
    Then the VPC should have 3 subnets

  Scenario: VPC has DNS support
    Then the VPC should have DNS support enabled

  Scenario: VPC has DNS hostnames
    Then the VPC should have DNS hostnames enabled
