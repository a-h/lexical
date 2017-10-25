package parse

func InOrder(pi Input, functions ...Function) Result {
	start := pi.Index()
	for _, f := range functions {
		r := f(pi)
		if !r.Success {
			rewind(pi, int(pi.Index()-start))
			return r
		}
	}
	//TODO: Decide what to do with the captured stuff.
	return Success("inorder", nil, nil)
}
