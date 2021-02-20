package main

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
  var a = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
  var b = []string{"2", "4", "5", "6", "8"}
  var result = difference(a, b)
  expext := []string{"1", "3", "7"}
  if !reflect.DeepEqual(result, expext) {
    t.Error("\nresult： ", result, "\nexpext： ", expext)
  }

  t.Log("TestDifference終了")
}

func TestChunked(t *testing.T) {
  var a = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
  var result = chunked(a, 3)
  expext := [][]string{{"1","2", "3"}, {"4", "5", "6"}, {"7", "8"}}
  if !reflect.DeepEqual(result, expext) {
    t.Error("\nresult： ", result, "\nexpext： ", expext)
  }

  t.Log("TestChunked終了")
}

func TestJoinStringSlice( t *testing.T) {
  var a = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
  var result = joinStringSlice(a)
  expext := "1,2,3,4,5,6,7,8"
  if result != expext {
    t.Error("\nresult： ", result, "\nexpext： ", expext)
  }

  t.Log("TestJoinStringSlice終了")
}


