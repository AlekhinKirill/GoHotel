// Пакет restaurant реализует учет заказов, сделанных в ресторане
package restaurant

import (
	"context"
)

// Restraunt реализует взаимодействие с приотельным рестораном: позволяет сделать заказ и разместить его в базе данных ресторана,
// а также выписать счет от ресторана при выселении
// Предполагается, что постояльцы, решившие поесть в отеле, записывают заказ на счет своего номера в отеле и оплачивают при выселении
type Restaurant interface {
	Bill(ctx context.Context, roomNumber int) (int, error)
	PlaceOrder(ctx context.Context, roomNumber int, dishes []string) (id int, err error)
	PlaceBreakfast(ctx context.Context, roomNumber int, count int) (id int, err error)
	ShowMenu(ctx context.Context) error
	Close() error
}

// Dinner хранит в себе информацию о сделанном заказе: названия заказанных блюд и их суммарную стоимость
type Dinner struct {
	Dishes []string
	Price  int
}

// Menu описывет меню ресторана: позволяет выгрузить список доступных блюд с ценами, автоматически получить цену на интересующее блюдо,
// а также получить информацию о цене и наличии завтрака в формате шведский стол
type Menu interface {
	Show(ctx context.Context) error
	Price(ctx context.Context, dish string) (int, error)
	Breakfast(ctx context.Context) (int, error)
	Close() error
}
