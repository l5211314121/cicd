package main

import (
	"fmt"
	"strings"
)

func main(){
	abc := `./test/test.py
./utils/utils.py
./utils/__pycache__
./utils/__pycache__/utils.cpython-36.pyc
./ansible_test_2.py
./ansible_test.py`
	s1 := strings.Split(abc, "\n")
	for _, s2 := range(s1){
		fmt.Println(strings.Split(s2, "/"))
	}
}