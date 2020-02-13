#!/bin/bash

function buildLambdas() {
  rootDir=${1}
  infraDir=${rootDir}/infra
  names=( "${@:3:$2}" ); shift "$(( $2 + 1 ))"
  versions=( "${@:3:$2}" ); shift "$(( $2 + 1 ))"

  for i in "${!names[@]}";
  do
    name="${names[$i]}";
    version="${versions[$i]}";
    cd ${rootDir}/cmd/${name}
    env GOOS=linux GOARCH=amd64 go build -o ${infraDir}/repository/${name}
    echo "Finished building ${name} go binary"

    cd ${infraDir}
    zip -j ./repository/${name}-${version}.zip ./repository/${name}
    echo "Finished zipping ${name} go binary"
  done
}

rootDir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
infraDir=${rootDir}/infra
cd ${infraDir}

mkdir -p repository
rm -fr repository/*

lambdaNames=("verification-link" "clone-user" "user")
lambdaVersions=("0.1.1" "0.1.5" "0.2.3")

buildLambdas ${rootDir} "${#lambdaNames[@]}" "${lambdaNames[@]}" "${#lambdaVersions[@]}" "${lambdaVersions[@]}"

export AWS_PROFILE=rethesis_personal
export TF_VAR_ENV=dev
export TF_VAR_WEBSITE_URL=http://localhost:9000
export TF_VAR_VERIFICATION_LINK_VERSION=${lambdaVersions[0]}
export TF_VAR_CLONE_USER_VERSION=${lambdaVersions[1]}
export TF_VAR_USER_VERSION=${lambdaVersions[2]}
echo "Start building infrastructure"
terraform init
#terraform plan
terraform apply -auto-approve
echo "Finished building infrastructure"
