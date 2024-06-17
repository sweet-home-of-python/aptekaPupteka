package basicAuth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

var Users = map[string]string{  // пользователи и пароли
	"yar": "asslove",
	"dim": "ass",
}

func BasicAuth(next http.Handler) http.Handler { // кастомная функция для авторизаци
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization") // получаем header авторизации
		if auth == "" { // если нет авторизации
			w.WriteHeader(http.StatusUnauthorized) // возвращаем ошибку
			w.Write([]byte("Unauthorized")) 
			return
		}

		parts := strings.SplitN(auth, " ", 2) // разбиваем на 2 части
		if len(parts) != 2 || parts[0] != "Basic" { // если нет Basic
			w.WriteHeader(http.StatusUnauthorized) // возвращаем ошибку
			w.Write([]byte("Unauthorized"))
			return
		}

		payload, err := base64.StdEncoding.DecodeString(parts[1]) // декодируем
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		pair := strings.SplitN(string(payload), ":", 2) // разбиваем на 2 части
		if len(pair) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		user, pass := pair[0], pair[1]

		if password, ok := Users[user]; !ok || password != pass { // если пароли не совпадают
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r) // передаем контроль следующему обработчику middleware
	})
}
