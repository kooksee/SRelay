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

	"github.com/ethereum/go-ethereum/crypto"
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
