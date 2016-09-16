
junostorage is th simple implementation of Redis-like in-memory cache.




## Features

Redis keys commands

- `DEL` this command deletes the key, if exists
- `EXPIRE` expires the key after the specified time

Redis strings commands

- `SET` sets the value at the specified key.
- `GET` get the value of a key.

Redis lists commands

- `LPUSH`  prepend one or multiple values to a list
- `LLEN`   get the length of a list
- `LINDEX` get an element from a list by its index
- `LPOP`   remove and get the first element in a list

Redis dict commands

- `HSET`    set the string value of a hash field
- `HGET`    get the value of a hash field stored at specified key
- `HGETALL` get all the fields and values stored in a hash at specified key
- `HDEL`    delete one or more hash fields


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
