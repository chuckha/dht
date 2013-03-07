package db

import "testing"

func TestAdd(t *testing.T) {
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	db.Add("Fifth", 7)
	val := db.Get("Fifth")
	if val != 7 {
		t.Errorf("Should have been 7, got: %d", val)
	}
}

func TestConcurrentAdd(t *testing.T) {
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	results := make(chan []int)
	// [expected, actual]
	for i := 0; i < 100; i++ {
		go func () {
			db.Add("First", i)
			db.Get("First")
			results <- []int{i, db.Get("First")}
		}()
	}
	for i := 0; i < 100; i++ {
		vals := <-results
		if vals[0] != vals[1] {
			t.Errorf("Add is not threadsafe")
		}
	}
}

func TestGet(t *testing.T) {
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	val := db.Get("First")
	if val != 3 {
		t.Errorf("First should have been 3, got: %d", val)
	}
}

func TestConcurrentGet(t *testing.T) {
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	results := make(chan int)
	for i := 0; i < 100; i++ {
		go func () {
			results<-db.Get("First")
			results<-db.Get("Second")
			results<-db.Get("Third")
			results<-db.Get("Fourth")
		}()
	}
	for i := 0; i < 100; i++ {
		val := <-results
		if val != 3 && val != 4 && val != 5 && val != 6 {
			t.Errorf("Should have been 3 4 5 or 6 but got: %d", val)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for k := range m {
			db.Add(k, i)
		}
	}
}

func BenchmarkConcurrentAdd(b *testing.B) {
	b.StopTimer()
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	b.StartTimer()
	results := make(chan struct{})
	for j := 0; j < b.N; j++ {
		for i := 0; i < 100; i++ {
			go func () {
				db.Add("First", i)
				results <- struct{}{}
			}()
		}
	}
	for i := 0; i < 100 * b.N; i++ {
		<-results
	}
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for k := range m {
			db.Get(k)
		}
	}
}

func BenchmarkConcurrentGet(b *testing.B) {
	b.StopTimer()
	var m = map[string]int{
		"First": 3,
		"Second": 4,
		"Third": 5,
		"Fourth": 6,
	}
	db := NewDb("test", m)
	results := make(chan struct{})
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		for i := 0; i < 100; i++ {
			go func () {
				db.Get("First")
				db.Get("Second")
				db.Get("Third")
				db.Get("Fourth")
				results <- struct{}{}
			}()
		}
	}
	for i := 0; i < 100 * b.N; i++ {
		<-results
	}
}
