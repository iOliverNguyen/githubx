#!/usr/bin/env bash
set -eo pipefail

env_file=extension/env/dev.js
out_dir=build/extension
if [[ -n $GITHUBX_ENV_FILE ]]; then env_file="$GITHUBX_ENV_FILE" ; fi
if [[ -n $GITHUBX_OUT_DIR  ]]; then out_dir="$GITHUBX_OUT_DIR"   ; fi

rm -rf "$out_dir" || true
mkdir -p "$out_dir"
cp -r extension/static/* "$out_dir"
cat extension/lib/vue.min.js "$env_file" extension/src/*.js > "${out_dir}/inject.js"
