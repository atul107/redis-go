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
EXPIRE mykey 10 XX
EXPIRE mykey 10 NX 
TTL mykey
SET key:1 val-1
SET key:2 val-2
KEYS key:*
HSET user:1 name John
HSET user:1 age 30
HGET user:1 name
ZADD myzset 1 one
ZADD myzset 1 uno
ZADD z 1 a 2 b -1 c
ZRANGE z 0 -1
```

### Running tests
1. Start server by running `make server` command inside the project root directory
2. Run test cases by using `make test` command in another terminal in the project root directory

## TODOs:
- Fix data persistence flow
- Add logs
