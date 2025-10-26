package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

var errInvalidInput = errors.New("invalid input")

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	// разбиваем строку формата "шаги,длительность" на компоненты
	// и возвращаем распарсенные значения либо ошибку валидации.

	data = strings.TrimSpace(data)
	if data == "" {
		fmt.Println("ожидалось два значения через запятую")
		return 0, 0, errInvalidInput
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		fmt.Println("ожидалось два значения через запятую")
		return 0, 0, errInvalidInput
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		fmt.Println("количество шагов должно быть больше 0")
		return 0, 0, errInvalidInput
	}

	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		fmt.Println("продолжительность должна быть больше 0")
		return 0, 0, errInvalidInput
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	// Проверяем физические параметры пользователя.
	if height <= 0 {
		fmt.Println("Рост должен быть больше 0")
		return ""
	}
	if weight <= 0 {
		fmt.Println("Вес должен быть больше 0")
		return ""
	}

	distM := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distM, calories)
}
