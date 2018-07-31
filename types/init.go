package types

import (
	"net/url"

	"github.com/kooksee/cmn"
)

var NewKBuffer = cmn.NewKBuffer
var JsonUnmarshal = cmn.Json.Unmarshal
var JsonMarshal = cmn.Json.Marshal
var errs = cmn.Err.Err

func GetNodeID(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	if u.Scheme != "sp2p" {
		return "", errs("invalid URL scheme, want sp2p ")
	}
	// Parse the node ID from the user portion.
	if u.User == nil {
		return "", errs("does not contain node ID")
	}
	return u.User.String(), nil
}
