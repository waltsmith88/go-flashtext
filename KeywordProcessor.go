package flashtext

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const __key__ = int32(197)

//KeywordRes Result of Extraction
type KeywordRes struct {
	CleanName string
	StartPos  int
	EndPos    int
}

// Keyword Processor
type keywordProcessor struct {
	caseSensitive     bool
	uniqueKeyword     bool
	_keyword          string
	whiteSpaceChars   []string
	nonWordBoundaries string
	keywordTrieDict   map[int32]interface{}
	termsInTrie       int
	delimiter         string
}

type KeywordProcessor interface {
	GetKeywordTrieDict() map[int32]interface{}
	GetCaseSensitive() bool
	SetCaseSensitive(bool)
	GetUniqueKeyword() bool
	SetUniqueKeyword(bool)
	GetDelimiter() string
	SetDelimiter(string)
	Len() int
	IsContains(string) bool
	GetKeyword(string) (string, bool)
	GetAllKeywords() map[string]string
	AddKeyword(string, ...string) bool
	AddKeywordsFromMap(map[string]string)
	AddKeywordsFromList([]string)
	AddKeywordsFromFile(string)
	RemoveKeywordFromList([]string)
	RemoveKeyword(string) bool
	ExtractKeywords(string) []string
	ExtractKeywordsWithSpanInfo(string) []KeywordRes
	ReplaceKeywords(string) string
}

func NewKeywordProcessor() KeywordProcessor {

	kp := &keywordProcessor{
		caseSensitive:     true,
		uniqueKeyword:     false,
		_keyword:          "_keyword_",
		whiteSpaceChars:   []string{".", "\t", "\n", "\a", " ", ","},
		nonWordBoundaries: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz",
		keywordTrieDict:   map[int32]interface{}{},
		termsInTrie:       0,
		delimiter:         "|",
	}
	return kp
}

// Default returns an keywordProcessor instance with default values.
func Default() KeywordProcessor {
	kp := NewKeywordProcessor()
	return kp
}

// Get keyword TrieTree
func (KeywordProcessor *keywordProcessor) GetKeywordTrieDict() map[int32]interface{} {
	return KeywordProcessor.keywordTrieDict
}

// Get case sensitive
func (KeywordProcessor *keywordProcessor) GetCaseSensitive() bool {
	return KeywordProcessor.caseSensitive
}

// Set case sensitive
func (KeywordProcessor *keywordProcessor) SetCaseSensitive(sensitive bool) {
	KeywordProcessor.caseSensitive = sensitive
}

// Get case sensitive
func (KeywordProcessor *keywordProcessor) GetUniqueKeyword() bool {
	return KeywordProcessor.uniqueKeyword
}

// Set case sensitive
func (KeywordProcessor *keywordProcessor) SetUniqueKeyword(sensitive bool) {
	KeywordProcessor.uniqueKeyword = sensitive
}

// Get the delimiter which is used to joining two different cleanNames of same keyword.
func (KeywordProcessor *keywordProcessor) GetDelimiter() string {
	return KeywordProcessor.delimiter
}

// Set the delimiter which is used to joining two different cleanNames of same keyword.
// Be careful of setting the delimiter
// and make sure the delimiter be the unique identifier of all cleanNames.
func (KeywordProcessor *keywordProcessor) SetDelimiter(delimiter string) {
	KeywordProcessor.delimiter = delimiter
}

// Return number of terms present in the keyword_trie_dict
func (KeywordProcessor *keywordProcessor) Len() int {
	return KeywordProcessor.termsInTrie
}

// If TrieTree contains keyword
func (KeywordProcessor *keywordProcessor) IsContains(keyword string) bool {
	currentDict := KeywordProcessor.keywordTrieDict
	for _, val := range keyword {

		if tmp_, err := currentDict[val]; !err {
			return false
		} else {
			currentDict = tmp_.(map[int32]interface{})
		}
	}
	_, err := currentDict[__key__]
	return err
}

// Get clean Name of keyword from TrieTree.
// Return (cleanName, bool).
func (KeywordProcessor *keywordProcessor) GetKeyword(keyword string) (string, bool) {
	currentDict := KeywordProcessor.keywordTrieDict
	for _, val := range keyword {

		if tmp_, err := currentDict[val]; !err {
			return "nil", false
		} else {
			currentDict = tmp_.(map[int32]interface{})
		}
	}
	if res, err := currentDict[__key__]; err {
		return res.(string), true
	} else {
		return "nil", false
	}
}

// Get all clean Name of keyword from TrieTree.
// Return map(keyword, cleanName).
func (KeywordProcessor *keywordProcessor) GetAllKeywords() map[string]string {
	return KeywordProcessor.__getAllKeywords__("", KeywordProcessor.keywordTrieDict)
}

// Get all clean Name of keyword from TrieTree.
// Return map(keyword, cleanName).
func (KeywordProcessor *keywordProcessor) __getAllKeywords__(termSoFar string, currendDict map[int32]interface{}) map[string]string {
	var allKeywords = map[string]string{}
	for key := range currendDict {
		if key == __key__ {
			allKeywords[termSoFar] = currendDict[__key__].(string)
		} else {
			subValues := KeywordProcessor.__getAllKeywords__(
				termSoFar+fmt.Sprintf("%c", key),
				currendDict[key].(map[int32]interface{}))
			for key := range subValues {
				allKeywords[key] = subValues[key]
			}
		}
	}
	return allKeywords
}

// Build the map relation between keyword and cleanName in TrieTree
func (KeywordProcessor *keywordProcessor) __setItem__(keyword string, cleanName string) map[int32]interface{} {
	keyReverse := __reverseString__(keyword)
	var resMap map[int32]interface{}
	for i, j := range keyReverse {
		if i == 0 {
			tmp := map[int32]interface{}{__key__: cleanName}
			resMap = map[int32]interface{}{j: tmp}
		} else {
			resMap = map[int32]interface{}{j: resMap}
		}
	}
	return resMap
}

// Add a keyword and its' cleanName to TrieTree
func (KeywordProcessor *keywordProcessor) AddKeyword(keyword string, cleanNames ...string) bool {
	var (
		cleanName string
		diffDict  map[int32]interface{}
		commDict  map[int32]interface{}
	)
	if len(cleanNames) == 0 {
		cleanName = keyword
	} else {
		cleanName = cleanNames[0]
	}

	if !KeywordProcessor.caseSensitive {
		keyword = strings.ToLower(keyword)
	}
	currentDict := KeywordProcessor.keywordTrieDict

	for i, letter := range keyword {
		if currentDict_, err := currentDict[letter]; err {
			currentDict = currentDict_.(map[int32]interface{})
			commDict = currentDict
		} else {
			diffDict = KeywordProcessor.__setItem__(keyword[i:], cleanName)
			break
		}
	}

	if commDict == nil {
		if currentDict == nil {
			currentDict = diffDict
		} else {
			for k, v := range diffDict {
				currentDict[k] = v
			}
		}
		KeywordProcessor.keywordTrieDict = currentDict

		KeywordProcessor.termsInTrie++
		return true
	}

	if diffDict == nil {
		if tmpCleanName, err := commDict[__key__]; err {
			if tmpCleanName != cleanName && !KeywordProcessor.uniqueKeyword {
				commDict[__key__] = strings.Join([]string{tmpCleanName.(string), cleanName}, "|")
			} else { // not unique keyword
				commDict[__key__] = cleanName
			}

			KeywordProcessor.termsInTrie++
			return true
		} else {
			diffDict = map[int32]interface{}{__key__: cleanName}
		}
	}

	for k, v := range diffDict {
		commDict[k] = v
	}

	KeywordProcessor.termsInTrie++
	return true
}

// Add multiple keywords and its' respective cleanName from keyword Map to TrieTree
func (KeywordProcessor *keywordProcessor) AddKeywordsFromMap(keywordMap map[string]string) {
	for keyword, cleanName := range keywordMap {
		KeywordProcessor.AddKeyword(keyword, cleanName)
	}
}

// Add multiple keywords and cleanName same of keywords from keyword list to TrieTree
// (keyword) -> (keyword, keyword) -> (keyword, cleanName)
func (KeywordProcessor *keywordProcessor) AddKeywordsFromList(keywordList []string) {
	for _, keyword := range keywordList {
		KeywordProcessor.AddKeyword(keyword, keyword)
	}
}

// Add multiple keywords and cleanName same of keywords from keyword list to TrieTree
// each line in file:
// 		1. keyword => cleanName
//		2. keyword
func (KeywordProcessor *keywordProcessor) AddKeywordsFromFile(filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewReader(file)
	for {
		line, _, err := w.ReadLine()
		if err == io.EOF {
			break
		}
		if strings.Contains(string(line), "=>") { // 1. keyword => cleanName
			lineString := strings.Split(string(line), "=>")
			keyword := strings.TrimSpace(lineString[0])
			cleanName := strings.TrimSpace(lineString[1])
			KeywordProcessor.AddKeyword(keyword, cleanName)
		} else { // 2. keyword
			keyword := strings.TrimSpace(string(line))
			KeywordProcessor.AddKeyword(keyword, keyword)
		}
	}
}

// Delete keywords and their cleanName from TrieTree
func (KeywordProcessor *keywordProcessor) RemoveKeywordFromList(keywordList []string) {
	for _, keyword := range keywordList {
		KeywordProcessor.RemoveKeyword(keyword)
	}
}

// Delete a keyword and its' cleanName from TrieTree
func (KeywordProcessor *keywordProcessor) RemoveKeyword(keyword string) bool {
	var (
		commDictKey   []int32
		commDictValue []map[int32]interface{}
	)
	currentDict := KeywordProcessor.keywordTrieDict

	for _, letter := range keyword {
		commDictValue = append(commDictValue, currentDict)
		if tmp_, err := currentDict[letter]; !err {
			return false
		} else {
			commDictKey = append(commDictKey, letter)
			currentDict = tmp_.(map[int32]interface{})
		}
	}

	if _, err := currentDict[__key__]; err {
		delete(currentDict, __key__)
		for i := len(commDictKey) - 1; i >= 0; i-- {
			tmpDict := commDictValue[i]
			delDict := tmpDict[commDictKey[i]].(map[int32]interface{})
			if len(delDict) > 0 {
				return true
			} else {
				delete(tmpDict, commDictKey[i])
			}
		}
		return true
	}
	return false
}

// Extract keywords from sentence by searching TrieTree.
// And return the keywords' clean names.
func (KeywordProcessor *keywordProcessor) ExtractKeywords(sentence string) []string {
	var keywordList []string
	if len(sentence) == 0 {
		return keywordList
	}
	if !KeywordProcessor.caseSensitive {
		sentence = strings.ToLower(sentence)
	}

	var (
		start        []int
		sentenceRune = []rune(sentence)
		idx          = 0
		idy          = 0
		sentenceLen  = len(sentenceRune)
		cleanName    = ""
		currentDict  = KeywordProcessor.keywordTrieDict
	)

	for idx < sentenceLen {
		char := sentenceRune[idx]
		tmpCurrent, err := currentDict[int32(char)].(map[int32]interface{})
		if err {
			start = append(start, idx)
			idx++
			idy++
			if cleanNameTmp, err := tmpCurrent[__key__]; err {
				cleanName = cleanNameTmp.(string)
			}
			currentDict = tmpCurrent
			if idx < sentenceLen {
				continue
			}
		} else {
			idx++
			idy++
			currentDict = KeywordProcessor.keywordTrieDict
		}

		if cleanName != "" {
			keywordList = append(keywordList, cleanName)
			idx = start[len(start)-1] + 1
			idy = idx + 1
			cleanName = ""
			start = []int{}
		} else {
			if len(start) > 0 {
				idx = start[0] + 1
				idy = idx + 1
			}
			start = []int{}
		}
	}
	return keywordList
}

// Extract keywords from sentence by searching TrieTree.
// And return the keywords' clean names, the start position and the end position of keyword in sentence.
func (KeywordProcessor *keywordProcessor) ExtractKeywordsWithSpanInfo(sentence string) []KeywordRes {
	var keywordList []KeywordRes
	if len(sentence) == 0 {
		return keywordList
	}
	if !KeywordProcessor.caseSensitive {
		sentence = strings.ToLower(sentence)
	}

	var (
		start        []int
		idx          = 0
		idy          = 0
		cleanName    = ""
		sentenceRune = []rune(sentence)
		sentenceLen  = len(sentenceRune)
		currentDict  = KeywordProcessor.keywordTrieDict
	)

	for idx < sentenceLen {
		char := sentenceRune[idx]
		tmpCurrent, err := currentDict[int32(char)].(map[int32]interface{})
		if err {
			start = append(start, idx)
			idx++
			idy++
			if cleanNameTmp, err := tmpCurrent[__key__]; err {
				cleanName = cleanNameTmp.(string)
			}
			currentDict = tmpCurrent
			continue
		} else {
			idx++
			idy++
			currentDict = KeywordProcessor.keywordTrieDict
		}
		if cleanName != "" {
			startIndex := start[0]
			endIndex := start[0] + len(start)
			res := KeywordRes{cleanName, startIndex, endIndex}
			keywordList = append(keywordList, res)
			idx = start[len(start)-1] + 1
			idy = idx
			cleanName = ""
			start = []int{}
		} else {
			if len(start) > 0 {
				idx = start[0] + 1
				idy = idx
			}
			start = []int{}
		}
	}
	return keywordList
}

// Replace keywords in sentence with their cleanName
func (KeywordProcessor *keywordProcessor) ReplaceKeywords(sentence string) string {
	var (
		newSentence = ""
	)
	if len(sentence) == 0 {
		return newSentence
	}
	if !KeywordProcessor.caseSensitive {
		sentence = strings.ToLower(sentence)
	}

	var (
		start        []int
		replaceIndex int
		sentenceRune = []rune(sentence)
		idx          = 0
		idy          = 0
		sentenceLen  = len(sentenceRune)
		cleanName    = ""
		currentDict  = KeywordProcessor.keywordTrieDict
	)

	for idx < sentenceLen {
		char := sentenceRune[idx]
		tmpCurrent, err := currentDict[int32(char)].(map[int32]interface{})
		if err {
			start = append(start, idx)
			idx++
			idy++
			if cleanNameTmp, err := tmpCurrent[__key__]; err {
				cleanName = cleanNameTmp.(string)
			}
			currentDict = tmpCurrent
			if idx < sentenceLen {
				continue
			}
		} else {
			replaceIndex = idx
			idx++
			idy++
			currentDict = KeywordProcessor.keywordTrieDict
		}
		if cleanName != "" {
			newSentence += cleanName
			idx = start[len(start)-1] + 1
			idy = idx + 1
			cleanName = ""
			start = []int{}
		} else {
			if len(start) > 0 {
				idx = start[0] + 1
				idy = idx + 1
				replaceIndex = start[0]
			}
			start = []int{}
			if replaceIndex < sentenceLen {
				newSentence += string(sentenceRune[replaceIndex])
			}
		}
	}
	return newSentence
}

// Reverse the given string
func __reverseString__(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
