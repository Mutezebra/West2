#!/bin/bash

echo "ready to download ik"

ikpath="./repository/es/plugins/ik-download.sh"

if ! sh $ikpath ; then
  return 1;
fi

echo "ik download successfully"

docker-compose up -d