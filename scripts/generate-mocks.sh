#!/bin/bash

this_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $this_dir/../

if ! command -v mockery &> /dev/null; then
  echo "ERROR: couldn't find mockery in \$PATH, cannot generate mocks"
  exit 1
fi
find ./mocks -mindepth 1 -maxdepth 1 -exec rm -rf '{}' \;
mockery -all -keeptree
