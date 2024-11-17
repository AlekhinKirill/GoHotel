package main

import (
	"Go_projects/hotel"
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/accommodation/accommap"
	"Go_projects/hotel/accommodation/rooms/roomsmap"
	"Go_projects/hotel/restaurant"
	"Go_projects/hotel/restaurant/menu/menumap"
	"Go_projects/hotel/restaurant/restmap"
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
		"шеф-салат":       300,
		"хинкали":         300,
		"цезарь":          200,
		"стейк":           1000,
		"пицца-пепперони": 500,
		"маргарита":       300,
		"болоньезе":       300,
		"шницель":         500,
		"лосось-гриль":    700,
		"сациви":          400,
		"хачапури":        300,
		"вино":            500,
		"шампанское":      500,
		"мороженое":       200,
		"штрудель":        200,
	}
	accom := accommap.NewLocalStorage(make(map[int]accommodation.Room), roomsmap.NewLocalStorage(myRoomTypes, myPrices))
	rest := restmap.NewLocalStorage(menumap.NewStorage(myMenu, menumap.Breakfast{Provided: true, Price: 700}), make(map[int][]restaurant.Dinner))
	return hotel.NewHotel(rest, accom)
}

func main() {
	ctx := context.Background()
	h := createMyEmptyHotel()
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
	dinners := [][]string{{"вино", "хачапури", "стейк", "шницель"}, {"шампанское", "болоньезе", "пицца-пепперони", "штрудель"}, {"лосось-гриль", "сациви"}}
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
