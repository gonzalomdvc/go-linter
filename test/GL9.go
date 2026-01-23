package test

func main() {
	ch := make(chan int)

	select {
	case <-ch:
		return
	}
}
