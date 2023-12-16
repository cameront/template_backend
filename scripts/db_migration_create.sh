#!/bin/bash

set -e

if [ $# -eq 0 ]
  then
    echo "Usage: ./db_migration_create.sh [MIGRATION_NAME]"
fi

FAMILAIR_DIRS=("./ent" "./rpc")
for i in ${!FAMILAIR_DIRS[@]};
do
  dir=${FAMILAIR_DIRS[$i]}
  if [ ! -d $dir ]; then
    echo "Script running from `pwd`, and there is no ${dir} directory here, so we're probably not running from the root directory. Exiting..."
    exit 1
  fi
done

# Not strictly necessary, but will probably save someone's butt at leat once
echo "first running 'go generate ./ent'"
go generate ./ent

atlas migrate diff $1 --dir "file://ent/migrate/migrations" --to "ent://ent/schema" --dev-url "sqlite://file?mode=memory&_fk=1"

echo "Successfully created new migration"