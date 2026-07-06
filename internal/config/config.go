package config

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	EthNetwork   string
	SessionKey   []byte
	CookieSecure bool
}

// Load reads config from the environment, loading a .env file if present.
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:         getEnv("APP_PORT", "8080"),
		EthNetwork:   getEnv("ETH_NETWORK", "sepolia"),
		CookieSecure: getEnv("COOKIE_SECURE", "false") == "true",
	}

	if key := os.Getenv("SESSION_KEY"); key != "" {
		cfg.SessionKey = []byte(key)
	} else {
		// No persistent key: use an ephemeral one so dev still works, but
		// sessions reset on restart.
		cfg.SessionKey = randomKey(32)
		log.Println("WARNING: SESSION_KEY not set, using an ephemeral random key (sessions reset on restart)")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func randomKey(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to generate session key: %v", err)
	}
	return b
}
