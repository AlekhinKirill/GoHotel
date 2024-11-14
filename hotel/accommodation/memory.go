package accommodation

import (
	"Go_projects/hotel/oops"
	"context"
	"fmt"
	"sync"
	"time"
)

type LocalStorage struct {
	Database    map[int]Room
	Description RoomsDescription
	Mu          *sync.Mutex
}

func NewLocalStorage(database map[int]Room, description RoomsDescription) *LocalStorage {
	return &LocalStorage{
		Database:    database,
		Description: description,
		Mu:          &sync.Mutex{},
	}
}

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
	return price * room.stayTime, nil
}

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
	s.Database[number] = Room{number, tenants, stayTime}
	s.Mu.Unlock()
	return number, nil
}

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
