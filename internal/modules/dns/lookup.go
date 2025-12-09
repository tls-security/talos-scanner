package dns

import (
	"net"
)

// Analyze faz o lookup de IP e MX records
func Analyze(hostname string) map[string]interface{} {
	result := make(map[string]interface{})

	// 1. Busca IPs (A Records)
	// Usamos o resolvedor padrão do sistema, que é mais simples e não precisa da var 'r'
	ips, err := net.LookupIP(hostname)
	if err == nil && len(ips) > 0 {
		// Converte IPs para string
		var ipStrings []string
		for _, ip := range ips {
			ipStrings = append(ipStrings, ip.String())
		}
		
		result["ip"] = ipStrings[0]
		result["all_ips"] = ipStrings
	} else {
		result["error"] = "Não foi possível resolver o DNS"
	}

	// 2. Busca MX (Email)
	mx, err := net.LookupMX(hostname)
	if err == nil {
		var mxList []string
		for _, m := range mx {
			mxList = append(mxList, m.Host)
		}
		result["mx_records"] = mxList
	}

	return result
}