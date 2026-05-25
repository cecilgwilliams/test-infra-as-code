# Run the following commands from the demos folder

# Initialize Python Virutal Environment
source .venv/bin/activate

# Demo 1: cfn-lint Static Analysis
cfn-lint ./cloudformation/bad-template.yaml
cfn-lint ./cloudformation/good-template.yaml

# Demo 2: checkov Security Scanning
checkov -f ./cloudformation/bad-template.yaml --framework cloudformation
checkov -f ./cloudformation/good-template.yaml --framework cloudformation

# Demo 3: cfn-guard Policy as Code
cfn-guard validate -r ./cloudformation/s3-rules.guard -d ./cloudformation/bad-template.yaml
cfn-guard validate -r ./cloudformation/s3-rules.guard -d ./cloudformation/good-template.yaml

# Demo 4: taskcat Static Analysis
taskcat lint

# Demo 5: taskcat integration test (view AWS Cloudformation console)
cd cloudformation
AWS_PROFILE=cecil-personal taskcat test run

# Demo 6: terraform static analysis
cd terraform-bad
terraform init -backend=false
terraform validate
terraform fmt -check
tflint --init && tflint
tfsec .

# Demo 7: terratest integration test
cd terraform-s3/test
AWS_PROFILE=cecil-personal go test -v -timeout 1m

# Demo 8: terratest integration test with multiple resources
cd terraform-vpc/test
AWS_PROFILE=cecil-personal go test -v -timeout 1m

# Demo 9: Test Driven Development for Infrastructure
cd terraform-tdd
AWS_PROFILE=cecil-personal go test -v -run TestS3BucketHasVersioning -timeout 3m

# Demo 10: 
cd terraform-bdd
AWS_PROFILE=cecil-personal go test -v -timeout 5m


