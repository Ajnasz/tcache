// Package tcache implements a simple library, to store []byte in memory for a
// given amount of time
package tcache

import(
	"time"
)

// TCacheItem represents an item which is stored in the cace
// Name will be used to acces to a item
// Values should be the actual value
// Expire should be a time when the cache expired
type TCacheItem struct {
	Name string
	Value []byte
	Expire time.Time
}

// IsExpired returns true, if current Now() is later than the time was defined
// in Expire
func (item *TCacheItem) IsExpired() bool {
	return item.Expire.UnixNano() / int64(time.Millisecond) < time.Now().UnixNano() / int64(time.Millisecond)
}

type TCollection interface {
	Add(TCacheItem) TCacheItem
	Get(string) (TCacheItem, bool)
	GetAll() []TCacheItem
}

// TCacheCollection is a struct which stores the cache items
type TCacheCollection struct {
	HitCount int
	MissCount int
	Items map[string]TCacheItem
}

// Add adds a new item to the cache
func (c *TCacheCollection) Add(item TCacheItem) {
	c.Items[item.Name] = item
}

// Returns all items in the cache
func (c *TCacheCollection) GetAll() map[string]TCacheItem {
	return c.Items
}

// Returns one specific item from the cache
// Also increments HitCount or MissCount
func (c *TCacheCollection) Get(name string) (TCacheItem, bool) {
	var found bool = false
	var item TCacheItem

	item, found = c.Items[name]

	if (found) {
		c.HitCount++
	} else {
		c.MissCount++
	}

	return item, found
}

// Removes item from cache
func (c *TCacheCollection) Remove(name string) {
	delete(c.Items, name)
}

// Removes expired cache items from TCacheCollection
func RemoveExpired(collection *TCacheCollection) {
	for key, item := range collection.Items {
		if item.IsExpired() {
			collection.Remove(key)
		}
	}
}

// Schedules a call of RemoveExpired as frequent as it's specified in freq atribute
func ScheduleRemoveExpired(collection *TCacheCollection, freq time.Duration) {
	for {
		RemoveExpired(collection)
		time.Sleep(freq)
	}
}

// Creates a cache collections
// Calls ScheduleRemoveExpired to remove expired items in every 100
// milliseconds
func CreateCache() *TCacheCollection {
	items := make(map[string]TCacheItem)

	collection := TCacheCollection{
		Items: items,
	}

	go ScheduleRemoveExpired(&collection, 100 * time.Millisecond)

	return &collection
}
