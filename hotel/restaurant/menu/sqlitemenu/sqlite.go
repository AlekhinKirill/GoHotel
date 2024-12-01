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

type Dish struct {
	Id    int
	Name  string
	Price int
}

type Storage struct {
	database *sql.DB
}

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

/*
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Exec("insert into products (model, company, price) values ('iPhone X', $1, $2)",
		"Apple", 72000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
}
*/
