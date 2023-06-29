package util

func IotaToValues[T comparable](kv map[T]string) []string {
	var values []string
	for _, v := range kv {
		values = append(values, v)
	}
	return values
}
