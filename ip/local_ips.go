package ip

// 本地可用IP
var LocalIps *SortedSet

func InitLocalIps() {
	LocalIps = NewSortedSet()
}
