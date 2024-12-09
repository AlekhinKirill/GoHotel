// Пакет sqliterooms предоставляет описание аппартаментов отеля посредством работы с базой данных
package sqliterooms

import (
	"Go_projects/hotel/oops"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Storage реализует интерфейс accommodation.RoomsDescription на основе сохранения данных в базе данных
type Storage struct {
	database *sql.DB
}

// NewStorage является конструктором для Storage
func NewStorage(path string) *Storage {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Errorf("NewStorage error: %w", err))
	}
	return &Storage{
		database: db,
	}
}

// Request является классом для хранения информации из сторок базы данных
type Request struct {
	Number   int
	Capacity int
	Class    string
	Price    int
}

// Show выводит информацию о номерах в отеле
func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "SELECT * FROM Apartments")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Number, &req.Capacity, &req.Class, &req.Price)
		if err != nil {
			return err
		}
		fmt.Printf("Комната №%d: %d-хместный номер класса %s. Стоимость: %d рублей за ночь\n",
			req.Number, req.Capacity, req.Class, req.Price)
	}
	return nil
}

// Capacity возвращает вместимость номера, то есть число человек, на которое он расчитан
func (s Storage) Capacity(ctx context.Context, roomNumber int) (int, error) {
	row := s.database.QueryRowContext(ctx, "SELECT Capacity FROM Apartments WHERE Number = $1", roomNumber)
	var capacity int
	err := row.Scan(&capacity)
	if errors.Is(err, sql.ErrNoRows) {
		err = oops.ErrNoRoom{Number: roomNumber}
	}
	if err != nil {
		return 0, fmt.Errorf("sqliterooms.Storage.Capacity error: %w", err)
	}
	return capacity, nil
}

// Type возвращает уровень комфортности номера
func (s Storage) Type(ctx context.Context, roomNumber int) (string, error) {
	row := s.database.QueryRowContext(ctx, "SELECT Class FROM Apartments WHERE Number = $1", roomNumber)
	var class string
	err := row.Scan(&class)
	if errors.Is(err, sql.ErrNoRows) {
		err = oops.ErrNoRoom{Number: roomNumber}
	}
	if err != nil {
		return "", fmt.Errorf("sqliterooms.Storage.Type error: %w", err)
	}
	return class, nil
}

// Price возвращает стоимость номера за одну ночь
func (s Storage) Price(ctx context.Context, roomNumber int) (int, error) {
	row := s.database.QueryRowContext(ctx, "select Price from Apartments where Number = $1", roomNumber)
	var price int
	err := row.Scan(&price)
	if errors.Is(err, sql.ErrNoRows) {
		err = oops.ErrNoRoom{Number: roomNumber}
	}
	if err != nil {
		return 0, fmt.Errorf("sqliterooms.Storage.Price error: %w", err)
	}
	return price, nil
}

func (s Storage) Close() error {
	err := s.database.Close()
	if err != nil {
		return fmt.Errorf("sqliterooms.Storage.Close error: %w", err)
	}
	return s.database.Close()
}
