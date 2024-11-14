package restaurant

import (
	"context"
)

// Dinner хранит в себе информацию о сделанном заказе: названия заказанных блюд и их суммарную стоимость
type Dinner struct {
	dishes []string
	price  int
}

type Menu interface {
	Price(ctx context.Context, dish string) (int, error)
	Breakfast(ctx context.Context) (int, error)
}
