package Solvers

import (
	"Redeemer/Core/Helpers"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Ab5Solver(ApiKey, SiteKey, RqData string) string {
	Proxy, err := Helpers.GetProxy()
	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return ""
	}

	client := http.Client{Timeout: time.Second * 60}
	type Response struct {
		Pass string `json:"pass"`
	}
	var passResponse Response

	for {
		baseURL := "https://api.ab5.wtf/solve"
		u, _ := url.Parse(baseURL)
		q := u.Query()

		q.Add("url", "https://discord.com")
		q.Add("sitekey", SiteKey)
		q.Add("proxy", "http://"+Proxy)
		q.Add("rqdata", RqData)
		q.Add("userAgent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9017 Chrome/108.0.5359.215 Electron/22.3.12 Safari/537.36")

		u.RawQuery = q.Encode()

		req, _ := http.NewRequest("GET", u.String(), nil)
		req.Header.Set("authorization", ApiKey)

		resp, _ := client.Do(req)
		body, _ := io.ReadAll(resp.Body)

		if strings.Contains(string(body), "pass") {
			_ = json.Unmarshal(body, &passResponse)
			return passResponse.Pass

		} else if strings.Contains(string(body), "error") {
			Helpers.LogError("Couldn't Solve Captcha, Error", string(body))
			return ""
		}
	}
}
