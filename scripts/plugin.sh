#!/bin/sh

pluginInstallPath="$HOME/.local/lib/nasin/plugins"

mkdir -p build
mkdir -p "$pluginInstallPath"

install() {
	go build -buildmode=plugin -o "build/$1.so" "./plugins/$1" && \
	cp "build/$1.so" $pluginInstallPath
}
