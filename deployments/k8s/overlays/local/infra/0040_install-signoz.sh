#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -i [-v] [-?] -- program to install SigNoz into a local cluster

where:
    -?                  show this help text
    -i, --version       Signoz version
                        Default is '1.20'
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
  "--help") set -- "$@" "-?" ;;
  "--verbose") set -- "$@" "-v" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

while getopts 'v?' option; do
  case "$option" in
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

if ! command -v helm &>/dev/null; then
  fail "helm is not installed.  It's not optional"
fi

########################
##### print inputs #####
########################
echo
debug "inputs:"
debug "-"

#################################
##### interact with cluster #####
#################################

debug "Checking if cluster context is local"
if [ "$(is_local_context)" != 'true' ]; then
    fail "Cluster context not a local context"
fi

cluster_context="$(kubectl config current-context)"

info "Installing SigNoz, this may take 3-5 minutes... Use 'kubectl -n signoz get pods' to check status"
helm repo add signoz https://charts.signoz.io
helm repo update

helm install signoz signoz/signoz \
   --namespace signoz \
   --create-namespace \
   --wait \
   --timeout 1h

success "Installed Signoz"
