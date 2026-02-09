package seven

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func seven() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

/*
asChan запускает горутину, отправляющую числа в канал и затем закрывающую его.
merge читает из a и b и пересылает в выходной c.
main (seven) читает из c в for range и печатает значения.
select мультиплексирует чтение из a и b, выбирая любую готовую операцию; если оба готовы — выбор произвольный; иначе блокируется.
При чтении закрытого канала ok == false - присваивают a = nil (или b = nil), чтобы отключить соответствующий case (чтение из nil канала всегда блокирует), пока не завершатся оба.
Когда a == nil && b == nil, merge вызывает close(c) и возвращается — main завершается после for v := range c.
*/
