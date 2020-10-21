package iputil

import (
	"fmt"
	"github.com/Zfzlv/ipInfo/geo"
	"github.com/Zfzlv/ipInfo/reader"
	"github.com/Zfzlv/ipInfo/util"
	"net"
	//"os"
	"log"
	"math/big"
	//"time"
	//"os/exec"
	"path"
	//"path/filepath"
	"runtime"
	"strings"
)

type IpInfo struct {
	IP         net.IP   `json:"ip"`         //origin ip
	IPDecimal  *big.Int `json:"ip_decimal"` //ip to bigInt
	Country    string   `json:"country,omitempty"`
	CountryEU  *bool    `json:"country_eu,omitempty"`  //is Eu Country？
	CountryISO string   `json:"country_iso,omitempty"` //ISO country code
	Province   string   `json:"province,omitempty"`
	City       string   `json:"city,omitempty"`
	Hostname   string   `json:"hostname,omitempty"`
	Latitude   float64  `json:"latitude,omitempty"`
	Longitude  float64  `json:"longitude,omitempty"`
	ASN        string   `json:"asn,omitempty"` //asn network code
	ASNOrg     string   `json:"asn_org,omitempty"`
}

var (
	Reader geo.Reader
	err    error
)

func init() {
	log.Println("-init-geoReader-")
	dir := getCurrentDirectory()
	Reader, err = geo.Open(dir+"/GeoLite2-Country.mmdb", dir+"/GeoLite2-City.mmdb", dir+"/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatalln("-init-geoReader-err:-" + err.Error())
	}
	region, err := reader.New(dir + "/Cnphy.db")
	if err != nil {
		log.Fatalln("-init-cnPhy-err:-" + err.Error())
	} else {
		defer region.Close()
		//init to memory
		ip, _ := region.MemorySearch("8.8.8.8")
		log.Info(ip)
	}
}

func getCurrentDirectory() string {
	_, filename, _, ok := runtime.Caller(1)
	var cwdPath string
	if ok {
		cwdPath = path.Join(path.Dir(filename), "")
	} else {
		cwdPath = "./"
	}
	cwdPath = strings.Replace(cwdPath, "iputil", "data", -1)
	return cwdPath
}

func GetIpInfo(remoteIP string, outputUseChinese bool) (IpInfo, error) {
	if err != nil {
		return IpInfo{}, err
	}
	ip := net.ParseIP(remoteIP)
	ipDecimal := util.ToDecimal(ip)
	country, _ := Reader.Country(ip)
	city, _ := Reader.City(ip)
	asn, _ := Reader.ASN(ip)
	hostname, _ := util.LookupAddr(ip)
	province := ""
	var autonomousSystemNumber string
	if asn.AutonomousSystemNumber > 0 {
		autonomousSystemNumber = fmt.Sprintf("AS%d", asn.AutonomousSystemNumber)
	}
	if outputUseChinese {
		cnIp, cnErr := region.MemorySearch(remoteIP)
		if cnErr == nil && cnIp != "" {
			str := strings.Split(cnIp, "|")
			if len(str) > 5 {
				if str[1] != "0" {
					country.Name = str[1]
				}
				if str[3] != "0" {
					province = str[3]
				}
				if str[4] != "0" {
					city.Name = str[4]
				}
				if str[5] != "0" {
					asn.AutonomousSystemOrganization = str[5]
				}
			}
		}
	}
	return IpInfo{
		IP:         ip,
		IPDecimal:  ipDecimal,
		Country:    country.Name,
		CountryISO: country.ISO,
		CountryEU:  country.IsEU,
		Province:   province,
		City:       city.Name,
		Hostname:   hostname,
		Latitude:   city.Latitude,
		Longitude:  city.Longitude,
		ASN:        autonomousSystemNumber,
		ASNOrg:     asn.AutonomousSystemOrganization,
	}, nil
}
