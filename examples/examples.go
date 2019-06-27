package main

import (
	"fmt"
	gf "github.com/waltsmith88/go-flashtext"
	"log"
	"path/filepath"
	"runtime"
)

func main() {

	file := "examples/data.txt"
	filePath, _ := filepath.Abs(file)

	keywordProcessor := gf.NewKeywordProcessor()

	// Example 1: add multiple keywords from file
	keywordProcessor.AddKeywordsFromFile(filePath)

	// Example 2: add a keyword
	keywordProcessor.AddKeyword("abc", "abc")

	// Example 3: add keywords from Map
	keywordMap := map[string]string{
		"abcd": "abcd",
		"student": "stu",
	}
	keywordProcessor.AddKeywordsFromMap(keywordMap)

	// Example 4: add same keyword "abc" with different cleanName
	keywordProcessor.AddKeyword("abc", "abc1")

	allKeywords := keywordProcessor.GetAllKeywords()
	fmt.Println(allKeywords)

	// Example 5: get cleanName of keyword from TrieTree
	if cleanName, err := keywordProcessor.GetKeyword("abc"); err {
		fmt.Println("Success：", cleanName, err)
	} else {
		fmt.Println("Failed：", cleanName, err)
	}

	// Example 6: 添加中文关键词并提取 add Chinese keywords
	keywordProcessor.AddKeyword("中文", "中文")
	// Extract keywords from sentence by searching TrieTree
	sentence := "中文支持1bceAbcd支持中文student abc"
	cleanNameList := keywordProcessor.ExtractKeywords(sentence)
	fmt.Println(cleanNameList)

	// Example 7: set properties of keywordProcessor
	fmt.Println(keywordProcessor.GetCaseSensitive())
	keywordProcessor.SetCaseSensitive(false)
	fmt.Println(keywordProcessor.GetCaseSensitive())
	cleanNameList1 := keywordProcessor.ExtractKeywords(sentence)
	fmt.Println(cleanNameList1)

	// Example 8: Extract keywords from sentence by searching TrieTree and return keywords' span
	cleanNameRes := keywordProcessor.ExtractKeywordsWithSpanInfo(sentence)
	sentence1 := []rune(sentence)
	fmt.Println(cleanNameRes, sentence1)
	for _, resSpan := range cleanNameRes {
		fmt.Println(resSpan.CleanName, resSpan.StartPos, resSpan.EndPos, fmt.Sprintf("%c", sentence1[resSpan.StartPos:resSpan.EndPos]))
	}

	// Example 9: delete keyword
	keywordProcessor.RemoveKeyword("abc")

	// Example 10: delete keywords in list
	keywordProcessor.RemoveKeywordFromList([]string{"student", "abcd", "abc", "中文"})
	fmt.Println(keywordProcessor.GetAllKeywords())

	// Example 11: replace keywords in sentence with their cleanName
	sourceSentence := "hello中国helloabc"
	newSentence := keywordProcessor.ReplaceKeywords(sourceSentence)
	fmt.Println(fmt.Sprintf("source sentence: %s; \nnew sentence: %s", sourceSentence, newSentence))

	printMemStats()
}

func printMemStats() {
	_, _, line, _ := runtime.Caller(1)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Line %v : Alloc = %v TotalAlloc = %v Sys = %v NumGC = %v\n", line, m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024,
		m.NumGC)
}
