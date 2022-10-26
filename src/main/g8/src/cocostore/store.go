// Package cocostore offers an Echo's cookie store with data compression.
//
// Available since template-r3
package cocostore

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	myRander          = rand.New(rand.NewSource(time.Now().UnixNano()))
	ErrorValueTooLong = errors.New("the value is too long")
)

// CompressionLevel defines how much session data is compressed
type CompressionLevel int

const (
	CompressionLevelNone CompressionLevel = iota
	CompressionLevelBestSpeed
	CompressionLevelBalance
	CompressionLevelBestCompression
	// CompressionLevelRandom
)

// NewCompressedCookieStore creates a new NewCompressedCookieStore.
//
// Keys are defined in pairs to allow key rotation, but the common case is
// to set a single authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for
// encryption. The encryption key can be set to nil or omitted in the last
// pair, but the authentication key is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes.
// The encryption key, if set, must be either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256 modes.
func NewCompressedCookieStore(compressionLevel CompressionLevel, keyPairs ...[]byte) *CompressedCookieStore {
	ccs := &CompressedCookieStore{
		compressionLevel: compressionLevel,
		maxLength:        4096,
		CookieStore:      sessions.NewCookieStore(keyPairs...),
	}
	for _, codec := range ccs.Codecs {
		// disable codec max length
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxLength(0)
		}
	}
	return ccs
}

// CompressedCookieStore stores sessions using secure cookies in compressed format.
type CompressedCookieStore struct {
	compressionLevel CompressionLevel
	maxLength        int
	*sessions.CookieStore
}

// GetCompressionLevel returns the current compression level value.
func (s *CompressedCookieStore) GetCompressionLevel() CompressionLevel {
	return s.compressionLevel
}

// SetCompressionLevel sets the compression level used to compress session data.
func (s *CompressedCookieStore) SetCompressionLevel(value CompressionLevel) *CompressedCookieStore {
	s.compressionLevel = value
	return s
}

// GetMaxLength returns the current maximum length of session data (after compression).
func (s *CompressedCookieStore) GetMaxLength() int {
	return s.maxLength
}

// SetMaxLength limits the maximum length of session data (after compression).
func (s *CompressedCookieStore) SetMaxLength(value int) *CompressedCookieStore {
	s.maxLength = value
	return s
}

// Get overrides CookieStore.Get.
func (s *CompressedCookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New overrides CookieStore.New.
func (s *CompressedCookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		var compressed, decompressed []byte
		if compressed, err = base64.StdEncoding.DecodeString(c.Value); err == nil {
			if decompressed, err = zlibDecompress(compressed); err == nil {
				raw := string(decompressed)
				if err = securecookie.DecodeMulti(name, raw, &session.Values, s.Codecs...); err == nil {
					session.IsNew = false
				}
			}
		}
	}
	return session, err
}

// Save overrides CookieStore.Save.
//
//	- Save returns ErrorValueTooLong if the length of encoded session values (after compression) exceeds limit
func (s *CompressedCookieStore) Save(_ *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	raw, err := securecookie.EncodeMulti(session.Name(), session.Values, s.CookieStore.Codecs...)
	if err != nil {
		return err
	}
	compressed := zlibCompress(s.compressionLevel, []byte(raw))
	encoded := base64.StdEncoding.EncodeToString(compressed)
	if s.maxLength > 0 && len(encoded) > s.maxLength {
		return ErrorValueTooLong
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func zlibCompress(compressionLevel CompressionLevel, data []byte) []byte {
	var level = zlib.DefaultCompression
	switch compressionLevel {
	case CompressionLevelNone:
		level = zlib.NoCompression
	case CompressionLevelBestSpeed:
		level = zlib.BestSpeed
	case CompressionLevelBalance:
		level = zlib.DefaultCompression
	case CompressionLevelBestCompression:
		level = zlib.BestCompression
		// case CompressionLevelRandom:
		// 	level = myRander.Intn(zlib.BestCompression) + 1
	}
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, level)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

func zlibDecompress(compressedData []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	r.Close()
	return b.Bytes(), err
}
