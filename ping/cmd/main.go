package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"sync"

	"github.com/leejansq/pkg/ping"
)

func main() {
	ip, _, err := net.ParseCIDR(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip.String())
	ip = ip.To4()
	// ip[3] = 1
	// fmt.Printf(ip.String())
	// return
	resOk := []string{}
	var group sync.WaitGroup
	for i := byte(1); i < 255; i++ {
		group.Add(1)
		go func(j byte) {
			tip := net.ParseIP("0.0.0.0")
			tip = tip.To4()
			copy(tip, ip)
			// tip[0] = ip[0]
			// tip[1] = ip[1]
			// tip[2] = ip[2]
			tip[3] = j
			//tip = tip.To4()
			//tip[3] = i
			log.Println("ping>>:", tip.String())
			if ok, err := ping.Ping(tip.String(), 5); ok && err == nil {
				resOk = append(resOk, tip.String())
			}
			group.Done()
		}(i)

	}
	group.Wait()
	fmt.Println("======================================")
	fmt.Println(resOk, len(resOk))
}
