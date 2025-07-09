package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"
)

/*
登录时：验证用户名和密码
登录成功后：创建 session 并设置 cookie
后续请求中：通过 cookie 中的会话 ID 验证 session
数据修改操作时：额外验证 CSRF 令牌
CSRF:
1.当用户访问包含表单的页面（如用户编辑、删除页面）时，
  系统会自动检查会话中是否已有CSRF令牌，如果没有则创建新令牌。
  存储到服务器对应的Session中和通过HTML表单中的隐藏字段或JavaScript变量传递给客户端
2.当用户提交表单时，系统会检查(客户端)表单中的CSRF令牌是否与(服务端)会话中的令牌匹配。
  如果匹配，则认为请求是合法的，否则认为请求是非法的。
3.如果请求是非法的，则系统会返回错误信息，并要求用户重新提交表单。
4.如果请求是合法的，则系统会执行相应的操作。
5.如果用户在会话过期后重新登录，则系统会自动创建新的CSRF令牌。
*/

/*
Session.Data 存储:
CSRF令牌
用户偏好设置（通过辅助函数设置)
可能的其他用途（虽然当前项目未明确使用）
用户角色和权限缓存  购物车信息（在电商应用中）  UI偏好   临时消息（Flash消息
*/

/*
当服务器重启后  所有内存中的会话数据（包括CSRF令牌）全部丢失
当用户带着旧Cookie访问重启后的服务器 都会验证失败，需要重新登录
解决方案
通常会采用持久化存储会话数据的方法
数据库存储：将会话存储在MySQL、PostgreSQL等数据库中
Redis存储：使用Redis等内存数据库存储会话
件系统存储：将会话数据序列化到文件中
*/

// Session 表示一个用户会话
type Session struct {
	ID        string                 // 会话唯一标识符 随机的sid
	UserID    int                    // 关联的用户ID 数据库中的id
	Data      map[string]interface{} // 会话数据存储 存储数据+CSRF令牌
	CreatedAt time.Time              // 会话创建时间
	ExpiresAt time.Time              // 会话过期时间
}

// Manager 会话管理器，负责创建、获取和销毁会话
type Manager struct {
	cookieName  string              // 表示这个Manager实例是管理session的 固定为session_id
	lock        sync.RWMutex        // 读写锁，保证并发安全
	sessions    map[string]*Session // 会话存储，键为会话ID(sid)
	maxLifetime time.Duration       // 默认session会话最大生存时间,如果没设置就使用这个 2h
}

// 全局会话管理器实例
var GlobalManager *Manager

// NewManager 创建一个新的会话管理器
func NewManager(cookieName string, maxLifetime time.Duration) *Manager {
	return &Manager{
		cookieName:  cookieName,
		sessions:    make(map[string]*Session),
		maxLifetime: maxLifetime,
	}
}

// Init 初始化全局会话管理器
func Init(cookieName string, maxLifetime time.Duration) {
	GlobalManager = NewManager(cookieName, maxLifetime)
	// 启动定期清理过期会话的goroutine
	// 延迟启动GC，避免影响启动速度
	go func() {
		time.Sleep(10 * time.Second) // 延迟10秒后启动GC
		GlobalManager.GC()
	}()
}

// generateSessionID 生成一个唯一的会话ID
func (manager *Manager) generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateSession 创建一个新会话
func (manager *Manager) CreateSession(w http.ResponseWriter, userID int, remember bool) (*Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	// 生成会话ID
	sid, err := manager.generateSessionID()
	if err != nil {
		return nil, err
	}

	// 计算过期时间
	var expiresAt time.Time
	var maxAge int
	if remember {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
		maxAge = 30 * 24 * 60 * 60
	} else {
		expiresAt = time.Now().Add(manager.maxLifetime)
		maxAge = int(manager.maxLifetime.Seconds())
	}

	// 创建新会话
	session := &Session{
		ID:        sid,
		UserID:    userID,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	// 立即生成 CSRF token
	token := generateCSRFTokenDirect()
	session.Data[CSRFTokenKey] = token

	// 存储会话
	manager.sessions[sid] = session

	// 设置Cookie
	cookie := http.Cookie{
		Name:     manager.cookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
		Secure:   false,                // 生产环境应该设为 true
		SameSite: http.SameSiteLaxMode, // 新增：防止 CSRF
	}
	http.SetCookie(w, &cookie)

	return session, nil
}

// GetSession 从请求中获取会话
func (manager *Manager) GetSession(r *http.Request) (*Session, error) {
	// 从Cookie中获取会话ID
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		return nil, err
	}

	sid := cookie.Value

	// 使用写锁，避免并发问题
	manager.lock.Lock()
	defer manager.lock.Unlock()

	session, exists := manager.sessions[sid]
	if !exists {
		return nil, errors.New("会话不存在或已过期")
	}

	// 检查会话是否过期
	if session.ExpiresAt.Before(time.Now()) {
		// 直接删除过期会话，因为我们已经持有写锁
		delete(manager.sessions, sid)
		return nil, errors.New("会话已过期")
	}

	return session, nil
}

// DestroySession 销毁会话
func (manager *Manager) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil {
		return
	}

	sid := cookie.Value

	//删除会话
	manager.lock.Lock()
	delete(manager.sessions, sid)
	manager.lock.Unlock()

	// 使Cookie过期
	expiredCookie := http.Cookie{
		Name:     manager.cookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, &expiredCookie)
}

// GC 垃圾收集，清理过期的会话
func (manager *Manager) GC() {
	for {
		time.Sleep(time.Minute) // 每分钟检查一次

		manager.lock.Lock()
		for sid, session := range manager.sessions {
			if session.ExpiresAt.Before(time.Now()) {
				delete(manager.sessions, sid)
			}
		}
		manager.lock.Unlock()
	}
}

func generateCSRFTokenDirect() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
