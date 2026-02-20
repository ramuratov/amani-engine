package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	dbConnStr = "host=localhost port=5432 user=postgres password=qwerty123 dbname=postgres sslmode=disable"
	botToken  = "8236828498:AAHgcFlXaab-lqp8Z-5Oom5JVCgb2CanDqM"
	csvURL    = "https://docs.google.com/spreadsheets/d/e/2PACX-1vQk0u-g6Q0Y9EoqRshxLZiCPGr8Nulg971jZvIZ5XhDQUmqDygLm4CnJ6SkZwLLtO0LU_L2SkKNdHZg/pub?gid=1503408859&single=true&output=csv"
)

// parseRange разделяет "88-92" на min=88, max=92
func parseRange(s string) (int, int) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0, 0
	}
	if !strings.Contains(s, "-") {
		val, _ := strconv.Atoi(s)
		return val, val
	}
	parts := strings.Split(s, "-")
	min, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	max, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return min, max
}

func syncData(db *sql.DB) {
	fmt.Printf("🔄 [%s] Авто-синхронизация с Google Таблицей...\n", time.Now().Format("15:04:05"))

	resp, err := http.Get(csvURL)
	if err != nil {
		log.Println("❌ Ошибка загрузки CSV:", err)
		return
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("❌ Ошибка чтения CSV:", err)
		return
	}

	_, _ = db.Exec("DELETE FROM product_metadata")

	for i, record := range records {
		if i == 0 || len(record) < 11 || record[0] == "" {
			continue
		}

		sku := record[0]
		category := record[1]
		sizeName := record[2]

		// Парсим диапазоны из колонок D, E, F (индексы 3, 4, 5)
		bMin, bMax := parseRange(record[3])
		wMin, wMax := parseRange(record[4])
		hMin, hMax := parseRange(record[5])

		// Свобода из колонки K (индекс 10)
		ease, _ := strconv.Atoi(strings.TrimSpace(record[10]))

		_, err := db.Exec(`
			INSERT INTO product_metadata (
				sku, category, size_name, ease_allowance_cm, 
				bust_min, bust_max, waist_min, waist_max, hips_min, hips_max
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			sku, category, sizeName, ease, bMin, bMax, wMin, wMax, hMin, hMax)
		if err != nil {
			log.Println("❌ Ошибка вставки артикула", sku, ":", err)
		}
	}
	fmt.Println("✅ База данных успешно обновлена!")
}

func main() {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Запуск фоновой синхронизации (раз в 24 часа)
	go func() {
		for {
			syncData(db)
			time.Sleep(24 * time.Hour)
		}
	}()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🚀 AMANI ENGINE запущена. Бот онлайн.")

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Привет! Я AI-стилист AMANI. ✨\n\nПришлите: `АРТИКУЛ ОГ-ОТ-ОБ` (напр. `04042 92-74-100`)", telebot.ModeMarkdown)
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		parts := strings.Fields(m.Text)
		if len(parts) < 2 {
			bot.Send(m.Sender, "Формат: `АРТИКУЛ ОГ-ОТ-ОБ`")
			return
		}

		sku := parts[0]
		params := strings.Split(parts[1], "-")
		if len(params) < 3 {
			bot.Send(m.Sender, "Нужно 3 параметра: ОГ-ОТ-ОБ")
			return
		}

		var uB, uW, uH int
		fmt.Sscanf(params[0], "%d", &uB)
		fmt.Sscanf(params[1], "%d", &uW)
		fmt.Sscanf(params[2], "%d", &uH)

		rows, err := db.Query("SELECT size_name, ease_allowance_cm, bust_min, bust_max, waist_max, hips_max FROM product_metadata WHERE sku = $1", sku)
		if err != nil {
			bot.Send(m.Sender, "Ошибка базы данных.")
			return
		}
		defer rows.Close()

		var bestSize string
		var bestEase int
		minDiff := 999.0

		for rows.Next() {
			var sn string
			var eB, bMin, bMax, wMax, hMax int
			rows.Scan(&sn, &eB, &bMin, &bMax, &wMax, &hMax)

			// Экспертная логика (Золотое сечение)
			if uB >= (bMin-6) && uB <= (bMax+4) && (wMax == 0 || uW <= wMax+8) && (hMax == 0 || uH <= hMax+8) {
				currEase := (bMax + eB) - uB
				diff := float64(currEase - eB)
				if diff < 0 {
					diff = -diff
				}

				if diff < minDiff {
					minDiff = diff
					bestSize = sn
					bestEase = currEase
				}
			}
		}

		if bestSize != "" {
			bot.Send(m.Sender, fmt.Sprintf("✅ **Ваш идеальный размер: %s**\n\nЗапас воздуха по груди: %d см.\nПосадка будет соответствовать задумке дизайнера.", bestSize, bestEase), telebot.ModeMarkdown)
		} else {
			bot.Send(m.Sender, "❌ К сожалению, модель не подходит под ваши параметры. Рекомендуем присмотреться к другому крою.")
		}
	})

	bot.Start()
}
