#! /usr/bin/env bash

set -e

GOP_ROOT=$(pwd)
GOP_HOME_DIR="$HOME/gop"

ADD_GOPATH_COMMAND="export PATH=\$PATH:\$GOPATH/bin"
ADD_GO_BIN_COMMAND="export PATH=\$PATH:\$HOME/go/bin"
MANUAL_EXPORT_COMMAND=""

GIT_COMMIT_HASH=$(git rev-parse --verify HEAD)
BUILD_DATE=$(date '+%Y-%m-%d_%H-%M-%S')
GIT_BRANCH=$(git symbolic-ref --short -q HEAD)
GO_FLAGS="-X github.com/goplus/gop/build.Date=${BUILD_DATE} \
  -X github.com/goplus/gop/build.Commit=${GIT_COMMIT_HASH} \
  -X github.com/goplus/gop/build.Branch=${GIT_BRANCH}"
GOP_CACHE_DIR="$GOP_ROOT/.gop"

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

build_go_plus_tutorials() {
  command_exists gop || {
    if [ -n "$GOPATH" ]; then
      MANUAL_EXPORT_COMMAND="$ADD_GOPATH_COMMAND"
      echo "Execute command: $ADD_GOPATH_COMMAND"
      export PATH=$PATH:$GOPATH/bin
    else
      MANUAL_EXPORT_COMMAND="$ADD_GO_BIN_COMMAND"
      echo "Execute command: $ADD_GO_BIN_COMMAND"
      export PATH=$PATH:$HOME/go/bin
    fi
  }

  command_exists gop || {
    echo "Error: something wrong, you could create a new issue on https://github.com/goplus/gop/issues, we will help you."
    exit 1
  }

  cd $GOP_ROOT

  echo "Building all Go+ tutorials."
  gop install -ldflags "${GO_FLAGS}" ./...
  echo "Go+ tutorials builded successfully!"
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

hello_world() {
  HELLO_WORLD_COMMAND="gop run tutorial/01-Hello-world/hello.gop"

  cd $GOP_ROOT

  EXPORT_CMD=""
  if [ -n "$MANUAL_EXPORT_COMMAND" ]; then
    EXPORT_CMD="${MANUAL_EXPORT_COMMAND} && "
  fi

  echo '-----------------------------'
  cat <<-EOF
Let's have a try now, Copy below command to get classic Hello, world!

${EXPORT_CMD}${HELLO_WORLD_COMMAND}

Besides, there are another more tutorials listed under ./tutorial/ directory.

Have fun!
EOF
}

# Build all Go+ tools
build_go_plus_tools

# Clear gop cache files
clear_gop_cache

# Link Gop root directory to home/ dir
link_gop_root_dir

# Build all Go+ tutorials
build_go_plus_tutorials

# Summary
summary

# hello world
hello_world
