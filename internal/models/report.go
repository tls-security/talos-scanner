package models

// RequestPayload é o que recebemos do Frontend
type RequestPayload struct {
    URL string `json:"url"`
}

// FullReport é o JSON final que devolvemos
type FullReport struct {
    URL           string                 `json:"url"`
    RiskScore     int                    `json:"risk_score"`
    Verdict       string                 `json:"verdict"`
    DNS           map[string]interface{} `json:"dns_data"`
    SSL           map[string]interface{} `json:"ssl_data"`
    Reputation    map[string]interface{} `json:"reputation"`
    Screenshot    string                 `json:"screenshot_base64"` // A string base64
    ProcessingTime string                `json:"processing_time"`
}