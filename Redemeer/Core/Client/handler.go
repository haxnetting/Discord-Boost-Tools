package Client

import (
	"Redeemer/Core/Helpers"
	"fmt"
	"os"
	"sync"
	"time"
)

var VccClaims, Success, Failed, Captcha = 0, 0, 0, 0
var mutex sync.Mutex

func Redeemer() {
	var wg sync.WaitGroup

	Helpers.CheckResources()

	timestamp := time.Now().Format("[2006-01-02] [3-04 PM]")
	folder := "./Data/Output/" + timestamp
	err := os.MkdirAll(folder, 0644)
	if err != nil {
		Helpers.LogError("Failed Creating Output Path, Error", err.Error())
	}

	for {
		Promos, _ := Helpers.GetResources("./Data/Input/Promos.txt")
		Vcc, _ := Helpers.GetResources("./Data/Input/Vcc's.txt")
		Tokens, _ := Helpers.GetResources("./Data/Input/Tokens.txt")
		Helpers.UpdateTitle(fmt.Sprintf("Success: %v - Failed: %v - Captcha: %v [VCC: %v | Promos: %v | Tokens: %v]", Success, Failed, Captcha, len(Vcc), len(Promos), len(Tokens)))

		done := make(chan bool)

		for i := 1; i <= settings.MiscellaneousSettings.Threads; i++ {
			wg.Add(1)

			if len(Vcc) == 0 || len(Promos) == 0 || len(Tokens) == 0 {
				//Helpers.LogFinished(fmt.Sprintf("Stopping Threads, Ran Out of Resources (Success: %v | Failed: %v)", Success, Failed))
				break
			}

			go worker(&wg, done, folder)
		}

		close(done)
		wg.Wait()
		done = make(chan bool)

	}
}

func worker(wg *sync.WaitGroup, done chan bool, folder string) {
	defer wg.Done()

	for {
		mutex.Lock()
		token := Helpers.GetResource("./Data/Input/Tokens.txt", true)
		promo := Helpers.GetResource("./Data/Input/Promos.txt", false)
		vcc := Helpers.GetResource("./Data/Input/Vcc's.txt", false)
		mutex.Unlock()

		if token != "" && promo != "" && vcc != "" {
			r, err := NewClient(token, promo, vcc, folder)
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.CheckElements()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.StripeCookies()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.GetStripeToken()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.IntentSetup()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.ValidateBilling()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.StripeIntents()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.PaymentSource()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.GetPaymentSourceID()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			err = r.Redeem()
			if err != nil {
				Failed++
				Helpers.LogError(err.Error(), Helpers.Replacelast(r.Token))
				return
			}

			if VccClaims == settings.VccSettings.MaxClaims {
				go Helpers.RemoveLine(vcc, "./Data/Input/Vcc's.txt")
				VccClaims = 0
			}
		} else {
			Helpers.LogFinished(fmt.Sprintf("Stopping Threads, Ran Out of Resources (Success: %v | Failed: %v)", Success, Failed))
			return
		}

	}

	done <- true
}
