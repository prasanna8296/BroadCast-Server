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

- Go 1.24 or newer
- Dependencies: `github.com/gin-gonic/gin` and `github.com/gorilla/websocket`

# Clone and setup
git clone https://github.com/prasanna8296/BroadCast-Server.git
cd Broadcast-Server

# Download dependencies
go mod tidy

# Run server (in one terminal)
go run . start 8080

# Run client (in another terminal)
go run . connect localhost:8080

Any line typed into a client is sent to the server and broadcast to every
connected client.

## Conclusion

For more details about the project, check out the [Task Tracker Project on Roadmap.sh](https://roadmap.sh/projects/broadcast-server).
