package matcher

import (
	"amani-engine/internal/domain"
	"fmt"
)

// MatchSize сравнивает параметры клиента с конкретным изделием
func MatchSize(user domain.UserParams, product domain.ProductSpec) domain.MatchResult {
	// Считаем разницу (свободу облегания) в груди
	diff := product.BustGarment - user.Bust

	// Если изделие меньше тела — не подходит
	if diff < 0 {
		return domain.MatchResult{
			IsFit:           false,
			RecommendedSize: product.SizeLabel,
			Comment:         fmt.Sprintf("Размер %s будет мал в груди", product.SizeLabel),
		}
	}

	// Описываем посадку
	comment := "Идеальная посадка"
	if diff > 15 {
		comment = fmt.Sprintf("Свободный оверсайз (запас %.1f см)", diff)
	}

	return domain.MatchResult{
		IsFit:           true,
		RecommendedSize: product.SizeLabel,
		Comment:         comment,
	}
}
