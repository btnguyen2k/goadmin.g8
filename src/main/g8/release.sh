#!/bin/sh

## Utility script to release project with a tag
## Usage:
##   ./release.sh <tag-name>

if [ "$1" == "" ]; then
	echo "Usage: $0 tag-name"
	exit -1
fi

echo "$1"
git commit -m "$1"
git tag -f -a "$1" -m "$1"
git push origin "$1" -f
git push

