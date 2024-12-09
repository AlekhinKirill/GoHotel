// Пакет sqlitaccom реализует учет занимаемых номеров в отеле посредством работы с базой данных
package sqlitaccom

import (
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/oops"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Storage реализует интерфейс accommodation.Accommodation на основе база данных
type Storage struct {
	database *sql.DB
	rooms    accommodation.RoomsDescription
}

// NewStorage является конструктором для Storage
func NewStorage(path string, rooms accommodation.RoomsDescription) *Storage {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Errorf("NewStorage error: %w", err))
	}
	return &Storage{
		database: db,
		rooms:    rooms,
	}
}

// Request является классом для хранения информации из сторок базы данных
type Request struct {
	Room     int
	Tenants  string
	StayTime int
}

// Show выводит информацию о всех заселенных номерах отеля и их постояльцах
func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "select * from Accommodation")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Room, &req.Tenants, &req.StayTime)
		if err != nil {
			return err
		}
		fmt.Printf("Комната №%d: %s. Время проживания: %d ночи\n", req.Room, req.Tenants, req.StayTime)
	}
	return nil
}

// Place размещает новых постояльцев в структуре отеля
func (s Storage) Place(ctx context.Context, number int, tenants []string, stayTime int) (int, error) {
	var (
		t string
	)
	capacity, err := s.rooms.Capacity(ctx, number)
	if err != nil {
		return -1, fmt.Errorf("sqlitaccom.Storage.Place error: %w", err)
	}
	if capacity != len(tenants) {
		return -1, fmt.Errorf("sqlitaccom.Storage.Place error: %w", oops.ErrRoomInconsistency{Number: number, Capacity: capacity, GuestsNumber: len(tenants)})
	}
	row := s.database.QueryRowContext(ctx, "SELECT Room FROM Accommodation WHERE Room = $1", number)
	var num int
	err = row.Scan(&num)
	if !errors.Is(err, sql.ErrNoRows) {
		return -1, fmt.Errorf("sqlitaccom.Storage.Place error: %w", oops.ErrOccupiedAlready{Number: number})
	}
	for _, person := range tenants {
		t += person + ", "
	}
	t = strings.TrimRight(t, ", ")
	result, err := s.database.Exec("insert into Accommodation (Room, Tenants, StayTime) values ($1, $2, $3)", number, t, stayTime)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// Bill выставляет счет за проживание в номере при выселении гостей из отеля
func (s Storage) Bill(ctx context.Context, roomNumber int) (int, error) {
	row := s.database.QueryRowContext(ctx, "SELECT StayTime FROM Accommodation WHERE Room = $1", roomNumber)
	var stayTime int
	err := row.Scan(&stayTime)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("sqlitaccom.Storage.Bill error: %w", oops.ErrEmptyRoom{Number: roomNumber})
	}
	if err != nil {
		return 0, fmt.Errorf("sqlitaccom.Storage.Bill error: %w", err)
	}
	price, err := s.rooms.Price(ctx, roomNumber)
	if err != nil {
		return 0, fmt.Errorf("sqlitaccom.Storage.Bill error: %w", err)
	}
	_, err = s.database.Exec("DELETE FROM Accommodation WHERE Room = $1", roomNumber)
	if err != nil {
		return 0, fmt.Errorf("sqlitaccom.Storage.Bill error: %w", err)
	}
	return price * stayTime, nil
}

func (s Storage) Description(ctx context.Context) error {
	return s.rooms.Show(ctx)
}

func (s Storage) Close() error {
	err := s.database.Close()
	if err != nil {
		return fmt.Errorf("sqliteaccom.Storage.Close error: %w", err)
	}
	err = s.rooms.Close()
	if err != nil {
		return fmt.Errorf("sqliteaccom.Storage.Close error: %w", err)
	}
	return nil
}
