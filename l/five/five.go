package five

/*
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	// S1021: should merge variable declaration with assignment on next line (staticcheck) go-golangci-lint-v2
	var err error
	// SA4023(related information): the lhs of the comparison gets its value from here and has a concrete type (staticcheck) go-golangci-lint-v2
	err = test()
	// tautological condition: non-nil != nil nilnesscond
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
*/

// Вывод: error
// Функция test() возвращает nil указатель на *customError.
// Переменная err объявлена как интерфейс error.
// В Go интерфейс состоит из типа и данных. При присваивании nil указателя конкретного типа интерфейсу, интерфейс не становится nil, так как тип (*customError) записан.
// Поэтому err != nil возвращает true, и программа печатает "error" и завершается, не доходя до "ok".
