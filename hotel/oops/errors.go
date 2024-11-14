package oops

import (
	"fmt"
)

type ErrEmptyRoom struct {
	Number int
}

func (err ErrEmptyRoom) Error() string {
	return fmt.Sprintf("в номере %d никто не проживает", err.Number)
}

type ErrRoomInconsistency struct {
	Number       int
	Capacity     int
	GuestsNumber int
}

func (err ErrRoomInconsistency) Error() string {
	return fmt.Sprintf("номер %d рассчитан на %d-х человек -- заселение %d-х человек невозможно", err.Number, err.Capacity, err.GuestsNumber)
}

type ErrOccupiedAlready struct {
	Number int
}

func (err ErrOccupiedAlready) Error() string {
	return fmt.Sprintf("номер %d уже занят", err.Number)
}

type ErrNoRoom struct {
	Number int
}

func (err ErrNoRoom) Error() string {
	return fmt.Sprintf("номера %d в отеле нет", err.Number)
}

type ErrNoPrice struct {
	Capacity int
	Class    string
}

func (err ErrNoPrice) Error() string {
	return fmt.Sprintf("на категорию %d-xместный номер класса %s не установлена стоимость проживания", err.Capacity, err.Class)
}

type ErrOutOfMenu struct {
	Dish string
}

func (err ErrOutOfMenu) Error() string {
	return fmt.Sprintf("блюда %s пока нет в нашем меню", err.Dish)
}
