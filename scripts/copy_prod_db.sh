#!/bin/bash

# Meant to be run from the main repo directory

set -e

# Pull in the DB_FILE var
source env-dev.sh

read -r -p "Do you want to delete the existing (local) sqlite database? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY])
        rm -rf data/${DB_FILE}*
        ;;
    *)
        echo "not deleting db"
        ;;
esac

litestream restore -if-db-not-exists -o data/${DB_FILE} gcs://some-remote-path/some-db
