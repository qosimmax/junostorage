
**junostorage** is th simple implementation of Redis-like in-memory cache.




## Features

Redis keys commands

- `DEL` this command deletes the key, if exists
- `EXPIRE` expires the key after the specified time
- `KEYS` Find all keys matching the specified pattern

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

**junostorage** can be compiled and used on Linux, OSX, Windows

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
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.

SET storage redis
+OK

GET storage
$5
redis

ping        
+PONG

LPUSH list 0 1 2 3 4 5
:6

LINDEX list
-ERR wrong number of arguments for 'lindex' command

LINDEX list 0
$1
5

HSET person name nemo
:1

HSET person age 25
:1

quit
+OK
Connection closed by foreign host.


```



junostorage Client API
=============


junostorage Client API is a [Go](http://golang.org/)

## Examples

#### Connection
```go
package main

import (
	"log"
	"time"

	"github.com/junostorage/client"
)

func main() {
	con, err := client.DialTimeout("localhost:6380", 60*time.Duration(time.Second))
	if err != nil {
		log.Fatalf("Dial error:%v", err)
	}
	val, err := con.Do("SET", "storage", "redis")
	if err != nil {
		log.Fatalf("On SET error:%v", err)
	}

	if val.Error() != nil {
		log.Fatalf("resp error:%v", val.Error())
	}
	log.Println(val)

	val, err = con.Do("GET", "storage")
	if err != nil {
		log.Fatalf("On GET error:%v", err)
	}
	log.Println(val)

}

```
