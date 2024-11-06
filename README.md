# WAV to FLAC Streaming Audio Converter

The WAV to FLAC Streaming Audio Converter is a robust backend service built in Go that provides real-time conversion of WAV audio streams to FLAC format. This service is designed to integrate with a frontend application, enabling users to test audio submissions through a base URL.

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Local Setup](#local-setup)
  - [Docker Deployment](#docker-deployment)
- [Testing](#testing)
  - [Unit Tests](#unit-tests)
  - [Integration Tests](#integration-tests)
- [WebSocket API Usage](#websocket-api-usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Efficient streaming of WAV audio data to the server
- Real-time conversion of WAV to FLAC format
- Streaming of FLAC data back to the client
- Handles multiple simultaneous connections
- Graceful error handling and resilient to connection issues
- Optimized for low-latency audio processing
- Comprehensive test suite for unit and integration testing

## Architecture

The project follows a modular structure with the following key components:

- **cmd/server/main.go**: Entry point of the application, sets up the Fiber web framework and handles server startup and shutdown.
- **internal/handlers**: Defines the WebSocket route and handles the audio conversion logic.
- **internal/services**: Implements the core audio conversion functionality.
- **internal/middleware**: Provides logging and rate-limiting middleware.
- **internal/models**: Defines data structures related to audio processing and conversion.
- **pkg/utils**: Implements utility functions for audio data manipulation and validation.
- **tests/unit**: Contains unit tests for the core conversion logic.
- **tests/integration**: Includes integration tests for end-to-end testing of the WebSocket API and conversion process.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)

### Local Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/audio-converter.git
   cd audio-converter
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the application:

   ```bash
   go run cmd/server/main.go
   ```

### Docker Deployment

1. Build and run using Docker Compose:

   ```bash
   docker-compose up --build
   ```

## Testing

### Unit Tests

To run the unit tests:

```bash
go test ./tests/unit/...
```

### Integration Tests

To run the integration tests:

```bash
go test ./tests/integration/...
```

The integration test suite includes:

1. Basic WebSocket connection and conversion test
2. Multiple simultaneous connections test
3. Large file conversion test
4. Error handling tests for various scenarios
5. Connection resilience test
6. Performance test

## WebSocket API Usage

Connect to the WebSocket endpoint:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/convert');

// Send WAV data
ws.send(wavData);

// Receive FLAC data
ws.onmessage = function(event) {
    const flacData = event.data;
    // Handle FLAC data
};
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the [MIT License](LICENSE).
