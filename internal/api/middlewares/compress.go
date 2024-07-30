package middlewares

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares/compress"
	"net/http"
	"strings"
)

var whiteList = map[string]int{
	"application/json":   1,
	"application/x-gzip": 2,
	"text/html":          3,
}

func Compress(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Функция сжатия должна работать для контента с типами ContentTypeJson, ContentTypeGzip или ContentTypeText
		if !isTargetContentType(r.Header.Get("Content-Type")) {
			next.ServeHTTP(w, r)
			return
		}
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		originalWriter := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			gzWriter := compress.NewGzipWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			originalWriter = gzWriter
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer gzWriter.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			gzReader, err := compress.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = gzReader
			defer gzReader.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(originalWriter, r)
	}

	return http.HandlerFunc(fn)
}

func isTargetContentType(contentType string) bool {
	_, ok := whiteList[contentType]

	return ok
}
