package convert

func PageToOffset(page int) int {
	switch page {
	case 0:
	case 1:
		page = 0
	default:
		page *= 10
		page -= 10
	}

	return page
}
