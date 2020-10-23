package geo

import (
	"github.com/Zfzlv/ipInfo/reader"
	"math"
	"net"
)

type Reader interface {
	Country(net.IP, string) (Country, error)
	City(net.IP, string) (City, error)
	ASN(net.IP) (ASN, error)
	IsEmpty() bool
}

type Country struct {
	Continent           string
	ContinentCode       string
	Name                string
	ISO                 string
	IsEU                *bool
	IsAnonymousProxy    bool
	IsSatelliteProvider bool
}

type City struct {
	Province            string
	ProvinceCode        string
	Name                string
	Latitude            float64
	Longitude           float64
	TimeZone            string
	IsAnonymousProxy    bool
	IsSatelliteProvider bool
}

type ASN struct {
	AutonomousSystemNumber       uint
	AutonomousSystemOrganization string
}

type geoip struct {
	country *reader.Reader
	city    *reader.Reader
	asn     *reader.Reader
}

func Open(countryDB, cityDB string, asnDB string) (Reader, error) {
	var country, city, asn *reader.Reader
	if countryDB != "" {
		r, err := reader.Open(countryDB)
		if err != nil {
			return nil, err
		}
		country = r
	}
	if cityDB != "" {
		r, err := reader.Open(cityDB)
		if err != nil {
			return nil, err
		}
		city = r
	}
	if asnDB != "" {
		r, err := reader.Open(asnDB)
		if err != nil {
			return nil, err
		}
		asn = r
	}
	return &geoip{country: country, city: city, asn: asn}, nil
}

func (g *geoip) Country(ip net.IP, lang string) (Country, error) {
	country := Country{}
	if g.country == nil {
		return country, nil
	}
	record, err := g.country.Country(ip)
	if err != nil {
		return country, err
	}
	if c, exists := record.Country.Names[lang]; exists {
		country.Name = c
	}
	if c, exists := record.RegisteredCountry.Names[lang]; exists && country.Name == "" {
		country.Name = c
	}
	if record.Country.IsoCode != "" {
		country.ISO = record.Country.IsoCode
	}
	if record.RegisteredCountry.IsoCode != "" && country.ISO == "" {
		country.ISO = record.RegisteredCountry.IsoCode
	}
	isEU := record.Country.IsInEuropeanUnion || record.RegisteredCountry.IsInEuropeanUnion
	country.IsEU = &isEU
	if c, exists := record.Continent.Names[lang]; exists {
		country.Continent = c
		country.ContinentCode = record.Continent.Code
	}
	country.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	country.IsSatelliteProvider = record.Traits.IsSatelliteProvider
	return country, nil
}

func (g *geoip) City(ip net.IP, lang string) (City, error) {
	city := City{}
	if g.city == nil {
		return city, nil
	}
	record, err := g.city.City(ip)
	if err != nil {
		return city, err
	}
	if c, exists := record.City.Names[lang]; exists {
		city.Name = c
	}
	if !math.IsNaN(record.Location.Latitude) {
		city.Latitude = record.Location.Latitude
	}
	if !math.IsNaN(record.Location.Longitude) {
		city.Longitude = record.Location.Longitude
	}
	city.TimeZone = record.Location.TimeZone
	if len(record.Subdivisions) > 0 {
		if c, exists := record.Subdivisions[0].Names[lang]; exists {
			city.Province = c
			city.ProvinceCode = record.Subdivisions[0].IsoCode
		}
	}
	city.IsAnonymousProxy = record.Traits.IsAnonymousProxy
	city.IsSatelliteProvider = record.Traits.IsSatelliteProvider
	return city, nil
}

func (g *geoip) ASN(ip net.IP) (ASN, error) {
	asn := ASN{}
	if g.asn == nil {
		return asn, nil
	}
	record, err := g.asn.ASN(ip)
	if err != nil {
		return asn, err
	}
	if record.AutonomousSystemNumber > 0 {
		asn.AutonomousSystemNumber = record.AutonomousSystemNumber
	}
	if record.AutonomousSystemOrganization != "" {
		asn.AutonomousSystemOrganization = record.AutonomousSystemOrganization
	}
	return asn, nil
}

func (g *geoip) IsEmpty() bool {
	return g.country == nil && g.city == nil
}
