# JS Runtime

This is a custom JavaScript runtime built using Go and V8 (via v8go). It supports running JavaScript code and loading custom modules.

## Folder Structure

- `cmd/`: Entry point for the CLI application.
- `pkg/runtime/`: Core runtime logic (V8 initialization, JS bindings, etc.).
- `pkg/utils/`: Utility functions (e.g., logging).
- `js/modules/`: Custom JavaScript modules.
- `js/lib/`: Built-in JavaScript functionality (like `console.log`).
- `scripts/`: Build or other utility scripts.
- `tests/`: Unit and integration tests.
