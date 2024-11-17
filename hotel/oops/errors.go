// Пакет oops содержит описание ошибок
package oops

import (
	"fmt"
)

// ErrEmptyRoom возникает при попытке выселения из незаселенного номера или оформления заказа на этот номер
type ErrEmptyRoom struct {
	Number int
}

func (err ErrEmptyRoom) Error() string {
	return fmt.Sprintf("в номере %d никто не проживает", err.Number)
}

// ErrRoomInconsistency возникает при попытке заселения в номер неправильного числа гостей, например, троих в двухместный номер
type ErrRoomInconsistency struct {
	Number       int
	Capacity     int
	GuestsNumber int
}

func (err ErrRoomInconsistency) Error() string {
	return fmt.Sprintf("номер %d рассчитан на %d-х человек -- заселение %d-х человек невозможно", err.Number, err.Capacity, err.GuestsNumber)
}

// ErrOccupiedAlready возникает при попытке заселения гостей в уже занятый номер
type ErrOccupiedAlready struct {
	Number int
}

func (err ErrOccupiedAlready) Error() string {
	return fmt.Sprintf("номер %d уже занят", err.Number)
}

// ErrNoRoom возникает при попытке взаимодействовать с комнатой, которой в отеле в принципе нет
type ErrNoRoom struct {
	Number int
}

func (err ErrNoRoom) Error() string {
	return fmt.Sprintf("номера %d в отеле нет", err.Number)
}

// ErrNoPrice возникает если по каким-то причинам цена на номер оказалась неустановленнной (пропуск в данных)
type ErrNoPrice struct {
	Capacity int
	Class    string
}

func (err ErrNoPrice) Error() string {
	return fmt.Sprintf("на категорию %d-xместный номер класса %s не установлена стоимость проживания", err.Capacity, err.Class)
}

// ErrOutOfMenu возникает при попытке заказать блюдо, которое не представлено в меню ресторана
type ErrOutOfMenu struct {
	Dish string
}

func (err ErrOutOfMenu) Error() string {
	return fmt.Sprintf("блюда %s пока нет в нашем меню", err.Dish)
}
