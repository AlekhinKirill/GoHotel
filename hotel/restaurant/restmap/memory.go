// Пакет restmap реализует учет заказов, сделанных в ресторане на основе сохранения информации о них в локальной памяти
// компьютера (в map)
package restmap

import (
	"Go_projects/hotel/restaurant"
	"context"
	"fmt"
	"sync"
	"time"
)

// LocalStorage реализует интерфейс restaurant.Restaurant на основе сохранения информации о заказах в локальной памяти компьютера (в map)
type LocalStorage struct {
	menu     restaurant.Menu
	mu       sync.Mutex
	database map[int][]restaurant.Dinner
}

// NewLocalStorage является конструктором для LocalStorage
func NewLocalStorage(menu restaurant.Menu, database map[int][]restaurant.Dinner) *LocalStorage {
	return &LocalStorage{
		menu:     menu,
		database: database,
		mu:       sync.Mutex{},
	}
}

// Bill выставляет счет от ресторана при выселении постояльцем из отеля с учетом всех сделанных ими заказов и посещенных завтраков
func (s *LocalStorage) Bill(ctx context.Context, roomNumber int) (int, error) {
	defer time.Sleep(time.Second)
	var sum int
	for _, dinner := range s.database[roomNumber] {
		sum += dinner.Price
	}
	s.mu.Lock()
	delete(s.database, roomNumber)
	s.mu.Unlock()
	return sum, nil
}

// PlaceOrder размещает заказ в базе данных ресторана
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
	s.database[roomNumber] = append(s.database[roomNumber], restaurant.Dinner{Dishes: dishes, Price: sum})
	s.mu.Unlock()
	return roomNumber, nil
}

// PlaceBreakfast размещает информацию о завтраках в базе данных отеля
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
	s.database[roomNumber] = append(s.database[roomNumber], restaurant.Dinner{Dishes: breakfasts, Price: count * price})
	s.mu.Unlock()
	return roomNumber, nil
}

func (s *LocalStorage) ShowMenu(ctx context.Context) error {
	return s.menu.Show(ctx)
}

func (s *LocalStorage) Close() error {
	err := s.menu.Close()
	if err != nil {
		return fmt.Errorf("restmap.LocalStorage.Close error: %w", err)
	}
	return s.menu.Close()
}
