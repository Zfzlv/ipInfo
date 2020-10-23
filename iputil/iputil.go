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
	//"os/exec"
	"path"
	//"path/filepath"
	"runtime"
	"strings"
)

type IpInfo struct {
	IP                  net.IP   `json:"ip"`         //origin ip
	IPDecimal           *big.Int `json:"ip_decimal"` //ip to bigInt
	Continent           string   `json:"continent"`
	ContinentIso        string   `json:"continent_iso"`
	Country             string   `json:"country"`
	CountryEU           *bool    `json:"country_eu"`  //is Eu Country？
	CountryISO          string   `json:"country_iso"` //ISO country code
	Province            string   `json:"province"`
	ProvinceIso         string   `json:"province_iso"`
	City                string   `json:"city"`
	Hostname            string   `json:"hostname,omitempty"`
	Latitude            float64  `json:"latitude,omitempty"`
	Longitude           float64  `json:"longitude,omitempty"`
	TimeZone            string   `json:"time_zone,omitempty"`
	ASN                 string   `json:"asn"` //asn network code
	ASNOrg              string   `json:"asn_org"`
	IsAnonymousProxy    bool     `json:"is_anonymous_proxy"`
	IsSatelliteProvider bool     `json:"is_satellite_provider"`
}

var (
	Reader geo.Reader
	region *reader.GetCnPhy
	err    error
)

func init() {
	log.Println("-init-geoReader-")
	dir := getCurrentDirectory()
	Reader, err = geo.Open(dir+"/GeoLite2-Country.mmdb", dir+"/GeoLite2-City.mmdb", dir+"/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatalln("-init-geoReader-err:-" + err.Error())
	}
	region, err = reader.New(dir + "/Cnphy.db")
	if err != nil {
		log.Fatalln("-init-cnPhy-err:-" + err.Error())
	} else {
		defer region.Close()
		//init load db to memory
		ip, _ := region.MemorySearch("8.8.8.8")
		log.Println(ip)
	}
	log.Println("-init-geoReader-success-")
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

func GetIpInfo(remoteIP, outputLang string) (IpInfo, error) {
	if err != nil {
		return IpInfo{}, err
	}
	if outputLang == "" {
		outputLang = "zh-CN" //support en/de/es/fr/ja/pt-BR/ru/zh-CN
	}
	ip := net.ParseIP(remoteIP)
	ipDecimal := util.ToDecimal(ip)
	country, _ := Reader.Country(ip, outputLang)
	city, _ := Reader.City(ip, outputLang)
	asn, _ := Reader.ASN(ip)
	//hostname, _ := util.LookupAddr(ip)
	var autonomousSystemNumber string
	if asn.AutonomousSystemNumber > 0 {
		autonomousSystemNumber = fmt.Sprintf("AS%d", asn.AutonomousSystemNumber)
	}
	if outputLang == "zh-CN" {
		cnIp, cnErr := region.MemorySearch(remoteIP)
		if cnErr == nil && cnIp.ISP != "0" {
			if cnIp.ISP != "0" {
				asn.AutonomousSystemOrganization = cnIp.ISP
			}
			if country.Name == "" && cnIp.Country != "0" {
				country.Name = cnIp.Country
			}
			if city.Province == "" && cnIp.Province != "0" {
				city.Province = cnIp.Province
			}
			if city.Name == "" && cnIp.City != "0" {
				city.Name = cnIp.City
			}
		}
		if strings.Contains(city.Province, "省") {
			city.Province = city.Province[0:strings.LastIndex(city.Province, "省")]
		}
		if strings.Contains(city.Province, "自治区") {
			city.Province = city.Province[0:strings.LastIndex(city.Province, "自治区")]
		}
		if strings.Contains(city.Province, "特别行政区") {
			city.Province = city.Province[0:strings.LastIndex(city.Province, "特别行政区")]
		}
		if strings.Contains(city.Name, "市") {
			city.Name = city.Name[0:strings.LastIndex(city.Name, "市")]
		}
	}
	return IpInfo{
		IP:                  ip,
		IPDecimal:           ipDecimal,
		Continent:           country.Continent,
		ContinentIso:        country.ContinentCode,
		Country:             country.Name,
		CountryISO:          country.ISO,
		CountryEU:           country.IsEU,
		Province:            city.Province,
		ProvinceIso:         city.ProvinceCode,
		City:                city.Name,
		Latitude:            city.Latitude,
		Longitude:           city.Longitude,
		TimeZone:            city.TimeZone,
		ASN:                 autonomousSystemNumber,
		ASNOrg:              asn.AutonomousSystemOrganization,
		IsAnonymousProxy:    country.IsAnonymousProxy || city.IsAnonymousProxy,
		IsSatelliteProvider: country.IsSatelliteProvider || city.IsSatelliteProvider,
	}, nil
}
