// Пакет menumap реализует меню на основе сохраниения информации в локальной памяти компьютера (в map)
package menumap

import (
	"Go_projects/hotel/oops"
	"context"
	"fmt"
)

// Breakfast хранит информацию о том, доступна ли опция завтрака в приотельном ресторане и какова ее цена
type Breakfast struct {
	Provided bool
	Price    int
}

// Storage реализует интерфейс restaurnat.Menu на основе сохраниения информации в локальной памяти компьютера (в map)
type Storage struct {
	table     map[string]int
	breakfast Breakfast
}

// NewStorage является конструктором Storage
func NewStorage(table map[string]int, breakfast Breakfast) *Storage {
	return &Storage{
		table:     table,
		breakfast: breakfast,
	}
}

// Price возвращает стоимость конкретного блюда в меню
func (s *Storage) Price(ctx context.Context, dish string) (int, error) {
	value, exists := s.table[dish]
	if exists {
		return value, nil
	}
	return 0, oops.ErrOutOfMenu{Dish: dish}
}

// Breakfast возвращает стоимость завтрака на одного человека
func (s *Storage) Breakfast(ctx context.Context) (int, error) {
	if s.breakfast.Provided {
		return s.breakfast.Price, nil
	}
	return 0, fmt.Errorf("завтрак в отеле не предусмотрен")
}

// Show выводит меню на экран
func (s *Storage) Show(ctx context.Context) {
	for dish, price := range s.table {
		fmt.Printf("%s : %d рублей\n", dish, price)
	}
	if s.breakfast.Provided {
		fmt.Printf("Завтрак : %d рублей", s.breakfast.Price)
	}
}
