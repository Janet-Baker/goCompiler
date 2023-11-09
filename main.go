package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//content, err := os.ReadFile("./test.txt")
	//if err == nil {
	//	fmt.Printf("%+v", tokenize(content))
	//} else {
	//	fmt.Println(err)
	//}

	flag.Parse()
	target := flag.Args()
	fmt.Println(target)
	for i := 0; i < len(target); i++ {
		content, err := os.ReadFile(target[i])
		if err == nil {
			// 词法分析
			// fmt.Printf("%+v", tokenize(content))
			_ = os.WriteFile(target[i]+".token", []byte(fmt.Sprintf("%+v", tokenize(content))), 0o644)
		} else {
			fmt.Println(err)
		}
	}
}
