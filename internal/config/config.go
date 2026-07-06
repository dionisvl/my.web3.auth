package config

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	Port       string
	EthNetwork string
	SessionKey []byte
}

// Load reads configuration from the environment. It loads a .env file if
// present (dev convenience); in production real environment variables win.
func Load() *Config {
	// Best-effort: ignore error so missing .env is not fatal (like PHP safeLoad).
	_ = godotenv.Load()

	cfg := &Config{
		Port:       getEnv("APP_PORT", "8080"),
		EthNetwork: getEnv("ETH_NETWORK", "sepolia"),
	}

	if key := os.Getenv("SESSION_KEY"); key != "" {
		cfg.SessionKey = []byte(key)
	} else {
		// No persistent key configured: generate an ephemeral one so the app
		// still runs in dev. Sessions won't survive a restart — warn loudly.
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
