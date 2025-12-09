package dns

import (
	"net"
	"time"
)

// Analyze faz o lookup de IP e MX records
func Analyze(hostname string) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Define um timeout para não travar
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 2 * time.Second,
			}
			return d.DialContext(ctx, network, address)
		},
	}
    // (Nota: Para simplificar, vamos usar o default resolver do sistema no container)
    
	// 1. Busca IPs (A Records)
	ips, err := net.LookupIP(hostname)
	if err == nil && len(ips) > 0 {
		result["ip"] = ips[0].String()
		result["all_ips"] = ips
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