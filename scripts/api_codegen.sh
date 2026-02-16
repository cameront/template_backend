#/bin/bash

set -e

FAMILAIR_DIRS=("./_ui" "./rpc" "./rpc_internal")
for i in ${!FAMILAIR_DIRS[@]};
do
  dir=${FAMILAIR_DIRS[$i]}
  if [ ! -d $dir ]; then
    echo "Script running from `pwd`, and there is no ${dir} directory here, so we're probably not running from the root directory. Exiting..."
    exit 1
  fi
done

CODEGEN_DIR_GO="."
CODEGEN_DIR_TS="./_ui/src/codegen"
PROTOC_GEN_TWIRP_BIN="./_ui/node_modules/.bin/protoc-gen-twirp_ts"
PROTOC_GEN_TS_BIN="./_ui/node_modules/.bin/protoc-gen-ts"

function codegen() {
  PROTO_DIR="$1"
  PROTO_PATH="${PROTO_DIR}/$2"
  echo "Generating code for ${PROTO_PATH}..."
  # Go
  protoc --twirp_out=. --twirp_opt=paths=source_relative --go_out=${CODEGEN_DIR_GO} --go_opt=paths=source_relative ${PROTO_PATH}
  # TS
  protoc -I ${PROTO_DIR} --plugin=protoc-gen-ts=${PROTOC_GEN_TS_BIN} --plugin=protoc-gen-twirp_ts=${PROTOC_GEN_TWIRP_BIN} --ts_out=${CODEGEN_DIR_TS} --twirp_ts_out=./${CODEGEN_DIR_TS} ${PROTO_PATH}
}

codegen "rpc/count" "countservice.proto"

