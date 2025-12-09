package sandbox

import (
    "context"
    "encoding/base64"
    "time"
    "github.com/chromedp/chromedp"
)

func TakeScreenshot(urlStr string) (string, error) {
    // Cria um contexto com timeout (segurança para não travar o servidor)
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // Configura o Chrome para rodar Headless (sem janela) e compatível com Docker
    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("no-sandbox", true), // Obrigatório para Docker/Fly.io
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("ignore-certificate-errors", true),
    )

    allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
    defer cancelAlloc()

    taskCtx, cancelTask := chromedp.NewContext(allocCtx)
    defer cancelTask()

    var buf []byte
    
    // A automação em si
    err := chromedp.Run(taskCtx,
        chromedp.Navigate(urlStr),
        chromedp.Sleep(2*time.Second), // Espera carregar um pouco
        chromedp.FullScreenshot(&buf, 90),
    )

    if err != nil {
        return "", err
    }

    // Converte para Base64 para enviar no JSON
    return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf), nil
}