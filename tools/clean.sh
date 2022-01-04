#!/usr/bin/env bash

# clean deletes generated sources from the source tree

# delete generated go code from the api directory
find api -type f -iname "*.go" -exec rm -f {} \;