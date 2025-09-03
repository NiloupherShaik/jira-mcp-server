#!/usr/bin/env bash

usage() {
  echo "Usage: $0 [-a] -- program to setup k8s contexts

where:
    -a                  account shorthand name, eg FAI or CDO
    " 1>&2
  exit 1
}

aws_account=''

while getopts 'a:v' flag; do
  case "${flag}" in
  a) aws_account=${OPTARG} ;;
  *) usage
    exit 1 ;;
  esac
done

if [ -z "$aws_account" ];
  then
    echo "aws-account not set, please use -account <account-name>"
    exit 1
fi

echo "aws-account: ${aws_account}"
working_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

# Ensure we know where we are before starting
cd "${working_dir}" || exit 1
project_root=$(cd "../" && pwd)

# Add some helper functions
. "${project_root}/scripts/common.sh"
. "${project_root}/scripts/common.k8.sh"

info "setup-k8s-contexts.sh $*"

environments=("dev" "stg" "prd")
regions=("ap-southeast-2" "eu-central-1" "eu-west-1" "us-west-2")
cluster_versions=("01")

for env in "${environments[@]}"; do
    info "Environment: ${env}"
    aws_profile="${aws_account}-${env}-developer"

    for region in "${regions[@]}"; do

        info "Region: ${region}"

        if [[ "${region}" == "ap-southeast-2" ]]; then
            regionName="apse2"
        else
            regionName="${region:0:2}${region:3:2}${region: -1}"
        fi

        for cluster_version in "${cluster_versions[@]}"; do
            cluster="slr-${aws_account}-${env}-${cluster_version}-${regionName}"
            debug "Cluster: ${cluster}"

             info "adding cluster to config: ${cluster}"
             if ! AWS_PROFILE="${aws_profile}" aws eks update-kubeconfig --name "${cluster}" --region "${region}"; then
                error "Failed to add access to cluster ${cluster}"
                continue;
             fi

             success "Added access to cluster ${cluster}"
         done
    done
done

success "Completed Cluster authentication setup"
