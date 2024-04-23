#!/bin/bash

export GOOSE_DRIVER=clickhouse
export GOOSE_DBSTRING="tcp://admin:admin@localhost:9000/test"
export GOOSE_MIGRATION_DIR="migrations"

if [ "$1" = "--dryrun" ]; then
  goose status -v
elif [ "$1" = "--down" ]; then
   goose down -v
else
   goose up -v
fi
