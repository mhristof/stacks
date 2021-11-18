#! /usr/bin/env bash
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

function commit {
    git checkout "$1" || git checkout -b "$1"
    date >>"$1.md"
    git add "$1.md"
    git commit "$1.md" -m "$1: commiting"
}

rm -rf test-repo
mkdir -p test-repo
cd test-repo
git init
touch main.md
git add main.md
git commit main.md -m 'initial import'
commit main
commit main
commit feat1
commit feat1.1
commit main
commit feat1.1

exit 0
