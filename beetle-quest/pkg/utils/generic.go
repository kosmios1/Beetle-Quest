package utils

import "os"

func PanicIfError[Type any](i Type, err error) Type {
	if err != nil {
		panic(err)
	}
	return i
}

func FindEnv(env string) string {
	value, exists := os.LookupEnv(env)
	if !exists {
		panic(env + " is not set!")
	}
	return value
}
