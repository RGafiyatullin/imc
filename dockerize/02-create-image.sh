#!/bin/sh

DOCKER="docker"
TMP_D="tmp"

cp ./Dockerfile "${TMP_D}/"
$DOCKER build "${TMP_D}"

git_hist_len=$(git log --pretty=oneline | wc -l | awk '{print $1}')
git_hash=$(git log --pretty="format:%h" | head -n1)

echo "Suggested tag: imcd:$git_hist_len.$git_hash"
