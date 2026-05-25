#!/bin/bash

###############################################################################
# presentation-demo.sh — IaC Testing Tools Demo Script
#
# Usage:  cd demos && bash ./presentation-demo.sh        (runs all demos)
#         cd demos && bash ./presentation-demo.sh 5      (starts at demo 5)
# Press Enter to advance through each command.
###############################################################################

START_DEMO=${1:-1}
DEMO_DIR=$(pwd)

clear

source .venv/bin/activate

# Load demo-magic
source ./demo-magic.sh

# Customize demo-magic settings
DEMO_PROMPT="(.venv) > "
DEMO_COMMENT_COLOR="\033[0;37m"
DEMO_CMD_COLOR="\033[0;37m"

# Set the typing speed
TYPE_SPEED=0.04

# ============================================================================
# Demo 1: cfn-lint Static Analysis
# ============================================================================
if [ $START_DEMO -le 1 ]; then

echo "# === Demo 1: cfn-lint Static Analysis ==="

pe "cfn-lint ./demo1-4/bad-template.yaml"

pe "cfn-lint ./demo1-4/good-template.yaml"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 2: checkov Security Scanning
# ============================================================================
if [ $START_DEMO -le 2 ]; then

echo "# === Demo 2: checkov Security Scanning ==="

pe "checkov -f ./demo1-4/bad-template.yaml --framework cloudformation"

pe "checkov -f ./demo1-4/good-template.yaml --framework cloudformation"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 3: cfn-guard Policy as Code
# ============================================================================
if [ $START_DEMO -le 3 ]; then

echo "# === Demo 3: cfn-guard Policy as Code ==="

pe "cfn-guard validate -r ./demo1-4/s3-rules.guard -d ./demo1-4/bad-template.yaml"

pe "cfn-guard validate -r ./demo1-4/s3-rules.guard -d ./demo1-4/good-template.yaml"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 4: taskcat Static Analysis
# ============================================================================
if [ $START_DEMO -le 4 ]; then

cd "$DEMO_DIR"

echo "# === Demo 4: taskcat Static Analysis ==="

pe "cd demo1-4"

pe "taskcat lint"

echo "Edit Taskcat Config File"

wait

pe "taskcat lint"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 5: taskcat Integration Test
# ============================================================================
if [ $START_DEMO -le 5 ]; then

cd "$DEMO_DIR"

echo "# === Demo 5: taskcat Integration Test ==="

echo "Start Docker"
wait
echo "Open AWS Console"
wait

pe "cd demo5"

pe "AWS_PROFILE=cecil-personal taskcat test run"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 6: Terraform Static Analysis
# ============================================================================
if [ $START_DEMO -le 6 ]; then

cd "$DEMO_DIR"

echo "# === Demo 6: Terraform Static Analysis ==="

pe "cd demo6"

pe "terraform init -backend=false"

pe "terraform validate"

pe "terraform fmt -check"

pe "tflint --init && tflint"

pe "tfsec ."

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 7: OPA Policy Testing
# ============================================================================
if [ $START_DEMO -le 7 ]; then

cd "$DEMO_DIR"

echo "# === Demo 7: OPA Policy Testing ==="

pe "cd demo6"

pe "AWS_PROFILE=cecil-personal terraform plan -out=tfplan"

pe "AWS_PROFILE=cecil-personal terraform show -json tfplan > tfplan.json"

pe "opa eval -d ../demo7/s3_encrypted.rego -i tfplan.json \"data.terraform.deny\""

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 8: Terratest Integration Test
# ============================================================================
if [ $START_DEMO -le 8 ]; then

cd "$DEMO_DIR"

echo "# === Demo 8: Terratest Integration Test ==="

pe "cd demo8/test"

pe "AWS_PROFILE=cecil-personal go test -v -timeout 1m"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 9: Terratest Integration Test with Multiple Resources
# ============================================================================
if [ $START_DEMO -le 9 ]; then

cd "$DEMO_DIR"

echo "# === Demo 9: Terratest Integration Test with Multiple Resources ==="

pe "cd demo9/test"

pe "AWS_PROFILE=cecil-personal go test -v -timeout 1m"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 10: Test Driven Development for Infrastructure
# ============================================================================
if [ $START_DEMO -le 10 ]; then

cd "$DEMO_DIR"

echo "# === Demo 10: Test Driven Development for Infrastructure ==="

pe "cd demo10/test"

pe "AWS_PROFILE=cecil-personal go test -v -run TestS3BucketHasVersioning -timeout 3m"

echo "Implement the desired behavior"

wait

pe "AWS_PROFILE=cecil-personal go test -v -run TestS3BucketHasVersioning -timeout 3m"

echo "Return to Slides"

wait

clear
fi
# ============================================================================
# Demo 11: Behavior Driven Development for Infrastructure
# ============================================================================
if [ $START_DEMO -le 11 ]; then

cd "$DEMO_DIR"

echo "# === Demo 11: Behavior Driven Development for Infrastructure ==="

pe "cd demo11"

pe "AWS_PROFILE=cecil-personal go test -v -timeout 5m"

cd "$DEMO_DIR"

echo "# === All demos complete! ==="
fi