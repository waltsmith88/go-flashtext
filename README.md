# go-flashtext

This module can be used to replace keywords in sentences or extract keywords from sentences. It is based on the [FlashText algorithm](https://arxiv.org/abs/1711.00046).

Compared with standard [FlashText algorithm](https://arxiv.org/abs/1711.00046)， there are some differences which make [go-flashtext](https://github.com/waltsmith88/go-flashtext) more powerful:

* **Chinese** is support fully. [Python implement](https://github.com/vi3k6i5/flashtext#flashtext) supports Chinese not well.
* We break **nonWordBoundaries** in FlashText algorithm to make it more powerful, which means that keyword could contains char not in [_0-9a-zA-Z]. 
* We allow the **same keyword with different cleanNames** exists, which means keywords are not unique. We found this is very useful in Industry envs.



## Installation

To install GoFlashText package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed, then you can use the below Go command to install GoFlashText.

```shell
$ go get -u github.com/waltsmith88/go-flashtext
```

2. Import it in your code:

```go
imoprt gf "github.com/waltsmith88/go-flashtext"
```



## Usage

- Extract keywords

```go
package main

import (
	"fmt"
	gf "github.com/waltsmith88/go-flashtext"
)

func main() {
	// add keywords from Map
	keywordMap := map[string]string{
		"love": "love",
		"hello": "hello",
	}
	keywordProcessor := gf.NewKeywordProcessor()
	keywordProcessor.AddKeywordsFromMap(keywordMap)
	foundList := keywordProcessor.ExtractKeywords("I love coding.")
	fmt.Println(foundList)
}
// [love]
```

- Extract keywords With Chinese Support

```go
package main

import (
	"fmt"
    	gf "github.com/waltsmith88/go-flashtext"
)

func main() {
	// add keywords from Map
	keywordMap := map[string]string{
		"love": "love",
		"中国": "中文",
	}
	keywordProcessor := gf.NewKeywordProcessor()
	keywordProcessor.AddKeywordsFromMap(keywordMap)
	keywordProcessor.AddKeyword("love", "ove")
	foundList := keywordProcessor.ExtractKeywords("I Love 中国.")
	fmt.Println(foundList)
}
// [中文]
```

- Case Sensitive example

```go
package main

import (
	"fmt"
    	gf "github.com/waltsmith88/go-flashtext"
)

func main() {
	// add keywords from Map
	keywordMap := map[string]string{
		"love": "love",
		"中国": "中文",
	}
	keywordProcessor := gf.NewKeywordProcessor()
	keywordProcessor.SetCaseSensitive(false)
	keywordProcessor.AddKeywordsFromMap(keywordMap)
	keywordProcessor.AddKeyword("love", "ove")
	foundList := keywordProcessor.ExtractKeywords("I Love 中国.")
	fmt.Println(foundList)
}
// [love|ove 中文]
```

- Unique Keywords example

```go
func main() {
	// add keywords from Map
	keywordMap := map[string]string{
		"love": "love",
		"中国": "中文",
	}
	keywordProcessor := gf.NewKeywordProcessor()
	keywordProcessor.SetUniqueKeyword(true)
	keywordProcessor.SetCaseSensitive(false)
	keywordProcessor.AddKeywordsFromMap(keywordMap)
	keywordProcessor.AddKeyword("love", "ove")
	foundList := keywordProcessor.ExtractKeywords("I Love 中国.")
	fmt.Println(foundList)
}
// [ove 中文]
```

- Span of keywords extracted

```go
func main() {
	// add keywords from Map
	keywordMap := map[string]string{
		"love": "love",
		"中国": "中文",
	}
	keywordProcessor := gf.NewKeywordProcessor()
	keywordProcessor.AddKeywordsFromMap(keywordMap)
	sentence := "I love 中国."
	cleanNameRes := keywordProcessor.ExtractKeywordsWithSpanInfo(sentence)
	sentence1 := []rune(sentence)
	for _, resSpan := range cleanNameRes {
		fmt.Println(resSpan.CleanName, resSpan.StartPos, resSpan.EndPos, fmt.Sprintf("%c", sentence1[resSpan.StartPos:resSpan.EndPos]))
	}
}
// love 2 6 [l o v e]
// 中文 7 9 [中 国]
```

- Add Multiple Keywords simultaneously

```go
// way 1: from Map
keywordMap := map[string]string{
		"abcd": "abcd",
		"student": "stu",
	}
keywordProcessor.AddKeywordsFromMap(keywordMap)
// way 2: from Slice
keywordProcessor.AddKeywordsFromList([]string{"student", "abcd", "abc", "中文"})
// way 3: from file. Line: keyword => cleanName
keywordProcessor.AddKeywordsFromFile(filePath)
```

- To Remove keywords

```go
keywordProcessor.RemoveKeyword("abc")
keywordProcessor.RemoveKeywordFromList([]string{"student", "abcd", "abc", "中文"})
```

- To Replace keywords

```go
newSentence := keywordProcessor.ReplaceKeywords(sourceSentence)
```

- To check Number of terms in KeywordProcessor

```go
keywordProcessor.Len()
```

- To check if term is present in KeywordProcessor

```go
keywordProcessor.IsContains("abc")
```

- Get all keywords in dictionary

```go
keywordProcessor.GetAllKeywords()
```


More Examples about Usage in go-flashtext/examples/examples.go and you could have a taste by using following command:

```shell
$ go run examples/examples.go
```



## Test

```shell
$ git clone github.com/waltsmith88/go-flashtext
$ cd go-flashtext
$ go test -v
```



## Why not Regex?

It's a custom algorithm based on [Aho-Corasick algorithm](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm) and [Trie Dictionary](https://en.wikipedia.org/wiki/TrieDictionary).

![Benchmark](https://github.com/vi3k6i5/flashtext/raw/master/benchmark.png)



Time taken by FlashText to find terms in comparison to Regex.

[![https://thepracticaldev.s3.amazonaws.com/i/xruf50n6z1r37ti8rd89.png](https://camo.githubusercontent.com/53e63b19336a7dfbe5d2874b70d73e37c4cd744d/68747470733a2f2f74686570726163746963616c6465762e73332e616d617a6f6e6177732e636f6d2f692f7872756635306e367a31723337746938726438392e706e67)](https://camo.githubusercontent.com/53e63b19336a7dfbe5d2874b70d73e37c4cd744d/68747470733a2f2f74686570726163746963616c6465762e73332e616d617a6f6e6177732e636f6d2f692f7872756635306e367a31723337746938726438392e706e67)

Time taken by FlashText to replace terms in comparison to Regex.

[![https://thepracticaldev.s3.amazonaws.com/i/k44ghwp8o712dm58debj.png](https://camo.githubusercontent.com/28e8b327359b6f93bf3ac4733b92c5dec0576851/68747470733a2f2f74686570726163746963616c6465762e73332e616d617a6f6e6177732e636f6d2f692f6b343467687770386f373132646d35386465626a2e706e67)](https://camo.githubusercontent.com/28e8b327359b6f93bf3ac4733b92c5dec0576851/68747470733a2f2f74686570726163746963616c6465762e73332e616d617a6f6e6177732e636f6d2f692f6b343467687770386f373132646d35386465626a2e706e67)

Link to code for benchmarking the [Find Feature](https://gist.github.com/vi3k6i5/604eefd92866d081cfa19f862224e4a0) and [Replace Feature](https://gist.github.com/vi3k6i5/dc3335ee46ab9f650b19885e8ade6c7a).

The idea for this library came from the following [StackOverflow question](https://stackoverflow.com/questions/44178449/regex-replace-is-taking-time-for-millions-of-documents-how-to-make-it-faster).

## Citation

The original paper published on [FlashText algorithm](https://arxiv.org/abs/1711.00046).

```
@ARTICLE{2017arXiv171100046S,
   author = {{Singh}, V.},
    title = "{Replace or Retrieve Keywords In Documents at Scale}",
  journal = {ArXiv e-prints},
archivePrefix = "arXiv",
   eprint = {1711.00046},
 primaryClass = "cs.DS",
 keywords = {Computer Science - Data Structures and Algorithms},
     year = 2017,
    month = oct,
   adsurl = {http://adsabs.harvard.edu/abs/2017arXiv171100046S},
  adsnote = {Provided by the SAO/NASA Astrophysics Data System}
}
```

The article published on [Medium freeCodeCamp](https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f).



## Contribute

- Issue Tracker: <https://github.com/waltsmith88/go-flashtext/issues>
- Source Code: <https://github.com/waltsmith88/go-flashtext>



## License

The project is licensed under the MIT license.
