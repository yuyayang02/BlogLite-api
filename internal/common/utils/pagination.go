package utils

func Paginator(totalItems, currentPage, itemsPerPage int) (hasPrevPage, hasNextPage bool) {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	return currentPage > 1, currentPage < totalPages
}
