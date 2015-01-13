package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"github.com/bradfitz/gomemcache/memcache"
)

type cacheItem struct {
	Value []byte
	Created time.Time
	Expire time.Duration
}

type requestItem struct {
	Key string
	Val cacheItem
}

func cleanCache(cache map[string]cacheItem, c chan requestItem) {
	for {
		for key, item := range cache {
			if item.Created.Add(item.Expire).Unix() < time.Now().Unix() {
				delete (cache, key)
				fmt.Println("Deleted %s", key)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func handlePostRequest(r *http.Request,cs chan requestItem) {
	body, _ := ioutil.ReadAll(r.Body)

	itemStruct := requestItem{
		Key: r.URL.Path,
		Val: cacheItem{
			Value: body,
			Created: time.Now(),
			Expire: 60 * time.Second,
		},
	}

	fmt.Println(itemStruct.Val.Value)

	cs <- itemStruct
}

var hitcount int = 0
var memhitcount int = 0

func handleGetRequest(w http.ResponseWriter, r *http.Request, cache map[string]cacheItem, mc *memcache.Client) {
	fmt.Println("handle get request", r.URL.Path)
	cachedValue, found := cache[r.URL.Path]

	if found {
		hitcount++
		w.Write([]byte(cachedValue.Value))
	} else {
		mcValue, err := mc.Get(r.URL.Path)

		if err != nil {
			http.Error(w, "Not found", 404)
		} else {
			memhitcount++
			w.Write(mcValue.Value)
			cache[r.URL.Path] = cacheItem{
				Value: mcValue.Value,
				Created: time.Now(),
				Expire: 60 * time.Second,
			}
		}
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request, cs chan requestItem, cache map[string]cacheItem, mc *memcache.Client) {
	if (r.Method == "POST") {
		handlePostRequest(r, cs)
		fmt.Fprint(w, "{\"status\":\"OK\"}")
	} else if (r.Method == "GET") {
		handleGetRequest(w, r, cache, mc)
	}
}

func listenToRequests(cs chan requestItem, cache map[string]cacheItem, mc *memcache.Client) {
	fmt.Println("LISTEN")
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, cs, cache, mc)
	})
	http.ListenAndServe(":8081", nil)
}

func addMsgToCache(cache map[string]cacheItem, cs chan requestItem) {
	for {
		select {
			case msg := <-cs:
				cache[msg.Key] = msg.Val
				fmt.Println("added message to cache", msg.Key)
		}
	}
}

func main() {
	cs := make(chan requestItem)

	cache := make(map[string]cacheItem)

	mc := memcache.New("172.17.0.11:11211")

	go cleanCache(cache, cs)

	go listenToRequests(cs, cache, mc)

	mc.Add(&memcache.Item{
		Key: "/foo/bar/baz3",
		Value: []byte("google.com"),
	})

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Ticker at", t)
			fmt.Println("Hit Count", hitcount)
			fmt.Println("Memcache Hit Count", memhitcount)
		}
	}()

	addMsgToCache(cache, cs)
}
