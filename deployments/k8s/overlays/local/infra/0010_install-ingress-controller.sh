#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -i [-v] [-?] -- program to install Istio and addons

where:
    -?                  show this help text
    -i, --version       Istio version
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

# If host if not passed in via parameter, pester the user
if [ -z "${version}" ]; then
  # check if the terminal is interactive
  fail_if_non_interactive

  debug "No istio version passed, requesting from user"
  read -p "Enter the Istio version to install [ENTER]: 1.20 " version
  version="${version:-1.20}"
fi

###################################################
##### confirm required programs are available #####
###################################################

if ! command -v istioctl &>/dev/null; then
  fail "istioctl is not installed.  It's not optional"
fi

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

cluster_context="$(kubectl config current-context)"

info "Creating istio-system namespace"
kubectl create namespace istio-system

#if [ "${cluster_context}" == "kind-ai" ]; then
#  info "Installing Istio with kind-ai profile"
#  istioctl install -f "${project_root}/deployments/k8s/overlays/local/infra/kind-istio.yaml" -y
#else
#  info "Installing Istio without kind-ai profile"
  info "Installing Istio"
  istioctl install -y
#fi

success "Installed Istio"

kubectl apply -f "https://raw.githubusercontent.com/istio/istio/release-${version}/samples/addons/kiali.yaml"
success "Installed kiali addons (${version})"

kubectl apply -f "https://raw.githubusercontent.com/istio/istio/release-${version}/samples/addons/jaeger.yaml"
success "Installed jaeger addons (${version})"

kubectl apply -f "https://raw.githubusercontent.com/istio/istio/release-${version}/samples/addons/prometheus.yaml"
success "Installed prometheus addons (${version})"
