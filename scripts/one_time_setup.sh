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


read -p 'What is your git repo name (e.g. cameront/foo-bar-repo)? ' REPO_NAME

echo ""
echo Repo Name: ${REPO_NAME}

read -p "Proceed? [y/n] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Not proceeding"
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

# Redirecting otuput to /dev/null prevents the line from being printed twice
title() {
    { set +x; } 2>/dev/null #echo off
    printf "\n###\n"
    printf "$@\n"
    { set -x; } 2>/dev/null # echo back on
}

# Note, repo names likely have "/" in them, so we need to use a non-default
# sed delimeter. Looks like @s are not allowed in repo names.

{ title "Step 1: renaming Go module and Go import paths"; } 2>/dev/null
# go module name
sed -i '' "s@go-svelte-sqlite-template@${REPO_NAME}@g" go.mod
# go import paths
find ./ -name '*.go' -exec sed -i '' "s@go-svelte-sqlite-template@${REPO_NAME}@g" {} \;

{ title "Step 2: verifying backend"; } 2>/dev/null
go get ./...
go build ./...

{ title "Step 3: verifying frontend"; } 2>/dev/null
pushd _ui 
npm install
npm run build
popd

{ title "Step 4: creating a database at data/database.db"; } 2>/dev/null
mkdir -p data
touch data/database.db

{ title "Step 5: migrating database"; } 2>/dev/null
./scripts/db_migrations_apply.sh

{ title "Script ran successfully."; } 2>/dev/null
{ set +x; } 2>/dev/null

printf "\nYou can run the backend with:\n  (source ./env-dev.sh && air)\n"
printf "And the frontend with: \n  npm run dev\n"
printf "\nThis script can be deleted, and you should probably set the new upstream: \"git remote set-url origin https://github.com/${REPO_NAME}.git\"\n"
