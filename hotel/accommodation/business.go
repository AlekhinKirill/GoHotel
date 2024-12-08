// Пакет accommodation релизует учет занимаемых номеров в отеле
package accommodation

import (
	"context"
)

// Accommodation осуществляет взаимодействие с помещениями отеля: дает возможность заселить постояльцев и
// выставить им счет за номер при выселении
type Accommodation interface {
	Bill(ctx context.Context, roomNumber int) (int, error)
	Place(ctx context.Context, number int, tenants []string, stayTime int) (id int, err error)
	Close() error
}

// Room хранит в себе информацию о данном номере в отеле и его текущих жильцах
// предполагается, что объект этого класса создается при заселении в него жильцов
type Room struct {
	Number   int
	Tenants  []string
	StayTime int
}

// RoomsDescription содержит информацию об имеющихся в отеле номерах: об их вместительности, уровне комфорта, цене
type RoomsDescription interface {
	Show(ctx context.Context) error
	Price(ctx context.Context, roomNumber int) (int, error)
	Capacity(ctx context.Context, roomNumber int) (int, error)
	Type(ctx context.Context, roomNumber int) (string, error)
	Close() error
}
