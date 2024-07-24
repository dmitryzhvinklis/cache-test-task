# Задание: "Написать сервис (класс/структура) кэширования".

# Основные условия:
- Кэш ограниченной емкости, метод вытеснения ключей LRU;
- Сервис должен быть потокобезопасный;
- Сервис должен принимать любые значения;
- Реализовать unit-тесты.

# Сервис должен реализовывать следующий интерфейс:
```
    type ICache interface {
        Cap() int
        Len() int
        Clear() // удаляет все ключи
        Add(key, value interface{})
        AddWithTTL(key, value interface{}, ttl time.Duration) // добавляет ключ со сроком жизни ttl
        Get(key interface{}) (value interface{}, ok bool)
        Remove(key interface{})
    }
```