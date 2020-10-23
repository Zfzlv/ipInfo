/**
* @author zfzlv
* 2019/09/26
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/Zfzlv/ipInfo/iputil"
	"testing"
	"time"
)

func TestGetIpInfo(t *testing.T) {
	fmt.Println(time.Now(), "--test-ipInfo--")
	r, e := iputil.GetIpInfo("114.84.151.36", "")
	if e != nil {
		fmt.Println(e.Error())
	} else {
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	}
}
