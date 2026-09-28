#/bin/bash
protoc --go_out=./ --go-grpc_out=./ ./task_tracker.proto