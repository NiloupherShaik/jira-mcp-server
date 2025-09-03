#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -i [-v] [-?] -- program to install a self-signer issuer

where:
    -?                  show this help text
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
  "--version") set -- "$@" "-i" ;;
  "--help") set -- "$@" "-?" ;;
  "--verbose") set -- "$@" "-v" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

while getopts 'i:v?' option; do
  case "$option" in
  i) version="${OPTARG}" ;;
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

# nothing to validate

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

#################################
##### interact with cluster #####
#################################

debug "Checking if cluster context is local"
if [ "$(is_local_context)" != 'true' ]; then
    fail "Cluster context not a local context"
fi

# This script installs self signed issuer
info "Installing self signed issuer"
kubectl apply -f ./issuer.yaml

debug "Waiting to stabilize..."
sleep 10

success "Installed self signed issuer"
