package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
)

const (
	// CSRFTokenKey 是存储在会话(session)中的CSRF令牌的键名
	CSRFTokenKey = "csrf_token"
	// CSRFTokenLength 是CSRF令牌的字节长度
	CSRFTokenLength = 32
)

// GenerateCSRFToken 为给定的会话生成一个新的CSRF令牌
func GenerateCSRFToken(session *Session) (string, error) {
	//生成随机字节
	b := make([]byte, CSRFTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	//转换为base64字符串
	token := base64.StdEncoding.EncodeToString(b)

	//将令牌存储到会话中
	session.Data[CSRFTokenKey] = token

	return token, nil
}

// GetCSRFToken 从会话中获取CSRF令牌，如果不存在则生成一个新令牌
func GetCSRFToken(session *Session) (string, error) {
	// 检查会话中是否已存在令牌
	if token, ok := session.Data[CSRFTokenKey].(string); ok {
		return token, nil
	}

	// 如果不存在，生成一个新令牌
	return GenerateCSRFToken(session)
}

// ValidateCSRFToken 验证请求中的CSRF令牌是否与会话中的令牌匹配
func ValidateCSRFToken(r *http.Request, session *Session) error {
	//从会话中获取令牌
	sessionToken, ok := session.Data[CSRFTokenKey].(string)
	if !ok {
		return errors.New("会话中没有CSRF令牌")
	}

	// 从表单或请求头中获取令牌
	var requestToken string

	//首先尝试从表单中获取
	if r.Method == http.MethodPost {
		requestToken = r.FormValue("csrf_token")
	}

	//如果表单中没有,尝试从请求头中获取
	if requestToken == "" {
		requestToken = r.Header.Get("X-CSRF-Token")
	}

	//验证令牌是否匹配
	if requestToken == "" || requestToken != sessionToken {
		return errors.New("CSRF令牌无效")
	}

	return nil
}

// CSRFMiddleware 是一个中间件 , 用于验证CSRF令牌
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 只处理POST、PUT、DELETE和PATCH请求
		if r.Method == http.MethodPost || r.Method == http.MethodPut ||
			r.Method == http.MethodDelete || r.Method == http.MethodPatch {

			//获取会话
			session, err := GlobalManager.GetSession(r)
			if err != nil {
				http.Error(w, "未授权: 无效的会话", http.StatusUnauthorized)
				return
			}

			//验证CSRF 令牌
			if err := ValidateCSRFToken(r, session); err != nil {
				http.Error(w, "未授权："+err.Error(), http.StatusUnauthorized)
				return
			}
		}
		//继续处理请求
		next.ServeHTTP(w, r)
	})
}
