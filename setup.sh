#!/bin/sh

DIR="./db"

if test ! -d $DIR; then
	mkdir $DIR
	chmod -R o+w $DIR
fi

docker run -d --name oradb \
-p 1521:1521 \
-e ORACLE_PWD=lol \
-e ORACLE_CHARACTERSET=AL32UTF8 \
-v $DIR:/opt/oracle/oradata \
container-registry.oracle.com/database/free:latest
