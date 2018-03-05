package utilities

func createErrorMessage(key string, value string) interface{} {
	return map[string]string{key: value}
}