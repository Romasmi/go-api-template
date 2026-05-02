# Go Coding Guidelines for Junie

These guidelines are based on JetBrains Junie recommendations for Go projects, adapted for this Clean Architecture template.

## 1. Project Structure
- Follow the Clean Architecture layers: `domain`, `usecase`, `interface`, `infrastructure`.
- Keep the `api/` folder for source `.proto` files only.
- Generated code resides in `internal/api/`.

## 2. Dependency Injection
- Always use explicit construction.
- Pass dependencies as interfaces to constructors (`New...` functions).
- Avoid global state and `init()` functions for logic.

## 3. Error Handling
- Wrap errors with context: `fmt.Errorf("failed to do X: %w", err)`.
- Use `errors.Is` and `errors.As` for error checking.
- Don't just return errors; provide enough context to trace where they came from.

## 4. Concurrency and Context
- Always propagate `context.Context` through all layers.
- Use context for cancellation and timeouts.
- Be careful with goroutines; ensure they are properly managed and terminated.

## 5. Database Access
- Use `pgx` for PostgreSQL access.
- Implement repository interfaces in `internal/infrastructure/db/postgres/`.
- Use context-aware database calls.

## 6. Logging
- Use structured logging (e.g., `log` package or `slog`).
- Include relevant identifiers (like `userID`) in logs but avoid sensitive data.

## 7. Testing
- Use `t.Context()` for tests that need context.
- Write unit tests for use cases and domain logic.
- Use mocks/interfaces for external dependencies.

## 8. Modern Go Idioms (Go 1.25+)
- Use `any` instead of `interface{}`.
- Use `for i := range n` for simple loops.
- Use `slices.Contains`, `slices.Max`, etc., from the standard library.
- Use `omitzero` in JSON tags for structs/times/slices.
- Use `b.Loop()` in benchmarks.

## 9. API Design
- Follow the API-first approach with Protobuf.
- Use gRPC for internal service communication.
- Use gRPC-Gateway for providing a RESTful JSON API.
- Serve Swagger UI for documentation.
