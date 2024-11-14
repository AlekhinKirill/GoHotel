package rooms

import (
	"Go_projects/hotel/oops"
	"context"
	"fmt"
)

// Pair - вспомогательный класс для объединения в одном объекте двух полей с разными типами
// нам он будет полезен при работе с типами номеров, поскольку они отличаются как по вместимости, так и по уровню комфорта
type Pair struct {
	Capacity int
	Class    string
}

type LocalStorage struct {
	roomTypes map[int]Pair
	prices    map[Pair]int
}

func NewLocalStorage(roomTypes map[int]Pair, prices map[Pair]int) *LocalStorage {
	return &LocalStorage{
		roomTypes: roomTypes,
		prices:    prices,
	}
}

func (s *LocalStorage) Price(ctx context.Context, roomNumber int) (int, error) {
	pair, roomExists := s.roomTypes[roomNumber]
	if !roomExists {
		return 0, fmt.Errorf("LocalStorage.Price error %w", oops.ErrNoRoom{Number: roomNumber})
	}
	price, priceExists := s.prices[pair]
	if !priceExists {
		return 0, fmt.Errorf("LocalStorage.Price error %w", oops.ErrNoPrice{Capacity: pair.Capacity, Class: pair.Class})
	}
	return price, nil
}

func (s *LocalStorage) Capacity(ctx context.Context, roomNumber int) (int, error) {
	pair, exists := s.roomTypes[roomNumber]
	if !exists {
		return 0, fmt.Errorf("LocalStorage.Capacity error %w", oops.ErrNoRoom{Number: roomNumber})
	}
	return pair.Capacity, nil
}

func (s *LocalStorage) Type(ctx context.Context, roomNumber int) (string, error) {
	pair, exists := s.roomTypes[roomNumber]
	if !exists {
		return "", fmt.Errorf("LocalStorage.Type error %w", oops.ErrNoRoom{Number: roomNumber})
	}
	return pair.Class, nil
}
