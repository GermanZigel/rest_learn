package proto

//go:generate protoc --go_out=paths=source_relative:. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative  message.proto
