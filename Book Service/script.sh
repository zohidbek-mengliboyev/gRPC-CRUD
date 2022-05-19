protoc --proto_path=proto/book_service --gofast_out=plugins=grpc:. book.proto
protoc --proto_path=proto/user_service --gofast_out=plugins=grpc:. user.proto
go install github.com/gogo/protobuf/protoc-gen-gofast@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go get -u github.com/golang/protobuf/protoc-gen-go