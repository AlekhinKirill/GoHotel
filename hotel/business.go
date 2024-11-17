// Пакет hotel реализует микросервис, который позволяет автоматизировать процессы выселения и заселения гостей из отеля, а также
// осуществляет автоматический учет их взаимодействия с приотельным рестораном
package hotel

import (
	"Go_projects/hotel/accommodation"
	"Go_projects/hotel/restaurant"
	"context"
	"fmt"
	"time"
)

// Hotel хранит и обрабатывает все данные отеля: о ресторане, номерах, постояльцах, а также о выручке отеля
type Hotel struct {
	rest    restaurant.Restaurant
	accom   accommodation.Accommodation
	revenue int
}

// NewHotel является кострутором объекта класса Hotel
func NewHotel(rest restaurant.Restaurant, accom accommodation.Accommodation) *Hotel {
	return &Hotel{rest: rest, accom: accom, revenue: 0}
}

// CheckIn осуществляет регистрацию гостей в базе данных отеля, а также передает информацию о том, будут ли гости
// завтракать в отеле, в ресторан
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

// CheckIn реализует выселение гостей из отеля, возвращая их счет и увеличивая выручку отеля на соответствующую величину
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

// PlaceOrder позволяет офрмить заказ в приотельном ресторане
func (h *Hotel) PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (id int, err error) {
	return h.rest.PlaceOrder(ctx, roomNumber, dishes)
}

// Money возвращает информацию о текущей выручке отеля
func (h *Hotel) Money() int {
	return h.revenue
}
