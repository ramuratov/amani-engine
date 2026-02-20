package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

// parseSafeInt обрабатывает пустые строки, запятые и округляет дробные числа
func parseSafeInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", ".")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int(math.Round(val))
}

// parseRange разделяет диапазоны вида "96-104" на минимальное и максимальное значения
func parseRange(s string) (int, int) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0
	}
	if !strings.Contains(s, "-") {
		v := parseSafeInt(s)
		return v, v
	}
	parts := strings.Split(s, "-")
	return parseSafeInt(parts[0]), parseSafeInt(parts[1])
}

func main() {
	// --- ВСТАВЬ СВОЙ ПАРОЛЬ НИЖЕ ВМЕСТО 'YOUR_PASSWORD' ---
	connStr := "host=localhost port=5432 user=postgres password=qwerty123 dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Ссылка на опубликованный CSV твоей Google Таблицы
	url := "https://docs.google.com/spreadsheets/d/e/2PACX-1vQk0u-g6Q0Y9EoqRshxLZiCPGr8Nulg971jZvIZ5XhDQUmqDygLm4CnJ6SkZwLLtO0LU_L2SkKNdHZg/pub?gid=1503408859&single=true&output=csv"

	fmt.Println("⏳ Скачиваю данные из Google Таблиц...")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Ошибка сети:", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Ошибка чтения CSV:", err)
	}

	fmt.Println("🚀 Начинаю синхронизацию с PostgreSQL...")

	// Очищаем таблицу перед новой загрузкой
	_, err = db.Exec("TRUNCATE TABLE product_metadata")
	if err != nil {
		log.Fatal("Ошибка очистки таблицы:", err)
	}

	for i, row := range records {
		// Пропускаем заголовок (строка 1)
		if i == 0 || len(row) < 15 || row[0] == "" {
			continue
		}

		sku := row[0]      // Артикул
		category := row[1] // Категория
		sizeName := row[2] // Размер

		bMin, bMax := parseRange(row[3]) // Обхват груди (D)
		wMin, wMax := parseRange(row[4]) // Обхват талии (E)
		hMin, hMax := parseRange(row[5]) // Обхват бедер (F)

		prodLen := parseSafeInt(row[6])   // Длина изделия (G)
		sleeveLen := parseSafeInt(row[7]) // Длина рукава (H)

		silhouette := row[9]          // Силуэт (J)
		ease := parseSafeInt(row[10]) // Свобода грудь (K)

		hMinRec, hMaxRec := parseRange(row[13]) // Рост (N-O)

		// Запрос на вставку данных в твою таблицу product_metadata
		query := `INSERT INTO product_metadata 
			(sku, category, size_name, bust_min, bust_max, waist_min, waist_max, hips_min, hips_max, 
			product_length, sleeve_length, silhouette, ease_allowance_cm, rec_height_min, rec_height_max) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

		_, err := db.Exec(query, sku, category, sizeName, bMin, bMax, wMin, wMax, hMin, hMax,
			prodLen, sleeveLen, silhouette, ease, hMinRec, hMaxRec)
		if err != nil {
			fmt.Printf("⚠️ Ошибка в строке %d (арт %s): %v\n", i+1, sku, err)
		}
	}

	fmt.Println("✅ Синхронизация успешно завершена! Проверь данные в DBeaver.")
}
