package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// структура для тестового сравнения составных типов
type testStruct struct {
	A int
	B string
}

// Тест на проверку корректного вытеснения элементов при превышении емкости кэша
func TestLRUCacheCapacity(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	// Проверяем, что добавленные элементы доступны
	value, ok := cache.Get("key1")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key1'")
	assert.Equal(t, "value1", value)

	value, ok = cache.Get("key2")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key2'")
	assert.Equal(t, "value2", value)

	cache.Add("key3", "value3") // Добавляем новый элемент, который вытеснит один из существующих

	// Проверяем, что ключ 'key1' был вытеснен, а ключи 'key2' и 'key3' доступны
	_, ok = cache.Get("key1")
	assert.False(t, ok, "Ключ 'key1' не был вытеснен")

	value, ok = cache.Get("key2")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key2'")
	assert.Equal(t, "value2", value)

	value, ok = cache.Get("key3")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key3'")
	assert.Equal(t, "value3", value)
}

// Тест на потокобезопасность
// Создаем несколько горутин, которые будут одновременно добавлять, удалять и получать элементы из кэша
func TestLRUCacheThreadSafety(t *testing.T) {
	cache := NewLRUCache(10)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := i % 5 // используем ключи от 0 до 4
			value := i

			cache.Add(key, value)

			v, ok := cache.Get(key)
			assert.True(t, ok, "Проблема с потокобезопасностью при добавлении ключа %d", key)
			assert.Equal(t, value, v, "Проблема с потокобезопасностью при добавлении ключа %d", key)

			if i%2 == 0 {
				cache.Remove(key)

				_, ok := cache.Get(key)
				assert.False(t, ok, "Проблема с потокобезопасностью при удалении ключа %d", key)
			}
		}(i)
	}
	wg.Wait()
}

// Тест на работу кэша с разными типами значений
// Проверяем, что кэш корректно работает с int, string, slice и struct
func TestLRUCacheAnyValue(t *testing.T) {
	cache := NewLRUCache(4) // Создаем кэш с емкостью 4

	// Добавляем значения разных типов
	cache.Add("key1", 123)
	cache.Add("key2", "string value")
	cache.Add("key3", []int{1, 2, 3})
	cache.Add("key4", testStruct{A: 1, B: "struct value"})

	// Проверяем, что ключи 'key1', 'key2', 'key3', 'key4' доступны
	value, ok := cache.Get("key1")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key1'")
	assert.Equal(t, 123, value)
	assert.IsType(t, 123, value)

	value, ok = cache.Get("key2")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key2'")
	assert.Equal(t, "string value", value)
	assert.IsType(t, "", value)

	value, ok = cache.Get("key3")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key3'")
	assert.Equal(t, []int{1, 2, 3}, value)
	assert.IsType(t, []int{}, value)

	value, ok = cache.Get("key4")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key4'")
	assert.Equal(t, testStruct{A: 1, B: "struct value"}, value)
	assert.IsType(t, testStruct{}, value)

	// Добавляем значение с типом времени, чтобы вытеснить наименее используемый элемент
	cache.Add("key5", time.Now())

	// Проверяем, что ключ 'key1' был вытеснен
	_, ok = cache.Get("key1")
	assert.False(t, ok, "Ключ 'key1' не был вытеснен")

	// Проверяем, что ключи 'key2', 'key3', 'key4', и 'key5' доступны
	value, ok = cache.Get("key2")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key2'")
	assert.Equal(t, "string value", value)

	value, ok = cache.Get("key3")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key3'")
	assert.Equal(t, []int{1, 2, 3}, value)

	value, ok = cache.Get("key4")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key4'")
	assert.Equal(t, testStruct{A: 1, B: "struct value"}, value)

	value, ok = cache.Get("key5")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key5'")
	assert.IsType(t, time.Time{}, value)
}

// Тест на добавление элементов с ограниченным сроком действия (TTL)
// Проверяем, корректность работы удаления, после истечения срока действия
func TestLRUCacheAddWithTTL(t *testing.T) {
	cache := NewLRUCache(2)

	cache.AddWithTTL("key1", "value1", 2*time.Second)

	// Проверяем, что ключ доступен до истечения TTL
	value, ok := cache.Get("key1")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key1' до истечения TTL")
	assert.Equal(t, "value1", value)
	assert.IsType(t, "", value)

	// Ждем 3 секунды для истечения TTL
	time.Sleep(3 * time.Second)

	// Проверяем, что ключ был удален после истечения TTL
	_, ok = cache.Get("key1")
	assert.False(t, ok, "Ключ 'key1' не был удален после истечения TTL")
}

// Тест на очистку кэша
// Проверяем, что все элементы удаляются при вызове метода Clear
func TestLRUCacheClear(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Clear() // Очищаем кэш

	// Проверяем, что все ключи были удалены
	_, ok := cache.Get("key1")
	assert.False(t, ok, "Ключ 'key1' не был удален после очистки")

	_, ok = cache.Get("key2")
	assert.False(t, ok, "Ключ 'key2' не был удален после очистки")

	_, ok = cache.Get("key3")
	assert.False(t, ok, "Ключ 'key3' не был удален после очистки")
}

// Тест на удаление элемента из кэша
// Проверяем, что элемент корректно удаляется и недоступен после удаления
func TestLRUCacheRemove(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	cache.Remove("key2") // Удаляем ключ 'key2'

	// Проверяем, что ключ 'key2' был удален
	_, ok := cache.Get("key2")
	assert.False(t, ok, "Ключ 'key2' не был удален")

	// Проверяем, что другие ключи все еще доступны
	value, ok := cache.Get("key1")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key1'")
	assert.Equal(t, "value1", value)

	value, ok = cache.Get("key3")
	assert.True(t, ok, "Не удалось получить значение по ключу 'key3'")
	assert.Equal(t, "value3", value)
}
