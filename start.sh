#!/bin/sh

DIR="db"

if test ! -d $DIR; then
	mkdir $DIR
	chmod -R o+w $DIR

	docker run -d --name postgresdb \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD=lol \
  postgres
fi

pushd parsing/setup-db

go run .

popd

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