package utils

import "crypto/rand"

func GenerateRandomID(byte_size int) ([]byte, error) {
	bytes := make([]byte, byte_size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
