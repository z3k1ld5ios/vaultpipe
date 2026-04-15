package token

import "strings"

const maskChar = "*"

// Mask returns a partially redacted representation of a token safe for logging.
// Tokens shorter than 8 characters are fully masked.
func Mask(t string) string {
	t = strings.TrimSpace(t)
	if len(t) < 8 {
		return strings.Repeat(maskChar, len(t))
	}
	visible := 4
	return t[:visible] + strings.Repeat(maskChar, len(t)-visible)
}

// MaskToken is a convenience wrapper that operates on a *Token.
// It returns an empty string if tok is nil.
func MaskToken(tok *Token) string {
	if tok == nil {
		return ""
	}
	return Mask(tok.Value)
}
