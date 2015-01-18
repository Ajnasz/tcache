package tcache

import (
	"testing"
	"fmt"
	"time"
	"reflect"
)

func fillCollection(collection *TCacheCollection) {
	now := time.Now()

	collection.Add(TCacheItem{
		"First",
		[]byte("Expires right now"),
		now,
	})

	collection.Add(TCacheItem{
		"Second",
		[]byte("Expires 1 second later"),
		now.Add(time.Second),
	})

	collection.Add(TCacheItem{
		"Second",
		[]byte("Expires 2 second later"),
		now.Add(2 * time.Second),
	})

	collection.Add(TCacheItem{
		"Third",
		[]byte("Expires 3 second later"),
		now.Add(3 * time.Second),
	})

	collection.Add(TCacheItem{
		"Fourth",
		[]byte("Expires 4 second later"),
		now.Add(4 * time.Second),
	})

	collection.Add(TCacheItem{
		"Fifth",
		[]byte("Expires 5 second later"),
		now.Add(5 * time.Second),
	})
}

func TestGetItem(t *testing.T) {
	collection := CreateCache()

	fillCollection(collection)

	item, found := collection.Get("First")

	if !found {
		t.Error("First item not found")
	}

	if (item.Name != "First" || string(item.Value) != "Expires right now") {
		t.Error("First item is not the elem I expected")
	}
}

func TestGetAll(t *testing.T) {
	collection := CreateCache()

	fillCollection(collection)

	if (!reflect.DeepEqual(collection.Items, collection.GetAll())) {
		t.Error("collection.Items not the same as collection.GetAll()")
	}
}

func TestAddItem(t *testing.T) {
	collection := CreateCache()

	originalLength := len(collection.GetAll())

	collection.Add(TCacheItem{
		"My Best stuff on earth",
		[]byte("VAO!!!"),
		time.Now().Add(time.Second),
	})

	if (originalLength >= len(collection.GetAll())) {
		t.Error("Length of items did not change")
	}
}

func TestTCacheItemIsExpired(t *testing.T) {
	item := TCacheItem{
		"Foo",
		[]byte("Value"),
		time.Now(),
	}

	if (item.IsExpired()) {
		t.Error("Item expired, but just created it", item)
	}

	time.Sleep(time.Second)

	if (!item.IsExpired()) {
		t.Error("Item not expired, after 1 second", item)
	}

	fmt.Println("Done")
}

func TestRemoveExpired(t *testing.T) {
	collection := CreateCache()

	fillCollection(collection)

	if _, found := collection.Get("First"); !found {
		t.Error("First item not found after RemoveExpired started")
	}

	if len(collection.GetAll()) != 5 {
		t.Error("Collection length should be 5")
	}

	time.Sleep(time.Second)

	if _, found := collection.Get("First"); found {
		t.Error("First item not removed")
	}

	if len(collection.GetAll()) != 4 {
		t.Error("Collection length should be 4")
	}

	time.Sleep(1100 * time.Millisecond)

	if _, found := collection.Get("Second"); found {
		t.Error("Second item not removed")
	}

	if len(collection.GetAll()) != 3 {
		t.Error("Collection length should be 3")
	}
}

func TestConcurrentCaches(t *testing.T) {
	collection1 := CreateCache()
	collection2 := CreateCache()

	collection1.Add(TCacheItem{
		"C1-1",
		[]byte("C1-1 Value"),
		time.Now().Add(time.Second),
	})

	collection1.Add(TCacheItem{
		"C1-2",
		[]byte("C1-2 Value"),
		time.Now().Add(time.Second),
	})

	collection2.Add(TCacheItem{
		"C2-1",
		[]byte("C2-1 Value"),
		time.Now().Add(time.Second),
	})

	if (len(collection1.GetAll()) == len(collection2.GetAll())) {
		t.Error("collection1 and collection2 has the same length")
	}

	if item, found := collection1.Get("C1-1"); !found || item.Name != "C1-1" {
		t.Error("Somehow couldn't get the right item for C1-1", item, found)
	}

	if item, found := collection1.Get("C1-2"); !found || item.Name != "C1-2" {
		t.Error("Somehow couldn't get the right item for C1-2", item, found)
	}

	if item, found := collection2.Get("C2-1"); !found || item.Name != "C2-1" {
		t.Error("Somehow couldn't get the right item for C2-1", item, found)
	}
}
