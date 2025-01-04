package utils

func Offset(page, limit int) int {
	return (page - 1) * limit
}
