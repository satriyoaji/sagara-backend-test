package helper

func OutputPanicError(err error) {
	if err != nil {
		panic(err)
	}
}
