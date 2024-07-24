package cache

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache реализует кэш с ограниченной емкостью и алгоритмом LRU.
type LRUCache struct {
	capacity int
	cache    map[interface{}]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

// entry представляет элемент кэша.
type entry struct {
	key        interface{}
	value      interface{}
	expiration time.Time // Время истечения срока действия элемента
}

// NewLRUCache создаёт новый LRU-кэш с заданной емкостью.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[interface{}]*list.Element),
		list:     list.New(),
	}
}

// Cap возвращает емкость кэша.
func (c *LRUCache) Cap() int {
	return c.capacity
}

// Len возвращает текущий размер кэша.
func (c *LRUCache) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.list.Len()
}

// Clear очищает весь кэш.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[interface{}]*list.Element) // Очищаем хранилище кэша
	c.list.Init()                                 // Инициализируем список заново
}

// Add добавляет пару ключ-значение в кэш.
func (c *LRUCache) Add(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		e := elem.Value.(*entry)
		c.list.MoveToFront(elem)   // Перемещаем элемент в начало списка
		e.value = value            // Обновляем значение
		e.expiration = time.Time{} // Сбрасываем время истечения
		return
	}

	if c.list.Len() >= c.capacity {
		c.deleteLastUsedElement() // Вытесняем наименее используемый элемент
	}
	e := &entry{key: key, value: value, expiration: time.Time{}}
	elem := c.list.PushFront(e) // Добавляем новый элемент в начало списка
	c.cache[key] = elem         // Добавляем элемент в хранилище кэша
}

// AddWithTTL добавляет пару ключ-значение в кэш с TTL.
func (c *LRUCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		e := elem.Value.(*entry)
		c.list.MoveToFront(elem)
		e.value = value
		e.expiration = time.Now().Add(ttl) // Устанавливаем время истечения
		return
	}

	if c.list.Len() >= c.capacity {
		c.deleteLastUsedElement()
	}
	e := &entry{key: key, value: value, expiration: time.Now().Add(ttl)}
	elem := c.list.PushFront(e) // Добавляем новый элемент в начало списка с TTL
	c.cache[key] = elem         // Добавляем элемент в хранилище кэша
}

// Get получает значение по ключу.
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.cache[key]
	if !ok {
		return nil, false // Если ключ не найден
	}

	e := elem.Value.(*entry)
	if !e.expiration.IsZero() && time.Now().After(e.expiration) {
		c.list.Remove(elem)  // Удаляем элемент из списка
		delete(c.cache, key) // Удаляем элемент из хранилища кэша
		return nil, false
	}

	c.list.MoveToFront(elem) // Перемещаем элемент в начало списка
	return e.value, true     // Возвращаем значение и статус
}

// Remove удаляет пару ключ-значение из кэша.
func (c *LRUCache) Remove(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.cache[key]
	if !ok {
		return
	}
	c.list.Remove(elem)
	delete(c.cache, key)
}

// deleteLastUsedElement вытесняет наименее используемый элемент из кэша.
func (c *LRUCache) deleteLastUsedElement() {
	elem := c.list.Back()
	if elem == nil {
		return
	}
	c.list.Remove(elem)
	delete(c.cache, elem.Value.(*entry).key)
}
