# template_backend

This is a template application I use frequently to spin simple applications up quickly. The stack is essentially:

* DB: Sqlite (schema managed by [Atlas](https://atlasgo.io/), entities managed by [Ent](https://entgo.io/))
* Backend: Go
* Frontend: React (in Typescript)
* RPC API: [Twirp](https://github.com/twitchtv/twirp)

The golang backend implements the twirp API server and also serves static FE files (index.html, css, js, etc). The protocol buffer API definition generates both go code for the server implementation and the client-side typescript code the frontend to use. The sqlite db is replicated to S3 (in production) by litestream. 

# How to deploy

TODO

# Things I'd like to improve

1. Serve static files via the go static file server without giving up Hot Module Replacement, so that you don't have to use the node webserver in development.
1. Split the one_time_setup script into 2 different scripts: onetime/rename_repo.sh and onetime/verify.sh, so that the latter can be used when changes need to be made to this (template) repo but you want to try out the working demo before comitting changes.

# Steps

1. `git clone https://github.com/cameront/template_backend [your_directory]`
1. `./scripts/one_time_setup.sh`
1. `(source env_dev.sh && air)`
1. `pushd _ui && pnpm run dev`

# Optional Additional Steps

## Rename ish

Rename the "count" service to the name of the service ([servicename]) you want to build.

1. replace rpc/count/countservice.proto package names with references to [servicename]
1. replace scripts/protogen.sh with references to [servicename]
1. `mv rpc/count/countservice.proto rpc/[servicename]/[servicename]service.proto`
1. `mv rpc/count -> rpc/[servicename]`
1. `rm -rf rpc/count`
1.  ./scripts/api_codegen.sh
1. optionally replace the server port (5001) and the UI dev port (5000) with those of your liking
1. remove this readme and the scripts/one_time_setup.sh script