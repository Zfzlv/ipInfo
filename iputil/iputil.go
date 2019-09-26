package iputil

import(
	"net"
	"github.com/Zfzlv/ipInfo/geo"
	"github.com/Zfzlv/ipInfo/util"
	"fmt"
	"os"
	"math/big"
	"log"
	//"time"
	"os/exec"
	"path"
)

type IpInfo struct {
	IP         net.IP   `json:"ip"`   //origin ip
	IPDecimal  *big.Int `json:"ip_decimal"` //ip to bigInt
	Country    string   `json:"country,omitempty"` 
	CountryEU  *bool    `json:"country_eu,omitempty"` //is Eu Country？
	CountryISO string   `json:"country_iso,omitempty"` //ISO country code
	City       string   `json:"city,omitempty"`
	Hostname   string   `json:"hostname,omitempty"`
	Latitude   float64  `json:"latitude,omitempty"`
	Longitude  float64  `json:"longitude,omitempty"`
	ASN        string   `json:"asn,omitempty"` //asn network code
	ASNOrg     string   `json:"asn_org,omitempty"`
}

var(
	Reader geo.Reader
	err error
)

func init(){
	log.Println("-init-geoReader-")
	file, _ := exec.LookPath(os.Args[0])
    dir,_ := path.Split(file)
	Reader, err = geo.Open(dir+"/data/GeoLite2-Country.mmdb",dir+"/data/GeoLite2-City.mmdb",dir+"/data/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatalln("-init-geoReader-err:-"+err.Error())
	}
}

func GetIpInfo(remoteIP string) (IpInfo,error){
	if err != nil {
		return IpInfo{}, err
	}
	ip := net.ParseIP(remoteIP)
	ipDecimal := util.ToDecimal(ip)
	country, _ := Reader.Country(ip)
	city, _ := Reader.City(ip)
	asn, _ := Reader.ASN(ip)
	hostname, _ := util.LookupAddr(ip)
	var autonomousSystemNumber string
	if asn.AutonomousSystemNumber > 0 {
		autonomousSystemNumber = fmt.Sprintf("AS%d", asn.AutonomousSystemNumber)
	}
	return IpInfo{
		IP:         ip,
		IPDecimal:  ipDecimal,
		Country:    country.Name,
		CountryISO: country.ISO,
		CountryEU:  country.IsEU,
		City:       city.Name,
		Hostname:   hostname,
		Latitude:   city.Latitude,
		Longitude:  city.Longitude,
		ASN:        autonomousSystemNumber,
		ASNOrg:     asn.AutonomousSystemOrganization,
	}, nil
}
