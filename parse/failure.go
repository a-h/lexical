package parse

func Failure() InterfaceResult {
	return InterfaceResult{
		success: false,
	}
}
