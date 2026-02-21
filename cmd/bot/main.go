package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
	msToken   = "7b10aaf4f6c38ab25c9930699a3d3de09e82d25b"
)

func getMSData(articul string) (string, string) {
	client := &http.Client{}
	url := fmt.Sprintf("https://online.moysklad.ru/api/remap/1.2/entity/product?filter=article=%s", articul)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+msToken)
	resp, err := client.Do(req)
	if err != nil {
		return "Модель AMANI", ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Rows []struct {
			Name       string `json:"name"`
			SalePrices []struct {
				Value float64 `json:"value"`
			} `json:"salePrices"`
		} `json:"rows"`
	}
	json.Unmarshal(body, &data)
	if len(data.Rows) > 0 {
		price := data.Rows[0].SalePrices[0].Value / 100
		return data.Rows[0].Name, fmt.Sprintf("%.0f ₸", price)
	}
	return "Модель AMANI", ""
}

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
	if db == nil {
		log.Println("❌ Ошибка: База данных не инициализирована")
		return
	}
	fmt.Printf("🔄 [%s] Синхронизация с таблицей...\n", time.Now().Format("15:04:05"))

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

	_, err = db.Exec("DELETE FROM product_metadata")
	if err != nil {
		log.Println("❌ Ошибка очистки базы:", err)
		return
	}

	for i, record := range records {
		if i == 0 || len(record) < 15 || record[0] == "" {
			continue
		}

		sku := record[0]
		cat := record[1]
		size := record[2]
		bMin, bMax := parseRange(record[3])
		wMin, wMax := parseRange(record[4])
		hMin, hMax := parseRange(record[5])
		ease, _ := strconv.Atoi(record[10])
		rMin, _ := strconv.Atoi(record[13])
		rMax, _ := strconv.Atoi(record[14])

		_, err = db.Exec(`INSERT INTO product_metadata (sku, category, size_name, ease_allowance_cm, bust_min, bust_max, waist_min, waist_max, hips_min, hips_max, height_min, height_max) 
                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			sku, cat, size, ease, bMin, bMax, wMin, wMax, hMin, hMax, rMin, rMax)
		if err != nil {
			log.Println("❌ Ошибка вставки строки:", err)
		}
	}
	fmt.Println("✅ База обновлена.")
}

func main() {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	menu := &telebot.ReplyMarkup{}
	btnHelp := menu.Text("📏 Как сделать замеры?")
	menu.Reply(menu.Row(btnHelp))

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Привет! Я AI-стилист AMANI. ✨\n\nПришлите: `АРТИКУЛ ОГ-ОТ-ОБ-РОСТ` \n(например: `04042 92-74-100-168`)", menu)
	})

	bot.Handle(&btnHelp, func(m *telebot.Message) {
		photo := &telebot.Photo{File: telebot.FromDisk("guide.jpg")}
		photo.Caption = "📐 **Инструкция по замерам AMANI**\n\nПожалуйста, используйте схему выше. Замеры лучше делать в тонком белье на выдохе."
		_, err := bot.Send(m.Sender, photo, telebot.ModeMarkdown)
		if err != nil {
			log.Println("❌ Ошибка отправки фото:", err)
			bot.Send(m.Sender, "Ошибка: файл guide.jpg не найден в папке проекта.")
		}
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		parts := strings.Fields(m.Text)
		if len(parts) < 2 {
			return
		}
		sku := parts[0]
		params := strings.Split(parts[1], "-")
		if len(params) < 3 {
			return
		}
		var uB, uW, uH, uR int
		_, _ = fmt.Sscanf(params[0], "%d", &uB)
		_, _ = fmt.Sscanf(params[1], "%d", &uW)
		_, _ = fmt.Sscanf(params[2], "%d", &uH)
		if len(params) > 3 {
			_, _ = fmt.Sscanf(params[3], "%d", &uR)
		}

		prodName, price := getMSData(sku)
		rows, _ := db.Query("SELECT size_name, ease_allowance_cm, bust_min, bust_max, waist_min, waist_max, hips_min, hips_max, height_min, height_max FROM product_metadata WHERE sku = $1", sku)
		defer rows.Close()

		var bestSize string
		var hWarn string

		for rows.Next() {
			var sn string
			var eb, bMin, bMax, wMin, wMax, hMin, hMax, rMin, rMax int
			_ = rows.Scan(&sn, &eb, &bMin, &bMax, &wMin, &wMax, &hMin, &hMax, &rMin, &rMax)

			if uB >= (bMin-4) && uB <= (bMax+4) {
				bestSize = sn
				if wMax > 0 && uW > (wMax+6) {
					hWarn += "\n⚠️ *Модель может быть плотной в талии.*"
				}
				if hMax > 0 && uH > (hMax+6) {
					hWarn += "\n⚠️ *Модель может быть плотной в бедрах.*"
				}
				if uR > 0 && rMin > 0 && (uR < rMin || uR > rMax) {
					hWarn += fmt.Sprintf("\n⚠️ *На ваш рост (%d см) модель может сесть иначе (стандарт: %d-%d см).* ", uR, rMin, rMax)
				}
				break
			}
		}

		if bestSize != "" {
			res := fmt.Sprintf("👗 **%s**\n💰 Цена: %s\n\n✅ Ваш рекомендуемый размер: **%s**\n%s\n\nЖелаете оформить заказ?", prodName, price, bestSize, hWarn)
			shopMenu := &telebot.ReplyMarkup{}
			btnOrder := shopMenu.URL("🛍️ Написать менеджеру", "https://t.me/amani_manager")
			shopMenu.Inline(shopMenu.Row(btnOrder))
			bot.Send(m.Sender, res, telebot.ModeMarkdown, shopMenu)
		} else {
			bot.Send(m.Sender, "К сожалению, модель не подходит под ваши параметры.")
		}
	})

	fmt.Println("🚀 AMANI ENGINE запущена. Бот онлайн.")
	bot.Start()
}
