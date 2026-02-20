package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=qwerty123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Данные клиентки Plus Size
	targetArticul := "04042"
	userBust := 106
	userWaist := 88
	userHips := 115

	fmt.Printf("🤖 AI-СТИЛИСТ AMANI (Экспертный вердикт)\n")
	fmt.Println(strings.Repeat("=", 50))

	query := `
		SELECT size_name, category, ease_allowance_cm, 
		       bust_min, bust_max, waist_min, waist_max, hips_min, hips_max
		FROM product_metadata 
		WHERE sku = $1
		ORDER BY bust_max ASC
	`
	rows, err := db.Query(query, targetArticul)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type BestMatch struct {
		Size  string
		Ease  int
		Score float64 // Насколько размер близок к идеалу
	}
	var best *BestMatch

	for rows.Next() {
		var sizeName, category string
		var easeBase, bMin, bMax, wMin, wMax, hMin, hMax int
		rows.Scan(&sizeName, &category, &easeBase, &bMin, &bMax, &wMin, &wMax, &hMin, &hMax)

		// Проверки на физическое соответствие
		bustOk := userBust >= (bMin-4) && userBust <= (bMax+4)
		waistOk := (wMax == 0) || (userWaist <= wMax+8)
		hipsOk := (hMax == 0) || (userHips <= hMax+8)

		if bustOk && waistOk && hipsOk {
			currentEase := (bMax + easeBase) - userBust

			// Считаем отклонение от "задуманной" свободы облегания
			// Чем меньше разница между реальной свободой и базовой — тем лучшеScore
			diff := float64(currentEase - easeBase)
			if diff < 0 {
				diff = -diff
			} // Берем модуль

			if best == nil || diff < best.Score {
				best = &BestMatch{
					Size:  sizeName,
					Ease:  currentEase,
					Score: diff,
				}
			}
		}
	}

	if best != nil {
		fmt.Println("💬 ВЕРДИКТ СТИЛИСТА:")
		fmt.Printf("Для ваших параметров идеально подходит размер: **%s**.\n", best.Size)
		fmt.Printf("Посадка будет именно такой, как задумано дизайнером (свобода %d см).\n", best.Ease)
		fmt.Println("Этот размер гарантирует комфорт и сохранение правильного силуэта модели.")
	} else {
		fmt.Println("💬 ВЕРДИКТ СТИЛИСТА:")
		fmt.Println("К сожалению, эта модель не сядет на ваши параметры так, как это предусмотрено стандартами бренда.")
		fmt.Println("Рекомендую обратить внимание на модели другого кроя.")
	}
}
