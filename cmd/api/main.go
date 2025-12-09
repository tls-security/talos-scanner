package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/tls-security/talos-scanner/internal/scanner"
    "github.com/tls-security/talos-scanner/internal/models"
    "net/http"
)

func main() {
    r := gin.Default()

    // Configuração de CORS (Para o Frontend na Vercel poder chamar depois)
    r.Use(cors.Default())

    // Servir o Frontend Provisório (Arquivos estáticos)
    r.Static("/public", "./public")
    r.StaticFile("/", "./public/index.html") // Raiz abre o HTML

    // API Endpoint
    r.POST("/analyze", func(c *gin.Context) {
        var req models.RequestPayload
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
            return
        }

        // Chama o orquestrador
        result := scanner.ScanURL(req.URL)
        
        c.JSON(http.StatusOK, result)
    })

    // Roda na porta definida pelo Fly.io ou 8080 local
    r.Run() 
}