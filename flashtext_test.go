package flashtext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddKeywords(t *testing.T) {

	var testSet = []struct{
		in			map[string]string
		expected	map[string]string
	} {
		{map[string]string{"teacher": "tea"}, 				map[string]string{"teacher": "tea"}},
		{map[string]string{"student": "stu", "中国": "中文"}, map[string]string{"student": "stu", "中国": "中文"}},
	}
	// Add keywords from map
	for _, testItem := range testSet {
		keywordProcessor := NewKeywordProcessor()
		keywordProcessor.AddKeywordsFromMap(testItem.in)
		assert.Equal(t, testItem.expected, keywordProcessor.GetAllKeywords())
	}
}

func TestAddKeywordsFromFile(t *testing.T) {

	var testSet = []struct{
		in			string
		expected	map[string]string
	} {
		{"examples/data.txt", map[string]string{"abc": "abc", "中国": "中文",}},
	}
	// Extract keywords from sentence by searching TrieTree
	for _, testItem := range testSet {
		keywordProcessor := NewKeywordProcessor()
		keywordProcessor.AddKeywordsFromFile(testItem.in)
		assert.Equal(t, testItem.expected, keywordProcessor.GetAllKeywords())
	}
}

func TestRemoveKeywordsFromList(t *testing.T) {

	// add keywords from Map
	keywordMap := map[string]string{
		"teacher": "tea",
		"student": "stu",
		"中国": "中文",
	}

	var testSet = []struct{
		in			[]string
		expected	map[string]string
	} {
		{[]string{"teacher"},		map[string]string{"student": "stu", "中国": "中文"}},
		{[]string{"student", "teacher"},map[string]string{"中国": "中文"}},
	}
	// Remove keywords from list
	for _, testItem := range testSet {
		keywordProcessor := NewKeywordProcessor()
		keywordProcessor.AddKeywordsFromMap(keywordMap)
		keywordProcessor.RemoveKeywordFromList(testItem.in)
		assert.Equal(t, testItem.expected, keywordProcessor.GetAllKeywords())
	}
}

func TestRemoveKeyword(t *testing.T) {

	// add keywords from Map
	keywordMap := map[string]string{
		"teacher": "tea",
		"student": "stu",
		"中国": "中文",
	}

	var testSet = []struct{
		in			string
		expected	map[string]string
	} {
		{"teacher",		map[string]string{"student": "stu", "中国": "中文",}},
		{"student",		map[string]string{"teacher": "tea", "中国": "中文",}},
	}
	// remove a keyword
	for _, testItem := range testSet {
		keywordProcessor := NewKeywordProcessor()
		keywordProcessor.AddKeywordsFromMap(keywordMap)
		keywordProcessor.RemoveKeyword(testItem.in)
		assert.Equal(t, testItem.expected, keywordProcessor.GetAllKeywords())
	}
}

func TestExtractKeywords(t *testing.T) {
	keywordProcessor := NewKeywordProcessor()

	// add keywords from Map
	keywordMap := map[string]string{
		"teacher": "tea",
		"student": "stu",
	}
	keywordProcessor.AddKeywordsFromMap(keywordMap)

	// 添加中文关键词
	keywordProcessor.AddKeyword("中文", "中文")
	keywordProcessor.AddKeyword("abc")

	var testSet = []struct{
		in			string
		expected	[]string
	} {
		{"hello abc, what up", 				[]string{"abc"}},
		{"hello, 你会说中文吗？", 			[]string{"中文"}},
		{"hello, abc 你会说中文吗？ oHabc", []string{"abc", "中文", "abc"}},
	}
	// Extract keywords from sentence by searching TrieTree
	for _, testItem := range testSet {
		cleanNameList := keywordProcessor.ExtractKeywords(testItem.in)
		assert.Equal(t, testItem.expected, cleanNameList)
	}
}
