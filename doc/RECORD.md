Day1: gRPC foundation
Done
- gRPC server running
- Client receiving response
- HealthCheck working
- Understood proto → generation → registration flow
- Next: add real UserService proto
Covered
- Repo structure
- Proto generation workflow
- HealthService
- gRPC server wiring
- Reflection / testing flow

Day2: User Service with batching
Done
- Added GetUser, BatchGetUsers, error codes
- Data: in-memory slice
Covered
- UserService
- GetUser
- Proper gRPC error codes
- BatchGetUsers
- Behaviour decisions (preserve order, ignore unknown IDs)
- Thread safety (RWMutex)
- Working client call

Day3: Replace with Postgres repo, then Redis cache
- Introduce a repository interface and swap in Postgres, keeping the gRPC contract unchanged
- Postgres backing
- Redis caching
- Architecture: Each layer has one responsibility for `Separation of Concerns`.
  - Service layer: The service layer handles business logic and request orchestration. It validates input, applies business rules, calls repositories or other services, and converts results or domain errors into API responses such as gRPC status codes. It does not know how data is stored, only what data it needs.
  - Repository layer: The repository layer manages data access. It retrieves and stores data, hides the underlying storage implementation (database, cache, in memory structures, etc.), and returns domain objects. It does not decide request validity or API responses.
  - Why repositories do not return gRPC errors: Repositories should remain independent from transport layers such as gRPC or HTTP. They return domain errors (for example ErrUserNotFound). The service layer then translates those errors into transport specific responses, such as a gRPC NotFound or an HTTP 404.
```
Client
  ↓
Transport layer (gRPC)
  ↓
Service Layer: Business rules + request validation + coordination of operations + API behaviour.
  ↓
Repository Layer: Abstraction for retrieving and persisting data.
  ↓
Storage (map or Postgres)
```