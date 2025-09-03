#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -i [-v] [-?] -- program to install sealed-secrets in a local cluster

where:
    -?                  show this help text
    -i, --version       Sealed Secrets version
                        Default is '0.17.5'
    " 1>&2
  exit 1
}

echo
working_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
project_root=$(cd "../../../../../" && pwd)

. "${project_root}/scripts/common.sh"
. "${project_root}/scripts/common.k8.sh"

info "Script: $(basename "$0")" "$*"

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--version") set -- "$@" "-s" ;;
  "--help") set -- "$@" "-?" ;;
  "--verbose") set -- "$@" "-v" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

while getopts 's:v?' option; do
  case "$option" in
  s) version="${OPTARG}" ;;
  v) enable_debug ;;
  \?)
    usage
    exit 1
    ;;
  :)
    printf "%s","Option -${OPTARG} requires an argument."
    usage
    exit 1
    ;;
  esac
done
shift $((OPTIND - 1))

###########################
##### Validate inputs #####
###########################

# If host if not passed in via parameter, pester the user
if [ -z "${version}" ]; then
  # check if the terminal is interactive
  fail_if_non_interactive

  debug "No istio version passed, requesting from user"
  read -p "Enter the Istio version to install [ENTER]: 0.17.5 " version
  version="${version:-0.17.5}"
fi

###################################################
##### confirm required programs are available #####
###################################################

if ! command -v kubectl &>/dev/null; then
  fail "kubectl is not installed.  It's not optional"
fi

########################
##### print inputs #####
########################
echo
debug "inputs:"
debug "version : ${version}"

#################################
##### interact with cluster #####
#################################

debug "Checking if cluster context is local"
if [ "$(is_local_context)" != 'true' ]; then
    fail "Cluster context not a local context"
fi

# This script installs sealed secrets controller in a pre-existing cluster
kubectl apply -f "https://github.com/bitnami-labs/sealed-secrets/releases/download/v${version}/controller.yaml"

debug "Waiting to stabilize..."
sleep 10
kubectl wait --namespace kube-system \
  --for=condition=ready pod \
  -l name=sealed-secrets-controller \
  --timeout=90s

success "Installed sealed-secrets (${version})"
