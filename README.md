Experimental library in go

```go
package main

import (
	"fmt"
	"time"
	"github.com/ajnasz/tcache"
)

func main() {
	var cache *tcache.TCacheCollection
	cache = tcache.CreateCache()

	cache.Add(tcache.TCacheItem{
		Name: "Foobar",
		Value: []byte("BazQux"),
		Expire: time.Now().Add(5 * time.Second),
	})

	cacheItem, found := cache.Get("Foobar")

	if found {
		fmt.Println(string(cacheItem.Value))
	} else {
		fmt.Println("Not found in cache")
	}
}
```
