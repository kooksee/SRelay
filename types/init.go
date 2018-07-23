package types

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ClientPort struct {
	Port     int64  `json:"port,omitempty"`
	Protocol string `json:"proto,omitempty"`
}
