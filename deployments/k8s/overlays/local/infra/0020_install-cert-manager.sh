#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -i [-v] [-?] -- program to install cert-manager in a local cluster

where:
    -?                  show this help text
    -c, --version       Cert manager version
                        Default is '1.5.4'
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
  "--version") set -- "$@" "-c" ;;
  "--help") set -- "$@" "-?" ;;
  "--verbose") set -- "$@" "-v" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

while getopts 'c:v?' option; do
  case "$option" in
  c) version="${OPTARG}" ;;
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

  debug "No Cert Manager version passed, requesting from user"
  read -p "Enter the Cert Manager to install [ENTER]: 1.5.4 " version
  version="${version:-1.5.4}"
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

# This script installs cert-manager controller in a pre-existing cluster
kubectl apply -f "https://github.com/jetstack/cert-manager/releases/download/v${version}/cert-manager.yaml"

debug "Waiting to stabilize..."
sleep 10
kubectl wait --namespace cert-manager \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

success "Installed cert-manager (${version})"
