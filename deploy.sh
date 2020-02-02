#!/bin/bash

function buildLambdas() {
  rootDir=${1}
  terraformDir=${rootDir}/terraform
  names=( "${@:3:$2}" ); shift "$(( $2 + 1 ))"
  versions=( "${@:3:$2}" ); shift "$(( $2 + 1 ))"

  for i in "${!names[@]}";
  do
    name="${names[$i]}";
    version="${versions[$i]}";
    cd ${rootDir}/cmd/${name}
    env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ${terraformDir}/repository/${name}
    echo "Finished building ${name} go binary"

    cd ${terraformDir}
    zip -j ./repository/${name}-${version}.zip ./repository/${name}
    echo "Finished zipping ${name} go binary"
  done
}

rootDir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
terraformDir=${rootDir}/terraform
cd ${terraformDir}

mkdir -p repository
rm -fr repository/*

lambdaNames=("verification-link" "clone-user")
lambdaVersions=("0.1.0" "0.1.0")

declare -a array=("one" "two" "three")

buildLambdas ${rootDir} "${#lambdaNames[@]}" "${lambdaNames[@]}" "${#lambdaVersions[@]}" "${lambdaVersions[@]}"

export AWS_PROFILE=rethesis_personal
export TF_VAR_ENV=dev
export TF_VAR_WEBSITE_URL=http://localhost:1234
export TF_VAR_VERIFICATION_LINK_VERSION=${lambdaVersions[0]}
export TF_VAR_CLONE_USER_VERSION=${lambdaVersions[1]}
echo "Start building infrastructure"
terraform init
terraform apply -auto-approve
echo "Finished building infrastructure"
