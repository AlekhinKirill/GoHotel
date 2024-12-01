package sqlitaccom

import (
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/oops"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	//"Go_projects/databases"
)

type Storage struct {
	database *sql.DB
	rooms    accommodation.RoomsDescription
}

func NewStorage(path string, rooms accommodation.RoomsDescription) *Storage {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Errorf("NewStorage error: %w", err))
	}
	db.Close()
	return &Storage{
		database: db,
		rooms:    rooms,
	}
}

type Request struct {
	Id       int
	Room     int
	Tenants  string
	StayTime int
}

func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "select * from Accommodation")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Id, &req.Room, &req.Tenants, &req.StayTime)
		if err != nil {
			return err
		}
		fmt.Printf("%d. Комната №%d: %s. Время проживания: %d ночей", req.Id, req.Room, req.Tenants, req.StayTime)
	}
	return nil
}

func (s Storage) Place(ctx context.Context, number int, tenants []string, stayTime int) (int, error) {
	var (
		t string
	)
	capacity, err := s.rooms.Capacity(ctx, number)
	if err != nil {
		return -1, fmt.Errorf("localStorage.Place error: %w", err)
	}
	if capacity != len(tenants) {
		return -1, fmt.Errorf("localStorage.Place error: %w", oops.ErrRoomInconsistency{Number: number, Capacity: capacity, GuestsNumber: len(tenants)})
	}
	rows, err := s.database.QueryContext(ctx, "select Room from Accommodation where Room = ?", number)
	if err != nil {
		return -1, err
	}
	if rows != nil {
		return -1, fmt.Errorf("Storage.Place error: %w", oops.ErrOccupiedAlready{Number: number})
	}
	for _, person := range tenants {
		t += person + ", "
	}
	t = strings.TrimRight(t, ", ")
	result, err := s.database.Exec("insert into Accommodation (Room, Tenants, Duration) values (?, ?, ?)", number, t, stayTime)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (s Storage) Bill(ctx context.Context, roomNumber int) (int, error) {
	row, err := s.database.QueryContext(ctx, "select Duration from Menu where Room = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, fmt.Errorf("Storage.Bill error: %w", oops.ErrEmptyRoom{Number: roomNumber})
	}
	var stayTime int
	err = row.Scan(&stayTime)
	if err != nil {
		return 0, err
	}
	price, err := s.rooms.Price(ctx, roomNumber)
	if err != nil {
		return 0, fmt.Errorf("Storage.Bill error: %w", err)
	}
	_, err = s.database.Exec("delete from Accommodation where Room = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	return price * stayTime, nil
}
