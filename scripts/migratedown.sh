#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd db/schema
goose postgres $connection_string down
