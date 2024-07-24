package models

import "time"

// ICache определяет интерфейс для кэша.
type ICache interface {
	Cap() int
	Len() int
	Clear() // удаляет все ключи
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration) // добавляет ключ со сроком жизни ttl
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}
