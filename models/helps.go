package models

// 辅助函数群
func stringInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
func removeFromSliceUsingCopy(slice []string, target string) []string {
	index := -1
	for i, s := range slice {
		if s == target {
			index = i
			break
		}
	}
	if index == -1 {
		return slice // 如果目标字符串不在切片中，直接返回原始切片
	}
	// 使用 copy 函数将后面的元素向前移动，覆盖掉目标元素
	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1] // 删除最后一个元素
}
