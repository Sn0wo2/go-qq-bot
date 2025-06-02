package token

import (
	"bytes"
	"encoding/json"
	"github.com/Sn0wo2/go-qq-bot/config"
	"github.com/Sn0wo2/go-qq-bot/types"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

var Instance = NewTokenProvider(func() (*Token, error) {
	accessPayload, err := json.Marshal(types.APIAccessPayload{AppId: config.AppID, ClientSecret: config.BotSecret})
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Post("https://bots.qq.com/app/getAppAccessToken", "application/json", bytes.NewBuffer(accessPayload))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	payload := gjson.ParseBytes(body)
	return &Token{
		AccessToken: payload.Get("access_token").Str,
		ExpireAt:    time.Now().Add(time.Duration(payload.Get("expires_in").Int()) - 30*time.Second),
	}, nil
})

type Token struct {
	AccessToken string
	ExpireAt    time.Time
}

type Provider struct {
	mu           sync.Mutex
	cachedToken  *Token
	fetchTokenFn func() (*Token, error)
}

func NewTokenProvider(fetchFn func() (*Token, error)) *Provider {
	return &Provider{
		fetchTokenFn: fetchFn,
	}
}

func (p *Provider) GetToken() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.cachedToken == nil || now.After(p.cachedToken.ExpireAt) {
		token, err := p.fetchTokenFn()
		if err != nil {
			return "", err
		}
		p.cachedToken = token
	}

	return p.cachedToken.AccessToken, nil
}
