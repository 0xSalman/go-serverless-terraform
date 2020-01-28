#!/bin/bash

dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
cd "$dir/terraform"

export AWS_PROFILE=rethesis_personal
export TF_VAR_ENV=dev
echo "Start building infrastructure"
terraform init
terraform apply -auto-approve
echo "Finished building infrastructure"

