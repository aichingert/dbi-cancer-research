cd setup-db || echo "fatal error"
go build setup-db || echo "consider installing go"
./setup-db || echo "tf"
cd ..