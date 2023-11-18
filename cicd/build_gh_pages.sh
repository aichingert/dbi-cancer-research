#!/bin/sh

repo="dbi-cancer-research"

pushd "visualize"
	trunk build --release --public-url "aichingert/$repo"

	mkdir -p "../../dist"

	ls -lha dist
	mv dist/visu* ../../dist/
	ls -lha ../../dist
popd

ls -lha

ls -lha "../dist"

pushd "../dist"
	js=$(ls | grep front*\.js)
	wasm=$(ls | grep front*\.wasm)
	cp "../$repo/visualize/index.html cp.html"

	echo $(cat cp.html | 
		sed -e "s/<link data-trunk rel=\"rust\" data-wasm-opt=\"z\"\/>/<script type=\"module\">import init from '.\/$js';init('$wasm');<\/script>/g") > index.html

	rm cp.html
	echo "$js - $wasm"
popd

git switch gh-pages

rm -rf fron*
rm -rf inde*

mv ../dist .
mv dist/* .
rm -rf dist

