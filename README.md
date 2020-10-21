# ipInfo

<p align="center">
	<img src="https://camo.githubusercontent.com/5b13bf8be0d98cf8e2764ce07ca68ee02a273f63/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f676f6c616e672d312e31332d626c75652e7376673f7374796c653d666c6174">
	<a href="https://raw.githubusercontent.com/onevcat/Kingfisher/master/LICENSE"><img src="https://img.shields.io/cocoapods/l/Kingfisher.svg?style=flat"></a>
</p>

ipInfo is a tiny tool to get physic address and asn info from ip.

## Installation

```golang
go get github.com/Zfzlv/ipInfo

import "github.com/Zfzlv/ipInfo/iputil"
```

## Usage

Example

An example file can be found at [main_test.go](https://github.com/Zfzlv/ipInfo/blob/master/main_test.go)

It is as simple as doing this

```golang
import "github.com/Zfzlv/ipInfo/iputil"

...

r,e := iputil.GetIpInfo("180.155.22.93",true)
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
[ip2region](https://github.com/lionsoul2014/ip2region)

## License
[MIT](https://choosealicense.com/licenses/mit/)
