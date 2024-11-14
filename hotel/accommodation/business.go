package accommodation

import (
	"context"
)

// Room хранит в себе информацию о данном номере в отеле и его текущих жильцах
// предполагается, что объект этого класса создается при заселении в него жильцов
type Room struct {
	number   int
	tenants  []string
	stayTime int
}

type RoomsDescription interface {
	Price(ctx context.Context, roomNumber int) (int, error)
	Capacity(ctx context.Context, roomNumber int) (int, error)
	Type(ctx context.Context, roomNumber int) (string, error)
}
