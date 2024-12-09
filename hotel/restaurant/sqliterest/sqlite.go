// Пакет sqliterest реализует учет заказов, сделанных в ресторане посредством работы с базой данных
package sqliterest

import (
	"Go_projects/hotel/restaurant"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Storage реализует интерфейс restaurant.Restaurant посредством работы с базой данных
type Storage struct {
	database *sql.DB
	menu     restaurant.Menu
}

// NewStorage является конструктором для Storage
func NewStorage(path string, menu restaurant.Menu) *Storage {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Errorf("sqliterest.NewStorage error: %w", err))
	}
	return &Storage{
		database: db,
		menu:     menu,
	}
}

// Request хранит информацию о заказах, полученную из строк базы данных
type Request struct {
	Room   int
	Dinner string
	Price  int
}

// Show выводит информацию о сделанных заказах на экран
func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "select * from Restaurant")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Room, &req.Dinner, &req.Price)
		if err != nil {
			return fmt.Errorf("sqliterest.Storage.Show error: %w", err)
		}
		fmt.Printf("Комната №%d: %s. Стоимость: %d рублей\n", req.Room, req.Dinner, req.Price)
	}
	return nil
}

// PlaceOrder размещает заказ в базе данных ресторана
func (s Storage) PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (int, error) {
	var (
		order string
		sum   int
	)
	for _, dish := range dishes {
		order += dish + ", "
		price, err := s.menu.Price(ctx, dish)
		if err != nil {
			return -1, fmt.Errorf("sqliterest.Storage.PlaceOrder error: %w", err)
		}
		sum += price
	}
	order = strings.TrimRight(order, ", ")
	result, err := s.database.Exec("INSERT INTO Restaurant (Room, Dinner, Price) VALUES ($1, $2, $3);", roomNumber, order, sum)
	if err != nil {
		return -1, fmt.Errorf("sqliterest.Storage.PlaceOrder error in db.Exec: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("sqliterest.Storage.PlaceOrder error: %w", err)
	}
	return int(id), nil
}

// PlaceBreakfast размещает информацию о завтраках в базе данных отеля
func (s Storage) PlaceBreakfast(ctx context.Context, roomNumber int, count int) (int, error) {
	price, err := s.menu.Breakfast(ctx)
	if err != nil {
		return -1, fmt.Errorf("sqliterest.Storage.PlaceBreakfast error: %w", err)
	}
	order := strings.Repeat("Завтрак, ", count)
	order = strings.TrimRight(order, ", ")
	result, err := s.database.Exec("INSERT INTO Restaurant (Room, Dinner, Price) VALUES ($1, $2, $3);", roomNumber, order, price*count)
	if err != nil {
		return -1, fmt.Errorf("sqliterest.Storage.PlaceBreakfast error in db.Exec: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("sqliterest.Storage.PlaceBreakfast error: %w", err)
	}
	return int(id), nil
}

// Bill выставляет счет от ресторана при выселении постояльцем из отеля с учетом всех сделанных ими заказов и посещенных завтраков
func (s Storage) Bill(ctx context.Context, roomNumber int) (int, error) {
	rows, err := s.database.QueryContext(ctx, "Select Price FROM Restaurant WHERE Room = $1", roomNumber)
	if err != nil {
		return 0, fmt.Errorf("sqliterest.Storage.Bill error: %w", err)
	}
	var sum int
	for rows.Next() {
		var price int
		err = rows.Scan(&price)
		if err != nil {
			return 0, fmt.Errorf("sqliterest.Storage.Bill error: %w", err)
		}
		sum += price
	}
	_, err = s.database.Exec("DELETE FROM Restaurant WHERE Room = $1", roomNumber)
	if err != nil {
		return 0, fmt.Errorf("sqliterest.Storage.Bill error: %w", err)
	}
	return sum, nil
}

func (s Storage) ShowMenu(ctx context.Context) error {
	return s.menu.Show(ctx)
}

func (s Storage) Close() error {
	err := s.database.Close()
	if err != nil {
		return fmt.Errorf("sqliterest.Storage.Close error: %w", err)
	}
	err = s.menu.Close()
	if err != nil {
		return fmt.Errorf("sqliterest.Storage.Close error: %w", err)
	}
	return nil
}
