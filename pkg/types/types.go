package pkg

import "strconv"

func HasValue(v any) bool {
	switch v := v.(type) {
	case int:
		if v != 0 {
			return true
		}
	case string:
		if v != "" {
			return true
		}
	case float64:
		if v != float64(0) {
			return true
		}
	case bool:
		return v
	case nil:
		return false
	default:
		return false
	}
	return false
}

func ParseIntOrDefault(str string, defaultValue int) int {
	// Tenta converter a string para int
	value, err := strconv.Atoi(str)
	if err != nil {
		// Se houver erro, retorna o valor padrão
		return defaultValue
	}
	// Se a conversão for bem-sucedida, retorna o valor convertido
	return value
}
