package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/huichen/sego"
	"github.com/ipfs/go-ipfs-chunker"
)

// 测试数据的切片和数据合并以及数据备份和metadata的生成
func TestName4(t *testing.T) {
	f := "main"
	//bf := "mainexe"

	d, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err.Error())
	}

	rootHash := crypto.Keccak256(d)
	fmt.Println(hex.EncodeToString(crypto.Keccak256(d)))

	f1, _ := os.Open(f)
	r := chunk.NewRabin(f1, 1024*256)

	//var chunks []byte

	for i := 1; ; i++ {
		ck, err := r.NextBytes()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}

		fmt.Println(i)

		//fmt.Println(hex.EncodeToString(crypto.Keccak256(ck)))
		bf := hex.EncodeToString(crypto.Keccak256(append(rootHash, big.NewInt(int64(i)).Bytes()...)))
		if err := ioutil.WriteFile(filepath.Join("bk", bf), ck, 0755); err != nil {
			panic(err.Error())
		}

		//fmt.Println(len(ck))

		//chunks = append(chunks, ck...)
		//fmt.Println(string(ck))
	}
	//if err := ioutil.WriteFile(bf, chunks, 0755); err != nil {
	//	panic(err.Error())
	//}
}

func TestName5(t *testing.T) {
	const n = 81
	const bk = "bk"
	rootHash, _ := hex.DecodeString("85e073d11db2a8d8e1fe5acb3ce0176076586f8c48b6a392545b999fb4229989")

	var cks []byte
	for i := 1; i <= n; i++ {
		bf := hex.EncodeToString(crypto.Keccak256(append(rootHash, big.NewInt(int64(i)).Bytes()...)))
		if d, err := ioutil.ReadFile(filepath.Join("bk", bf)); err != nil {
			panic(err.Error())
		} else {
			cks = append(cks, d...)
		}
	}

	if err := ioutil.WriteFile("kkkk", cks, 0755); err != nil {
		panic(err.Error())
	}
}

// 全文检索功能测试
func TestName11(t *testing.T) {
	message := struct {
		Id   string
		From string
		Body string
	}{
		Id:   "example",
		From: "marty.schoch@gmail.com",
		Body: "bleve indexing is easy",
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		panic(err)
	}
	if err := index.Index(message.Id, message); err != nil {
		panic(err.Error())
	}
}

// 全文检索查询
func TestName22(t *testing.T) {
	index, _ := bleve.Open("example.bleve")
	query := bleve.NewQueryStringQuery("bleve")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult.String())
	fmt.Println("ok")
}

// sego 中文分词
func TestName44(t *testing.T) {
	// 载入词典
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("/Users/barry/gopath/src/github.com/huichen/sego/data/dictionary.txt")

	// 分词
	text := []byte("中华人民共和国中央人民政府")
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	fmt.Println(sego.SegmentsToString(segments, true))
	fmt.Println(sego.SegmentsToSlice(segments, true))

	message := struct {
		Id   string
		Body string
	}{
		Id:   "example12345",
		Body: "中华人民共和国中央人民政府",
	}

	index, err := bleve.Open("example.bleve")
	if err != nil {
		panic(err)
	}
	if err := index.Index(message.Id, message); err != nil {
		panic(err.Error())
	}

	query := bleve.NewQueryStringQuery("人民 共和国")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	fmt.Println(searchResult.String())
	fmt.Println("ok")
}
