junostorage Client
=============


junostorage Client is a [Go](http://golang.org/)

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
		log.Fatalf("On cmd send error:%v", err)
	}

	if val.Error() != nil {
		log.Fatalf("resp error:%v", val.Error())
	}
	log.Println(val)
}

```
