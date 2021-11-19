#! /usr/bin/env bash

#
# Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

GOP_ROOT=$(pwd)
GOP_HOME_DIR="$HOME/gop"
GOP_CACHE_DIR="$GOP_ROOT/.gop"

ADD_GOPATH_COMMAND="export PATH=\$PATH:\$GOPATH/bin"
ADD_GO_BIN_COMMAND="export PATH=\$PATH:\$HOME/go/bin"
MANUAL_EXPORT_COMMAND=""

GIT_BRANCH=$(git branch --show-current)
GIT_COMMIT_HASH=$(git rev-parse --verify HEAD)
BUILD_DATE=$(date '+%Y-%m-%d_%H-%M-%S')
GO_FLAGS="-X github.com/goplus/gop/build.Date=${BUILD_DATE} \
  -X github.com/goplus/gop/build.Commit=${GIT_COMMIT_HASH} \
  -X github.com/goplus/gop/build.Branch=${GIT_BRANCH}"

command_exists() {
	command -v "$@" >/dev/null 2>&1
}

build_go_plus_tools() {
	command_exists go || {
    echo "Error: go is not installed but required, please visit https://golang.org/doc/install for help."
		exit 1
	}

  COMMANDS_DIR="$GOP_ROOT/cmd"
  if [ ! -e "$COMMANDS_DIR" ]; then
    echo "Error: This shell script should be run at root directory of gop repository."
    exit 1
  fi

  echo "Installing Go+ tools..."
  cd $COMMANDS_DIR

  # will be overwritten by gop build
  go install -v -ldflags "${GO_FLAGS}" ./...

  echo "Go+ tools installed successfully!"
}

clear_gop_cache() {
  echo "Clearing gop cache"
  cd $GOP_ROOT
  if [ -d "$GOP_CACHE_DIR" ]; then
    rm -r $GOP_CACHE_DIR
    echo "Gop cache files cleared"
  else
    echo "No gop cache files found"
  fi
}

link_gop_root_dir() {
  echo "Linking $GOP_ROOT to $GOP_HOME_DIR"
  cd $GOP_ROOT
  if [ ! -e "$GOP_HOME_DIR" ] && [ "$GOP_ROOT" != "$GOP_HOME_DIR" ]; then
    ln -s $GOP_ROOT $GOP_HOME_DIR
  fi
  echo "$GOP_ROOT linked to $GOP_HOME_DIR successfully!"
}

summary() {
  echo "Installation Summary:"
  echo "Go+ is now installed."
  echo ""
  if [ -n "$MANUAL_EXPORT_COMMAND" ]; then
    cat <<-EOF
Notice, we just temporarily add gop command and other tools into your PATH.
To make it permanent effect, you should manually add below command:

${MANUAL_EXPORT_COMMAND}

to your shell startup file, such as: ~/.bashrc, ~/.zshrc...
type 'go help install' for more details.

EOF
  fi
}

gop_test() {
  echo "Running gop test"
  cd $GOP_ROOT
  PATH=$PATH:$GOPATH/bin gop test -v -coverprofile=coverage.txt -covermode=atomic ./...
  echo "Finished running gop test"
}

default() {
  # Build all Go+ tools
  build_go_plus_tools

  # Clear gop cache files
  clear_gop_cache

  # Link Gop root directory to home/ dir
  link_gop_root_dir

  # Summary
  summary
}

if [ "$#" -eq 0 ]; then
  default
  exit 0
fi

# To add more options below, juse add another case.
while [ "$#" -gt 0 ]; do
  case "$1" in
    -c|--compile)
      build_go_plus_tools
      ;;
    -t|--test)
      gop_test
      ;;
    -*)
      echo "Unknown option: $1"
      echo "Valid options:"
      echo "  -t, --test     Running testcases with gop test"
      echo "  -c, --compile  Compile gop and related tools"
      exit 1
      ;;
  esac
  shift
done
