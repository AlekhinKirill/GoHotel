package restaurant

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type LocalStorage struct {
	menu     Menu
	mu       sync.Mutex
	database map[int][]Dinner
}

func NewLocalStorage(menu Menu, database map[int][]Dinner) *LocalStorage {
	return &LocalStorage{
		menu:     menu,
		database: database,
	}
}

func (s *LocalStorage) Bill(ctx context.Context, roomNumber int) (int, error) {
	defer time.Sleep(time.Second)
	var sum int
	for _, dinner := range s.database[roomNumber] {
		sum += dinner.price
	}
	s.mu.Lock()
	delete(s.database, roomNumber)
	s.mu.Unlock()
	return sum, nil
}

func (s *LocalStorage) PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (id int, err error) {
	defer time.Sleep(time.Second)
	var sum int
	for _, dish := range dishes {
		price, err := s.menu.Price(ctx, dish)
		if err != nil {
			return 0, fmt.Errorf("localStorage.PlaceOrder error: %w", err)
		}
		sum += price
	}
	s.mu.Lock()
	s.database[roomNumber] = append(s.database[roomNumber], Dinner{dishes, sum})
	s.mu.Unlock()
	return roomNumber, nil
}

func (s *LocalStorage) PlaceBreakfast(ctx context.Context, roomNumber int, count int) (id int, err error) {
	price, err := s.menu.Breakfast(ctx)
	if err != nil {
		return 0, fmt.Errorf("LocalStorge.PlaceBreakfast error: %w", err)
	}
	breakfasts := make([]string, count)
	for i := 0; i < count; i++ {
		breakfasts[i] = "Завтрак"
	}
	s.mu.Lock()
	s.database[roomNumber] = append(s.database[roomNumber], Dinner{breakfasts, count * price})
	s.mu.Unlock()
	return roomNumber, nil
}
