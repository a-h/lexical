package parse

var _ Result = (*InterfaceResult)(nil)

type InterfaceResult struct {
	success bool
	item    ResultItem
	next    Function
	err     error
}

func (r InterfaceResult) Success() bool {
	return r.success
}

func (r InterfaceResult) Item() ResultItem {
	return r.item
}

func (r InterfaceResult) Next() Function {
	return r.next
}

func (r InterfaceResult) Error() error {
	return r.err
}
