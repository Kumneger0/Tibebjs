# JavaScript Runtime Project Changelog

## Recent Changes and Improvements

We've recently refactored and improved the JavaScript runtime project. Here's a summary of the changes:

### 1. Code Restructuring

- **main.go**: Simplified and reduced in size. It now focuses on initializing the runtime and executing the script.
- **pkg/runtime/runtime.go**: New file that encapsulates the core runtime functionality.
- **pkg/modules/modules.go**: Updated to support ES modules instead of CommonJS.

### 2. ES Module Support

- Implemented basic ES module support, replacing the previous CommonJS-style `require` function.
- Added an `__import__` function that allows dynamic importing of modules.
- Modules are now executed in their own context, with support for `import.meta.url`.

### 3. Asynchronous Execution

- Both the main script and imported modules are now wrapped in async functions.
- This change allows for top-level `await` usage in scripts.

### 4. Improved Error Handling

- Consolidated error handling in the `main.go` file.
- More descriptive error messages throughout the codebase.

### 5. Runtime Encapsulation

- Created a `Runtime` struct in `pkg/runtime/runtime.go` that encapsulates the V8 isolate and context.
- Provides methods for setting up globals, executing scripts, and cleaning up resources.

### 6. Global Object Setup

- Moved the setup of global objects (like `console` and `__import__`) into the `Runtime.SetupGlobals` method.
- This change makes it easier to add or modify global objects in the future.

### 7. Script Execution

- The main script is now executed using the `Runtime.ExecuteScript` method.
- Scripts are wrapped in an async IIFE (Immediately Invoked Function Expression) to support top-level await.

## How to Use

1. Ensure you have the necessary dependencies installed.
2. Run your JavaScript file using: `go run main.go <path_to_your_script.js>`

## Future Improvements

- Implement module caching to improve performance for repeated imports.
- Add support for more ES module features and edge cases.
- Expand the global object to include more built-in JavaScript functionalities.

## Notes

This implementation provides basic ES module support. While it covers many common use cases, it may not fully comply with all aspects of the ES module specification. Future updates may address any limitations or edge cases.
