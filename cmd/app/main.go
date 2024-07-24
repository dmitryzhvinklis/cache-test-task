package main

import (
	"cache-test/internal/cache"
	"cache-test/internal/models"
	"log"
	"time"
)

func main() {
	lruCache := cache.NewLRUCache(3) // Создаем новый LRUCache с емкостью 3
	testCache(lruCache)
}

func testCache(lruCache models.ICache) {
	log.Println("Добавляем ключи a, b, c")
	lruCache.Add("a", 1)
	lruCache.Add("b", 2)
	lruCache.Add("c", 3)
	printCacheState(lruCache)

	log.Println("Получаем значение для ключа 'a'")
	if value, ok := lruCache.Get("a"); ok {
		log.Printf("Ключ: a, Значение: %v\n", value)
	} else {
		log.Println("Ключ 'a' не найден")
	}

	log.Println("Добавляем ключ 'd', что приводит к вытеснению одного из ключей")
	lruCache.Add("d", 4)
	printCacheState(lruCache)

	log.Println("Проверяем наличие ключа 'b'")
	if _, ok := lruCache.Get("b"); !ok {
		log.Println("Ключ 'b' был вытеснен")
	} else {
		log.Println("Ключ 'b' все еще в кэше")
	}

	log.Println("Добавляем ключ 'e' с TTL 2 секунды")
	lruCache.AddWithTTL("e", 5, 2*time.Second)
	printCacheState(lruCache)

	log.Println("Ждем 3 секунды...")
	time.Sleep(3 * time.Second)
	log.Println("Проверяем наличие ключа 'e' после истечения TTL")
	if _, ok := lruCache.Get("e"); !ok {
		log.Println("Ключ 'e' был удален по истечении TTL")
	} else {
		log.Println("Ключ 'e' все еще в кэше")
	}

	printCacheState(lruCache)

	log.Println("Очищаем кэш")
	lruCache.Clear()
	printCacheState(lruCache)
}

// текущее состояние кэша
func printCacheState(lruCache models.ICache) {
	log.Printf("Текущее состояние кэша:\n")
	log.Printf("Емкость: %d\n", lruCache.Cap())
	log.Printf("Количество ключей: %d\n", lruCache.Len())
}
