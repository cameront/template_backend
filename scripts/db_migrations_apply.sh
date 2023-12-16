#!/bin/bash

set -e

source env-dev.sh

FAMILAIR_DIRS=("./ent" "./rpc")
for i in ${!FAMILAIR_DIRS[@]};
do
  dir=${FAMILAIR_DIRS[$i]}
  if [ ! -d $dir ]; then
    echo "Script running from `pwd`, and there is no ${dir} directory here, so we're probably not running from the root directory. Exiting..."
    exit 1
  fi
done

atlas migrate apply --dir "file://ent/migrate/migrations" --url "sqlite://${DB_FILE}?_fk=1"

echo "Successfully ran migrations"