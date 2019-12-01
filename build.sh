#!/bin/bash
clear

usage="build.sh TAG"

if [ -z "$1" ]; then
  echo $usage
  exit 1
fi;

docker build --rm -t figassis/mysql-backup:$1 . && docker push figassis/mysql-backup:$1