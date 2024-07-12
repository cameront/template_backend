# Prerequisites

### Install Go (backend language), Node (typescript frontend), Atlas (managing db schema), Protoc (protocol buffers )
```
brew install go
brew install node
brew install atlas
brew install protobuf
```

Install go tools: ProtoGen helps generate code from our API definition, ent manages our database entities, and air does live compile and reload during dev.

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go get -d entgo.io/ent/cmd/ent
go install github.com/air-verse/air@latest
```

Finally, run the one-time setup script to update repo paths and verify builds.

```
./scripts/one_time_setup.sh
```

# Setting up DB Replication

[Litestream Setup](/LITESTREAM_SETUP.md)

# How this was created

```
go mod init https://github.com/cameront/template_backend
```

```
mkdir template_backend
cd template_backend
```

```
npx degit sveltejs/template _ui
node _ui/scripts/setupTypeScript.js
pushd _ui
npm install
popd
```

And then lots of other shit... haha.