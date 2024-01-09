package Solvers

import (
	"Redeemer/Core/Helpers"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func Hcoptcha(ApiKey, SiteKey, RqData string) string {
	Proxy, err := Helpers.GetProxy()
	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return ""
	}

	client := http.Client{Timeout: 30 * time.Second}

	payload := map[string]interface{}{
		"task_type": "hcaptchaEnterprise",
		"api_key":   ApiKey,
		"data": map[string]string{
			"sitekey": SiteKey,
			"url":     "https://discord.com",
			"proxy":   Proxy,
			"rqdata":  RqData,
		},
	}
	jsonPayload, err := json.Marshal(payload)

	req1, _ := http.NewRequest("POST", "https://api.hcoptcha.online/api/createTask", bytes.NewBuffer(jsonPayload))
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := client.Do(req1)
	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return ""
	}
	bodytext1, _ := io.ReadAll(resp1.Body)

	defer resp1.Body.Close()

	type hcopResponse struct {
		Error  bool   `json:"error"`
		TaskID string `json:"task_id"`
	}

	var hcopresp hcopResponse
	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return ""
	}

	err = json.Unmarshal(bodytext1, &hcopresp)

	if hcopresp.TaskID != "" {
		Helpers.LogDebug("Captcha Task ID: " + hcopresp.TaskID)
		for i := 0; i < 10; i++ {
			json2 := map[string]interface{}{
				"api_key": ApiKey,
				"task_id": hcopresp.TaskID,
			}

			p2, _ := json.Marshal(json2)

			req2, _ := http.NewRequest("POST", "https://api.hcoptcha.online/api/getTaskData", bytes.NewReader(p2))
			req2.Header.Set("Content-Type", "application/json")

			resp2, err := client.Do(req2)
			if err != nil {
				Helpers.LogError(err.Error(), "N/A")
				return ""
			}

			defer resp2.Body.Close()

			type hcopSol struct {
				Error bool `json:"error"`
				Task  struct {
					CaptchaKey string `json:"captcha_key"`
					Refunded   bool   `json:"refunded"`
					State      string `json:"state"`
				} `json:"task"`
			}

			var hcopsolution hcopSol
			bodytext2, _ := io.ReadAll(resp2.Body)

			err = json.Unmarshal(bodytext2, &hcopsolution)
			if err != nil {
				Helpers.LogError(err.Error(), "N/A")
				return ""
			}

			if hcopsolution.Error {
				return "error"
			}

			if hcopsolution.Task.State != "completed" && hcopsolution.Task.State != "processing" {
				Helpers.LogError("Couldn't Solved Captcha, Retrying", "N/A")
				return ""
			} else if hcopsolution.Task.State == "completed" {
				return hcopsolution.Task.CaptchaKey
			} else {
				time.Sleep(time.Second * 2)
				continue
			}
		}
	} else {
		Helpers.LogError("Couldn't Get Captcha Task ID, Check API Key or Contact Support", "N/A")
	}

	return ""
}
