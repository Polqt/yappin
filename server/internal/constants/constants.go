package constants

import "time"

// JWT and Authentication
const (
	JWTCookieName     = "jwt"
	JWTCookieDuration = 24 * time.Hour
	JWTTokenExpiry    = 24 * time.Hour
)

// Room Configuration
const (
	RoomDefaultExpiry   = 24 * time.Hour
	RoomCleanupInterval = 5 * time.Minute
	DefaultRoomLimit    = 50
	MaxRoomHistory      = 100
)

// Rate Limiting
const (
	DefaultRateLimit  = 100
	AuthRateLimit     = 10
	RoomCreationLimit = 10
	RateLimitWindow   = time.Minute
)

// WebSocket
const (
	WebSocketReadBufferSize  = 1024
	WebSocketWriteBufferSize = 1024
)

// Database
const (
	DBMaxOpenConns    = 25
	DBMaxIdleConns    = 5
	DBConnMaxLifetime = 5 * time.Minute
)

// Service Timeouts
const (
	DefaultServiceTimeout = 2 * time.Second
	HTTPServerTimeout     = 30 * time.Second
)

// AllowedReactionEmojis defines the valid emoji reactions for messages
var AllowedReactionEmojis = []string{"👍", "❤️", "😂", "😮", "😢", "👏", "🎉"}

// IsValidReactionEmoji checks if the given emoji is in the allowed list
func IsValidReactionEmoji(emoji string) bool {
	for _, e := range AllowedReactionEmojis {
		if e == emoji {
			return true
		}
	}
	return false
}
