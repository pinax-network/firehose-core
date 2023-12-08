#!/usr/bin/env bash
# Copyright 2019 dfuse Platform Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

main() {
  current_dir="`pwd`"
  trap "cd \"$current_dir\"" EXIT
  pushd "$ROOT" &> /dev/null

  cmd=$1; shift

  if [[ $cmd == "" ]]; then
    usage_error "argument <cmd> is required"
  fi

  if [[ ! -d "./cmd/$cmd" ]]; then
    usage_error "argument <cmd> is invalid, valid ones are: \"`ls ./cmd | xargs | tr ' ' ','`\""
  fi

  echo "Building $cmd..."
  go build -o "./$cmd" "./cmd/$cmd"
  if [[ $? != 0 ]]; then
    echo "Build failed"
    exit 1
  fi

  command="./$cmd `cmd_args $cmd $@`"
  echo "Starting $command"
  exec $command
}

cmd_args() {
  cmd=$1; shift

  backup_dir="$ROOT/data/$cmd/backups"
  config_dir="$ROOT/data/$cmd/config"
  data_dir="$ROOT/data/$cmd/storage"
  snapshot_dir="$ROOT/data/$cmd/snapshots"
  mindreader_dir="$ROOT/data/$cmd/deep-mind"

  geth_args="--geth-binary geth --data-dir $data_dir --backup-store-url $backup_dir"
  nodeos_args="--nodeos-path nodeos --config-dir $config_dir --data-dir $data_dir --backup-store-url $backup_dir --snapshot-store-url $snapshot_dir"

  case "$cmd" in
  geth_manager)
      args="$geth_args"
      extra_flags="--deep-mind-block-progress"
      ;;
  geth_mindreader)
      args="$geth_args --oneblock-store-url=$mindreader_dir/blocks --working-dir=$mindreader_dir/tmp"
      extra_flags="--deep-mind"
      ;;
  nodeos_manager)
      args="$nodeos_args "
      ;;
  nodeos_mindreader)
      args="$nodeos_args --oneblock-store-url=$mindreader_dir/blocks --working-dir=$mindreader_dir/tmp"
      ;;
  esac

  extra_flags_processed=""
  for arg in "$@"; do
    if [[ "$arg" == "--" ]]; then
      args+=" -- $extra_flags"
      extra_flags_processed="true"
    else
      args+=" $arg"
    fi
  done

  if [[ $extra_flags_processed == "" ]]; then
    args+=" -- $extra_flags"
  fi

  echo "$args"
}

usage_error() {
  message="$1"
  exit_code="$2"

  echo "ERROR: $message"
  echo ""
  usage
  exit ${exit_code:-1}
}

usage() {
  echo "usage: serve.sh <cmd> [<cmd arguments>]"
  echo ""
  echo "Build & serve the appropriate manager/mindreader operator"
  echo ""
  echo "Valid <cmd>"
  ls "$ROOT/cmd" | xargs | tr " " "\n" | sed 's/^/ * /'
  echo ""
}

main $@