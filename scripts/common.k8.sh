#!/usr/bin/env bash

# Begin Standard 'imports'
#set -e
set -o pipefail

#######################################
# returns true or false if context
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
is_same_context(){
    if [ "$(kubectx -c)" == "${context}" ]; then
        echo 'true'
    else
        echo 'false'
    fi
}

is_same_namespace(){
    if [ "$(kubens -c)" == "${namespace}" ]; then
        echo 'true'
    else
        echo 'false'
    fi
}

revert_config(){
    if [ "$(is_same_context)" == "false" ]; then
        info "Switching to previous context"
        kubectx -
    fi

    if [ "$(is_same_namespace)" == "false" ]; then
        info "Switching to previous namespace"
        kubens -
    fi
}

is_local_context(){
    if [ "$(kubectl config current-context)" != "minikube" ] &&\
     [ "$(kubectl config current-context)" != "docker-desktop" ] &&\
     [[ "$(kubectl config current-context)" != kind-* ]]; then
        echo "false"
    else
        echo "true"
    fi
}
