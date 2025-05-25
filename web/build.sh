#!/usr/bin/env nix-shell
#!nix-shell -i bash -p nodejs kubo

set -ve

npm install

npx webpack

mkdir -p site
cp index.html view.css controller.js nacatgunma.svg nacatgunma.ico site/

ipfs add --recursive --pin=false site/
