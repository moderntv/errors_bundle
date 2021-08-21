package bundlerr

import "encoding/json"

func (b Bundle) Len() int {
	return len(b.errors)
}

func (b Bundle) Swap(i, j int) {
	b.errors[i], b.errors[j] = b.errors[j], b.errors[i]
}

func (b Bundle) Less(i, j int) bool {
	return b.errors[i].Error() < b.errors[j].Error()
}

func (b *Bundle) MarshalJSON() ([]byte, error) {
	errStrings := []string{}
	for _, err := range b.errors {
		errStrings = append(errStrings, err.Error())
	}

	return json.Marshal(errStrings)
}
