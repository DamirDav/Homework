package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// разбиваем строку вида "шаги,тип тренировки,длительность"
	// и возвращаем распарсенные значения или ошибку валидации.
	data = strings.TrimSpace(data)
	if data == "" {
		fmt.Println("ожидалось три значения через запятую")
		return 0, "", 0, errors.New("invalid input")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		fmt.Println("ожидалось три значения через запятую")
		return 0, "", 0, errors.New("invalid format")
	}

	stepsStr := strings.TrimSpace(parts[0])
	trainingType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		fmt.Println("количество шагов должно быть больше 0")
		return 0, "", 0, errors.New("invalid steps")
	}

	if trainingType == "" {
		return 0, "", 0, errors.New("empty activity")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		fmt.Println("продолжительность должна быть больше 0")
		return 0, "", 0, errors.New("invalid duration")
	}
	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	//возвращаем дистанцию в километрах, исходя из количества шагов и роста.
	lenStep := height * stepLengthCoefficient
	return float64(steps) * lenStep / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// вычисляем среднюю скорость (км/ч) по шагам, росту и длительности.
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	speed := dist / duration.Hours()
	return speed

}

// формируем отчёт о тренировке
// количество шагов, тип, длительность, дистанция, скорость и калории.
// Ошибки парсинга выводит в лог; при неизвестном типе возвращает ошибку.
func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Парсим строку и логируем ошибку, если есть
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 2. Проверяем параметры
	if steps <= 0 {
		return "", fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return "", fmt.Errorf("вес должен быть больше 0")
	}
	if duration <= 0 {
		return "", fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 3. Общие расчёты: дистанция и скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// 4. Определяем тип тренировки и считаем калории
	var calories float64
	switch trainingType {
	case "Бег", "бег", "running":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба", "ходьба", "walking":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 5. Формируем строку (даже «Сожгли калорий: …»)
	report := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration.Hours(), dist, speed, calories)
	return report, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//возвращаем количество калорий, потраченных при беге

	if steps <= 0 {
		fmt.Println("количество шагов должно быть больше 0")
		return 0, errors.New("invalid steps")
	}
	if weight <= 0 {
		fmt.Println("вес должен быть больше 0")
		return 0, errors.New("invalid weight")
	}
	if duration <= 0 {
		fmt.Println("продолжительность должна быть больше 0")
		return 0, errors.New("invalid duration")
	}

	speed := meanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Возвращаем количество калорий, потраченных при ходьбе,

	if steps <= 0 {
		fmt.Println("количество шагов должно быть больше 0")
		return 0, errors.New("invalid steps")
	}
	if weight <= 0 {
		fmt.Println("вес должен быть больше 0")
		return 0, errors.New("invalid weight")
	}
	if height <= 0 {
		fmt.Println("рост должен быть больше 0")
		return 0, errors.New("invalid height")
	}
	if duration <= 0 {
		fmt.Println("продолжительность должна быть больше 0")
		return 0, errors.New("invalid duration")
	}

	speed := meanSpeed(steps, height, duration)
	baseCalories := (weight * speed * duration.Minutes()) / minInH
	return baseCalories * walkingCaloriesCoefficient, nil
}
