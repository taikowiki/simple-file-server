package serverUtil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
)

type UserAuthData struct {
	Provider   string `json:"provider"`
	ProviderId string `json:"providerId"`
}

func Decipher(encryptedHex string, keyHex string) (*UserAuthData, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}

	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Node.js crypto's aes-256-gcm uses the key as IV if passed that way.
	// In GCM, the IV is usually 12 bytes.
	// Looking at the JS: createDecipheriv('aes-256-gcm', bufferKey, bufferKey)
	// bufferKey is 32 bytes (64 hex chars).
	// If Node.js receives a longer IV, it might be truncating or handling it specifically.
	// Actually, GCM IV should be 12 bytes. If the JS uses 32 bytes, it's non-standard.

	// Let's try to match Node.js behavior.
	// Node.js's createDecipheriv for GCM: if IV is not 12 bytes, it might be doing something else.
	iv := key[:12] // Standard GCM IV is 12 bytes.

	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(key)) // Try with 32 byte nonce if that's what JS did
	if err != nil {
		// If 32 byte nonce is not supported by default GCM, use standard 12 byte nonce from the key
		aesgcm, err = cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		iv = key[:12]
	} else {
		iv = key
	}

	decrypted, err := aesgcm.Open(nil, iv, encrypted, nil)
	if err != nil {
		return nil, err
	}

	var authData UserAuthData
	err = json.Unmarshal(decrypted, &authData)
	if err != nil {
		return nil, err
	}

	return &authData, nil
}
