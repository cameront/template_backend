#!/bin/bash

set -e

FAMILAIR_DIRS=("./ent" "./_ui")
for i in ${!FAMILAIR_DIRS[@]};
do
  dir=${FAMILAIR_DIRS[$i]}
  if [ ! -d $dir ]; then
    echo "Script running from `pwd`, and there is no ${dir} directory here, so we're probably not running from the root directory. Exiting..."
    exit 1
  fi
done


read -p 'What is your git repo name? (e.g. cameront/foo-bar-repo) ' REPO_NAME

echo ""
echo Repo Name: ${REPO_NAME}

read -p "Proceed? [y/n] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Not proceeding"
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

echo "Step 1, renaming Go module and Go import paths"
# go module name
sed -i '' "s/go-svelte-sqlite-template/${REPO_NAME}/g" go.mod
# go import paths
find ./ -name '*.go' -exec sed -i '' "s/go-svelte-sqlite-template/${REPO_NAME}/g" {} \;

echo "Step 2, verifying backend"
go get ./...
go build ./...

echo "Step 3, verifying frontend"
pushd _ui 
npm install
npm run build
popd

echo "Step 4, creating a database at data/database.db"
mkdir -p data
touch data/database.db

echo "Step 5, migrating database"
./scripts/db_migrations_apply.sh

echo "Success."
echo "You can run the backend with '(source ./env-dev.sh && air)'"
echo "Now would be a good time to set the new upstream too like \"git remote set-url origin https://github.com/OWNER/REPOSITORY.git\""
