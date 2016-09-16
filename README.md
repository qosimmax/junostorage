
junostorage is th simple implementation of Redis-like in-memory cache.




## Features

Redis keys commands

- `DEL` This command deletes the key, if exists
- `EXPIRE` Expires the key after the specified time

Redis strings commands

- `SET` Sets the value at the specified key.
- `GET` Get the value of a key.

Redis lists commands

- `LPUSH`  Prepend one or multiple values to a list
- `LLEN`   Get the length of a list
- `LINDEX` Get an element from a list by its index
- `LPOP`   Remove and get the first element in a list

Redis dict commands

- `HSET`    Set the string value of a hash field
- `HGET`    Get the value of a hash field stored at specified key
- `HGETALL` Get all the fields and values stored in a hash at specified key
- `HDEL`    Delete one or more hash fields


## Getting Started


### Building junostorage

junostorage can be compiled and used on Linux, OSX, Windows

To build project:
```
$ make
```

To test:
```
$ make test
```

## Running
For command line invoke:
```
$ ./juno-server -h
```

To run a server:

```
$ ./juno-server

```


#### Telnet
There is the option to use a plain telnet connection. The default output through telnet is [RESP](http://redis.io/topics/protocol).

```
telnet localhost 6380
ping
+PONG
set mykey hello
+OK
get mykey
$5
hello


```
