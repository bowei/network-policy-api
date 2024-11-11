#!/usr/bin/env sh

set -xv
set -euo pipefail

POLICY_ASSISTANT_ARGS=$@
RESULTS_DIR="${RESULTS_DIR:-/tmp/results}"


./policy-assistant $POLICY_ASSISTANT_ARGS > "${RESULTS_DIR}"/results.txt


cd "${RESULTS_DIR}"

# Sonobuoy worker expects a tar file.
tar czf results.tar.gz ./*

# Signal to the worker that we are done and where to find the results.
realpath results.tar.gz > ./done
