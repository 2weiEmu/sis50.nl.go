package main

import "math/rand/v2"


func MakeRandomString(length int) []byte {
	chars := "abcdefghijklmnopqrtstuvwxyzABCDEFGHIJKLMNOPQRTSTUVWXYZ"
	final := make([]byte, length)

	for i := 0; i < length; i++ {
		final[i] = chars[rand.IntN(len(chars))]
	}
	return final
}


