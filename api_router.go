package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	authURL  = "https://localhost:7176/api/auth/login"  // 登录接口，获取 Token
	dataURL  = "https://localhost:7027/api/gamingtable" // GamingTable API
	email    = "admin@city.com"                         // 测试账号
	password = "Admin@12345"                            // 测试密码
)

func main() {
	// 1️⃣ 获取 JWT Token
	token, err := getToken(email, password)
	if err != nil {
		fmt.Println("❌ 获取 Token 失败:", err)
		return
	}

	fmt.Println("✅ 成功获取 Token:", token)

	// 2️⃣ 带上 Token 访问 GamingTable API
	err = getGamingTables(token)
	if err != nil {
		fmt.Println("❌ 访问 GamingTable 失败:", err)
		return
	}
}

// 獲取 JWT Token
func getToken(email, password string) (string, error) {
	// 請求體
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	// 發送 POST 請求
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 響應
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("登录失败: %s", string(body))
	}

	// 解析 JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// 返回 Token
	token, exists := result["token"]
	if !exists {
		return "", fmt.Errorf("响应中未找到 token: %s", string(body))
	}
	return token.(string), nil
}

// 訪問 GamingTable API
func getGamingTables(token string) error {
	// Request for gaming table queryable
	req, _ := http.NewRequest("GET", dataURL, nil)
	req.Header.Set("Authorization", "Bearer "+token) // 設置 JWT Token 請求頭

	client := &http.Client{}    // Build Client
	resp, err := client.Do(req) // 發送請求
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 將 Gaming Table Info List 從 Body 中讀取出來
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("🎯 GamingTable API 响应:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("访问失败: %s", string(body))
	}

	return nil
}
