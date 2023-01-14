package convert

func PageToOffset(page int) int {
	switch page {
	case 0:
	case 1:
		page = 0
	default:
		page = (page * 10) - 10 //nolint:gomnd // .
	}

	return page
}
