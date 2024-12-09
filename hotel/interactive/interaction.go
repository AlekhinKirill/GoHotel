package interactive

import (
	"Go_projects/hotel"
	"Go_projects/hotel/oops"
	"context"
	"fmt"
)

type Interactive struct {
	Hotel *hotel.Hotel
}

func receiveRoomNumber() int {
	var roomNumber int
	fmt.Print("введите номер комнаты: ")
	fmt.Scan(&roomNumber)
	return roomNumber
}

func (in *Interactive) checkIn(ctx context.Context) error {
	var (
		name         string
		stayTime     int
		yesOrNo      string
		ifBreakfast  bool
		i            int
		guestsNumber int
	)
	number := receiveRoomNumber()
	fmt.Print("введите число гостей: ")
	fmt.Scan(&guestsNumber)
	tenants := make([]string, guestsNumber)
	fmt.Println("введите фамилии постояльцев: ")
	for i = 0; i < guestsNumber; i++ {
		fmt.Scan(&name)
		tenants[i] = name
	}
	fmt.Print("планируемое время пребывания (число ночей): ")
	fmt.Scan(&stayTime)
	fmt.Print("будут ли завтракать в отеле? ")
	fmt.Scan(&yesOrNo)
	switch yesOrNo {
	case "да":
		ifBreakfast = true
	case "нет":
		ifBreakfast = false
	default:
		return fmt.Errorf("ответ на вопрос, будут ли завтракать в отеле, должен быть либо да, либо нет")
	}
	_, err := in.Hotel.CheckIn(ctx, number, tenants, stayTime, ifBreakfast)
	if err != nil {
		return fmt.Errorf("interactive.Interactive.CheckIn error: %w", err)
	}
	return nil
}

func receiveAnOrder() []string {
	var (
		dish   string
		number int
		i      int
		arr    = make([]string, 0)
	)
	fmt.Println("введите заказ в формате: блюдо и через пробел их количество, для завершения вбейте слово завершить")
	for dish != "завершить" {
		fmt.Scan(&dish)
		if dish == "завершить" {
			break
		}
		fmt.Scan(&number)
		for i = 0; i < number; i++ {
			arr = append(arr, dish)
		}
	}
	return arr
}

func printBill(sum int) {
	fmt.Printf("к оплате: %d рублей \n", sum)
}

func (in *Interactive) information(ctx context.Context) error {
	var mode string
	fmt.Println("выберите, что Вас интересует: комнаты, постояльцы, меню, выручка, заказы")
	fmt.Scan(&mode)
	var err error
	switch mode {
	case "комнаты":
		fmt.Println("информация о номерах:")
		err = in.Hotel.Accom.Description(ctx)
	case "постояльцы":
		fmt.Println("информация о постояльцах:")
		err = in.Hotel.Accom.Show(ctx)
	case "меню":
		err = in.Hotel.Rest.ShowMenu(ctx)
		fmt.Println("приятного аппетита!")
	case "выручка":
		fmt.Printf("выручка составляет %d рублей\n", in.Hotel.Money())
	case "заказы":
		err = in.Hotel.Rest.Show(ctx)
	default:
		err = oops.ErrOperationNameMistake{Input: mode}
	}
	if err != nil {
		return fmt.Errorf("interactive.Interactive.information error: %w", err)
	}
	return nil
}

// Request реализует всевозможные взаимодействия пользователя со струкурой отеля посредством ввода с командной строки
func (in *Interactive) Request(ctx context.Context) error {
	var mode string
	var err error
	finish := false
	for !finish {
		fmt.Println("выберите операцию: регистрация, выселение, заказ, информация или завершение")
		fmt.Scan(&mode)
		switch mode {
		case "регистрация":
			err = in.checkIn(ctx)
		case "выселение":
			var sum int
			sum, err = in.Hotel.CheckOut(ctx, receiveRoomNumber())
			printBill(sum)
		case "заказ":
			_, err = in.Hotel.Rest.PlaceOrder(ctx, receiveRoomNumber(), receiveAnOrder())
		case "информация":
			err = in.information(ctx)
		case "завершение":
			finish = true
			err = in.Hotel.Close()
		default:
			err = oops.ErrOperationNameMistake{Input: mode}
		}
		if err != nil {
			err = fmt.Errorf("при выполнении запроса %s возникла ошибка: %w", mode, err)
			fmt.Println(err)
		}
	}
	return err
}
