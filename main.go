package main

import (
	"fmt"
	"io/ioutil"
  "time"
	"net/http"
	"net/url"
	"encoding/json"
  "github.com/koron/go-dproxy"
)

const openbdEndpoint = "https://api.openbd.jp/v1"

var count int = 0
var total int = 0

func main() {
  startTime := time.Now()
  logMessage(fmt.Sprintf("start %v", startTime))
  var coverage = getCoverage()
  logMessage(fmt.Sprintf("coverage %v", len(coverage)))
  var knownIsbnList = GetIsbnList()
  logMessage(fmt.Sprintf("knownIsbnList %v", len(knownIsbnList)))
  coverage = difference(coverage, knownIsbnList)
  logMessage(fmt.Sprintf("defference %v", len(coverage)))
  total = len(coverage)

  var chunkedCoverage = chunked(coverage, 10000)
  var isbns string
  for _, chunk := range chunkedCoverage {
    isbns = joinStringSlice(chunk)
    getBookInfo(isbns)
  }
  logMessage(fmt.Sprintf("開始時刻 %v", startTime))
  logMessage(fmt.Sprintf("終了時刻 %v", time.Now()))
}

// go - How to find the difference between two slices of strings - Stack Overflow https://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings
// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}
// リストをn個ずつに分割する
func chunked(list []string, n int) [][]string {
  var chunkedList = [][]string{}

  for x := 0; x < len(list); x = x + n {
    var chunk []string
    if (x + n) < len(list) {
      chunk = list[x:x+n]
    } else {
      chunk = list[x:]
    }
    var next = x + n
    if next > len(list) {
      next = len(list)
    }
    chunkedList = append(chunkedList, chunk)
  }

  return chunkedList
}

// 文字列のスライスを連結する
// 参照: Goでは文字列連結はコストの高い操作 - Qiita https://qiita.com/ruiu/items/2bb83b29baeae2433a79
func joinStringSlice(a []string) string {
  b := make([]byte, 0, 10)

  for i := 0; i < len(a); i++ {
    b = append(b, a[i]...)
    b = append(b, ',')
  }

  // 末尾のカンマ削除
  b = b[:len(b)-1]

  return string(b)
}

// 指定したISBNの書誌情報を取得する
func getBookInfo(isbn string) {
  var book BookInfo
  params := url.Values{}
  params.Add("isbn", isbn)

  logMessage("send http request")
  resp, err := http.PostForm(openbdEndpoint + "/get", params)
  if err != nil {
	  logMessage(fmt.Sprintf("%v", err))
    return
  }
  defer resp.Body.Close()
  
  jsonBlob, err := ioutil.ReadAll(resp.Body)
  var obj interface{}
  json.Unmarshal(jsonBlob, &obj)

  length := len(obj.([]interface{}))
  for i := 0; i < length; i++ {
    // TODO dproxy なしの場合と性能比較する
    // isbn := obj.([]interface{})[0].(map[string]interface{})["summary"].(map[string]interface{})["isbn"].(string)
    p := dproxy.New(obj).A(i)
    var d dproxy.Drain

    book.Isbn        = d.String(p.M("summary").M("isbn"))
    book.Title       = d.String(p.M("summary").M("title"))
    book.Author      = d.String(p.M("summary").M("author"))
    book.Publisher   = d.String(p.M("summary").M("publisher"))
    book.Pubdate     = d.String(p.M("summary").M("pubdate"))
    book.Cover       = d.String(p.M("summary").M("cover"))
    book.Keywords    = ""
    book.Ccode       = ""
    book.Genre       = ""
    book.Description = ""
    book.Contents    = ""

    var ps dproxy.ProxySet
    ps = p.M("onix").M("DescriptiveDetail").M("Subject").ProxySet()
    if !ps.Empty() {
      for j := 0; j < ps.Len(); j++ {
        switch d.String(ps.A(j).M("SubjectSchemeIdentifier")) { 
          case "20":
            book.Keywords = d.String(ps.A(j).M("SubjectHeadingText"))
          case "78":
            book.Ccode = d.String(ps.A(j).M("SubjectCode"))
          case "79":
            book.Genre = d.String(ps.A(j).M("SubjectCode"))
        }
      }
    }
    
    ps = p.M("onix").M("CollateralDetail").M("TextContent").ProxySet()
    if !ps.Empty() {
      for j := 0; j < ps.Len(); j++ {
        switch d.String(ps.A(j).M("TextType")) { 
          case "03":
            book.Description = d.String(ps.A(j).M("Text"))
          case "04":
            book.Contents = d.String(ps.A(j).M("Text"))
        }
      }
    }

    if err := d.CombineErrors(); err != nil {
      logMessage(fmt.Sprintf("%v", err))
      ps = p.M("onix").M("DescriptiveDetail").M("Subject").ProxySet()
      fmt.Println("ps", ps)
      //return
      continue
    }
    InsertBook(book)
    count++
    //logMessage(count, total)
    //logMessage(book)
  }

}

// OpenBDからcoverageを取得する
func getCoverage() []string {
  var coverage []string
  resp, err := http.Get(openbdEndpoint + "/coverage")
  if err != nil {
		panic(err.Error())
  }
  defer resp.Body.Close()

  jsonBlob, err := ioutil.ReadAll(resp.Body)
  if err != nil {
		panic(err.Error())
  }
  if err := json.Unmarshal(jsonBlob, &coverage); err != nil {
		panic(err.Error())
	}
  return coverage
}

func logMessage(msg ...string) {
  now := time.Now()
  fmt.Println(now, msg)
}
