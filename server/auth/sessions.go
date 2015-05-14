package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	sql "github.com/aodin/aspect"
	"github.com/aodin/volta/config"

	db "github.com/aodin/listofthings/db"
)

// TODO Check session expiration

type SessionManager struct {
	conn   sql.Connection
	cookie config.CookieConfig
	keyGen func() string
}

// Create creates a new session with a random key
func (m *SessionManager) Create(user db.User) db.Session {
	return m.create(time.Now().UTC(), user)
}

func (m *SessionManager) create(now time.Time, user db.User) (session db.Session) {
	// Set the expires from the cookie config
	session.UserID = user.ID
	session.Expires = now.Add(m.cookie.Age)

	// TODO Single transaction

	// Generate a new random session key
	for {
		session.Key = m.keyGen()

		// No duplicates - generate a new key if this key already exists
		var duplicate string
		stmt := sql.Select(
			db.Sessions.C["key"],
		).Where(db.Sessions.C["key"].Equals(session.Key)).Limit(1)
		if !m.conn.MustQueryOne(stmt, &duplicate) {
			break
		}
	}

	// Insert the session
	m.conn.MustExecute(db.Sessions.Insert().Values(session))
	return session
}

func (m *SessionManager) Get(key string) (session db.Session) {
	stmt := db.Sessions.Select().Where(db.Sessions.C["key"].Equals(key))
	m.conn.MustQueryOne(stmt, &session)
	return
}

func (m *SessionManager) GetUser(key string) (user db.User) {
	stmt := db.Users.Select().JoinOn(
		db.Sessions,
		db.Sessions.C["user_id"].Equals(db.Users.C["id"]),
	).Where(
		db.Sessions.C["key"].Equals(key),
	).Limit(1)
	m.conn.MustQueryOne(stmt, &user)
	return
}

// Sessions creates a new session manager
func Sessions(conf config.Config, conn sql.Connection) *SessionManager {
	return &SessionManager{
		conn:   conn,
		cookie: conf.Cookie,
		keyGen: RandomKey,
	}
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
