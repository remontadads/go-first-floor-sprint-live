package spentalories

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.

	stepLengthCoefficient      = 0.5 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5 // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	// Парсим количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат количества шагов")
	}

	// Парсим продолжительность
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	// Возвращаем количество шагов, вид активности, продолжительность и nil
	return steps, parts[1], duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLen := height * stepLengthCoefficient

	// Вычисляем дистанцию в метрах и переводим в километры
	return (float64(steps) * stepLen) / mInKm
}

// meanSpeed вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Переводим продолжительность в часы
	durationHours := duration.Hours()

	// Вычисляем и возвращаем среднюю скорость
	return dist / durationHours
}

// RunningSpentCalories вычисляет калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность параметров
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Вычисляем калории
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories вычисляет калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность параметров
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Вычисляем калории с корректирующим коэффициентом
	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo возвращает информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем данные из строки
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Определяем тип тренировки и вычисляем калории
	var calories float64
	var trainingTypeName string

	switch trainingType {
	case "Бег", "Running":
		trainingTypeName = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба", "Walking":
		trainingTypeName = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	// Формируем строку результата
	return "Тип тренировки: " + trainingTypeName + "\n" +
		"Длительность: " + strconv.FormatFloat(duration.Hours(), 'f', 2, 64) + " ч.\n" +
		"Дистанция: " + strconv.FormatFloat(dist, 'f', 2, 64) + " км.\n" +
		"Скорость: " + strconv.FormatFloat(speed, 'f', 2, 64) + " км/ч\n" +
		"Сожгли калорий: " + strconv.FormatFloat(calories, 'f', 2, 64), nil
}
