package matcher

import (
	"amani-engine/internal/domain"
)

// MatchResult описывает итог проверки
type MatchResult struct {
	IsFit           bool
	RecommendedSize string
	Comment         string
}

// MatchSize — тот самый "мозг", который сравнивает мерки клиента и изделия
func MatchSize(client domain.UserParams, product domain.ProductSpec) MatchResult {
	// Простейшая проверка по груди для теста
	// Если грудь изделия больше груди клиента — значит влезет
	if product.BustGarment >= client.Bust {
		return MatchResult{
			IsFit:           true,
			RecommendedSize: product.SizeLabel,
			Comment:         "Размер подходит, отличная посадка!",
		}
	}

	return MatchResult{
		IsFit:   false,
		Comment: "Изделие будет мало в груди",
	}
}
