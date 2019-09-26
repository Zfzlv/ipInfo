/**
* thanks to
* github.com/mpolden/echoip
* github.com/oschwald/geoip2-golang
* github.com/oschwald/maxminddb-golang
* GeoLite2 maxmind Db
* @author zfzlv
* 2019/09/26
*/

package main

import(
	"encoding/json"
	"fmt"
	"github.com/Zfzlv/ipInfo/iputil"
	"time"
)

func main(){
	fmt.Println(time.Now(),"--test-ipInfo--")
	r,e := iputil.GetIpInfo("83.145.209.213")
	if e!=nil{
		fmt.Println(e.Error())
	}else{
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	}
	r,e = iputil.GetIpInfo("66.249.64.103")
	if e!=nil{
		fmt.Println(e.Error())
	}else{
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	}
}