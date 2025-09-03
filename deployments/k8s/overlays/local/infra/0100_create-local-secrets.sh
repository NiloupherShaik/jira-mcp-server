#!/usr/bin/env bash

usage() {
  echo "Usage: $0 -n [-v] [-?] -- program to create local sealed secrets

where:
    -?                            show this help text
    -n, --namespace <namespace>  namespace to use
    " 1>&2
  exit 1
}

echo

project_root="$(git rev-parse --show-toplevel)"

## Add some helper functions
. "${project_root}/scripts/common.sh"
. "${project_root}/scripts/common.k8.sh"

info "Script: $(basename "$0")" "$*"

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
  "--help") set -- "$@" "-?" ;;
  "--verbose") set -- "$@" "-v" ;;
  "--namespace") set -- "$@" "-n" ;;
  *) set -- "$@" "$arg" ;;
  esac
done

while getopts 'n:v?' option; do
  case "$option" in
  n) namespace="$OPTARG" ;;
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
if [ -z "${namespace}" ]; then
  # check if the terminal is interactive
  fail_if_non_interactive

  debug "No namespace passed, requesting from user"
  read -p "Enter the Namespace to generate the secrets in [ENTER]: fusion-ai " namespace
  namespace="${namespace:-fusion-ai}"
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
debug "  namespace: ${namespace}"

#################################
##### interact with cluster #####
#################################

debug "Checking if cluster context is local"
if [ "$(is_local_context)" != 'true' ]; then
    fail "Cluster context not a local context"
fi

secretFiles=("apollolicense.txt" )

echo "checking for secret text files"
for file in "${secretFiles[@]}"; do
  info "Checking file: ${file}"
  filePath="${project_root}/deployments/k8s/overlays/local/resources/${file}"
  if ! test -f "${filePath}"; then
    warn "File ${file} does not exist, creating blank file"
    touch "${filePath}"
    success "${file} created"
  else
    success "${file} already exists!"
  fi
done

cd "${project_root}"

info "creating Apollo license secret"
kubectl create secret generic apollo-license-secret \
		-n "${namespace}" \
		--from-file=APOLLO_KEY="${project_root}"/deployments/k8s/overlays/local/resources/apollolicense.txt \
		--dry-run=client \
		-o yaml | kubeseal -o yaml > deployments/k8s/overlays/local/resources/apollo-license-sealed-secret.yaml
success "done"
