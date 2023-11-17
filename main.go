package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	content, err := os.ReadFile("./test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 词法分析
	tokens, err := tokenize(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Printf("%+v\n", tokens)
	_ = os.WriteFile("test.token", []byte(fmt.Sprintf("%+v", tokens)), 0o644)
	// 语法分析
	ast, err := parser(&tokens)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Printf("%+v\n", ast)
	_ = os.WriteFile("test.ast", []byte(fmt.Sprintf("%+v", ast)), 0o644)
	j, e := json.Marshal(ast)
	if e == nil {
		_ = os.WriteFile("test.ast.json", j, 0o644)
	}

	// 中间代码执行
	_ = ast.run()

	//flag.Parse()
	//target := flag.Args()
	//fmt.Println(target)
	//for i := 0; i < len(target); i++ {
	//	content, err := os.ReadFile(target[i])
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		continue
	//	}
	//	// 词法分析
	//	// fmt.Printf("%+v", tokenize(content))
	//	tokens, err := tokenize(content)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		continue
	//	}
	//	_ = os.WriteFile(target[i]+".token", []byte(fmt.Sprintf("%+v", tokens)), 0o644)
	//	// 语法分析
	//	ast, err := parser(tokens)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		continue
	//	}
	//	_ = os.WriteFile(target[i]+".ast", []byte(fmt.Sprintf("%+v", ast)), 0o644)
	//}
	return
}
