package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func newPerson(name string, age int) Person {
	return Person{
		ID:   uuid.NewString(),
		Name: name,
		Age:  age,
	}
}

func main() {
	fmt.Println("Go Redis tutorial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	person := newPerson("Jay", 30)
	jsonString, _ := json.Marshal(person)

	err = client.Set(ctx, person.ID, jsonString, 0).Err()
	if err != nil {
		fmt.Errorf("error setting key %s", err)
	}

	name, err := client.Get(ctx, person.ID).Result()
	if err != nil {
		fmt.Errorf("error getting key %s", err)
	}
	fmt.Println("value retrieved:", name)
}
