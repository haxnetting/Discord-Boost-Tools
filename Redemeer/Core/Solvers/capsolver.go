package Solvers

import (
	"Redeemer/Core/Helpers"
	capsolver_go "github.com/capsolver/capsolver-go"
	"strings"
)

func Capsolver(ApiKey, SiteKey, RqData string) string {
	Proxy, err := Helpers.GetProxy()
	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return ""
	}

	capSolver := capsolver_go.CapSolver{ApiKey}
	s, err := capSolver.Solve(
		map[string]any{
			"type":        "HCaptchaTurboTask",
			"websiteURL":  "https://discord.com/",
			"websiteKey":  SiteKey,
			"isInvisible": false,
			"enterprisePayload": map[string]any{
				"rqdata": RqData,
			},
			"proxy":     "http:" + strings.Split(Proxy, "@")[1] + ":" + strings.Split(Proxy, "@")[0],
			"userAgent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9017 Chrome/108.0.5359.215 Electron/22.3.12 Safari/537.36",
		})

	if err != nil {
		Helpers.LogError(err.Error(), "N/A")
		return "error"
	}

	Helpers.LogDebug("Captcha Task ID: " + s.TaskId)

	return s.Solution.GRecaptchaResponse
}
