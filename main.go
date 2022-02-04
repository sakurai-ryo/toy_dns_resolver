package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

const RootNameServer = "198.41.0.4"

func getAnswer(reply *dns.Msg) []net.IP {
	if len(reply.Answer) == 0 {
		return nil
	}

	ans := make([]net.IP, len(reply.Answer))
	for i, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			ans[i] = record.(*dns.A).A
		}
	}
	return ans
}

// nameserverに登録されたAレコードを返す
func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

// 次の問い合わせ先のドメインを返す
func getNS(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("  ", record)
			return record.(*dns.NS).Ns
		}
	}
	return ""
}

func dnsQuery(name string, server net.IP) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server.String(), name)
	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	c := new(dns.Client)
	reply, _, _ := c.Exchange(msg, server.String()+":53")
	return reply
}

// Answerセクション => Glueレコード => NSレコードの順に問い合わせ
func resolve(name string) []net.IP {
	nameserver := net.ParseIP(RootNameServer)

	for {
		reply := dnsQuery(name, nameserver)

		if ip := getAnswer(reply); ip != nil {
			return ip
		} else if nsIP := getGlue(reply); nsIP != nil {
			nameserver = nsIP
		} else if domain := getNS(reply); domain != "" {
			nameserver = resolve(domain)[0]
		} else {
			panic("DNS failed")
		}
	}
}

// フルサービスリゾルバ
func main() {
	name := os.Args[1]
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}

	fmt.Println("Result:", resolve(name))
}
