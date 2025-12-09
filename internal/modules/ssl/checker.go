package ssl

import (
	"crypto/tls"
	"net"
	"time"
)

func Check(hostname string) map[string]interface{} {
	result := make(map[string]interface{})
	result["is_valid"] = false

	conf := &tls.Config{
		InsecureSkipVerify: true, // Aceitamos conectar para ler o erro depois
	}

	// Tenta conectar na porta 443 com timeout de 3s
	conn, err := net.DialTimeout("tcp", hostname+":443", 3*time.Second)
	if err != nil {
		result["error"] = "Porta 443 fechada ou timeout"
		return result
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, conf)
	err = tlsConn.Handshake()
	if err != nil {
		result["error"] = "Falha no Handshake SSL"
		return result
	}

	// Pega o certificado
	certs := tlsConn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		cert := certs[0]
		result["is_valid"] = true
		result["issuer"] = cert.Issuer.CommonName
		result["subject"] = cert.Subject.CommonName
		result["expires"] = cert.NotAfter.Format("2006-01-02")
		result["days_left"] = int(time.Until(cert.NotAfter).Hours() / 24)
	}

	return result
}