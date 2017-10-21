package parse

func String(name, value string, next Function, err error) InterfaceResult {
	return InterfaceResult{
		success: true,
		next:    next,
		item: InterfaceResultItem{
			name:  name,
			value: value,
		},
		err: err,
	}
}
