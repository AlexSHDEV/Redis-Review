package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func cache[T any](db *redis.Client, key string, value T) error {

	expiration := 10 * time.Minute // Время жизни кэша

	err := db.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Println("Ошибка при установке значения:", err)
		return err
	}
	fmt.Println("Данные успешно кэшированы")
	return nil
}

func getValue(db *redis.Client, key string) (string, error) {
	cachedValue, err := db.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("Данные не найдены в кэше")
		} else {
			fmt.Println("Ошибка при получении значения:", err)
		}
	}

	return cachedValue, nil
}

type User struct {
	ID        int
	Name      string
	Email     string
	Lastvisit time.Time
	Upgrade   bool
}

func main() {
	// Создаем клиент Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	fmt.Println("Подключение к Redis успешно")

	key := "user:3"
	value := User{
		ID:        1,
		Name:      "Max",
		Email:     "Max@inbox.ru",
		Lastvisit: time.Now(),
		Upgrade:   true}

	userJSON, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Ошибка при сериализации структуры:", err)
		return
	}

	fmt.Println(userJSON)
	_ = cache(rdb, key, userJSON)

	// Получаем данные из кэша
	cachedValue, _ := getValue(rdb, key)

	var cachedUser User
	err = json.Unmarshal([]byte(cachedValue), &cachedUser)
	if err != nil {
		fmt.Println("Ошибка при десериализации структуры:", err)
		return
	}

	fmt.Println("Данные из кэша:", cachedUser, "\n", cachedValue)
}
