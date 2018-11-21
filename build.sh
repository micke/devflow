mkdir -p build/mac
mkdir -p build/linux
go build --tags "static" -o build/mac/devflow main.go
cd build/mac
tar -czf ../mac.tar.gz devflow
cd ../..
shasum -a 256 build/mac.tar.gz
