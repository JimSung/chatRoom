package apps

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type FilterMgr struct {
	replaceWord string
	dirtyWords  map[string]bool
}

func NewFilterMgr() *FilterMgr {
	filterMgr := &FilterMgr{
		dirtyWords:  make(map[string]bool),
		replaceWord: "*",
	}
	dirty := "../../config/dirtyWords.txt"

	// 读取脏词库
	log.Printf("Loading Dirty Words...")
	path, _ := os.Getwd()
	log.Printf("os.pwd:%v", path)
	file, err := os.Open(dirty)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	defer file.Close()

	// 逐行扫描
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		words := strings.ToUpper(strings.TrimSpace(scanner.Text())) // 均处理为大写
		if words != "" {
			filterMgr.dirtyWords[words] = true
		}
	}
	log.Printf("Dirty Words Loaded")
	return filterMgr
}

func (f *FilterMgr) Filter(in string) string {
	indexMap := make(map[int]struct{})
	for key := range f.dirtyWords {
		for _, index := range checkWord(strings.ToUpper(in), key) {
			indexMap[index] = struct{}{}
		}
	}
	wordsRune := []rune(in)
	for i := range indexMap {
		wordsRune[i] = '*'
	}
	return string(wordsRune)
}

func checkWord(words, dirtyKey string) (indexList []int) {
	n := 0
	for i := 0; i+len(dirtyKey) <= len(words); i++ {
		if words[i] == dirtyKey[0] && (len(dirtyKey) == 1 || words[i:i+len(dirtyKey)] == dirtyKey) {
			wordsByte := []byte(words)
			if i > 0 {
				n = utf8.RuneCount(wordsByte[0:i])
			}
			for j := 0; j < utf8.RuneCount(wordsByte[i:i+len(dirtyKey)]); j++ {
				indexList = append(indexList, n+j)
			}
			i += len(dirtyKey) - 1
		}
	}
	return
}
