package util

import (
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
)

//find hostName
func LookupAddr(ip net.IP) (string, error) {
	names, err := net.LookupAddr(ip.String())
	if err != nil || len(names) == 0 {
		return "", err
	}
	return strings.TrimRight(names[0], "."), nil
}

//check port
func LookupPort(ip net.IP, port uint64) error {
	address := fmt.Sprintf("[%s]:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

//ip to bigInt
func ToDecimal(ip net.IP) *big.Int {
	i := big.NewInt(0)
	if to4 := ip.To4(); to4 != nil {
		i.SetBytes(to4)
	} else {
		i.SetBytes(ip)
	}
	return i
}
