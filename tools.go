package copple

import (
	"bytes"
	"encoding/json"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// timeOfDay day is "2006-01-02"
func TimeOfDay(day string) (int64, int64) {

	//time.Parse("2006-01-02 15:04:05", "2018-12-03 11:02:12") 记忆方式
	timeStrBegin := day + " 00:00:01"
	timeStrEnd := day + " 23:59:59"

	var err1, err2 error
	t1, err1 := time.Parse("2006-01-02 15:04:05", timeStrBegin)
	t2, err2 := time.Parse("2006-01-02 15:04:05", timeStrEnd)
	if err1 != nil || err2 != nil {
		return -1, -1
	}
	return t1.Unix() - 8*3600, t2.Unix() - 8*3600
}

func GetCommentInHtml(html string) []string {
	commentRegex := regexp.MustCompile(`<!--(.*?)-->`)
	matches := commentRegex.FindAllStringSubmatch(html, -1)
	var rv []string
	for _, match := range matches {
		if len(match) > 1 {
			rv = append(rv, strings.TrimSpace(match[1]))
		}
	}
	return rv
}

func Obj2Json(obj interface{}, indent bool) ([]byte, error) {
	var jsonBytes []byte
	var err error
	if indent {
		jsonBytes, err = json.MarshalIndent(obj, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(obj)
	}
	if err != nil {
		return nil, err
	}
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003c"), []byte("<"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003e"), []byte(">"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u0026"), []byte("&"), -1)
	return jsonBytes, nil
}

func SleepRandMS(m int, n int) {
	sleepDuration := time.Duration(rand.Intn(m)+n-m) * time.Millisecond
	time.Sleep(sleepDuration)
}

func findMetaRefresh(body string) string {
	re := regexp.MustCompile(`<meta[^>]+http-equiv=["']?refresh["']?[^>]+content=["']?[^;]*;\s*url=([^"'>]+)["']?`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func GBK2UTF8(gbkStr string) (string, error) {
	reader := transform.NewReader(strings.NewReader(gbkStr), simplifiedchinese.GBK.NewDecoder())
	utf8Bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Bytes), nil
}

func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}
