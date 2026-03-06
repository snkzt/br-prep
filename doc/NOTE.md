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


### Docker compose
- Key ideas: https://docs.docker.com/compose/intro/compose-application-model/#cli
- Starting docker container: `docker compose up -d`, `docker compose down`, `docker ps -a`, `docker compose logs db`
- Running and accessing to postgres DB: `docker compose exec db psql -U <USERNAME> -d <DBNAME> -c "SELECT 1;"`
- Execute sql migration file with psql command inside a container: `docker exec -i <CONATINERNAME> psql -U <USERNAME> -d <DBNAME like app_db> < migrations/001_users.sql`
- Start a DB session where -it means interactive terminal session: `docker exec -it <CONTAINERNAME> psql -U <USERNAME> -d <DBNAME>`
- Check config: `docker compose config`

### Run Redis with Docker
- `docker run -d --name redis -p 6379:6379 redis`
- https://hub.docker.com/_/redis?utm_source=chatgpt.com


### Interface and DI in Golang
Go Interfaces: Implicit Implementation
- Go uses implicit interface implementation.
- A type does not declare that it implements an interface. Instead, the Go compiler checks whether the type provides all required methods with exactly matching signatures.
  If it does, the type automatically satisfies the interface.
  Example:
  Interface:
  ```
  type UserRepository interface {
    GetUser(ctx context.Context, id int64) (*User, error)
    BatchGetUsers(ctx context.Context, ids []int64) ([]*User, error)
  }
  ```
  Implementation:
  ```
  type PostgresUserRepository struct {
    db *sql.DB
  }
  func (r *PostgresUserRepository) GetUser(ctx context.Context, id int64) (*User, error) {}
  func (r *PostgresUserRepository) BatchGetUsers(ctx context.Context, ids []int64) ([]*User, error) {}
  ```

- Even though PostgresUserRepository never mentions UserRepository, it implements the interface because it defines the same methods.
  This allows assignments like:
  ```
  var repo UserRepository
  repo = NewPostgresUserRepository(db)
  ```
- Important Rule: Method signatures must match exactly...
  - same method name
  - same parameters
  - same return types
  Otherwise compilation fails.

- Why Go uses implicit interfaces?
  1. Loose coupling: Implementations do not depend on interface definitions.
  2. Multiple implementations easily
    Example:
    ```
    InMemoryUserRepository
    PostgresUserRepository
    ```
  3. Easy testing
    You can create fake or mock implementations for tests.

- Optional compile-time check
  Sometimes developers add: `var _ UserRepository = (*PostgresUserRepository)(nil)`
  This forces the compiler to verify that the type implements the interface.

- Key idea
  An interface in Go describes behaviour (methods).
  Any type that provides those behaviours automatically satisfies the interface.

- DI
  Interface
  ```
  type UserRepository interface {
    GetUser(ctx context.Context, id int64) (*User, error)
  }
  ```
  Implementations
  ```
  InMemoryUserRepository
  PostgresUserRepository
  ```
  Dependency injection in main.go
  ```
  repo := NewPostgresUserRepository(db)
  userSvc := NewUserServiceServer(repo)
  ```
  Because both repositories satisfy the interface, the service does not care which one it receives.
  Mental model:
  ```
  Interface → defines behaviour
  Implementation → provides behaviour
  DI → decides which implementation is used
  ```