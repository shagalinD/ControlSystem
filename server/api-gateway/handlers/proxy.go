package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type ProxyHandler struct {
    AuthServiceURL      string
    ProjectDefectServiceURL string
    ContentServiceURL   string
    Client             *resty.Client
}

func NewProxyHandler(authURL, projectDefectURL, contentURL string) *ProxyHandler {
    return &ProxyHandler{
        AuthServiceURL:      authURL,
        ProjectDefectServiceURL: projectDefectURL,
        ContentServiceURL:   contentURL,
        Client:             resty.New(),
    }
}

func convertHeaders(header http.Header) map[string]string {
    result := make(map[string]string)
    for key, values := range header {
        // Объединяем множественные значения через запятую
        result[key] = strings.Join(values, ", ")
    }
    return result
}
func filterHeaders(header http.Header) http.Header {
    filtered := make(http.Header)
    
    // Копируем только нужные заголовки
    for key, values := range header {
        // Пропускаем заголовки, которые не нужно передавать
        if key == "Connection" || key == "Keep-Alive" {
            continue
        }
        filtered[key] = values
    }
    
    return filtered
}

func (h *ProxyHandler) ProxyRequest(serviceURL string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Создаем URL для целевого сервиса
        targetURL := serviceURL + c.Request.URL.Path

				 // Добавляем query параметры если есть
        if c.Request.URL.RawQuery != "" {
            targetURL += "?" + c.Request.URL.RawQuery
        }
        
        // Копируем тело запроса
        var bodyBytes []byte
        if c.Request.Body != nil {
            bodyBytes, _ = io.ReadAll(c.Request.Body)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
        }
        
				// Фильтруем заголовки
        filteredHeaders := filterHeaders(c.Request.Header)

        // Выполняем запрос к целевому сервису
        resp, err := h.Client.R().
            SetHeaders(convertHeaders(filteredHeaders)).
            SetBody(bodyBytes).
            Execute(c.Request.Method, targetURL)
        
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   "Service unavailable: " + err.Error(),
            })
            return
        }
        
        // Копируем заголовки ответа
        for key, values := range resp.Header() {
            for _, value := range values {
                c.Writer.Header().Add(key, value)
            }
        }
        
        // Устанавливаем статус и отправляем тело ответа
        c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
    }
}

// Специальные обработчики для разных сервисов
func (h *ProxyHandler) AuthProxy() gin.HandlerFunc {
    return h.ProxyRequest(h.AuthServiceURL)
}

func (h *ProxyHandler) ProjectDefectProxy() gin.HandlerFunc {
    return h.ProxyRequest(h.ProjectDefectServiceURL)
}

func (h *ProxyHandler) ContentProxy() gin.HandlerFunc {
    return h.ProxyRequest(h.ContentServiceURL)
}