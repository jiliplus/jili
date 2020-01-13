package clock

// var testTime = time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC)

// func TestMock_AddNext(t *testing.T) {
// 	m := clock.NewMock(testTime)
// 	test := func(C <-chan time.Time, want time.Duration) {
// 		now := <-C
// 		if got := now.Sub(testTime); got != want {
// 			t.Errorf("want timeout at t+%s, got: t+%s", want, got)
// 		}
// 	}
// 	next := func(want time.Duration) {
// 		_, got := m.AddNext()
// 		if want != got {
// 			t.Errorf("want c.AddNext(): %s, got: %s", want, got)
// 		}
// 	}
// 	tc := m.NewTicker(5 * time.Second)
// 	tm := m.NewTimer(10 * time.Second)

// 	next(5 * time.Second)
// 	test(tc.C, 5*time.Second)

// 	next(5 * time.Second)
// 	test(tc.C, 10*time.Second)
// 	test(tm.C, 10*time.Second)
// 	tc.Stop()
// 	tm.Reset(15 * time.Second)

// 	next(15 * time.Second)
// 	test(tm.C, 25*time.Second)
// 	next(0)

// 	tm.Reset(0) // fires immediately
// 	test(tm.C, 25*time.Second)
// 	next(0)
// 	tm.Reset(0) // fires immediately (again)
// 	test(tm.C, 25*time.Second)
// 	next(0)
// 	next(0)
// 	if got, want := m.Until(testTime), -25*time.Second; got != want {
// 		t.Fatalf("want: %s, got: %s", want, got)
// 	}

// 	done := make(chan struct{})

// 	m.AfterFunc(5*time.Second, func() {
// 		panic("unexpected")
// 	}).Stop()

// 	m.AfterFunc(5*time.Second, func() {
// 		panic("unexpected")
// 	}).Reset(35 * time.Second)

// 	m.AfterFunc(30*time.Second, func() {
// 		close(done)
// 	})
// 	next(30 * time.Second)
// 	<-done
// }

// func ExampleMock_AddNext() {
// 	start := time.Now()
// 	mock := clock.NewMock(start)
// 	mock.Tick(1 * time.Second)
// 	fizz := mock.Tick(3 * time.Second)
// 	buzz := mock.Tick(5 * time.Second)
// 	var items []string
// 	for i := 0; i < 20; i++ {
// 		mock.AddNext()
// 		var item string
// 		select {
// 		case <-fizz:
// 			select {
// 			case <-buzz:
// 				item = "FizzBuzz"
// 			default:
// 				item = "Fizz"
// 			}
// 		default:
// 			select {
// 			case <-buzz:
// 				item = "Buzz"
// 			default:
// 				item = strconv.Itoa(int(mock.Since(start) / time.Second))
// 			}
// 		}
// 		items = append(items, item)
// 	}
// 	fmt.Println(strings.Join(items, " "))
// 	// Output: 1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz 16 17 Fizz 19 Buzz
// }

// func ExampleNewMock() {
// 	// Use clock.Realtime() in production
// 	mock := clock.NewMock(time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC))
// 	fmt.Println("Time is now", mock.Now())
// 	timer := mock.NewTimer(15 * time.Second)
// 	mock.Add(25 * time.Second)
// 	fmt.Println("Time is now", mock.Now())
// 	fmt.Println("Timeout was", <-timer.C)
// 	// Output:
// 	// Time is now 2018-01-01 10:00:00 +0000 UTC
// 	// Time is now 2018-01-01 10:00:25 +0000 UTC
// 	// Timeout was 2018-01-01 10:00:15 +0000 UTC
// }
