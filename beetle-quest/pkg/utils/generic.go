package utils

func PanicIfError[Type any](i Type, err error) Type {
	if err != nil {
		panic(err)
	}
	return i
}
