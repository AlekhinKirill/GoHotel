package main

import (
	"Go_projects/hotel"
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/accommodation/rooms"
	"Go_projects/hotel/restaurant"
	"Go_projects/hotel/restaurant/menu"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func createMyEmptyHotel() *hotel.Hotel {
	myRoomTypes := map[int]rooms.Pair{
		101: {Capacity: 2, Class: "usual"},
		102: {Capacity: 4, Class: "usual"},
		103: {Capacity: 4, Class: "lux"},
		104: {Capacity: 3, Class: "comfort"},
		105: {Capacity: 2, Class: "lux"},
		201: {Capacity: 2, Class: "lux"},
		202: {Capacity: 2, Class: "lux"},
		203: {Capacity: 2, Class: "usual"},
		204: {Capacity: 4, Class: "comfort"},
		205: {Capacity: 3, Class: "comfort"},
		301: {Capacity: 1, Class: "comfort"},
		302: {Capacity: 1, Class: "usual"},
		303: {Capacity: 3, Class: "lux"},
		304: {Capacity: 4, Class: "lux"},
		305: {Capacity: 1, Class: "usual"},
	}

	myPrices := map[rooms.Pair]int{
		{Capacity: 1, Class: "usual"}:   3000,
		{Capacity: 1, Class: "comfort"}: 4000,
		{Capacity: 1, Class: "lux"}:     5000,
		{Capacity: 2, Class: "usual"}:   5000,
		{Capacity: 2, Class: "comfort"}: 6000,
		{Capacity: 2, Class: "lux"}:     7500,
		{Capacity: 3, Class: "usual"}:   7000,
		{Capacity: 3, Class: "comfort"}: 8000,
		{Capacity: 3, Class: "lux"}:     9000,
		{Capacity: 4, Class: "usual"}:   9000,
		{Capacity: 4, Class: "comfort"}: 10000,
		{Capacity: 4, Class: "lux"}:     12000,
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
	accom := accommodation.NewLocalStorage(make(map[int]accommodation.Room), rooms.NewLocalStorage(myRoomTypes, myPrices))
	rest := restaurant.NewLocalStorage(menu.NewStorage(myMenu, menu.Breakfast{Provided: true, Price: 700}), make(map[int][]restaurant.Dinner))
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
