#!/bin/sh

DOCKER="docker"

cp ./Dockerfile ./tmp/
$DOCKER build tmp

git_hist_len=$(git log --pretty=oneline | wc -l | awk '{print $1}')
git_hash=$(git log --pretty="format:%h" | head -n1)

echo "Suggested tag: imcd:$git_hist_len.$git_hash"