package menu

import (
	"Go_projects/hotel/oops"
	"context"
	"fmt"
)

type Breakfast struct {
	Provided bool
	Price    int
}

type Storage struct {
	table     map[string]int
	breakfast Breakfast
}

func NewStorage(table map[string]int, breakfast Breakfast) *Storage {
	return &Storage{
		table:     table,
		breakfast: breakfast,
	}
}

func (s *Storage) Price(ctx context.Context, dish string) (int, error) {
	value, exists := s.table[dish]
	if exists {
		return value, nil
	}
	return 0, oops.ErrOutOfMenu{Dish: dish}
}

func (s *Storage) Breakfast(ctx context.Context) (int, error) {
	if s.breakfast.Provided {
		return s.breakfast.Price, nil
	}
	return 0, fmt.Errorf("завтрак в отеле не предусмотрен")
}
