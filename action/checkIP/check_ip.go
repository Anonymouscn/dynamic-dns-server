package checkIP

import (
	"github.com/Anonymouscn/dynamic-dns-server/data/resp"
	"github.com/Anonymouscn/dynamic-dns-server/provider"
	"github.com/Anonymouscn/go-partner/restful"
	restfulmodel "github.com/Anonymouscn/go-partner/restful/model"
)

// GetMyIP 获取本地公网 get_my_ip 地址
func GetMyIP() (string, error) {
	rc := restful.NewRestClient().
		SetURL(provider.GetMyIPApiSecret.API).
		Get()
	r := &restfulmodel.Result[resp.IPInfo]{}
	if err := rc.Map(r); err != nil {
		return "", err
	}
	return r.Data.IP, nil
}
