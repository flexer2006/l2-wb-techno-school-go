package four

/*
func main() {
  ch := make(chan int)
  go func() {
    for i := 0; i < 10; i++ {
    ch <- i
  }
}()
  for n := range ch {
    println(n)
  }
}
*/

// 0
// 1
// 2
// 3
// 4
// 5
// 6
// 7
// 8
// 9
// fatal error: all goroutines are asleep - deadlock!
// goroutine 1 [chan receive]:
// main.main()
//         .../l2-wb-techno-school-go/main.go:10 +0xa8
// exit status 2

// После отправки 10 чисел горутина завершается, но канал не закрыт, поэтому цикл range будет ждать дополнительных данных.
// Это приведёт к deadlock.
