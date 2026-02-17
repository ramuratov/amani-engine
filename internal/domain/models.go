package domain

import "time"

// UserParams — параметры, которые вводит клиентка
type UserParams struct {
	Bust   float64 // Обхват груди
	Waist  float64 // Обхват талии
	Hips   float64 // Обхват бедер
	Height float64 // Рост
}

// ProductSpec — реальные замеры изделия из твоей таблицы (лекала)
type ProductSpec struct {
	SKU          string
	SizeLabel    string
	BustGarment  float64
	HipsGarment  float64
	BackLength   float64
	SleeveLength float64
}

// MatchResult — результат работы алгоритма
type MatchResult struct {
	IsFit           bool
	RecommendedSize string
	Comment         string
}

// User — для регистрации в базе данных
type User struct {
	ID        int
	Phone     string
	InstaNick string
	CreatedAt time.Time
}
