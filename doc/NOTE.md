### generate Go code with proto buffer compiler:
protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative health.proto
protoc --proto_path=proto --go-grpc_out=proto --go-grpc_opt=paths=source_relative health.proto

0. Command
1. Path to the input proto files
2. Path to where to put generated file
3. Path to where the ourput should be (source_relative -> output will be with the input file) This controls how file paths are calculated relative to the input file.
How should protoc determine the final file path? Without paths=source_relative, protoc may instead use the go_package import path to decide directory structure, which can generate nested directories you did not expect.
4. Input file name

```
service UserService {
  rpc GetUser (UserRequest) returns (UserResponse) {}
}
```
https://grpc.io/docs/languages/go/quickstart/
The "{}" is for detaield rpc configuration.
For simpl rpc, we can just use conventional ";" as "  rpc GetUser (UserRequest) returns (UserResponse);"


### When to check which generated code files:
`*.pb.go` - This file contains message types.
Refer when you...
- Want to know the exact Go struct generated from a proto message.
- Need to confirm the field names in Go.
- Want to see what getters exist, like GetId() or GetIds().
- Need to know the exact field name to set when constructing a response.
I used to refer `type BatchGetUsersRequest struct` to implement client request call for the batch.

`*_grpc.pb.go`: This file contains service-related code.
Refer when you...
- Want to see the exact server interface you must implement.
- Want to confirm the method signature for an RPC.
- Want to know how the client interface looks.
- Want to see the exact registration function name.
I referred it t get method signatures for server `func (u *userServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {`.
