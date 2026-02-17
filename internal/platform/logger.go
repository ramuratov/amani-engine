package platform

import (
	"log"
)

// InitLogger просто инициализирует стандартный логгер
func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Логгер инициализирован")
}
