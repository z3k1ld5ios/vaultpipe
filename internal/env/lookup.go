package env

import "os"

// lookupEnv is a thin wrapper around os.LookupEnv, extracted so that tests
// can verify behaviour without mutating the real process environment in ways
// that are hard to isolate. It is intentionally unexported — callers outside
// this package control fallback behaviour through Interpolator.useFallback.
func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
