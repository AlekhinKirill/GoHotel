package main

import (
	"Go_projects/hotel"
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/accommodation/accommap"
	"Go_projects/hotel/accommodation/rooms/roomsmap"
	"Go_projects/hotel/accommodation/rooms/sqliterooms"
	"Go_projects/hotel/accommodation/sqlitaccom"
	"Go_projects/hotel/interactive"
	"Go_projects/hotel/restaurant"
	"Go_projects/hotel/restaurant/menu/menumap"
	"Go_projects/hotel/restaurant/menu/sqlitemenu"
	"Go_projects/hotel/restaurant/restmap"
	"Go_projects/hotel/restaurant/sqliterest"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func createMyEmptyHotel() *hotel.Hotel {
	myRoomTypes := map[int]roomsmap.Pair{
		101: {Capacity: 2, Class: "эконом"},
		102: {Capacity: 4, Class: "эконом"},
		103: {Capacity: 4, Class: "люкс"},
		104: {Capacity: 3, Class: "комфорт"},
		105: {Capacity: 2, Class: "люкс"},
		201: {Capacity: 2, Class: "люкс"},
		202: {Capacity: 2, Class: "люкс"},
		203: {Capacity: 2, Class: "эконом"},
		204: {Capacity: 4, Class: "комфорт"},
		205: {Capacity: 3, Class: "комфорт"},
		301: {Capacity: 1, Class: "комфорт"},
		302: {Capacity: 1, Class: "эконом"},
		303: {Capacity: 3, Class: "люкс"},
		304: {Capacity: 4, Class: "люкс"},
		305: {Capacity: 1, Class: "эконом"},
	}

	myPrices := map[roomsmap.Pair]int{
		{Capacity: 1, Class: "эконом"}:  3000,
		{Capacity: 1, Class: "комфорт"}: 4000,
		{Capacity: 1, Class: "люкс"}:    5000,
		{Capacity: 2, Class: "эконом"}:  5000,
		{Capacity: 2, Class: "комфорт"}: 6000,
		{Capacity: 2, Class: "люкс"}:    7500,
		{Capacity: 3, Class: "эконом"}:  7000,
		{Capacity: 3, Class: "комфорт"}: 8000,
		{Capacity: 3, Class: "люкс"}:    9000,
		{Capacity: 4, Class: "эконом"}:  9000,
		{Capacity: 4, Class: "комфорт"}: 10000,
		{Capacity: 4, Class: "люкс"}:    12000,
	}

	myMenu := map[string]int{
		"Шеф-салат":       300,
		"Хинкали":         300,
		"Цезарь":          200,
		"Стейк":           1000,
		"Пицца-пепперони": 500,
		"Маргарита":       300,
		"Болоньезе":       300,
		"Шницель":         500,
		"Лосось-гриль":    700,
		"Сациви":          400,
		"Хачапури":        300,
		"Вино":            500,
		"Шампанское":      500,
		"Мороженое":       200,
		"Штрудель":        200,
	}
	accom := accommap.NewLocalStorage(make(map[int]accommodation.Room), roomsmap.NewLocalStorage(myRoomTypes, myPrices))
	rest := restmap.NewLocalStorage(menumap.NewStorage(myMenu, menumap.Breakfast{Provided: true, Price: 700}), make(map[int][]restaurant.Dinner))
	return hotel.NewHotel(rest, accom)
}

func createMySQLiteHotel() *hotel.Hotel {
	myMenu := sqlitemenu.NewStorage("D:/Go_projects/databases/menu.db")
	myRooms := sqliterooms.NewStorage("D:/Go_projects/databases/apartments.db")
	accom := sqlitaccom.NewStorage("D:/Go_projects/databases/accommodation.db", myRooms)
	rest := sqliterest.NewStorage("D:/Go_projects/databases/restaurant.db", myMenu)
	return hotel.NewHotel(rest, accom)
}

func demonstration() {
	ctx := context.Background()
	h := createMySQLiteHotel()
	defer h.Close()
	numbers := []int{104, 105, 302}
	names := [][]string{{"Иванов", "Петров", "Сидоров"}, {"Пирогов", "Пирогова"}, {"Козлов"}}
	guestsNumbers := []int{3, 2, 1}
	breakfasts := []bool{true, true, false}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			_, err := h.CheckIn(ctx, numbers[i], names[i], guestsNumbers[i], breakfasts[i])
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Println(duration)
	dinners := [][]string{{"Вино", "Хачапури", "Стейк", "Шницель"}, {"Шампанское", "Болоньезе", "Пицца-пепперони", "Штрудель"}, {"Лосось-гриль", "Сациви"}}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			_, err := h.PlaceOrder(ctx, numbers[i], dinners[i])
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	duration = time.Since(start)
	fmt.Println(duration)
	wg.Add(3)
	for _, number := range numbers {
		go func() {
			_, err := h.CheckOut(ctx, number)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%d\n", h.Money())
	duration = time.Since(start)
	fmt.Println(duration)
}

func main() {
	ctx := context.Background()
	h := createMySQLiteHotel()
	in := interactive.Interactive{Hotel: h}
	in.Request(ctx)
}
