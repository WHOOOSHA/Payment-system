package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

/*
	В задании сказано:
	"Необходимо разработать приложение, реализующее систему обработки транзакций платёжной системы".
	Исходя из этого сделан вывод, что разрабатываемое приложение является
	одним сервисом в микросервисной архитектуре, это означает, что запросы,
	которые приходят на данный сервер, должны исходить от другого сервера.

	Чтобы проверить источник запроса используется HMAC подпись, записанная в
	заголовке X-Signature. Подпись вычисляется с использованием алгоритма SHA-256
	на основании тела запроса и секретного ключа.

	Данный midlleware проверяет	подлиность подписи, повторяя вычисления,
	которые совершил клиентский сервис. Если подпись полученная от клиента
	совпадает с вычисленной подписью в middleware, то запрос пришел от сервиса
	и обрабатывается дальше, в ином случае, возвращается 401.

	Как получить подпись см. в Readme.me
*/

func HMACAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signature := r.Header.Get("X-Signature")
			if signature == "" {
				http.Error(w, "Missing signature", http.StatusUnauthorized)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(bodyBytes)
			expectedSig := hex.EncodeToString(mac.Sum(nil))

			if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
