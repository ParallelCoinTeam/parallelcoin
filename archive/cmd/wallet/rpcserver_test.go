package wallet

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	
	"github.com/p9c/qu"
)

func TestThrottle(t *testing.T) {
	const threshold = 1
	busy := qu.T()
	srv := httptest.NewServer(
		ThrottledFn(threshold,
			func(w http.ResponseWriter, r *http.Request) {
				<-busy
			},
		),
	)
	codes := make(chan int, 2)
	for i := 0; i < cap(codes); i++ {
		go func() {
			res, e := http.Get(srv.URL)
			if e != nil {
				t.Fatal(e)
			}
			codes <- res.StatusCode
		}()
	}
	got := make(map[int]int, cap(codes))
	for i := 0; i < cap(codes); i++ {
		got[<-codes]++
		if i == 0 {
			busy.Q()
		}
	}
	want := map[int]int{200: 1, 429: 1}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("status codes: want: %v, got: %v", want, got)
	}
}
