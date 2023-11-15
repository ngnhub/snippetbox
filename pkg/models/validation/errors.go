package validation

type Errors map[string][]string

func (err Errors) Add(field, message string) {
	err[field] = append(err[field], message)
}

func (err Errors) Get(field string) string {
	messages := err[field]
	if len(messages) == 0 {
		return ""
	}
	return messages[0]
}
