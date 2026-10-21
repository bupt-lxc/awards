package services

import (
	"awards-management/database"
	"fmt"
	"sort"
	"strings"
)

type NameSuggestionService struct{}

func NewNameSuggestionService() *NameSuggestionService {
	return &NameSuggestionService{}
}

// SuggestNames returns a small, name-only result set for autocomplete inputs.
// It deliberately excludes certificate numbers and all award details.
func (s *NameSuggestionService) SuggestNames(nameType, keyword string, limit int) ([]string, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []string{}, nil
	}

	if limit <= 0 || limit > 10 {
		limit = 10
	}

	if nameType != "group" {
		return nil, fmt.Errorf("不支持的名称类型")
	}
	table, column := "group_awards", "group_name"

	// Fetch a bounded candidate set, then rank prefix matches before substring matches.
	var names []string
	err := database.DB.Table(table).
		Distinct(column).
		Where(column+" <> '' AND LOCATE(?, "+column+") > 0", keyword).
		Pluck(column, &names).Error
	if err != nil {
		return nil, fmt.Errorf("查询名称建议失败: %w", err)
	}

	sort.SliceStable(names, func(i, j int) bool {
		iPrefix := strings.HasPrefix(strings.ToLower(names[i]), strings.ToLower(keyword))
		jPrefix := strings.HasPrefix(strings.ToLower(names[j]), strings.ToLower(keyword))
		if iPrefix != jPrefix {
			return iPrefix
		}
		if len([]rune(names[i])) != len([]rune(names[j])) {
			return len([]rune(names[i])) < len([]rune(names[j]))
		}
		return names[i] < names[j]
	})

	if len(names) > limit {
		names = names[:limit]
	}
	return names, nil
}
