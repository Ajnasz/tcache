package tcache

import (
	"fmt"
	"time"
)

const MSG_HIT byte = 1
const MSG_MISS byte = 0
const MSG_DELETED byte = 2

type CacheItem struct {
	Value []byte
	Created time.Time
	Expire time.Duration
}

type CacheItemCollection struct {
	channel chan byte
	HitCount int
	MissCount int
	Items map[string]CacheItem
}

func (collection *CacheItemCollection) Add(key string, item CacheItem) {
	collection.Items[key] = item
}

func (collection *CacheItemCollection) Remove(key string) {
	delete(collection.Items, key)
}

func (collection *CacheItemCollection) Get(key string) (CacheItem, bool) {
	val, ok := collection.Items[key]

	if (ok) {
		collection.channel <- MSG_HIT
		collection.HitCount++
	} else {
		collection.channel <- MSG_MISS
		collection.MissCount++
	}

	return val, ok
}

func cleanCache(collection CacheItemCollection, cs chan byte) {
	for {
		for key, item := range collection.Items {
			if item.Created.Add(item.Expire).Unix() < time.Now().Unix() {
				collection.Remove(key)
				cs <- MSG_DELETED
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func CreateCache(cs chan byte) CacheItemCollection{

	var c = CacheItemCollection{
		channel: cs,
		Items: make(map[string]CacheItem),
	}

	fmt.Println(c.Items)

	go cleanCache(c, cs)

	go func () {
		for {
			select {
				case <-cs:
			}
		}
	}()

	return c
}
