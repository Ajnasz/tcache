Experimental library in go

```go

cs := make(chan byte, 1000)
cache := tcache.CreateCache(cs)

cache.Add("Foobar", tcache.CacheItem{
	Value: []byte("BazQux"),
	Created: time.Now(),
	Expire: 5 * time.Second,
})
