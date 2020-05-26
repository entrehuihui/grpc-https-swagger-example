

protoc -I=. --go_out=plugins=grpc:. *.proto
protoc -I=. --grpc-gateway_out=logtostderr=true:. *.proto
protoc -I=. --swagger_out=logtostderr=true:. *.proto

go-bindata --nocompress -pkg swagger -o proto/swagger/datafile.go proto/third_party/swagger-ui/...