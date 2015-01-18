package tcache

import(
	"time"
)

type TCacheItem struct {
	Name string
	Value []byte
	Expire time.Time
}

func (item *TCacheItem) IsExpired() bool {
	return item.Expire.UnixNano() / int64(time.Millisecond) < time.Now().UnixNano() / int64(time.Millisecond)
}

type TCollection interface {
	Add(TCacheItem) TCacheItem
	Get(string) (TCacheItem, bool)
	GetAll() []TCacheItem
}

type TCacheCollection struct {
	HitCount int
	MissCount int
	Items map[string]TCacheItem
}

func (c *TCacheCollection) Add(item TCacheItem) {
	c.Items[item.Name] = item
}

func (c *TCacheCollection) GetAll() map[string]TCacheItem {
	return c.Items
}

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

func RemoveExpired(collection *TCacheCollection) {
	for {
		for key, item := range collection.Items {
			if item.IsExpired() {
				delete(collection.Items, key)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func CreateCache() *TCacheCollection {
	items := make(map[string]TCacheItem)

	collection := TCacheCollection{
		Items: items,
	}

	go RemoveExpired(&collection)

	return &collection
}
