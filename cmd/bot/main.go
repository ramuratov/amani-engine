package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/tucnak/telebot.v2" // Исправленный путь здесь
)

func main() {
	// 1. ПОДКЛЮЧЕНИЕ К БАЗЕ
	connStr := "host=localhost port=5432 user=postgres password=qwerty123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка базы:", err)
	}
	defer db.Close()

	// 2. НАСТРОЙКА ТЕЛЕГРАМ-БОТА
	pref := telebot.Settings{
		Token:  "8236828498:AAHgcFlXaab-lqp8Z-5Oom5JVCgb2CanDqM",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка бота:", err)
	}

	fmt.Println("🚀 AMANI ENGINE: Бот запущен и слушает сообщения...")

	// ОБРАБОТКА /start
	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Привет! Я AI-стилист AMANI. ✨\n\nПришлите мне артикул и ваши параметры через пробел:\n`АРТИКУЛ ОГ-ОТ-ОБ` \n\nНапример: `04042 92-74-100`", telebot.ModeMarkdown)
	})

	// ГЛАВНАЯ ЛОГИКА
	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		text := m.Text
		parts := strings.Fields(text)

		if len(parts) < 2 {
			bot.Send(m.Sender, "Пожалуйста, введите данные в формате: `04042 92-74-100`", telebot.ModeMarkdown)
			return
		}

		targetArticul := parts[0]
		params := strings.Split(parts[1], "-")
		if len(params) < 3 {
			bot.Send(m.Sender, "Укажите все 3 параметра через дефис: ОГ-ОТ-ОБ")
			return
		}

		var uBust, uWaist, uHips int
		fmt.Sscanf(params[0], "%d", &uBust)
		fmt.Sscanf(params[1], "%d", &uWaist)
		fmt.Sscanf(params[2], "%d", &uHips)

		// Запрос к твоей таблице
		query := `SELECT size_name, ease_allowance_cm, bust_min, bust_max, waist_max, hips_max 
		          FROM product_metadata WHERE sku = $1`
		rows, err := db.Query(query, targetArticul)
		if err != nil {
			log.Println("Ошибка SQL:", err)
			bot.Send(m.Sender, "Произошла ошибка при поиске модели.")
			return
		}
		defer rows.Close()

		type Match struct {
			Size string
			Ease int
			Diff float64
		}
		var best *Match

		for rows.Next() {
			var sn string
			var eBase, bMin, bMax, wMax, hMax int
			rows.Scan(&sn, &eBase, &bMin, &bMax, &wMax, &hMax)

			bustOk := uBust >= (bMin-6) && uBust <= (bMax+4)
			waistOk := (wMax == 0) || (uWaist <= wMax+8)
			hipsOk := (hMax == 0) || (uHips <= hMax+8)

			if bustOk && waistOk && hipsOk {
				currentEase := (bMax + eBase) - uBust
				diff := float64(currentEase - eBase)
				if diff < 0 {
					diff = -diff
				}

				if best == nil || diff < best.Diff {
					best = &Match{Size: sn, Ease: currentEase, Diff: diff}
				}
			}
		}

		if best != nil {
			msg := fmt.Sprintf("✅ **Ваш идеальный размер: %s**\n\nПосадка будет именно такой, как задумано дизайнером (запас воздуха %d см).", best.Size, best.Ease)
			bot.Send(m.Sender, msg, telebot.ModeMarkdown)
		} else {
			bot.Send(m.Sender, "💬 К сожалению, под эти параметры идеальный размер не найден.")
		}
	})

	bot.Start()
}
