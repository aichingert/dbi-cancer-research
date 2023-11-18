#!/bin/sh

set -e

repo="dbi-cancer-research"

pushd "visualize"
	trunk build --release --public-url "aichingert/$repo"

	mkdir -p "../../dist"

	ls -lha dist
	mv dist/visu* ../../dist/
	ls -lha ../../dist
popd

pushd "../dist"
	js=$(ls | grep visu*\.js)
	wasm=$(ls | grep visu*\.wasm)
	cp "../$repo/visualize/index.html" "cp.html"

	echo $(cat cp.html | 
		sed -e "s/<link data-trunk rel=\"rust\" data-wasm-opt=\"z\"\/>/<script type=\"module\">import init from '.\/$js';init('$wasm');<\/script>/g") > index.html

	rm cp.html
	echo "$js - $wasm"
popd

git switch gh-pages

mv ../dist .
mv dist/* .
rm -rf dist
