// Пакет sqliterest реализует учет заказов, сделанных в ресторане посредством работы с базой данных
package sqliterest

import (
	"Go_projects/hotel/restaurant"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	//"Go_projects/databases"
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
		log.Fatal(fmt.Errorf("NewStorage error: %w", err))
	}
	db.Close()
	return &Storage{
		database: db,
		menu:     menu,
	}
}

// Request хранит информацию о заказах, полученную из строк базы данных
type Request struct {
	Id    int
	Room  int
	Order string
	Price int
}

// Show выводит информацию о сделанных заказах на экран
func (s Storage) Show(ctx context.Context) error {
	rows, err := s.database.QueryContext(ctx, "select * from Restaurant")
	if err != nil {
		return err
	}
	for rows.Next() {
		var req Request
		err = rows.Scan(&req.Id, &req.Room, &req.Order, &req.Price)
		if err != nil {
			return err
		}
		fmt.Printf("Заказ №%d, комната №%d: %s. Стоимость: %d рублей", req.Id, req.Room, req.Order, req.Price)
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
			return -1, err
		}
		sum += price
	}
	order = strings.TrimRight(order, ", ")
	result, err := s.database.Exec("insert into Restaurant (Room, Order, Price) values (?, ?, ?)", roomNumber, order, sum)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// PlaceBreakfast размещает информацию о завтраках в базе данных отеля
func (s Storage) PlaceBreakfast(ctx context.Context, roomNumber int, count int) (int, error) {
	price, err := s.menu.Breakfast(ctx)
	if err != nil {
		return -1, err
	}
	order := strings.Repeat("Завтрак, ", count)
	order = strings.TrimRight(order, ", ")
	result, err := s.database.Exec("insert into dinners (room, order, price) values (?, ?, ?)", roomNumber, order, price*count)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// Bill выставляет счет от ресторана при выселении постояльцем из отеля с учетом всех сделанных ими заказов и посещенных завтраков
func (s Storage) Bill(ctx context.Context, roomNumber int) (int, error) {
	rows, err := s.database.QueryContext(ctx, "select Price from Menu where Room = ?", roomNumber)
	if err != nil {
		return 0, err
	}
	var sum int
	for rows.Next() {
		var price int
		err = rows.Scan(&price)
		if err != nil {
			return 0, err
		}
		sum += price
	}
	return sum, nil
}
