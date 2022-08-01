#!/bin/sh

if ! [[ $GITHUB_ACTION ]]; then
    echo "formatting"
    golines pkg resolver internal daos cmd models schema testutls -w --shorten-comments --reformat-tags -m 128
fi