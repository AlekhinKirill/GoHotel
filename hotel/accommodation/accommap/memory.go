// Пакет accommap реализует учет занимаемых номеров в отеле на основе сохранения данных в локальной памяти компьютера (в map)
package accommap

import (
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/oops"
	"context"
	"fmt"
	"sync"
	"time"
)

// LocalStorage реализует интерфейс accommodation.Accommodation на основе сохранения данных в локальной памяти компьютера (в map)
type LocalStorage struct {
	Database    map[int]accommodation.Room
	Description accommodation.RoomsDescription
	Mu          sync.Mutex
}

// NewLocalStorage является конструктором для LocalStorage
func NewLocalStorage(database map[int]accommodation.Room, description accommodation.RoomsDescription) *LocalStorage {
	return &LocalStorage{
		Database:    database,
		Description: description,
		Mu:          sync.Mutex{},
	}
}

// Bill выставляет счет за проживание в номере при выселении гостей из отеля
func (s *LocalStorage) Bill(ctx context.Context, roomNumber int) (int, error) {
	defer time.Sleep(time.Second)
	room, exists := s.Database[roomNumber]
	if !exists {
		return 0, fmt.Errorf("localStorage.Bill error: %w", oops.ErrEmptyRoom{Number: roomNumber})
	}
	price, err := s.Description.Price(ctx, roomNumber)
	if err != nil {
		return 0, fmt.Errorf("localStorage.Bill error: %w", err)
	}
	s.Mu.Lock()
	delete(s.Database, roomNumber)
	s.Mu.Unlock()
	return price * room.StayTime, nil
}

// Place размещает новых постояльцев в структуре отеля
func (s *LocalStorage) Place(ctx context.Context, number int, tenants []string, stayTime int) (id int, err error) {
	capacity, err := s.Description.Capacity(ctx, number)
	if err != nil {
		return 0, fmt.Errorf("localStorage.Place error: %w", err)
	}
	if capacity != len(tenants) {
		return 0, fmt.Errorf("localStorage.Place error: %w", oops.ErrRoomInconsistency{Number: number, Capacity: capacity, GuestsNumber: len(tenants)})
	}
	_, exists := s.Database[number]
	if exists {
		return 0, fmt.Errorf("localStorage.Place error: %w", oops.ErrOccupiedAlready{Number: number})
	}
	s.Mu.Lock()
	s.Database[number] = accommodation.Room{Number: number, Tenants: tenants, StayTime: stayTime}
	s.Mu.Unlock()
	return number, nil
}

// Replace удаляет данные о заселении комнаты из базы данных
func (s *LocalStorage) Replace(ctx context.Context, number int) error {
	_, exists := s.Database[number]
	if !exists {
		return fmt.Errorf("LocalStorage.Replace error %w", oops.ErrEmptyRoom{Number: number})
	}
	s.Mu.Lock()
	delete(s.Database, number)
	s.Mu.Unlock()
	return nil
}

func (s *LocalStorage) Close() error {
	return s.Description.Close()
}
