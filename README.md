# redis-go: A Simple Redis-like Server and Client in Go

 redis-go is a minimal Redis-like server and client implemented in Go. This project covers basic Redis commands and communication protocols.

## Features

- Basic Redis commands support: SET, GET, PING, KEYS, DEL, HGET, HSET, ZADD, ZRANGE.
- Simple client for interacting with the server.

## Getting Started

Follow these steps to set up and run the redis-go server and client:

### Prerequisites

- Go installed on your machine.

### Server

1. Clone this repository.

```bash
git clone https://github.com/atul107/redis-go.git
```

2. Navigate to the project root directory:

```bash
cd redis-go
```
3. Start the redis-go server:
```bash
make server
```

The server will listen on 127.0.0.1:6379.

### Client

1. Clone this repository (if not already done).

```bash
git clone https://github.com/atul107/redis-go.git
```
2. Navigate to the project root directory:
```bash
cd redis-go
```
3. Start the redis-go client:
```bash
make client
```

Now you can interact with the server by typing Redis-like commands.

## Example Commands
```bash
PING
PING hello
SET mykey myvalue
GET mykey
DEL mykey
SET key:1 val-1
SET key:2 val-2
KEYS key:*
HSET user:1 name "John"
HSET user:1 age 30
HGET user:1 name
```

## TODOs:
- Add support for TTL and EXPIRE command
- Fix data persistence flow
- Add logs
- Write test cases