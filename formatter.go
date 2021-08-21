package bundlerr

// Formatter is responsible for formatting all the bundled errors into a single string
type Formatter func(Bundle) string

func defaultFormatFn(b Bundle) string {
	res := ""
	delimiter := " █ "

	for _, e := range b.Errors() {
		if e == nil {
			continue
		}

		res += e.Error() + delimiter
	}

	if l := len(res); l > 0 {
		res = res[:l-len(delimiter)]
	}

	return res
}
