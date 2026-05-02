# Junie AI Guidelines

This directory contains instructions for Junie (the AI agent by JetBrains) to follow when working on this project.

## Core Instructions

- **Always follow Clean Architecture and DDD principles.**
- **Project Style**: Use modern Go idioms (Go 1.25+).
- **Documentation**: Refer to the detailed skill guides in the `skills/` directory for specific implementation patterns.

## Skills Reference

- [Architecture](../skills/architecture.md): Layers and dependency rules.
- [API](../skills/api.md): Protobuf and code generation.
- [UseCase](../skills/usecase.md): Generic UseCase pattern and registration.
- [Interface](../skills/interface.md): gRPC, CLI, and Kafka adapters.
- [Infrastructure](../skills/infrastructure.md): DB and Kafka implementations.
- [Domain](../skills/domain.md): Entities and pure business logic.

## Key Rules for Junie

1. **API-First**: Always edit `.proto` files in the `api/` directory first, then run `make generate` to update the generated code in `internal/api/`.
2. **UseCase Pattern**: Business logic MUST be implemented as a `UseCase[I, O]` struct in `internal/usecase/`. It MUST be registered in `internal/app/app.go`.
3. **No Direct Logic in Handlers**: gRPC handlers, CLI commands, and Kafka consumers must only act as adapters. They should retrieve a usecase from the registry and call its `Do` method.
4. **Dependency Injection**: Pass dependencies via constructors. Use interfaces for dependencies to allow for easy testing and decoupling.
5. **Modern Go**: Use modern Go features like `any`, `slices` package, `maps` package, and `for i := range n`.

## Specific Go Guidelines

See [.junie/go.md](./go.md) for detailed Go coding standards and practices.
