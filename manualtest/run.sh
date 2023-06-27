#!/usr/bin/env bash
set -o errexit
set -o nounset

if [[ -n ${DEBUG+undef} ]]; then
  set -o xtrace
fi

function log() {
  local level="${1}"
  local color="${2}"
  shift 2
  echo -e "\e[6;35m<$(date)> \e[0;1;4;${color}m[${level}]\e[0;${color}m ${@}\e[0m"
}

function info() { log "INFO" 32 "${@}"; }
function error() { log "ERROR" 31 "${@}"; }

this_dir="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"

build_dir="${this_dir}/.build"
plugin_cache_dir="${build_dir}/terraform-plugin-cache"
gobin_dir="${build_dir}/gobin"

info "Creating directories to store generated files in..."
mkdir -pv "${build_dir}"
mkdir -pv "${plugin_cache_dir}"
mkdir -pv "${gobin_dir}"

info "Installing provider into temporary GOBIN ${gobin_dir}..."
GOBIN="${gobin_dir}" go install "${this_dir}/../..."

info "Creating Terraform configuration file..."
tfrc_file="${build_dir}/config.tfrc"
tee "${tfrc_file}" <<EOF
  plugin_cache_dir = "$(realpath "${plugin_cache_dir}")"

  provider_installation {
    dev_overrides {
      "hashicorp.com/ascopes/ocicopy" = "${gobin_dir}"
    }

    direct {}
  }
EOF

info "Starting registry Docker container..."
registry_cid_file="${build_dir}/registry.cid"
docker run --detach --cidfile "${registry_cid_file}" -p 5000:5000 --rm registry:2
trap 'kill_container' EXIT INT TERM QUIT

function kill_container() {
  info "Stopping container..."
  docker stop "$(cat "${registry_cid_file}")"
  rm -v "${registry_cid_file}"
}

info "Performing a health check on the container..."

iterations=0
until [[ ${iterations} -gt 30 ]] || curl -I --fail http://localhost:5000; do
  info "Waiting for container to become healthy..."
  ((iterations++))
  sleep 1
done
if [[ "${iterations}" -gt 30 ]]; then
  error "Container failed to become healthy in time. Aborting..."
  exit 2
fi


info "Running Terraform"
(
  cd "${this_dir}"
  export GOBIN="${gobin_dir}"
  export TF_CLI_CONFIG_FILE="${tfrc_file}"
  terraform apply -auto-approve
)
