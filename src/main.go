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

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET at /tasks/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path == "/tasks/" {
		// Request all tasks
		var answers = make([][]byte, 4)
		var wg sync.WaitGroup
		wg.Add(4)
		go func() {
			defer wg.Done()
			var err error
			answers[0], err = processTask(task1)
			if err != nil {
				log.Fatalln("Something went wrong")
				return
			}
		}()
		go func() {
			defer wg.Done()
			var err error
			answers[1], err = processTask(task2)
			if err != nil {
				log.Fatalln("Something went wrong")
				return
			}
		}()
		go func() {
			defer wg.Done()
			var err error
			answers[2], err = processTask(task3)
			if err != nil {
				log.Fatalln("Something went wrong")
				return
			}
		}()
		go func() {
			defer wg.Done()
			var err error
			answers[3], err = processTask(task4)
			if err != nil {
				log.Fatalln("Something went wrong")
				return
			}
		}()
		wg.Wait()
		w.Write(bytes.Join(answers, []byte{}))

	} else {
		// Request "/task/<name>"
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 || pathParts[0] != "task" {
			http.Error(w, "expect /task/<name> in task handler", http.StatusBadRequest)
			return
		}

		if pathParts[1] != task1 && pathParts[1] != task2 && pathParts[1] != task3 && pathParts[1] != task4 {
			http.Error(w, "unexpected task name", http.StatusBadRequest)
		}

		answer, err := processTask(pathParts[1])
		if err != nil {
			log.Fatalln("Something went wrong")
			return
		}
		w.Write(answer)
	}
}

func processTask(currentTask string) ([]byte, error) {
	// get cases for the task
	var taskCases []json.RawMessage
	err := getCases(currentTask, &taskCases)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}

	// solve
	taskArray := [10]taskElement{}
	err = parseAndSolve(currentTask, taskCases, &taskArray)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}

	//check the results
	var data []byte
	data, err = checkResults(currentTask, &taskCases, &taskArray)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}

	return data, nil
}

func getCases(currentTask string, taskCases *[]json.RawMessage) error {
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
		return err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	err = json.Unmarshal(payload, &taskCases)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func parseAndSolve(currentTask string, taskCases []json.RawMessage, taskArray *[10]taskElement) error {
	var wg sync.WaitGroup
	wg.Add(10)

	for i, taskCase := range taskCases {
		var arguments []json.RawMessage
		err := json.Unmarshal(taskCase, &arguments)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		if currentTask == task1 {
			// rotation task takes two args
			json.Unmarshal(arguments[0], &taskArray[i].a)
			json.Unmarshal(arguments[1], &taskArray[i].k)
		} else {
			err = json.Unmarshal(arguments[0], &taskArray[i].a)
			if err != nil {
				log.Fatalln(err)
				return err
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
	}
	wg.Wait()
	return nil
}

func checkResults(currentTask string, taskCases *[]json.RawMessage, taskArray *[10]taskElement) ([]byte, error) {
	// marshaling results
	var raw []byte
	var rawArray []json.RawMessage
	var err error

	for _, element := range taskArray {
		if currentTask == task1 {
			raw, err = json.Marshal(element.result)
		} else {
			raw, err = json.Marshal(element.result[0])
		}
		if err != nil {
			log.Fatalln(err)
			return []byte{}, err
		}
		rawArray = append(rawArray, raw)
	}

	message := map[string]interface{}{
		"user_name": userName,
		"task":      currentTask,
		"results": map[string]interface{}{
			"payload": taskCases,
			"results": rawArray,
		},
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}

	// second request to Solution
	var reqAddr string
	if debugFlag {
		reqAddr = "http://" + solutionAddr + "/test/"
	} else {
		reqAddr = "http://" + solutionAddr + "/tasks/solution"
	}
	resp2, err := http.Post(reqAddr, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
		return []byte{}, err
	}
	defer resp2.Body.Close()

	var data []byte
	if resp2.StatusCode != http.StatusOK {
		data = []byte("Solution response status: " + resp2.Status)
	} else {
		data, err = ioutil.ReadAll(resp2.Body)
		if err != nil {
			log.Fatalln(err)
			return []byte{}, err
		}
	}
	return data, nil
}
