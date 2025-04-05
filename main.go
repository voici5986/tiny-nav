package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed public/*
var content embed.FS

// token 过期时间
const (
	defaultExpireTime = 30 * 24 * time.Hour
	defaulttokenCount = 10
)

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

// Token结构体用于存储token及其过期时间
type Token struct {
	Value    string
	ExpireAt time.Time
}

// TokenStore结构体用于管理token存储
type TokenStore struct {
	tokens map[string]Token
	mu     sync.Mutex
}

// NewTokenStore创建一个新的TokenStore
func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]Token),
	}
}

// AddToken添加一个新的token到存储中，如果超过数量则删除最早的token
func (ts *TokenStore) AddToken(token string, duration time.Duration) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.tokens) >= defaulttokenCount {
		oldestToken := ""
		oldestTime := time.Now()
		for k, v := range ts.tokens {
			if v.ExpireAt.Before(oldestTime) {
				oldestToken = k
				oldestTime = v.ExpireAt
			}
		}
		delete(ts.tokens, oldestToken)
	}

	expireAt := time.Now().Add(duration)
	ts.tokens[token] = Token{Value: token, ExpireAt: expireAt}
}

// ValidateToken验证token是否有效
func (ts *TokenStore) ValidateToken(token string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t, exists := ts.tokens[token]
	if !exists {
		return false
	}
	if time.Now().After(t.ExpireAt) {
		delete(ts.tokens, token)
		return false
	}
	newExpireAt := time.Now().Add(defaultExpireTime)
	ts.tokens[token] = Token{Value: token, ExpireAt: newExpireAt}
	return true
}

var tokenStore = NewTokenStore()

// 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
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
		// 生成新的令牌
		token, err := generateToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tokenStore.AddToken(token, defaultExpireTime)
		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// 中间件函数验证令牌
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !tokenStore.ValidateToken(token) {
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
	file, err := content.Open("public/index.html")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, file)
}

// 记录访问日志的中间件
func logAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 创建一个响应记录器来捕获状态码
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("%s - %s %s %d %v", r.RemoteAddr, r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// 自定义响应写入器，用于记录状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// 调试接口，输出所有的 token
func debugTokensHandler(w http.ResponseWriter, r *http.Request) {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	var tokens []string
	for _, token := range tokenStore.tokens {
		tokens = append(tokens, token.Value)
	}

	log.Printf("Current tokens: %v", tokens)
	w.WriteHeader(http.StatusOK)
}

func getIconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("Authorization")
	if !tokenStore.ValidateToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	iconURL := fmt.Sprintf("%s://%s/favicon.ico", parsedURL.Scheme, parsedURL.Host)
	resp, err := http.Get(iconURL)
	log.Printf("Fetching icon from URL:%s code:%d Content-Type:%s", iconURL, resp.StatusCode, resp.Header.Get("Content-Type"))
	if err != nil || resp.StatusCode != http.StatusOK ||
		resp.Header.Get("Content-Type") != "image/x-icon" ||
		resp.Header.Get("Content-Type") != "image/vnd.microsoft.icon" {

		// 尝试解析 HTML 来获取图标
		iconURL, err = fetchIconFromHTML(urlParam)
		if err != nil {
			log.Printf("Failed to fetch icon: %v", err)
			http.Error(w, "Failed to fetch icon", http.StatusInternalServerError)
			return
		}
		log.Printf("Fetching icon from html URL:%s", iconURL)
		resp, err = http.Get(iconURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to fetch icon from HTML: %v", err)
			http.Error(w, "Failed to fetch icon", http.StatusInternalServerError)
			return
		}
	}
	defer resp.Body.Close()

	// 创建 cache 目录
	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		err := os.Mkdir("cache", 0755)
		if err != nil {
			log.Printf("Failed to create cache directory: %v", err)
			http.Error(w, "Failed to create cache directory", http.StatusInternalServerError)
			return
		}
	}

	ext := getFileExtension(iconURL, resp)
	// 生成唯一的文件名
	fileName := fmt.Sprintf("%s%s", base64.URLEncoding.EncodeToString([]byte(parsedURL.Host)), ext)
	filePath := filepath.Join("cache", fileName)

	// 保存图标到 cache 目录
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Failed to save icon: %v", err)
		http.Error(w, "Failed to save icon", http.StatusInternalServerError)
		return
	}

	// 返回图标 URL
	iconResponse := map[string]string{
		"iconUrl": fmt.Sprintf("/cache/%s", fileName),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iconResponse)
}

func fetchIconFromHTML(pageURL string) (string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch page, status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var iconURL string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, attr := range n.Attr {
				if attr.Key == "rel" && (attr.Val == "icon" || attr.Val == "shortcut icon") {
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							iconURL = attr.Val
							log.Printf("Icon URL found: %s", iconURL)
							if !strings.HasPrefix(iconURL, "http") && !strings.HasPrefix(iconURL, "//") {
								baseURL, _ := url.Parse(pageURL)
								iconURL = baseURL.ResolveReference(&url.URL{Path: iconURL}).String()
							}
							if strings.HasPrefix(iconURL, "//") {
								baseURL, _ := url.Parse(pageURL)
								iconURL = baseURL.Scheme + ":" + iconURL
							}
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if iconURL == "" {
		return "", fmt.Errorf("icon not found in HTML")
	}

	return iconURL, nil
}

func getFileExtension(iconURL string, resp *http.Response) string {
	ext := filepath.Ext(iconURL)
	if ext == "" {
		// 尝试从 Content-Type 获取文件扩展名
		contentType := resp.Header.Get("Content-Type")
		if contentType != "" {
			exts, _ := mime.ExtensionsByType(contentType)
			if len(exts) > 0 {
				ext = exts[0]
			}
		}
	}
	if ext == "" {
		// 默认扩展名
		ext = ".ico"
	}
	return ext
}

func main() {
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/navigation", authMiddleware(getNavigationHandler))
	mux.HandleFunc("/navigation/add", authMiddleware(addLinkHandler))
	mux.HandleFunc("/navigation/update/", authMiddleware(updateLinkHandler))
	mux.HandleFunc("/navigation/delete/", authMiddleware(deleteLinkHandler))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/debug/tokens", debugTokensHandler)
	mux.HandleFunc("/get-icon", authMiddleware(getIconHandler))
	mux.Handle("/cache/", http.StripPrefix("/cache/", http.FileServer(http.Dir("cache"))))
	mux.Handle("/public/", http.FileServer(http.FS(content)))

	// 使用日志中间件包装 mux
	logHandler := logAccessMiddleware(mux)

	log.Printf("Server is running on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, logHandler)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
