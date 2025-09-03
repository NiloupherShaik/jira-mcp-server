#!/usr/bin/env bash

# Begin Standard 'imports'
set -e
set -o pipefail

## Bash v4.x and bash v.3x use different escape characters
if [ "${BASH_VERSINFO:-0}" -ge 4 ]; then
  ESCAPE_CHAR='\e'
else
  ESCAPE_CHAR='\033'
fi

gray="${ESCAPE_CHAR}[90m"
blue="${ESCAPE_CHAR}[36m"
red="${ESCAPE_CHAR}[31m"
yellow="${ESCAPE_CHAR}[33m"
green="${ESCAPE_CHAR}[32m"
reset="${ESCAPE_CHAR}[0m"


#######################################
# echoes a message in blue
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
info() { echo -e "${blue}INFO: $*${reset}"; }

#######################################
# echoes a message in red
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
error() { echo -e "${red}ERROR: $*${reset}"; }


#######################################
# echoes a message in yellow
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
warn() { echo -e "${yellow}WARN: $*${reset}"; }

#######################################
# echoes a message in grey. Only if debug mode is enabled
# Globals:
#   DEBUG
# Arguments:
#   Message
# Returns:
#   None
#######################################
debug() {
    echo -e "${gray}DEBUG: $*${reset}";
}

#######################################
# echoes a message in green
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
success() { echo -e "${green}✔ $*${reset}"; }

#######################################
# echoes a message in red and terminates the program
# Globals:
#   None
# Arguments:
#   Message
# Returns:
#   None
#######################################
fail() { echo -e "${red}✖ $*${reset}"; exit 1; }

## Enable debug mode.
enable_debug() {
  info "Enabling debug mode."
  set -x
}

#######################################
# echoes a message in blue
# Globals:
#   status: Exit status of the command that was executed.
#   output_file: Local path with captured output generated from the command.
# Arguments:
#   command: command to run
# Returns:
#   None
#######################################
#run() {
##  output_file="/var/tmp/pipe-$(date +%s)-$RANDOM"
#
##  echo "$@"
#  set +e
#  "$@"
#  status=$?
#  set -e
#}

# End standard 'imports'

#######################################
# Exits if the terminal isn't interactive
# Globals:
#   None
# Arguments:
#   Parameter
# Returns:
#   None
#######################################
fail_if_non_interactive(){
  # https://serverfault.com/a/991904/429389
  if [ "$(is_interactive)" == "false" ]; then
    fail "No ${1} id provided and terminal isn't interactive"
  fi
}

is_interactive(){
    # https://serverfault.com/a/991904/429389
  if [ "${-#*i}" != "$-" ]; then echo -n "false"; else echo -n "true"; fi
}

export -f debug
export -f info
export -f success
export -f error
export -f fail_if_non_interactive
