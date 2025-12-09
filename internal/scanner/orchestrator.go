package scanner

import (
	"net/url"
	"sync"
	"time"
    "strings"

	"github.com/tls-security/talos-scanner/internal/models"
	"github.com/tls-security/talos-scanner/internal/modules/sandbox"
	"github.com/tls-security/talos-scanner/internal/modules/dns"
	"github.com/tls-security/talos-scanner/internal/modules/ssl"
)

func ScanURL(targetURL string) models.FullReport {
	start := time.Now()
	var wg sync.WaitGroup

	// Normaliza URL para extrair hostname
	parsed, _ := url.Parse(targetURL)
	hostname := parsed.Hostname()
    if hostname == "" {
        // Fallback se o user não digitou http://
        if !strings.HasPrefix(targetURL, "http") {
             targetURL = "http://" + targetURL
             parsed, _ = url.Parse(targetURL)
             hostname = parsed.Hostname()
        }
    }

	report := models.FullReport{URL: targetURL}

	// --- Tarefa 1: Sandbox (Browser) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		screenshot, err := sandbox.TakeScreenshot(targetURL)
		if err == nil {
			report.Screenshot = screenshot
		} else {
            // AGORA VAMOS VER O ERRO NO JSON
			report.Screenshot = "ERRO: " + err.Error() 
		}
	}()

	// --- Tarefa 2: DNS ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		report.DNS = dns.Analyze(hostname)
	}()

    // --- Tarefa 3: SSL ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		report.SSL = ssl.Check(hostname)
	}()

	wg.Wait()

	report.ProcessingTime = time.Since(start).String()
    
    // Simples lógica de veredito baseada no que achamos
    report.RiskScore = 0
    report.Verdict = "SEGURO"
    
    if report.Screenshot != "" && strings.HasPrefix(report.Screenshot, "ERRO") {
         report.Verdict = "ERRO NA ANÁLISE"
    }

	return report
}