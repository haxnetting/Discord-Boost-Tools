package Helpers

import (
	"fmt"
	"github.com/gookit/color"
	title "github.com/lxi1400/GoTitle"
	"os"
	"time"
)

func GetTime() string {
	return time.Now().Format("3:04 PM")
}

func LogError(message, token string) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=FF5733>!</><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>%v -> </><fg=FF5733>%v\n</>", GetTime(), message, token)
}

func LogRedeemed(token string, vcc string, promo string, time time.Duration) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=32CD32>$</><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>Redeemed Promo -> </><fg=32CD32>%v</><fg=FFFFFF> (vcc: </><fg=32CD32>%v</><fg=FFFFFF> | promo: </><fg=32CD32>%v</><fg=FFFFFF> | elapsed time: </><fg=32CD32>%v</><fg=FFFFFF>)\n</>", GetTime(), Replacelast(token), vcc, promo, time)
}

func LogVcc(token, vcc string) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=32CD32>$</><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>Added VCC -> </><fg=32CD32>%v</><fg=FFFFFF>, vcc: </><fg=32CD32>%v\n</>", GetTime(), Replacelast(token), vcc)
}

func LogPromo(token, promo string) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=6050DC>></><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>Redeeming Promo -> </><fg=6050DC>%v</><fg=FFFFFF>, promo: </><fg=6050DC>%v\n</>", GetTime(), Replacelast(token), promo)
}

func LogInfo(message, token string) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=6050DC>></><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>%v -> </><fg=6050DC>%v\n</>", GetTime(), message, token)
}

func LogDebug(format string, content ...interface{}) {
	if settings.MiscellaneousSettings.Debug {
		color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=FFFF00>?</><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>%v\n</>", GetTime(), fmt.Sprintf(format, content...))
	}
}

func LogPanic(format string, content ...interface{}) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=7F00FF>!!</><fg=FFFFFF>]</><fg=FF00FF> - </><fg=7F00FF>%v\n</>", GetTime(), fmt.Sprintf(format, content...))
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func LogFinished(format string, content ...interface{}) {
	color.Printf("<fg=FF00FF>%v</><fg=FFFFFF> [</><fg=6050DC>></><fg=FFFFFF>]</><fg=FF00FF> - </><fg=FFFFFF>%v\n</>", GetTime(), fmt.Sprintf(format, content...))
}

func UpdateTitle(message string) {
	title.SetTitle("Nitro Redeemer @tempywempy. | " + message)
}
