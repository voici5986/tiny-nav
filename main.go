package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mat/besticon/v3/besticon"
	"gopkg.in/ini.v1"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:embed public/*
var embeddedFiles embed.FS

const (
	defaultExpireTime  = 30 * 24 * time.Hour // token 过期时间
	defaulttokenCount  = 10                  // 最多存储的 token 数量
	dataDir            = "data"
	tokenFileName      = "tokens.json"
	navigationFileName = "navigation.json"
	configFileName     = "config.ini"
)

var tokenStore *TokenStore

var envPort string           // 监听端口
var envUsername string       // 登录用户名
var envPassword string       // 登录密码
var envEnableNoAuth bool     // 是否启用无用户密码模式
var envEnableNoAuthView bool // 是否启用无用户密码浏览模式

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Link struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	Icon      string `json:"icon"`
	Category  string `json:"category"`
	SortIndex int    `json:"sortIndex"`
}

type Navigation struct {
	Links        []Link   `json:"links"`
	Categories   []string `json:"categories"`
	LastModified int64    `json:"lastModified"`
}

type NavigationLastModified struct {
	LastModified int64 `json:"lastModified"`
}

type Config struct {
	EnableNoAuth     bool `json:"enableNoAuth"`
	EnableNoAuthView bool `json:"enableNoAuthView"`
}

func loadConfig() {
	// Parse command line flags
	port := flag.String("port", "", "Port to listen on (e.g. 58080)")
	user := flag.String("user", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")
	noAuth := flag.Bool("no-auth", false, "Enable no-auth mode")
	noAuthView := flag.Bool("no-auth-view", false, "Enable no-auth-view mode")
	flag.Parse()

	// 确保 data 目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Errorf("failed to create data directory: %v", err)
	}

	configPath := filepath.Join(dataDir, configFileName)

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	// Priority: Command line args > Environment variables > Config file
	envPort = os.Getenv("LISTEN_PORT")
	if envPort == "" {
		if *port != "" {
			envPort = *port
		} else if cfg != nil {
			envPort = cfg.Section("").Key("LISTEN_PORT").MustString("58080")
		} else {
			envPort = "58080"
		}
	}

	envUsername = os.Getenv("NAV_USERNAME")
	if envUsername == "" {
		if *user != "" {
			envUsername = *user
		} else if cfg != nil {
			envUsername = cfg.Section("").Key("NAV_USERNAME").String()
		}
	}

	envPassword = os.Getenv("NAV_PASSWORD")
	if envPassword == "" {
		if *password != "" {
			envPassword = *password
		} else if cfg != nil {
			envPassword = cfg.Section("").Key("NAV_PASSWORD").String()
		}
	}

	noAuthStr := os.Getenv("ENABLE_NO_AUTH")
	if cfg != nil {
		noAuthStr = cfg.Section("").Key("ENABLE_NO_AUTH").MustString("false")
	}
	if *noAuth {
		envEnableNoAuth = true
	} else {
		envEnableNoAuth = noAuthStr == "true"
	}

	noAuthViewStr := os.Getenv("ENABLE_NO_AUTH_VIEW")
	if cfg != nil {
		noAuthViewStr = cfg.Section("").Key("ENABLE_NO_AUTH_VIEW").MustString("false")
	}
	if *noAuthView {
		envEnableNoAuthView = true
	} else {
		envEnableNoAuthView = noAuthViewStr == "true"
	}

	log.Printf("Config loaded: LISTEN_PORT=%s, NAV_USERNAME=%s, ENABLE_NO_AUTH=%v, ENABLE_NO_AUTH_VIEW=%v", envPort, envUsername, envEnableNoAuth, envEnableNoAuthView)
}

func loadNavigation() (Navigation, error) {
	// 确保 data 目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return Navigation{}, fmt.Errorf("failed to create data directory: %v", err)
	}

	navPath := filepath.Join(dataDir, navigationFileName)
	data, err := os.ReadFile(navPath)
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
	// 确保 data 目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	currentTime := time.Now()
	lastModified := currentTime.UnixNano() / int64(time.Millisecond)
	nav.LastModified = lastModified

	data, err := json.MarshalIndent(nav, "", "  ")
	if err != nil {
		return err
	}

	navPath := filepath.Join(dataDir, navigationFileName)
	return os.WriteFile(navPath, data, 0644)
}

// Token结构体用于存储token及其过期时间
type Token struct {
	Value    string
	ExpireAt time.Time
}

// TokenStore结构体用于管理token存储
type TokenStore struct {
	tokens   map[string]Token
	mu       sync.Mutex
	filePath string
}

func NewTokenStore() *TokenStore {
	// 确保data目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Failed to create data directory: %v", err)
	}

	ts := &TokenStore{
		tokens:   make(map[string]Token),
		filePath: filepath.Join(dataDir, tokenFileName),
	}

	// 从文件加载现有token
	ts.loadTokens()

	return ts
}

// 保存tokens到文件
func (ts *TokenStore) saveTokens() error {
	// 清理过期的token
	now := time.Now()
	for k, v := range ts.tokens {
		if now.After(v.ExpireAt) {
			delete(ts.tokens, k)
		}
	}

	// 将tokens转换为可序列化的格式
	data, err := json.MarshalIndent(ts.tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tokens: %v", err)
	}

	// 写入文件
	err = os.WriteFile(ts.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing tokens file: %v", err)
	}

	return nil
}

// 从文件加载tokens
func (ts *TokenStore) loadTokens() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	data, err := os.ReadFile(ts.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading tokens file: %v", err)
		}
		return
	}

	err = json.Unmarshal(data, &ts.tokens)
	if err != nil {
		log.Printf("Error unmarshaling tokens: %v", err)
		ts.tokens = make(map[string]Token) // 如果出错，使用空map
	}

	// 清理已过期的token
	now := time.Now()
	for k, v := range ts.tokens {
		if now.After(v.ExpireAt) {
			delete(ts.tokens, k)
		}
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

	// 保存到文件
	if err := ts.saveTokens(); err != nil {
		log.Printf("Error saving tokens: %v", err)
	}
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
	// 保存到文件
	if err := ts.saveTokens(); err != nil {
		log.Printf("Error saving tokens: %v", err)
	}
	return true
}

// 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// updateCategories 更新导航的分类列表，保持原有顺序，删除不存在的分类，添加新的分类
func updateCategories(nav *Navigation) {
	// 创建当前链接中存在的分类集合
	currentCategories := make(map[string]struct{})
	for _, link := range nav.Links {
		if link.Category != "" {
			currentCategories[link.Category] = struct{}{}
		}
	}

	// 如果 Categories 为空，初始化它
	if nav.Categories == nil {
		nav.Categories = make([]string, 0)
	}

	// 删除不存在的分类（保持原有顺序）
	newCategories := make([]string, 0, len(nav.Categories))
	for _, category := range nav.Categories {
		if _, exists := currentCategories[category]; exists {
			newCategories = append(newCategories, category)
			delete(currentCategories, category) // 从当前分类集合中删除已处理的分类
		}
	}

	// 添加新的分类（将剩余的分类追加到列表末尾）
	for category := range currentCategories {
		newCategories = append(newCategories, category)
	}

	nav.Categories = newCategories
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

	authCheck := false
	if envEnableNoAuth {
		log.Println("No authentication required (ENABLE_NO_AUTH=true)")
		authCheck = true
	} else {
		if user.Username == envUsername && user.Password == envPassword {
			authCheck = true
		}
	}
	if !authCheck {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	// 生成新的令牌
	token, err := generateToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tokenStore.AddToken(token, defaultExpireTime)
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
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

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := Config{
		EnableNoAuth:     envEnableNoAuth,
		EnableNoAuthView: envEnableNoAuthView,
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
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

func getNavigationLastModifiedHandler(w http.ResponseWriter, r *http.Request) {
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := NavigationLastModified{
		LastModified: nav.LastModified,
	}
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
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
	if len(newLink.Url) == 0 {
		http.Error(w, "Url required", http.StatusBadRequest)
		return
	}
	if len(newLink.Category) == 0 {
		http.Error(w, "Category required", http.StatusBadRequest)
		return
	}
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	nav.Links = append(nav.Links, newLink)
	updateCategories(&nav)
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
	if len(updatedLink.Url) == 0 {
		http.Error(w, "Url required", http.StatusBadRequest)
		return
	}
	if len(updatedLink.Category) == 0 {
		http.Error(w, "Category required", http.StatusBadRequest)
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
	updateCategories(&nav)
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
	updateCategories(&nav)
	err = saveNavigation(nav)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type UpdateSortIndexRequest struct {
	Updates []struct {
		Index     int    `json:"index"`              // 链接在数组中的索引
		SortIndex int    `json:"sortIndex"`          // 新的排序索引值
		Category  string `json:"category,omitempty"` // 可选的分类更新
	} `json:"updates"`
}

func updateSortIndicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateSortIndexRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 批量更新 sortIndex 和 category
	needUpdaeCategories := false
	for _, update := range req.Updates {
		if update.Index < 0 || update.Index >= len(nav.Links) {
			http.Error(w, fmt.Sprintf("Invalid index: %d", update.Index), http.StatusBadRequest)
			return
		}
		nav.Links[update.Index].SortIndex = update.SortIndex
		if update.Category != "" {
			nav.Links[update.Index].Category = update.Category
			needUpdaeCategories = true
		}
	}

	// 更新分类列表
	if needUpdaeCategories {
		updateCategories(&nav)
	}

	if err := saveNavigation(nav); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UpdateCategorysRequest struct {
	Categories []string `json:"categories"`
}

func updateCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req UpdateCategorysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 加载当前导航数据
	nav, err := loadNavigation()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 获取当前所有实际使用的分类
	currentCategories := make(map[string]struct{})
	for _, link := range nav.Links {
		if link.Category != "" {
			currentCategories[link.Category] = struct{}{}
		}
	}

	// 验证新的分类列表包含所有正在使用的分类
	for category := range currentCategories {
		found := false
		for _, newCategory := range req.Categories {
			if category == newCategory {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, fmt.Sprintf("Cannot remove category '%s' that is still in use", category), http.StatusBadRequest)
			return
		}
	}

	// 更新分类列表
	nav.Categories = req.Categories

	// 保存更新后的导航数据
	if err := saveNavigation(nav); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 返回更新后的完整导航数据
	data, err := json.MarshalIndent(nav, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
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
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	b := besticon.New(besticon.WithLogger(besticon.NewDefaultLogger(io.Discard)))
	finder := b.NewIconFinder()
	icons, err := finder.FetchIcons(url)
	if err != nil {
		log.Printf("Fetched icon from: %s err:%s", url, err)
		http.Error(w, "Failed to fetch icons", http.StatusBadRequest)
		return
	}

	if len(icons) == 0 {
		log.Printf("No icons from: %s", url)
		http.Error(w, "No icons", http.StatusBadRequest)
		return
	}
	best := icons[0]
	log.Printf("Fetched icon ok %s:  %s", url, best.URL)

	// 返回base64编码的图标数据
	iconResponse := getIconResponse(best)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iconResponse)
}

func getIconResponse(best besticon.Icon) map[string]string {
	// 1. 建立格式与 MIME 类型的映射表（覆盖常见图标格式）
	formatToContentType := map[string]string{
		"jpg":  "image/jpeg", // jpg 对应标准 MIME 类型 image/jpeg
		"jpeg": "image/jpeg", // 兼容 jpeg 格式标识
		"png":  "image/png",
		"svg":  "image/svg+xml", // svg 需指定为 image/svg+xml
		"gif":  "image/gif",     // 可选：扩展支持 gif 格式
		"ico":  "image/x-icon",  // 可选：扩展支持 ico 图标格式
	}

	// 2. 根据 best.Format 查找对应的 contentType，处理未知格式
	contentType, ok := formatToContentType[best.Format]
	if !ok {
		// 未知格式时，默认使用二进制流类型（或根据业务需求调整，如报错）
		log.Printf("Unknow image type %s", best.Format)
		contentType = "application/octet-stream"
	}

	base64Data := base64.StdEncoding.EncodeToString(best.ImageData)
	iconResponse := map[string]string{
		"iconData": fmt.Sprintf("data:%s;base64,%s", contentType, base64Data),
	}
	return iconResponse
}

// CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置允许的来源
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许的方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// 允许暴露的响应头
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		// 允许凭证
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 缓存预检请求结果
		w.Header().Set("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Add a simple usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --port=58080 --user=admin --password=123456\n", os.Args[0])
	}

	loadConfig()
	tokenStore = NewTokenStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	if envEnableNoAuthView {
		mux.HandleFunc("/navigation", getNavigationHandler)
		mux.HandleFunc("/navigation/last-modified", getNavigationLastModifiedHandler)
	} else {
		mux.HandleFunc("/navigation", authMiddleware(getNavigationHandler))
		mux.HandleFunc("/navigation/last-modified", authMiddleware(getNavigationLastModifiedHandler))
	}
	mux.HandleFunc("/navigation/add", authMiddleware(addLinkHandler))
	mux.HandleFunc("/navigation/update/", authMiddleware(updateLinkHandler))
	mux.HandleFunc("/navigation/delete/", authMiddleware(deleteLinkHandler))
	mux.HandleFunc("/navigation/sort", authMiddleware(updateSortIndicesHandler))
	mux.HandleFunc("/navigation/categories", authMiddleware(updateCategoriesHandler))
	mux.HandleFunc("/debug/tokens", debugTokensHandler)
	mux.HandleFunc("/get-icon", authMiddleware(getIconHandler))
	mux.HandleFunc("/config", getConfigHandler)
	mux.HandleFunc("/validate", authMiddleware(validateTokenHandler))

	// 静态文件
	staticFiles, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle("/", fileServer)

	// 使用 CORS 中间件和日志中间件
	handler := corsMiddleware(logAccessMiddleware(mux))

	log.Printf("Server is running on http://localhost:%s\n", envPort)
	err = http.ListenAndServe(":"+envPort, handler)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
