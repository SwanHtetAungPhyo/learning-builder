package common

func Must[T any](val T, err error, msg ...string) T {
	if err != nil {
		detail := ""
		if len(msg) > 0 {
			detail = msg[0] + ": "
		}
		panic(detail + err.Error())
	}
	return val
}

//
//func MustEnhance[T any](fn func() (T, error), msg string) T {
//	val, err := fn()
//	if err != nil {
//		panic(msg + err.Error())
//	}
//	return val
//}
