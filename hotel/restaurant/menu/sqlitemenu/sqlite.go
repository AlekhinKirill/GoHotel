// Пакет sqlitemenu реализует меню посредством работы с базой данных
package sqlitemenu

import (
	"Go_projects/hotel/oops"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	//"Go_projects/databases"
)

// Dish хранит информацию о блюде, полученную из строки базы данных
type Dish struct {
	Id    int
	Name  string
	Price int
}

// Storage реализует интерфейс restaurnat.Menu посредством работы с базой данных
type Storage struct {
	database *sql.DB
}

// NewStorage является конструктором Storage
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

// Load выгружает всю информацию о блюдах из меню в более удобной для дальнейшей работы форме
func (s Storage) Load(ctx context.Context) ([]Dish, error) {
	rows, err := s.database.QueryContext(ctx, "select * from Menu")
	if err != nil {
		return nil, err
	}
	var dishes = make([]Dish, 0)
	for rows.Next() {
		var dish Dish
		err = rows.Scan(&dish.Id, &dish.Name, &dish.Price)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}
	return dishes, nil
}

// Show выводит меню на экран
func (s Storage) Show(ctx context.Context) error {
	dishes, err := s.Load(ctx)
	if err != nil {
		return err
	}
	for _, dish := range dishes {
		fmt.Printf("%s: %d рублей", dish.Name, dish.Price)
	}
	return nil
}

// Price возвращает стоимость конкретного блюда в меню
func (s Storage) Price(ctx context.Context, dish string) (int, error) {
	dishes, err := s.Load(ctx)
	if err != nil {
		return 0, err
	}
	for _, value := range dishes {
		if value.Name != dish {
			continue
		}
		return value.Price, nil
	}
	return 0, oops.ErrOutOfMenu{Dish: dish}
}

// Breakfast возвращает стоимость завтрака на одного человека
func (s Storage) Breakfast(ctx context.Context) (int, error) {
	price, err := s.Price(ctx, "Завтрак")
	if err != nil {
		if errors.Is(err, oops.ErrOutOfMenu{Dish: "Завтрак"}) {
			return 0, fmt.Errorf("завтрак в отеле не предусмотрен")
		}
		return 0, err
	}
	return price, nil
}
