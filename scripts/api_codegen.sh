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

PROTO_DIR="rpc/count"
PROTO_PATH="${PROTO_DIR}/countservice.proto"

CODEGEN_DIR_GO="."
CODEGEN_DIR_TS="./src/codegen"

# Go code generation
protoc --twirp_out=. --twirp_opt=paths=source_relative --go_out=${CODEGEN_DIR_GO} --go_opt=paths=source_relative ${PROTO_PATH}

# TS code generation
pushd _ui
npx protoc --ts_out ${CODEGEN_DIR_TS} --proto_path ../${PROTO_DIR} ../${PROTO_PATH}
popd
