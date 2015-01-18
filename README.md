Experimental library in go

```go

cache := tcache.CreateCache()

cache.Add("Foobar", tcache.CacheItem{
	Name: "Foobar",
	Value: []byte("BazQux"),
	Expire: time.Now().Add(5 * time.Second),
})

cacheItem, found := cache.Get("Foobar")
```
