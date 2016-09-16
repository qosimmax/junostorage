
junostorage is th simple implementation of Redis-like in-memory cache.




## Features

- GET, SET, LPUSH methods

## Getting Started


### Building junostorage

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

## <a name="cli"></a>Telnet sample with juno-server

Basic operations:
```
$ ./juno-server

# .
> SET storage redis   # set storage key
> GET storage         # get value from storage key
```


#### Telnet
There is the option to use a plain telnet connection. The default output through telnet is [RESP](http://redis.io/topics/protocol).

```
telnet localhost 6380
set mykey hello
+OK

```
