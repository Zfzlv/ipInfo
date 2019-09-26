# ipInfo

ipInfo is a tiny tool to get physic address and asn info from ip.

## Installation

```golang
go get github.com/Zfzlv/ipInfo

import "github.com/Zfzlv/ipInfo/iputil"
```

## Usage

Example

An example file can be found at [main_test.go](https://github.com/Zfzlv/ipInfo/main_test.go)

It is as simple as doing this

```golang
import "github.com/Zfzlv/ipInfo/iputil"

...

r,e := iputil.GetIpInfo("180.155.22.93")
if e!=nil{
	log.Fatalln(e.Error())
}
b, _ := json.Marshal(r)
log.Println(string(b))

```

## Thanks
[maxmind](https://dev.maxmind.com/geoip/geoip2/geolite2/)
[geoip2-golang](https://github.com/oschwald/geoip2-golang)
[echoip](https://github.com/mpolden/echoip)

## License
[MIT](https://choosealicense.com/licenses/mit/)
