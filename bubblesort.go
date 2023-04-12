package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type PageData struct {
	SortedArray string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/sort", sortHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tpl.Execute(w, nil)
		return
	}
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		numElementsStr := r.FormValue("numElements")
		elementsStr := r.FormValue("elements")

		numElements, err := strconv.Atoi(numElementsStr)
		if err != nil {
			http.Error(w, "Invalid number of elements", http.StatusBadRequest)
			return
		}

		elements := strings.Split(elementsStr, ",")
		if len(elements) != numElements {
			http.Error(w, "Number of elements does not match", http.StatusBadRequest)
			return
		}

		array := make([]int, numElements)
		for i := 0; i < numElements; i++ {
			array[i], err = strconv.Atoi(elements[i])
			if err != nil {
				http.Error(w, "Invalid element", http.StatusBadRequest)
				return
			}
		}

		bubbleSort(array)

		pageData := PageData{
			SortedArray: fmt.Sprintf("%v", array),
		}

		err = tpl.Execute(w, pageData)
		if err != nil {
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}
	}
}

func bubbleSort(array []int) {
	n := len(array)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
}
