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

	// ТЕСТОВЫЕ ДАННЫЕ КЛИЕНТА
	userBust := 102
	fmt.Printf("\n👗 АМАНИ-ЭНДЖИН: РЕЗУЛЬТАТ ПОДБОРА (ОГ: %d см)\n", userBust)
	fmt.Println(strings.Repeat("-", 45))

	for i, row := range records {
		if i == 0 || len(row) < 11 || row[0] == "" {
			continue
		}

		articul := row[0]
		size := row[2]
		bustRangeStr := row[3]
		easeBust := row[10] // Колонка K: Свобода по груди

		minBust, maxBust := parseRange(bustRangeStr)

		if userBust >= minBust && userBust <= maxBust {
			fmt.Printf("✅ Артикул: %s | Размер: %s\n", articul, size)
			fmt.Printf("   Посадка: Свобода в груди +%s см\n\n", easeBust)
		} else if userBust > maxBust && userBust <= maxBust+4 {
			fmt.Printf("⚠️ Артикул: %s | Размер: %s\n", articul, size)
			fmt.Printf("   ВНИМАНИЕ: Будет сидеть плотно (впритык)\n\n")
		}
	}
	fmt.Println(strings.Repeat("-", 45))
}
