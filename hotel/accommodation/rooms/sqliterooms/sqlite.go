// Пакет sqliterooms предоставляет описание аппартаментов отеля посредством работы с базой данных
package sqliterooms

import (
	"Go_projects/hotel/oops"
	"context"
	"database/sql"
	"fmt"
	"log"
	//"Go_projects/databases"
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
	db.Close()
	return &Storage{
		database: db,
	}
}

// Request является классом для хранения информации из сторок базы данных
type Request struct {
	Id       int
	Number   int
	Capacity int
	Class    string
	Price    int
}

// Show выводит информацию о номерах в отеле
func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "select * from Rooms")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Id, &req.Number, &req.Capacity, &req.Class, &req.Price)
		if err != nil {
			return err
		}
		fmt.Printf("%d. Комната №%d: %d-хместный номер класса %s. Стоимость: %d рублей за ночь",
			req.Id, req.Number, req.Capacity, req.Class, req.Price)
	}
	return nil
}

// Capacity возвращает вместимость номера, то есть число человек, на которое он расчитан
func (s Storage) Capacity(ctx context.Context, roomNumber int) (int, error) {
	row, err := s.database.QueryContext(ctx, "select Capacity from Rooms where Number = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, fmt.Errorf("Storage.Bill error: %w", oops.ErrNoRoom{Number: roomNumber})
	}
	var capacity int
	err = row.Scan(&capacity)
	if err != nil {
		return 0, err
	}
	return capacity, nil
}

// Type возвращает уровень комфортности номера
func (s Storage) Type(ctx context.Context, roomNumber int) (int, error) {
	row, err := s.database.QueryContext(ctx, "select Class from Rooms where Number = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, fmt.Errorf("Storage.Bill error: %w", oops.ErrNoRoom{Number: roomNumber})
	}
	var class int
	err = row.Scan(&class)
	if err != nil {
		return 0, err
	}
	return class, nil
}

// Price возвращает стоимость номера за одну ночь
func (s Storage) Price(ctx context.Context, roomNumber int) (int, error) {
	row, err := s.database.QueryContext(ctx, "select Price from Menu where Room = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, fmt.Errorf("Storage.Bill error: %w", oops.ErrNoRoom{Number: roomNumber})
	}
	var price int
	err = row.Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
