package main

import (
	"bytes"
	"context"
	"elmaFinalTask/src/solver"
	"elmaFinalTask/src/testSolution"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	userName         string = "d_skit"
	port             string = ":8080"
	defSolutionAddr  string = "116.203.203.76:3000"
	testSolutionAddr string = "127.0.0.1:8081"
	task1            string = "Циклическая ротация"
	task2            string = "Чудные вхождения в массив"
	task3            string = "Проверка последовательности"
	task4            string = "Поиск отсутствующего элемента"
)

var solutionAddr string

var debugFlag bool = false

// var debugFlag bool = true

type taskElement struct {
	a      []int
	k      int
	result []int
}

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET at /tasks/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path == "/tasks/" {
		// Request all tasks
		fmt.Fprintf(w, "GET received at /tasks/")

	} else {
		// Request "/task/<name>"
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 || pathParts[0] != "task" {
			http.Error(w, "expect /task/<name> in task handler", http.StatusBadRequest)
			return
		}

		currentTask := pathParts[1]

		if currentTask != task1 && currentTask != task2 && currentTask != task3 && currentTask != task4 {
			http.Error(w, "unexpected task name", http.StatusBadRequest)
		}

		// request to Solution
		var reqAddr string
		if debugFlag {
			reqAddr = "http://" + solutionAddr + "/test/" + currentTask
		} else {
			reqAddr = "http://" + solutionAddr + "/tasks/" + currentTask
		}

		resp, err := http.Get(reqAddr)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer req.Body.Close()

		payload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var taskCases []json.RawMessage
		err = json.Unmarshal(payload, &taskCases)
		if err != nil {
			log.Fatalln(err)
			return
		}

		taskArray := [10]taskElement{}
		var wg sync.WaitGroup
		wg.Add(10)

		for i, taskCase := range taskCases {
			var arguments []json.RawMessage
			err = json.Unmarshal(taskCase, &arguments)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if currentTask == task1 {
				// rotation task takes two args
				json.Unmarshal(arguments[0], &taskArray[i].a)
				json.Unmarshal(arguments[1], &taskArray[i].k)
			} else {
				err = json.Unmarshal(arguments[0], &taskArray[i].a)
				if err != nil {
					log.Fatalln(err)
					return
				}

			}

			// solve
			go func(i int) {
				defer wg.Done()
				switch currentTask {
				case task1:
					taskArray[i].result = solver.Solution1(taskArray[i].a, taskArray[i].k)
				case task2:
					taskArray[i].result = append(taskArray[i].result, solver.Solution2(taskArray[i].a))
				case task3:
					taskArray[i].result = append(taskArray[i].result, solver.Solution3(taskArray[i].a))
				case task4:
					taskArray[i].result = append(taskArray[i].result, solver.Solution4(taskArray[i].a))
				default:
				}
			}(i)
		} // taskCases loop
		wg.Wait()

		// marshaling results
		var raw []byte
		var rawArray []json.RawMessage

		for _, element := range taskArray {
			if currentTask == task1 {
				raw, err = json.Marshal(element.result)
			} else {
				raw, err = json.Marshal(element.result[0])
			}
			if err != nil {
				log.Fatalln(err)
			}
			rawArray = append(rawArray, raw)
		}

		results, err := json.Marshal(rawArray)
		if err != nil {
			log.Fatalln(err)
		}

		message := map[string]interface{}{
			"user_name": userName,
			"task":      currentTask,
			"results": map[string]interface{}{
				"payload": payload,
				"results": results,
			},
		}
		bytesRepresentation, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		// second request to Solution
		if debugFlag {
			reqAddr = "http://" + solutionAddr + "/test/"
		} else {
			reqAddr = "http://" + solutionAddr + "/tasks/solution"
		}
		resp2, err := http.Post(reqAddr, "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatalln(err)
		}

		// var finalResult map[string]interface{}
		//json.NewDecoder(resp2.Body).Decode(&finalResult)
		// if resp2.Status == http.StatusBadRequest {

		// }
		if resp2.StatusCode != http.StatusOK {
			w.Write([]byte("Solution response status: " + resp2.Status))
		} else {
			data, err := ioutil.ReadAll(resp2.Body)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(data)
		}
	}
}

func main() {
	flag.StringVar(&solutionAddr, "addr", defSolutionAddr, "Net address of the server \"Solution\"")
	flag.Parse()
	if solutionAddr == "test" {
		debugFlag = true

	}
	if debugFlag {
		solutionAddr = testSolutionAddr
	}
	fmt.Println("Solution server at: " + solutionAddr)

	if debugFlag {
		ctx := context.Background()
		ctxTestServer, cancelTestServer := context.WithCancel(ctx)
		defer cancelTestServer()
		go testSolution.StartTestServer(ctxTestServer, testSolutionAddr)
	}

	http.HandleFunc("/task/", tasksHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	fmt.Println("Service is getting up")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
