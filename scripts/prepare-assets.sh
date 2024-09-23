#!/bin/bash

set -e

# ensure existence of release folder
if ! [ -d "./release" ]; then
    mkdir ./release
fi

# ensure zip is installed
if [ "$(which zip)" = "" ]; then
    apt-get update && apt-get install -y zip
fi

# add execution permission
chmod 750 ./build/pongo-darwin-amd64
chmod 750 ./build/pongo-darwin-arm64
chmod 750 ./build/pongo-linux-amd64
chmod 750 ./build/pongo-linux-arm64
chmod 750 ./build/pongo-linux-riscv64
chmod 750 ./build/pongo-windows-amd64.exe
chmod 750 ./build/pongo-windows-arm64.exe

# create archives
zip -j ./release/pongo-darwin-amd64.zip ./build/pongo-darwin-amd64
zip -j ./release/pongo-darwin-arm64.zip ./build/pongo-darwin-arm64
zip -j ./release/pongo-linux-amd64.zip ./build/pongo-linux-amd64
zip -j ./release/pongo-linux-arm64.zip ./build/pongo-linux-arm64
zip -j ./release/pongo-linux-riscv64.zip ./build/pongo-linux-riscv64
zip -j ./release/pongo-windows-amd64.zip ./build/pongo-windows-amd64.exe
zip -j ./release/pongo-windows-arm64.zip ./build/pongo-windows-arm64.exe

# calculate checksums
for file in  ./release/*; do
	checksum=$(sha256sum "${file}" | cut -d' ' -f1)
	filename=$(echo "${file}" | rev | cut -d/ -f1 | rev)
	echo "${checksum} ${filename}" >> ./release/checksums_sha256.txt
done
