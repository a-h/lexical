package parse

func Success(item ResultItem, next Function, err error) InterfaceResult {
	return InterfaceResult{
		success: true,
		next:    next,
		item:    item,
		err:     err,
	}
}
