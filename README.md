# Broadcast Server

## Overview

The Broadcast Server is a simple command-line application that allows clients to connect and send messages that are broadcast to all connected clients. This project demonstrates real-time communication using WebSockets, helping to understand the underlying principles of applications such as chat systems and live scoreboards.

## Features

- Start a broadcast server that listens for client connections.
- Connect multiple clients to the server.
- Broadcast messages from any connected client to all other clients.
- Handle client connections and disconnections gracefully.
- Simple command-line interface for server and client operations.

## Requirements

- Go 1.16^ (or any programming language of your choice)
- `websockets` library (if using Python)

# Clone and setup
git clone https://github.com/yourusername/broadcast-server.git
cd broadcast-server

# Initialize Go module
go mod init broadcast
go mod tidy

# Install dependencies
go get github.com/gin-gonic/gin
go get github.com/gorilla/websocket

# Run server (in one terminal)
go run main.go start 8080

# Run client (in another terminal)
go run main.go connect localhost:8080

## Conclusion

For more details about the project, check out the [Task Tracker Project on Roadmap.sh](https://roadmap.sh/projects/broadcast-server).
