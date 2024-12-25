# TibebJS

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A JavaScript runtime experiment inspired by [Roll your own JavaScript runtime](https://deno.com/blog/roll-your-own-javascript-runtime). This is just an experiment how JavaScript runtimes work and is not intended for production use.

</div>

## Features

- JavaScript runtime built on V8 engine
- File system operations
- Built-in HTTP server
- Console API implementation
- Timer functions (setTimeout, setInterval)

## Prerequisites

- Go 1.23 or higher
- GCC (for v8go compilation)

> Note: Windows is not supported as v8go [dropped Windows support](https://github.com/rogchap/v8go/pull/234)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/Kumneger0/Tibebjs.git
cd Tibebjs
```

2. Install dependencies:
```bash
go mod download
```

## Building

### On Linux
```bash
go build -o tibebjs
```

### On macOS
```bash
go build -o tibebjs
```

Note: Due to CGO dependencies (v8go), cross-compilation requires additional setup. It's recommended to build on the target platform directly.

## Running

Execute JavaScript files:
```bash
./tibebjs path/to/your/script.js
```

## Project Structure

- `pkg/`
  - `runtime/`: Core runtime implementation
  - `eventloop/`: Event loop and async operations
  - `console/`: Console API implementation
  - `timer/`: Timer functionality
  - `net/`: Network operations
  - `fs/`: File system operations
  - `fetch/`: HTTP client implementation
  - `utils/`: Utility functions
- `js/`: JavaScript example files
- `globals/`: Global objects and bindings

## Available APIs

### File System Operations
```javascript
// Write to file
await Tibeb.writeFile('file.txt', 'content');

// Read from file
const content = await Tibeb.readFile('file.txt');

// Remove file
await Tibeb.rmFile('file.txt');

// Rename file
await Tibeb.renameFile('old.txt', 'new.txt');
```

### HTTP Server
```javascript
// Create a simple HTTP server
Tibeb.serve((request) => {
  // Request object contains url, method, and headers
  const response = {
    url: request.url,        // URL path of the request
    method: request.method,  // HTTP method (GET, POST, etc.)
    headers: request.headers // Request headers
  };
  
  return response(JSON.stringify(response), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
}, 3000); // Listen on port 3000

// Example with routing
Tibeb.serve((request) => {
  switch(request.url) {
    case "/":
      return response(JSON.stringify({ 
        path: "home",
        method: request.method,
        headers: request.headers 
      }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      });
      
    case "/api":
      return response(JSON.stringify({ 
        path: "api",
        method: request.method,
        headers: request.headers 
      }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      });
      
    default:
      return response(JSON.stringify({ 
        error: "Not Found",
        path: request.url,
        method: request.method 
      }), {
        status: 404,
        headers: { "Content-Type": "application/json" }
      });
  }
}, 3000);
```

## Resources

- [Roll your own JavaScript runtime](https://deno.com/blog/roll-your-own-javascript-runtime) - Original Deno blog post
- [v8go Documentation](https://pkg.go.dev/rogchap.com/v8go) - Go V8 bindings documentation

## Contributing

Feel free to submit a Pull Request or open an issue for discussion.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
