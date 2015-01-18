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

type TCacheCCollection struct {
	Items map[string]TCacheItem
}

func (c *TCacheCCollection) Add(item TCacheItem) {
	c.Items[item.Name] = item
}

func (c *TCacheCCollection) GetAll() map[string]TCacheItem {
	return c.Items
}

func (c *TCacheCCollection) Get(name string) (TCacheItem, bool) {
	var found bool = false
	var item TCacheItem

	item, found = c.Items[name]

	return item, found
}

func RemoveExpired(collection *TCacheCCollection) {
	for {
		for key, item := range collection.Items {
			if item.IsExpired() {
				delete(collection.Items, key)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func CreateCache() *TCacheCCollection {
	items := make(map[string]TCacheItem)

	collection := TCacheCCollection{
		Items: items,
	}

	go RemoveExpired(&collection)

	return &collection
}
