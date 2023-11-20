#!/bin/sh

DIR="db"
WORKING=$(pwd)

if test ! -d $DIR; then
	mkdir $DIR
	chmod -R o+w $DIR

	docker run -d --name postgresdb \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD=lol \
  postgres
fi

cd parsing/setup-db
go run .

cd $WORKING

#pushd transform
#
#cargo r -r &
#
#popd
#
#pushd visualize
#
#trunk serve --open
#
#popd
