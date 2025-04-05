package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var content embed.FS

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Link struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Icon     string `json:"icon"`
	Category string `json:"category"`
}

type Navigation struct {
	Links []Link `json:"links"`
}

func loadNavigation() (Navigation, error) {
	data, err := os.ReadFile("navigation.json")
	if err != nil {
		if !os.IsNotExist(err) {
			return Navigation{}, err
		}
		return Navigation{}, nil
	}
	var nav Navigation
	err = json.Unmarshal(data, &nav)
	if err != nil {
		return Navigation{}, err
	}
	return nav, nil
}

func saveNavigation(nav Navigation) error {
	data, err := json.MarshalIndent(nav, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("navigation.json", data, 0644)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	envUsername := os.Getenv("NAV_USERNAME")
	envPassword := os.Getenv("NAV_PASSWORD")
	if user.Username == envUsername && user.Password == envPassword {
		// 简单返回一个令牌
		token := "valid_token"
		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// 中间件函数验证令牌
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedToken := r.Header.Get("Authorization")
		if providedToken != "valid_token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func getNavigationHandler(w http.ResponseWriter, r *http.Request) {
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(nav, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func addLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var newLink Link
	err := json.NewDecoder(r.Body).Decode(&newLink)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	nav.Links = append(nav.Links, newLink)
	err = saveNavigation(nav)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var updatedLink Link
	err := json.NewDecoder(r.Body).Decode(&updatedLink)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var index int
	fmt.Sscanf(r.URL.Path, "/navigation/update/%d", &index)
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if index < 0 || index >= len(nav.Links) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}
	nav.Links[index] = updatedLink
	err = saveNavigation(nav)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var index int
	fmt.Sscanf(r.URL.Path, "/navigation/delete/%d", &index)
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if index < 0 || index >= len(nav.Links) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}
	nav.Links = append(nav.Links[:index], nav.Links[index+1:]...)
	err = saveNavigation(nav)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := content.Open("templates/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, file)
}

func main() {
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/navigation", authMiddleware(getNavigationHandler))
	http.HandleFunc("/navigation/add", authMiddleware(addLinkHandler))
	http.HandleFunc("/navigation/update/", authMiddleware(updateLinkHandler))
	http.HandleFunc("/navigation/delete/", authMiddleware(deleteLinkHandler))
	http.HandleFunc("/", indexHandler)
	http.Handle("/templates/", http.FileServer(http.FS(content)))

	log.Printf("Server is running on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
