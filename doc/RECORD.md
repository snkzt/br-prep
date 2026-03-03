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