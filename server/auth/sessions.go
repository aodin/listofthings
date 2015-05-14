package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aodin/volta/config"
)

// SessionManager is an in-memory store of sessions
type SessionManager struct {
	sync.RWMutex
	sessions map[string]Session
	cookie   config.CookieConfig
	keyGen   func() string
}

// Create creates a new session with a random key
func (m *SessionManager) Create() Session {
	return m.create(time.Now().UTC())
}

func (m *SessionManager) create(now time.Time) (session Session) {
	m.Lock()
	defer m.Unlock()

	// Set the expires from the cookie config
	session = Session{
		Expires: now.Add(m.cookie.Age),
	}

	// Generate a new random session key
	for {
		session.Key = m.keyGen()
		if _, exists := m.sessions[session.Key]; !exists {
			break
		}
	}

	m.sessions[session.Key] = session
	return
}

func (m *SessionManager) Get(key string) Session {
	m.RLock()
	defer m.RUnlock()
	return m.sessions[key]
}

// Session creates a key with an optional user
type Session struct {
	Key     string
	UserID  int64
	Expires time.Time
}

func (session Session) Exists() bool {
	return session.Key != ""
}

// EncodeBase64String is a wrapper around the standard base64 encoding call.
func EncodeBase64String(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

// RandomBytes returns random bytes from the crypto/rand Reader or it panics.
func RandomBytes(n int) []byte {
	key := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		log.Panicf("auth: could not generate random bytes: %s", err)
	}
	return key
}

// RandomKey generates a new session key. It does so by producing 24
// random bytes that are encoded in URL safe base64, for output of 32 chars.
func RandomKey() string {
	return RandomKeyN(24)
}

// RandomKeyN generates a new Base 64 encoded random string. N is the length
// of the random bytes, not the final encoded string.
func RandomKeyN(n int) string {
	return EncodeBase64String(RandomBytes(n))
}

// Sessions creates a new session manager
func Sessions(conf config.Config) *SessionManager {
	return &SessionManager{
		cookie:   conf.Cookie,
		sessions: make(map[string]Session),
		keyGen:   RandomKey,
	}
}
