package utils

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReturnIfError(err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// func HandleIfError(err error, fn func() error) {
// 	if err := fn(); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }
