package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Загрузка сертификатов из Docker-тома /app/ssl/
	serverCert, err := tls.LoadX509KeyPair(
		"/app/ssl/server.crt",
		"/app/ssl/server.key",
	)
	if err != nil {
		log.Fatalf("Failed to load server certificates: %v", err)
	}

	// 2. Загрузка CA сертификата
	caCert, err := os.ReadFile("/app/ssl/ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to parse CA certificate")
	}

	// 3. Настройка TLS с Two-Way аутентификацией
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Требуем клиентский сертификат
		MinVersion:   tls.VersionTLS12,
	}

	// 4. Создание HTTP сервера
	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
		Handler:   enableCORS(setupRoutes()),
	}

	// 5. Логирование старта
	log.Println("Starting HTTPS server on :8080 with mTLS...")
	log.Printf("Using certificates from: %v", listFiles("/app/ssl"))

	// 6. Запуск сервера
	err = server.ListenAndServeTLS("", "") // Сертификаты уже в TLSConfig
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Проверка клиентского сертификата
		if len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "Client certificate required", http.StatusForbidden)
			return
		}

		clientCN := r.TLS.PeerCertificates[0].Subject.CommonName
		log.Printf("Request from client: %s", clientCN)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Secure data","client":"` + clientCN + `"}`))
	})

	return mux
}

// Вспомогательная функция для логгирования файлов
func listFiles(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return []string{"error reading directory"}
	}

	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return fileNames
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
