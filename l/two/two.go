package two

// package main

// import "fmt"

// func test() (x int) {
//   defer func() {
//     x++
//   }()
//   x = 1
//   return
// }

// func anotherTest() int {
//   var x int
//   defer func() {
//     x++
//   }()
//   x = 1
//   return x
// }

// func main() {
//   fmt.Println(test())
//   fmt.Println(anotherTest())
// }

// 2
// 1
// Выполняется x = 1.
// На return x сначала вычисляется выражение возврата  и это значение копируется/
// Затем выполняются defer: x++ делает локальный x равным 2.
// Но возвращаемое значение уже было вычислено как 1, и defer его не меняет.
