# Log Proxy and Deduplication Service

This project consists of two components: a caching log proxy server and a command-line client that fetches and deduplicates log data. The primary goal is to efficiently retrieve and process large Jenkins build logs for easier analysis.

## Overview

-   **Log Proxy Server**: An HTTP server that acts as a caching proxy for Jenkins build logs. It fetches logs from the public Jenkins CI server, caches them on disk to prevent redundant downloads, and serves them via a JSON API.
-   **Log Proxy Client**: A CLI tool that connects to the proxy server, fetches a specified build log, and applies deduplication techniques to minimize the output for human analysis.

---

## Features

-   **Efficient Caching**: The server downloads each log only once and caches it locally on disk.
-   **Concurrent Safe**: Handles multiple simultaneous requests for the same log without redundant downloads using a `sync.Map` to act as a per-resource lock.
-   **Upstream Server Error Handling**: Handles Internal errors as well jenkins server errors.
-   **Deduplication**: The client implements two forms of deduplication:
    1.  **Timestamp Removal**: Strips the leading `[timestamp]` from each log line.
    2.  **Consecutive Line Consolidation**: Identical consecutive lines are collapsed into a single line with a repeat count (e.g., `... (repeated 5 times)`).

---

## Project Structure

The project follows idiomatic Go standards for a project with multiple binaries:

```
/
├── cmd/
│   ├── client/   # Entry point for the log proxy client 
│   └── server/   # Entry point for the log proxy server
├── config/		  # Application config initialization
├── internal/
│   ├── cache/              # Caching and remote download logic for the server
│   ├── dedupe/             # Log deduplication processing logic
│   ├── logsapi/            # HTTP handlers and routing for the log proxy server
│   ├── logsclient/         # Log API client logic
│   └── server/             # Log proxy server initialization
├── pkg/                 	
│   ├── proxyerrors/        # Shared Error types
│   ├── httpclient/         # Shared HttpClient Wrapper
└── Makefile
```

---

## Usage

A `Makefile` is provided to simplify common operations.

1.  **Install Dependencies**:
    The project uses `gorilla/mux`. The dependencies will be fetched automatically on build.

    ```sh
    go mod tidy
    ```

2.  **Build Binaries**:
    This compiles both the server and client binaries into the project's root directory.

    ```sh
    make build
    ```

3.  **Run Tests**:
    This runs all unit tests for the project.

    ```sh
    make test
    ```

4.  **Run the Server**:
    The server will start on `http://localhost:8080`.

    ```sh
    make run-server
    ```

5.  **Run the Client**:
    In a separate terminal, run the client. `-build-id` sets to 7460 build by default.

    ```sh
    make run-client --build-id=7466
    ```

---

## Possible Improvements

- Better logger initialization
- Remote cached file to object store or database as blob data
- Better routing for GET and HEAD methods
- Redis for lock over go routine downloading logs for a build
- Better Error handling for api responses with error contract in json with log data
- Stop Signals for server
- Open API Spec for log proxy server
