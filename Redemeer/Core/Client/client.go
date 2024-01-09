package Client

import (
	"Redeemer/Core/Helpers"
	"Redeemer/Core/Solvers"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

var settings, _ = Helpers.LoadSettings()

const StripeKey = "pk_live_CUQtlpQUF0vufWpnpUmQvcdi"

func MakeSuperProperties() string {
	text := fmt.Sprintf(`{"os":"Windows","browser":"Chrome","device":"","system_locale":"en-US","browser_user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36","browser_version":"107.0.0.0","os_version":"10","referrer":"https://www.google.com/","referring_domain":"www.google.com","search_engine":"google", "referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number": 253927,"client_event_source":null}`)
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func NewClient(token string, promo string, vcc string, folder string) (RedeemerStruct, error) {
	var proxy string
	var err error

	jar := tls_client.NewCookieJar([]tls_client.CookieJarOption{}...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(60),
		tls_client.WithClientProfile(profiles.Chrome_117),
		tls_client.WithCookieJar(jar),
		tls_client.WithInsecureSkipVerify(),
	}

	if !settings.MiscellaneousSettings.Proxyless {
		proxy, err = Helpers.GetProxy()
		if err != nil {
			t := Helpers.FormatToken(token)
			Helpers.LogError(err.Error(), Helpers.Replacelast(t))
			return RedeemerStruct{}, err
		}
		options = append(options, tls_client.WithProxyUrl("http://"+proxy))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		go Helpers.AppendTextToFile(token+"\n", folder+"/Failed_To_Start.txt")
		go Helpers.RemoveLine(token, "./Data/Input/Tokens.txt")
		return RedeemerStruct{}, err
	}

	VccCard, VccCvv, VccMonth, VccYear := Helpers.ParseVcc(vcc)

	ParsedPromo, err := Helpers.ParsePromo(promo)
	if err != nil {
		go Helpers.AppendTextToFile(token+"\n", folder+"/Failed_To_Start.txt")
		go Helpers.RemoveLine(token, "./Data/Input/Tokens.txt")
		return RedeemerStruct{}, err
	}

	client.SetCookies(&url.URL{Path: "/", Host: "discord.com", Scheme: "https"}, []*http.Cookie{{Name: "locale", Value: "en-US"}})

	r := RedeemerStruct{
		Client:           client,
		Token:            Helpers.FormatToken(token),
		UnformattedToken: token,
		UnformattedPromo: promo,
		StartTime:        time.Now(),
		Promo:            ParsedPromo,
		VccInfo:          vcc,
		Proxy:            proxy,
		Folder:           folder,
		Vcc:              VccInfo{VccCard: VccCard, VccCVV: VccCvv, VccExpiryMonth: VccMonth, VccExpiryYear: VccYear},
		SuperProperties:  MakeSuperProperties(),
	}

	return r, nil
}

func (r *RedeemerStruct) GetHeaders() http.Header {
	headers := http.Header{
		"authority":          {"discord.com"},
		"scheme":             {"https"},
		"accept":             {"*/*"},
		"accept-encoding":    {"gzip, deflate"},
		"accept-language":    {"en-US"},
		"authorization":      {r.Token},
		"cookie":             {r.Cookies},
		"origin":             {"https://discord.com"},
		"sec-ch-ua":          {`"Chromium";v="107", "Google Chrome";v="107", "Not=A?Brand";v="99"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"Windows"},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-debug-options":    {"bugReporterEnabled"},
		"x-fingerprint":      {r.Fingerprints},
		"x-super-properties": {r.SuperProperties},
	}

	return headers

}

func (r *RedeemerStruct) GetStripeHeaders() http.Header {
	headers := http.Header{
		"authority":          {"api.stripe.com"},
		"scheme":             {"https"},
		"accept":             {"application/json"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en"},
		"content-type":       {"application/x-www-form-urlencoded"},
		"dnt":                {"1"},
		"origin":             {"https://js.stripe.com"},
		"referer":            {"https://js.stripe.com/"},
		"sec-ch-ua":          {`"Chromium";v="107", "Google Chrome";v="107", "Not=A?Brand";v="99"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"Windows"},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	}

	return headers
}

func (r *RedeemerStruct) CheckElements() error {
	var TokenResponse TokenCheckResponse
	var PromoResponse PromoCheckResponse

	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	req.Header.Set("authorization", r.Token)
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), r.Token)
		}
	}(resp.Body)

	if resp.StatusCode == 401 {
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Invalid_Token.txt")
		return errors.New("Invalid Token -> " + Helpers.Replacelast(r.Token))
	} else if resp.StatusCode == 429 {
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Rate_Limited.txt")
		return errors.New("Rate Limited Token -> " + Helpers.Replacelast(r.Token))
	}

	TokenBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(TokenBody, &TokenResponse)
	if err != nil {
		return errors.New("Failed to Get Token Check Response, Error: " + string(TokenBody))
	}

	if TokenResponse.PremiumType == 2 {
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Already_Nitro_On_Acc.txt")
		return errors.New("Token Already has Nitro -> " + Helpers.Replacelast(r.Token))
	}

	req2, err := http.NewRequest("GET", "https://discord.com/api/v9/entitlements/gift-codes/"+r.Promo+"?country_code=US&with_application=true&with_subscription_plan=true", nil)
	if err != nil {
		return err
	}

	resp2, err := r.Client.Do(req2)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp2.Body)

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &PromoResponse)
	if err != nil {
		return errors.New("Failed to Get Promo Check Response, Error: " + string(body))
	}

	if PromoResponse.Uses != 0 {
		go Helpers.AppendTextToFile(r.UnformattedPromo+"\n", r.Folder+"/Already_Redeemed_Promo.txt")
		go Helpers.RemoveLine(r.UnformattedPromo+"\n", "./Data/Input/Promos.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Promo Already Redeemed -> " + r.Promo)
	}

	return nil
}

func (r *RedeemerStruct) StripeCookies() error {
	req, err := http.NewRequest("POST", "https://m.stripe.com/6", nil)
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Get_Stripe_Cookies.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed to get Stripe Cookies, Error: " + string(body) + " -> " + Helpers.Replacelast(r.Token))
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Stripe.StripeGuid = jsonData["guid"].(string)
	r.Stripe.StripeMuid = jsonData["muid"].(string)
	r.Stripe.StripeSid = jsonData["sid"].(string)

	return nil
}

func (r *RedeemerStruct) GetStripeToken() error {
	headers := r.GetStripeHeaders()
	rand.Seed(time.Now().UnixNano())

	headers["authorization"] = []string{"Bearer " + StripeKey}

	parameters := url.Values{}
	parameters.Add("card[number]", r.Vcc.VccCard)
	parameters.Add("card[cvc]", r.Vcc.VccCVV)
	parameters.Add("card[exp_month]", r.Vcc.VccExpiryMonth)
	parameters.Add("card[exp_year]", r.Vcc.VccExpiryYear)
	parameters.Add("guid", r.Stripe.StripeGuid)
	parameters.Add("muid", r.Stripe.StripeMuid)
	parameters.Add("sid", r.Stripe.StripeSid)
	parameters.Add("payment_user_agent", "stripe.js/0651b07260; stripe-js-v3/0651b07260; split-card-element")
	parameters.Add("time_on_page", fmt.Sprintf("%d", rand.Intn(170000-130000+1)+130000))
	parameters.Add("key", StripeKey)

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/tokens?"+parameters.Encode(), nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Get_Stripe_Token.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed to get Stripe Token, Error: " + string(body) + " -> " + Helpers.Replacelast(r.Token))
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Stripe.StripeToken = jsonData["id"].(string)
	Helpers.LogDebug("Successfully Received Stripe Token: " + jsonData["id"].(string))

	return nil
}

func (r *RedeemerStruct) IntentSetup() error {
	headers := r.GetHeaders()

	req, err := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/billing/stripe/setup-intents", nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Setup_Client_Intents.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed Setup Client Intents, Error: " + string(body) + " -> " + Helpers.Replacelast(r.Token))
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Discord.ClientSecret = jsonData["client_secret"].(string)
	Helpers.LogDebug("Successfully Received Client Secret: " + jsonData["client_secret"].(string))

	return nil
}

func (r *RedeemerStruct) ValidateBilling() error {
	headers := r.GetHeaders()

	payload := map[string]interface{}{
		"billing_address": map[string]interface{}{
			"name":        "Aman",
			"line_1":      "20 Wheel Farm Drive",
			"line_2":      "",
			"city":        "Dagenham",
			"state":       "United Kingdom",
			"postal_code": "RM10 7AR",
			"country":     "GB",
			"email":       "",
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/billing/payment-sources/validate-billing-address", bytes.NewBuffer(jsonPayload))
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Validate_Billing.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed Validating Billing Address, Error: " + string(body) + " -> " + Helpers.Replacelast(r.Token))
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Discord.BillingToken = jsonData["token"].(string)
	Helpers.LogDebug("Successfully Received Billing Token: " + jsonData["token"].(string))

	return nil
}

func (r *RedeemerStruct) StripeIntents() error {
	headers := r.GetStripeHeaders()
	rand.Seed(time.Now().UnixNano())

	splitParts := strings.Split(r.Discord.ClientSecret, "_")
	slicedParts := splitParts[:2]
	clientSecretID := strings.Join(slicedParts, "_")

	headers["authorization"] = []string{"Bearer " + StripeKey}

	parameters := url.Values{}
	parameters.Add("payment_method_data[type]", "card")
	parameters.Add("payment_method_data[card][token]", r.Stripe.StripeToken)
	parameters.Add("payment_method_data[billing_details][address][line1]", "20 Wheel Farm Drive")
	parameters.Add("payment_method_data[billing_details][address][line2]", "")
	parameters.Add("payment_method_data[billing_details][address][city]", "Dagenham")
	parameters.Add("payment_method_data[billing_details][address][state]", "United Kingdom")
	parameters.Add("payment_method_data[billing_details][address][postal_code]", "RM10 7AR")
	parameters.Add("payment_method_data[billing_details][address][country]", "GB")
	parameters.Add("payment_method_data[billing_details][name]", "Aman")
	parameters.Add("payment_method_data[guid]", r.Stripe.StripeGuid)
	parameters.Add("payment_method_data[muid]", r.Stripe.StripeMuid)
	parameters.Add("payment_method_data[sid]", r.Stripe.StripeSid)
	parameters.Add("payment_method_data[payment_user_agent]", "stripe.js/0651b07260; stripe-js-v3/0651b07260")
	parameters.Add("payment_method_data[time_on_page]", fmt.Sprintf("%d", rand.Intn(60000-30000+1)+30000))
	parameters.Add("expected_payment_method_type", "card")
	parameters.Add("use_stripe_sdk", "true")
	parameters.Add("key", StripeKey)
	parameters.Add("client_secret", r.Discord.ClientSecret)

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/setup_intents/"+clientSecretID+"/confirm?"+parameters.Encode(), nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if strings.Contains(string(body), "error") {
			go Helpers.AppendTextToFile(r.VccInfo+"\n", r.Folder+"/Bad_VCC.txt")
			go Helpers.RemoveLine(r.VccInfo, "./Data/Input/Vcc's.txt")
			go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
			return errors.New("Bad Payment Method, Card is Expired or Was Declined")
		}

		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Setup_Stripe_Intents.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed Getting Stripe Intents, Error: " + string(body) + " -> " + Helpers.Replacelast(r.Token))
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Stripe.StripePaymentToken = jsonData["payment_method"].(string)
	Helpers.LogDebug("Successfully Received Payment Method Token: " + jsonData["payment_method"].(string))

	return nil
}

func (r *RedeemerStruct) PaymentSource() error {
	headers := r.GetHeaders()

	payload := map[string]interface{}{
		"payment_gateway": 1,
		"token":           r.Stripe.StripePaymentToken,
		"billing_address": map[string]interface{}{
			"name":        "Aman",
			"line_1":      "20 Wheel Farm Drive",
			"line_2":      "",
			"city":        "Dagenham",
			"state":       "United Kingdom",
			"postal_code": "RM10 7AR",
			"country":     "GB",
			"email":       "",
		},
		"billing_address_token": r.Discord.BillingToken,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/billing/payment-sources", bytes.NewBuffer(jsonPayload))
	req.Header = headers
	req.Header.Del("Content-Length")
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		Helpers.LogVcc(r.Token, r.VccInfo)
		return nil
	} else if strings.Contains(string(body), "captcha_key") {
		var solution string
		var ResponseData ServerJoinRQ

		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Captcha_Tokens.txt")
		Captcha++
		Helpers.LogInfo("Encountered Captcha", Helpers.Replacelast(r.Token))

		err = json.Unmarshal(body, &ResponseData)
		if err != nil {
			return err
		}

		if strings.ToLower(settings.CaptchaSettings.Service) == "capsolver" {
			for i := 0; i < 5; i++ {
				solution = Solvers.Capsolver(settings.CaptchaSettings.APIKey, ResponseData.CaptchaSitekey, ResponseData.CaptchaRqdata)
				if solution != "" {
					break
				} else if solution == "error" {
					return errors.New("Failed to Solve Captcha, Received Error While Solving")
				} else if solution == "" && i == 5 {
					go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Solve_Captcha.txt")
					return errors.New("Failed to Solve Captcha, Unable to add Payment Source")
				}
			}

		} else if strings.ToLower(settings.CaptchaSettings.Service) == "hcoptcha" {
			for i := 0; i < 5; i++ {
				solution = Solvers.Hcoptcha(settings.CaptchaSettings.APIKey, ResponseData.CaptchaSitekey, ResponseData.CaptchaRqdata)
				if solution != "" {
					break
				} else if solution == "error" {
					return errors.New("Failed to Solve Captcha, Received Error While Solving")

				} else if solution == "" && i == 5 {
					go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Solve_Captcha.txt")
					return errors.New("Failed to Solve Captcha, Unable to add Payment Source")
				}
			}

		} else if strings.ToLower(settings.CaptchaSettings.Service) == "ab5" {
			for i := 0; i < 5; i++ {
				solution = Solvers.Ab5Solver(settings.CaptchaSettings.APIKey, ResponseData.CaptchaSitekey, ResponseData.CaptchaRqdata)
				if solution != "" {
					break
				} else if solution == "" && i == 5 {
					go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Solve_Captcha.txt")
					return errors.New("Failed to Solve Captcha, Unable to add Payment Source")
				}
			}
		} else {
			return errors.New("Unsupported Captcha Service Listed in Config, Check Config")
		}

		CaptchaPayload := map[string]interface{}{
			"captcha_key":     solution,
			"captcha_rqtoken": ResponseData.CaptchaRqtoken,
			"payment_gateway": 1,
			"token":           r.Stripe.StripePaymentToken,
			"billing_address": map[string]interface{}{
				"name":        "Aman",
				"line_1":      "20 Wheel Farm Drive",
				"line_2":      "",
				"city":        "Dagenham",
				"state":       "United Kingdom",
				"postal_code": "RM10 7AR",
				"country":     "GB",
				"email":       "",
			},
			"billing_address_token": r.Discord.BillingToken,
		}

		CaptchaJsonPayload, err := json.Marshal(CaptchaPayload)
		if err != nil {
			return err
		}

		CaptchaReq, err := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/billing/payment-sources", bytes.NewBuffer(CaptchaJsonPayload))
		CaptchaReq.Header = headers
		CaptchaReq.Header.Del("Content-Length")

		if err != nil {
			return err
		}

		CaptchaResp, err := r.Client.Do(CaptchaReq)
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
			}
		}(CaptchaResp.Body)

		CaptchaBody, err := io.ReadAll(CaptchaResp.Body)
		if err != nil {
			return err
		}

		if CaptchaResp.StatusCode == 200 {
			Helpers.LogVcc(r.Token, r.VccInfo)
			return nil
		} else if CaptchaResp.StatusCode == 400 {
			go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Bad_Request_Error.txt")
			go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
			return errors.New("Failed to Add Payment Source, Received Bad Request Error")
		} else {
			go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Add_Payment_Source.txt")
			go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
			return errors.New("Failed to Add Payment Source, Error: " + string(CaptchaBody))
		}
	} else {
		go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Received_Unknown_Error.txt")
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Received Unknown Error, Error: " + string(body))
	}

}

func (r *RedeemerStruct) GetPaymentSourceID() error {
	headers := r.GetHeaders()

	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me/billing/payment-sources", nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	var jsonData []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	if len(jsonData) > 0 {
		r.Discord.PaymentSourceID = jsonData[0]["id"].(string)
		Helpers.LogDebug("Successfully Received Payment Source ID: " + jsonData[0]["id"].(string))
	} else {
		go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
		return errors.New("Failed to Receive Payment Source ID -> " + r.UnformattedToken)
	}

	return nil
}

func (r *RedeemerStruct) Redeem() error {
	if settings.VccSettings.AuthRetries > 0 {
		for i := 0; i < 1; i++ {
			Helpers.LogPromo(r.Token, r.Promo)

			headers := r.GetHeaders()

			headers["referer"] = []string{"https://discord.com/billing/promotions/" + r.Promo}

			payload := map[string]interface{}{
				"channel_id":               nil,
				"payment_source_id":        r.Discord.PaymentSourceID,
				"gateway_checkout_context": nil,
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			req, err := http.NewRequest("POST", "https://discord.com/api/v9/entitlements/gift-codes/"+r.Promo+"/redeem", bytes.NewBuffer(jsonPayload))
			req.Header = headers
			if err != nil {
				return err
			}

			resp, err := r.Client.Do(req)
			if err != nil {
				return err
			}

			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				}
			}(resp.Body)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if resp.StatusCode == 200 {
				Success++
				VccClaims++
				go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Redeemed.txt")
				go Helpers.RemoveLine(r.UnformattedPromo, "./Data/Input/Promos.txt")
				Helpers.LogRedeemed(r.Token, r.VccInfo, r.Promo, time.Since(r.StartTime))
				return nil
			} else if strings.Contains(string(body), "100029") {
				Helpers.LogInfo("Failed Authentication, Retrying", Helpers.Replacelast(r.Token))

				var jsonData map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&jsonData)
				if err != nil {
					return err
				}

				r.Discord.PaymentID = jsonData["payment_id"].(string)

				err = r.Setup3DPaymentIntents()
				if err != nil {
					return err
				}

				err = r.Confirm()
				if err != nil {
					return err
				}

				err = r.Authenticate()
				if err != nil {
					return err
				}

				if i == settings.VccSettings.AuthRetries {
					go Helpers.AppendTextToFile(r.VccInfo+"\n", r.Folder+"/Failed_To_Authenticate.txt")
					go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
					go Helpers.RemoveLine(r.VccInfo, "./Data/Input/Vcc's.txt")
					return errors.New("Failed Authentication, Couldn't Redeem Promo")
				}

			} else {
				go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Redeem.txt")
				go Helpers.RemoveLine(r.UnformattedToken, "./Data/Input/Tokens.txt")
				go Helpers.RemoveLine(r.VccInfo, "./Data/Input/Vcc's.txt")

				return errors.New("Failed to Redeem Promo, Error: " + string(body))
			}
		}
	} else {

		Helpers.LogPromo(r.Token, r.Promo)

		headers := r.GetHeaders()

		headers["referer"] = []string{"https://discord.com/billing/promotions/" + r.Promo}

		payload := map[string]interface{}{
			"channel_id":               nil,
			"payment_source_id":        r.Discord.PaymentSourceID,
			"gateway_checkout_context": nil,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "https://discord.com/api/v9/entitlements/gift-codes/"+r.Promo+"/redeem", bytes.NewBuffer(jsonPayload))
		req.Header = headers
		if err != nil {
			return err
		}

		resp, err := r.Client.Do(req)
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			Success++
			VccClaims++
			go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Redeemed.txt")
			go Helpers.RemoveLine(r.UnformattedPromo, "./Data/Input/Promos.txt")
			Helpers.LogRedeemed(r.Token, r.VccInfo, r.Promo, time.Since(r.StartTime))
			return nil
		} else if strings.Contains(string(body), "100029") {

			go Helpers.AppendTextToFile(r.VccInfo+"\n", r.Folder+"/Failed_To_Authenticate.txt")
			go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
			go Helpers.RemoveLine(r.VccInfo, "./Data/Input/Vcc's.txt")
			return errors.New("Failed Authentication, Couldn't Redeem Promo")

		} else {
			go Helpers.AppendTextToFile(r.UnformattedToken+"\n", r.Folder+"/Failed_To_Redeem.txt")
			go Helpers.AppendTextToFile(r.UnformattedToken, "./Data/Input/Tokens.txt")
			go Helpers.RemoveLine(r.VccInfo, "./Data/Input/Vcc's.txt")
			return errors.New("Failed to Redeem Promo, Error" + string(body))
		}
	}

	return nil
}

func (r *RedeemerStruct) Setup3DPaymentIntents() error {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me/billing/stripe/payment-intents/payments/"+r.Discord.PaymentID, nil)
	req.Header = r.GetHeaders()
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("Failed to Get Stripe Payment Intent Client Secret")
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	r.Discord.PaymentIntentClientSecret = jsonData["stripe_payment_intent_client_secret"].(string)
	Helpers.LogDebug("Successfully Received Stripe Payment Intent Client Secret: " + jsonData["stripe_payment_intent_client_secret"].(string))

	return nil
}

func (r *RedeemerStruct) Confirm() error {
	splitParts := strings.Split(r.Discord.PaymentIntentClientSecret, "_")
	slicedParts := splitParts[:2]
	clientSecretID := strings.Join(slicedParts, "_")

	headers := r.GetStripeHeaders()
	headers["authorization"] = []string{"Bearer " + StripeKey}

	parameters := url.Values{}
	parameters.Add("expected_payment_method_type", "card")
	parameters.Add("use_stripe_sdk", "true")
	parameters.Add("key", StripeKey)
	parameters.Add("client_secret", r.Discord.PaymentIntentClientSecret)

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/payment_intents/"+clientSecretID+"/confirm?"+parameters.Encode(), nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("Failed to Get Stripe 3D Source")
	}

	var response Stripe3DSource

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	r.Stripe.Stripe3D = response.NextAction.UseStripeSdk.ThreeDSecure2Source
	Helpers.LogDebug("Successfully Received Stripe 3D Source: " + response.NextAction.UseStripeSdk.ThreeDSecure2Source)

	return nil
}

func (r *RedeemerStruct) Authenticate() error {
	headers := r.GetStripeHeaders()
	headers["authorization"] = []string{"Bearer " + StripeKey}

	formData := url.Values{
		"source":                                 {r.Stripe.Stripe3D},
		"browser":                                {`{"fingerprintAttempted":false,"fingerprintData":null,"challengeWindowSize":null,"threeDSCompInd":"Y","browserJavaEnabled":false,"browserJavascriptEnabled":true,"browserLanguage":"en-US","browserColorDepth":"24","browserScreenHeight":"728","browserScreenWidth":"1366","browserTZ":"420","browserUserAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"}`},
		"one_click_authn_device_support[hosted]": {"false"},
		"one_click_authn_device_support[same_origin_frame]":                 {"false"},
		"one_click_authn_device_support[spc_eligible]":                      {"false"},
		"one_click_authn_device_support[webauthn_eligible]":                 {"false"},
		"one_click_authn_device_support[publickey_credentials_get_allowed]": {"true"},
		"key": {StripeKey},
	}

	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/3ds2/authenticate?"+formData.Encode(), nil)
	req.Header = headers
	if err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("Failed to Authenticating Stripe 3D Source")
	}

	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return err
	}

	if jsonData["state"].(string) == "succeeded" {
		Helpers.LogDebug("Successfully Authenticated Vcc")
	}

	return nil

}
