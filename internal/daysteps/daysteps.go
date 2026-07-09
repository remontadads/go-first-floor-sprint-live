package daysteps

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"spentalories" // импортируем пакет для расчёта калорий
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку с данными о шагах и продолжительности
func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Парсим количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("неверный формат количества шагов")
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// Парсим продолжительность
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	return steps, duration, nil
}

// DayActionInfo возвращает информацию о дневной активности
func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		// В случае ошибки выводим её и возвращаем пустую строку
		return ""
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории с помощью функции из пакета spentcalories
	calories, err := spentalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		// В случае ошибки возвращаем пустую строку
		return ""
	}

	// Формируем строку результата
	return "Количество шагов: " + strconv.Itoa(steps) + ".\n" +
		"Дистанция составила " + strconv.FormatFloat(distanceKm, 'f', 2, 64) + " км.\n" +
		"Вы сожгли " + strconv.FormatFloat(calories, 'f', 2, 64) + " ккал."
}
