package hotel

import (
	"context"
	"fmt"
	"time"
)

type Restaurant interface {
	Bill(ctx context.Context, roomNumber int) (int, error)
	PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (id int, err error)
	PlaceBreakfast(ctx context.Context, roomNumber int, count int) (id int, err error)
}

type Accommodation interface {
	Bill(ctx context.Context, roomNumber int) (int, error)
	Place(ctx context.Context, number int, tenants []string, stayTime int) (id int, err error)
}

type Hotel struct {
	rest    Restaurant
	accom   Accommodation
	revenue int
}

func NewHotel(rest Restaurant, accom Accommodation) *Hotel {
	return &Hotel{rest: rest, accom: accom, revenue: 0}
}

func (h *Hotel) CheckIn(ctx context.Context, number int, tenants []string, stayTime int, breakfast bool) (id int, err error) {
	defer time.Sleep(time.Second)
	_, err = h.accom.Place(ctx, number, tenants, stayTime)
	if err != nil {
		return 0, fmt.Errorf("hotel.checkIn error %w", err)
	}
	if breakfast {
		_, err = h.rest.PlaceBreakfast(ctx, number, stayTime*len(tenants))
		if err != nil {
			return 0, fmt.Errorf("hotel.checkIn error %w", err)
		}
	}
	return number, nil
}

func (h *Hotel) CheckOut(ctx context.Context, number int) (int, error) {
	accomValueChan, accomErrChan := make(chan int), make(chan error)
	go func() {
		value, err := h.accom.Bill(ctx, number)
		accomValueChan <- value
		accomErrChan <- err
	}()
	restValueChan, restErrChan := make(chan int), make(chan error)
	go func() {
		value, err := h.rest.Bill(ctx, number)
		restValueChan <- value
		restErrChan <- err
	}()
	accomValue, accomErr := <-accomValueChan, <-accomErrChan
	restValue, restErr := <-restValueChan, <-restErrChan
	if accomErr != nil {
		return 0, fmt.Errorf("hotel.checkOut error %w", accomErr)
	}
	if restErr != nil {
		return 0, fmt.Errorf("hotel.checkOut error %w", restErr)
	}
	h.revenue += accomValue + restValue
	return accomValue + restValue, nil
}

func (h *Hotel) PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (id int, err error) {
	return h.rest.PlaceOrder(ctx, roomNumber, dishes)
}

func (h *Hotel) Money() int {
	return h.revenue
}
