# go-svelte-sqlite-template

This is a template application I use frequently to spin simple applications up quickly. The stack is essentially:

* DB: Sqlite (schema managed by Atlas)
* Backend: Golang
* Frontend: Svelte (in Typescript)
* RPC API: Twirp

The golang backend implements the twirp API server and also serves static FE files (index.html, css, js, etc). The protocol buffer API definition generates both go code for the server implementation and the client-side typescript code the frontend to use. The sqlite db is replicated to S3 (in production) by litestream. 

# How to start

TODO: add gonew instructions
TODO: add instructions to rename the server, the db, and remove the default entities.

# Rename ish

Rename the "count" service to the name of the service ([servicename]) you want to build.

1. replace rpc/count/countservice.proto package names with references to [servicename]
1. replace scripts/protogen.sh with references to [servicename]
1. `mv rpc/count/countservice.proto rpc/[servicename]/[servicename]service.proto`
1. `mv rpc/count -> rpc/[servicename]`
1. `rm -rf rpc/count`
1.  ./scripts/protogen.sh
1. optionally replace the server port (5001) and the UI dev port (5001) with those of your liking

```

# How to develop

# How to deploy

