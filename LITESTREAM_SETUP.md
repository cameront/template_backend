# Setting Up Litestream

## Replicating on Google Cloud Storage

### Set up a replicator service account

1. Create the service account `gcloud iam service-accounts create database-replicator --description "litestream backup writer"`
1. Generate keys and store them locally `gcloud iam service-accounts keys create ./gcp_service_account_database_replicator.json --iam-account=database-replicator@$(gcloud config get project).iam.gserviceaccount.com`
1. Add `./gcp_service_account_database_replicator.json` to your .gitignore so you don't accidentally check in your service credentials.

### Set up the bucket to store db versions

1. Choose a bucket name (e.g. 'litestream-myproject') `BACKUPS_BUCKET_NAME='litestream-myproject'`
1. Create the bucket `gcloud storage buckets create gs://${BACKUPS_BUCKET_NAME}`
1. Grant access to the service account `gcloud storage buckets add-iam-policy-binding gs://${BACKUPS_BUCKET_NAME} --member=serviceAccount:database-replicator@$(gcloud config get project).iam.gserviceaccount.com --role=roles/storage.objectAdmin`

### Update your deployment

1. Update your dockerfile to add the credentials to your docker image, e.g. `COPY gcp_service_account_database_replicator.json /app/_creds/gcp_service_account_database_replicator.json`
1. In your production environment (e.g. fly.toml), set the GOOGLE_APPLICATION_CREDENTIALS env var (litestream looks for this) to the path of the creds in production, e.g. `GOOGLE_APPLICATION_CREDENTIALS="/app/_creds/gcp_service-acount_database_replicator.json"`
1. Set the correct db path and bucket location in your litestream.yml file

### Test (optional)

To ensure that permissions are correctly configured, you can attempt to restore the db to a temporary location
```
export GOOGLE_APPLICATION_CREDENTIALS="gcp_service_account_database_replicator.json"
litestream restore -parallelism 64 -if-db-not-exists -o /tmp/foo.db gcs://${BACKUPS_BUCKET_NAME}/some-backups-folder
```

You'll get the error "no matching backups found" if you haven't yet replicated the db to the bucket, but the permissions seem correct.

### Seed Database (optional)

To avoid the catch-22 situation described below, you may want to seed the production db with one you have on your local filesystem.

```
export GOOGLE_APPLICATION_CREDENTIALS="gcp_service_account_database_replicator.json"
litestream replicate data/path-to-database.db gcs://${BACKUPS_BUCKET_NAME}/some-backups-folder
```

### Deploy

Note that you start in kind of a catch-22 situation since when the first deployment runs, there's no db in google cloud storage to restore. So there's a magic environment variable used by the docker_entrypoint for this situation called "DB_FIRST_RUN_MODE", which will create a new database if one cannot be restored. So the steps are...

1. Set the DB_FIRST_RUN_MODE env var to anything non-empty e.g. `DB_FIRST_RUN_MODE="1"`
1. Deploy `fly deploy`
1. Verify via `fly logs` that the database is replicating and the application starts without issues
1. Unset the DB_FIRST_RUN_MODE and re-deploy