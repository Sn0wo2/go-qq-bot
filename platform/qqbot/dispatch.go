package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sn0wo2/go-qq-bot/platform/qqbot/token"
	"github.com/tidwall/gjson"
)

func Dispatch(w http.ResponseWriter, r *http.Request, payload gjson.Result) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := gjson.ParseBytes(body)

	type Message struct {
		Content string `json:"content,omitempty"`
		MsgType uint   `json:"msg_type,omitempty"`
		MsgID   string `json:"msg_id,omitempty"`
	}

	msgPayload, err := json.Marshal(Message{
		Content: data.Get("d.content").Str,
		MsgType: 0,
		MsgID:   data.Get("d.id").Str,
	})
	if err != nil {
		panic(err)
	}

	var url string

	switch data.Get("t").Str {
	case "GROUP_AT_MESSAGE_CREATE":
		url = fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/messages", data.Get("d.group_id").Str)
	default:
		url = fmt.Sprintf("https://api.sgroup.qq.com/v2/users/%s/messages", data.Get("d.author.union_openid"))
	}

	fmt.Println(url, "\n", string(msgPayload))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(msgPayload))
	if err != nil {
		panic(err)
	}
	t, err := token.Instance.GetToken()
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("QQBot %s", t))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
