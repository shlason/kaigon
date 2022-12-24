package utils

import "math/rand"

func RandStringBytes(n int) string {
	// 避免 0, 1, i, j, l, I 視覺上混淆，故把 '0', '1', 'i', 'j', 'l', 'I' 拿掉
	const letterBytes = "23456789abcdefghkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		// TODO: seed
		// TODO: nano ID
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
