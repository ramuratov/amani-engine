package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func parseRange(s string) (int, int) {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0
	}
	min, _ := strconv.Atoi(parts[0])
	max, _ := strconv.Atoi(parts[1])
	return min, max
}

func main() {
	url := "https://docs.google.com/spreadsheets/d/e/2PACX-1vQk0u-g6Q0Y9EoqRshxLZiCPGr8Nulg971jZvIZ5XhDQUmqDygLm4CnJ6SkZwLLtO0LU_L2SkKNdHZg/pub?gid=1503408859&single=true&output=csv"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка сети:", err)
		return
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	records, _ := reader.ReadAll()

	// ДАННЫЕ КЛИЕНТА
	userBust := 102
	fmt.Printf("\n👗 АМАНИ-ЭНДЖИН: УМНЫЙ ПОДБОР (Ваш ОГ: %d см)\n", userBust)
	fmt.Println(strings.Repeat("-", 50))

	found := false
	for i, row := range records {
		if i == 0 || len(row) < 11 || row[0] == "" {
			continue
		}

		articul := row[0]
		size := row[2]
		bustRangeStr := row[3]
		baseEaseStr := row[10] // Свобода для минимального порога (напр. для 96 см)

		minBust, maxBust := parseRange(bustRangeStr)
		baseEase, _ := strconv.Atoi(strings.TrimSpace(baseEaseStr))

		if userBust >= minBust && userBust <= maxBust {
			// --- ТВОЯ КРИТИЧЕСКАЯ ЛОГИКА ТУТ ---
			// 1. На сколько см клиент больше минимального порога?
			extraBody := userBust - minBust
			// 2. Сколько воздуха реально останется?
			realEase := baseEase - extraBody

			fmt.Printf("✅ Артикул: %s | Размер: %s\n", articul, size)
			fmt.Printf("   (Диапазон размера: %s см)\n", bustRangeStr)

			// Вердикт на основе РЕАЛЬНОГО остатка воздуха
			if realEase >= 20 {
				fmt.Printf("   ВЕРДИКТ: Свободный OVERSIZE (запас %d см воздуха).\n", realEase)
			} else if realEase >= 10 {
				fmt.Printf("   ВЕРДИКТ: Комфортная посадка (запас %d см воздуха).\n", realEase)
			} else if realEase > 0 {
				fmt.Printf("   ВЕРДИКТ: Плотная посадка (запас всего %d см).\n", realEase)
			} else {
				fmt.Printf("   ВЕРДИКТ: Экстра-облегание (впритык).\n")
			}
			fmt.Println()
			found = true
		}
	}

	if !found {
		fmt.Println("❌ К сожалению, подходящих размеров не найдено.")
	}
	fmt.Println(strings.Repeat("-", 50))
}
