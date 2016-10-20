#!/bin/bash

mkdir -p bin

for file in `ls *.go`
do
    binary=$(basename "$file" .go)
    go build -o bin/$binary $file
done
