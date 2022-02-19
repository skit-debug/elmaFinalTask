package testSolution

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	task1 string = "Циклическая ротация"
	task2 string = "Чудные вхождения в массив"
	task3 string = "Проверка последовательности"
	task4 string = "Поиск отсутствующего элемента"
)

func testTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		// filling test data
		var taskCases [10]json.RawMessage
		switch pathParts[1] {
		case task1:
			for i := 0; i < 10; i++ {
				a, _ := json.Marshal([]int{3, 8, 9, 7, 6})
				k, _ := json.Marshal(3)
				taskCases[i], _ = json.Marshal([]json.RawMessage{a, k})
			}
		case task2:
			for i := 0; i < 10; i++ {
				a, _ := json.Marshal([]int{9, 3, 9, 3, 9, 7, 9})
				taskCases[i], _ = json.Marshal([]json.RawMessage{a})
			}
		case task3:
			for i := 0; i < 10; i++ {
				a, _ := json.Marshal([]int{4, 1, 3, 2})
				taskCases[i], _ = json.Marshal([]json.RawMessage{a})
			}
		case task4:
			for i := 0; i < 10; i++ {
				a, _ := json.Marshal([]int{2, 3, 1, 5})
				taskCases[i], _ = json.Marshal([]json.RawMessage{a})
			}
		default:
		}
		bytesRepresentation, err := json.Marshal(taskCases)
		if err != nil {
			log.Fatalln(err)
		}
		w.Write(bytesRepresentation)
	}
	if req.Method == http.MethodPost {
		var finalResult map[string]interface{}
		json.NewDecoder(req.Body).Decode(&finalResult)

		message := map[string]interface{}{
			"percent": 90,
			"fails": map[string]int{
				"OriginalResult": 1,
				"ExternalResult": 0,
			},
		}

		bytesRepresentation, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}
		w.Write(bytesRepresentation)
	}
}

func StartTestServer(ctx context.Context, addr string) {
	http.HandleFunc("/test/", testTaskHandler)
	fmt.Println("Test server is getting up")
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Test server is down")
			return
		default:
		}
	}()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
