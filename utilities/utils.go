package utilities

func CreateErrorMessage(key string, value string) interface{} {
	return map[string]string{key: value}
}