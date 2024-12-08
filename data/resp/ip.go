package resp

// IPInfo IP 信息
type IPInfo struct {
	IP        string `json:"ip"`
	ISPrivate bool   `json:"is_private"`
	Time      int64  `json:"time"`
}
