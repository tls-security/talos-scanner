package sandbox

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/chromedp/chromedp"
)

func TakeScreenshot(urlStr string) (string, error) {
	// AUMENTADO: Timeout de 15s para 60s (O plano Free do Render precisa de tempo)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Configurações Otimizadas para Docker/Render
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
		// NOVAS FLAGS IMPORTANTE PARA DOCKER:
		chromedp.Flag("disable-dev-shm-usage", true), // Evita crash por falta de memória partilhada
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.WindowSize(1280, 720), // Tamanho fixo para economizar RAM
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	var buf []byte

	// A automação
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(urlStr),
		// Reduzimos o sleep para 1 segundo (se carregar antes, ele avança)
		chromedp.Sleep(1*time.Second), 
		chromedp.FullScreenshot(&buf, 80), // Qualidade 80% para ficar mais leve
	)

	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf), nil
}