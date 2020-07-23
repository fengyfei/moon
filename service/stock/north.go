package stock

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fengyfei/moon/service/request"
)

const (
	northURLFormat = "http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=SGT_SGTJYRB&txtDate=%s&random=0.07009646504151434"
)

var (
	errInvalidRequest = errors.New("[north] create request failed")
)

// GetNorthDailyReport -
func GetNorthDailyReport(date string) (string, error) {
	req := request.Get(fmt.Sprintf(northURLFormat, date), nil)

	if req == nil {
		return "", errInvalidRequest
	}

	request.SetAjaxHeader(req, "http://www.szse.cn/szhk/szhktradeinfo/szdaily/index.html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
