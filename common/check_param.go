package common

func CheckPageParam(pageSize, pageNum int) bool {
	if pageSize < 0 || pageNum < 0 {
		return false
	}
	if pageSize >= 1000 {
		return false
	}
	return true
}
