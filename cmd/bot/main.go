package main

import (
	"amani-engine/internal/domain"
	"amani-engine/internal/services/matcher"
	"fmt"
)

func main() {
	fmt.Println("=== AMANI Smart Engine: Тестирование алгоритма ===")

	// 1. Имитируем данные от клиентки (например, из квиза в боте)
	client := domain.UserParams{
		Bust:   92.0, // Обхват груди 92 см
		Waist:  70.0,
		Hips:   98.0,
		Height: 168.0,
	}

	// 2. Имитируем данные изделия (взяли первую колонку с твоего фото)
	product := domain.ProductSpec{
		SKU:         "SHIRT-SILK-01",
		SizeLabel:   "44-46",
		BustGarment: 125.0, // Тот самый оверсайз 125 см
		HipsGarment: 116.0,
	}

	// 3. Запускаем "Мозг" (алгоритм подбора)
	result := matcher.MatchSize(client, product)

	// 4. Выводим результат в консоль
	fmt.Printf("\nПараметры клиента: Грудь %.1f см", client.Bust)
	fmt.Printf("\nПараметры изделия (%s): Грудь %.1f см", product.SizeLabel, product.BustGarment)
	fmt.Println("\n-------------------------------------------")

	if result.IsFit {
		fmt.Printf("ВЕРДИКТ: РЕКОМЕНДОВАНО ✅\n")
		fmt.Printf("Размер: %s\n", result.RecommendedSize)
		fmt.Printf("Комментарий стилиста: %s\n", result.Comment)
	} else {
		fmt.Printf("ВЕРДИКТ: НЕ ПОДХОДИТ ❌\n")
		fmt.Printf("Причина: %s\n", result.Comment)
	}
	fmt.Println("-------------------------------------------")
}
