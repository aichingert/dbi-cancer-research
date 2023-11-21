#!/bin/sh

WORKING=$(pwd)

docker run -d --name postgresdb \
-p 5432:5432 \
-e POSTGRES_PASSWORD=lol \
postgres

sleep 3

cd parsing/setup-db
go run .

cd $WORKING

cd transform
cargo r -r &
cd ..

cd visualize
trunk serve --open &
cd ..
