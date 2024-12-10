// Пакет roomsmap предоставляет описание аппартаментов отеля, основанное на сохранении данных в локальной памяти компьютера (в map)
package roomsmap

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

// LocalStorage реализует интерфейс accommodation.RoomsDescription на основе сохранения данных в локальной памяти компьютера (в map)
type LocalStorage struct {
	roomTypes map[int]Pair
	prices    map[Pair]int
}

// NewLocalStorage является конструктором для LocalStorage
func NewLocalStorage(roomTypes map[int]Pair, prices map[Pair]int) *LocalStorage {
	return &LocalStorage{
		roomTypes: roomTypes,
		prices:    prices,
	}
}

// Price возвращает стоимость номера за одну ночь
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

// Capacity возвращает вместимость номера, то есть число человек, на которое он расчитан
func (s *LocalStorage) Capacity(ctx context.Context, roomNumber int) (int, error) {
	pair, exists := s.roomTypes[roomNumber]
	if !exists {
		return 0, fmt.Errorf("LocalStorage.Capacity error %w", oops.ErrNoRoom{Number: roomNumber})
	}
	return pair.Capacity, nil
}

// Type возвращает уровень комфортности номера
func (s *LocalStorage) Type(ctx context.Context, roomNumber int) (string, error) {
	pair, exists := s.roomTypes[roomNumber]
	if !exists {
		return "", fmt.Errorf("LocalStorage.Type error %w", oops.ErrNoRoom{Number: roomNumber})
	}
	return pair.Class, nil
}

// Show выводит информацию о номерах в отеле
func (s *LocalStorage) Show(ctx context.Context) error {
	for number, roomtype := range s.roomTypes {
		price, exists := s.prices[roomtype]
		if !exists {
			return oops.ErrNoPrice{Capacity: roomtype.Capacity, Class: roomtype.Class}
		}
		fmt.Printf("Номер %d : %d человека, класс %s, %d рублей\n", number, roomtype.Capacity, roomtype.Class, price)
	}
	return nil
}

// Close необходим для реализации интерфейса accommodation.RoomsDescription
// в случае использования map этот метод фиктивный
func (s *LocalStorage) Close() error {
	return nil
}
