#!/usr/bin/env bash

VERSION=$(cat ./VERSION)
MAJOR=$(echo $VERSION | cut -d. -f1)
MINOR=$(echo $VERSION | cut -d. -f2)
PATCH=$(echo $VERSION | cut -d. -f3)

if [ "$1" = "major" ]; then
    MAJOR=$(expr $MAJOR + 1)
    MINOR=0
    PATCH=0
fi

if [ "$1" = "minor" ]; then
    MINOR=$(expr $MINOR + 1)
    PATCH=0
fi

if [ "$1" = "patch" ]; then
    PATCH=$(expr $PATCH + 1)
fi

VER=$(printf "%d.%d.%d" $MAJOR $MINOR $PATCH)
echo $VER > ./VERSION

sed -i "" -e "s/VERSION=.*/VERSION=$VER/g" ./scripts/install.sh