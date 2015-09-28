#!/bin/bash

# Full clean build front+backend for the Z-Manager mulitenancy

E_BAD_GOGET=101
E_BAD_ASSETGEN=102
E_BAD_GOBUILD=103

#TODO: check for npm
hash "npm" 2>/dev/null || { echo >&2 "Node.js is not isntalled. Please install it and re-run."; exit 1; }
hash "go"  2>/dev/null || { echo >&2 "Golang is not isntalled. Please install it and re-run."; exit 1; }


#frontend
echo "cd web/ui"
pushd web/ui

echo "npm install"
npm install

echo "npm run build"
npm run build
#TODO(bzz): check for exit-code
popd


#backend
if [[ -z "${GOPATH}" ]] ; then
  export GOPATH="$PWD/third_party/go"
else
  export GOPATH="$PWD/third_party/go:$GOPATH"
fi

echo "cd ../server"
cd server

export GOPATH_BIN='../third_party/go/bin'
if [[ ! -f "$GOPATH_BIN/go-bindata" ]] ; then
  go get -u github.com/jteeuwen/go-bindata/...
  if [[ "$?" -ne 0 ]]; then
    echo "Unable to 'go get' github.com/jteeuwen/go-bindata" >&2
    exit "${E_BAD_GOGET}"
  fi
fi

echo "generating binary data assets"
"$GOPATH_BIN/go-bindata" -o ./assets.go ../web/ui/public/**
if [[ "$?" -ne 0 ]]; then
  echo "Unable to generate assets.go using go-bindata" >&2
  exit "${E_BAD_ASSETGEN}"
fi

echo "go build ./..."
go build ./...
if [[ "$?" -ne 0 ]]; then
  echo "Unable to 'go build' the server" >&2
  exit "${E_BAD_GOBUILD}"
fi
