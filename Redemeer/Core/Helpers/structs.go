package Helpers

type SettingsStruct struct {
	License     string `json:"license,omitempty"`
	VccSettings struct {
		MaxClaims   int `json:"maxClaims,omitempty"`
		AuthRetries int `json:"authRetries,omitempty"`
	} `json:"vccSettings,omitempty"`
	CaptchaSettings struct {
		Service string `json:"service,omitempty"`
		APIKey  string `json:"apiKey,omitempty"`
	} `json:"captchaSettings,omitempty"`
	MiscellaneousSettings struct {
		Debug     bool `json:"debug,omitempty"`
		Proxyless bool `json:"proxyless,omitempty"`
		Threads   int  `json:"threads,omitempty"`
	} `json:"miscellaneousSettings,omitempty"`
}
