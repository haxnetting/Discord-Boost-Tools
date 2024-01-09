package Client

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"time"
)

type RedeemerStruct struct {
	Client           tls_client.HttpClient
	StartTime        time.Time
	Token            string
	UnformattedToken string
	UnformattedPromo string
	Promo            string
	Proxy            string
	Cookies          string
	SuperProperties  string
	Fingerprints     string
	VccInfo          string
	Folder           string
	Vcc              VccInfo
	Stripe           StripeInfo
	Discord          DiscordInfo
}

type VccInfo struct {
	VccCard        string
	VccCVV         string
	VccExpiryMonth string
	VccExpiryYear  string
}

type StripeInfo struct {
	StripeMuid         string
	StripeGuid         string
	StripeSid          string
	StripeToken        string
	StripePaymentToken string
	StripeClientSecret string
	Stripe3D           string
}

type DiscordInfo struct {
	ClientSecret              string
	BillingToken              string
	PaymentSourceID           string
	PaymentID                 string
	PaymentIntentClientSecret string
}

type ServerJoinRQ struct {
	CaptchaKey     []string `json:"captcha_key"`
	CaptchaSitekey string   `json:"captcha_sitekey"`
	CaptchaService string   `json:"captcha_service"`
	CaptchaRqdata  string   `json:"captcha_rqdata"`
	CaptchaRqtoken string   `json:"captcha_rqtoken"`
}

type PromoCheckResponse struct {
	Code          string    `json:"code,omitempty"`
	SkuID         string    `json:"sku_id,omitempty"`
	ApplicationID string    `json:"application_id,omitempty"`
	Uses          int       `json:"uses,omitempty"`
	MaxUses       int       `json:"max_uses,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	Redeemed      bool      `json:"redeemed,omitempty"`
	Flags         int       `json:"flags,omitempty"`
	BatchID       string    `json:"batch_id,omitempty"`
}

type TokenCheckResponse struct {
	ID                   string `json:"id,omitempty"`
	Username             string `json:"username,omitempty"`
	Avatar               string `json:"avatar,omitempty"`
	Discriminator        string `json:"discriminator,omitempty"`
	PublicFlags          int    `json:"public_flags,omitempty"`
	PremiumType          int    `json:"premium_type,omitempty"`
	Flags                int    `json:"flags,omitempty"`
	Banner               any    `json:"banner,omitempty"`
	AccentColor          any    `json:"accent_color,omitempty"`
	GlobalName           any    `json:"global_name,omitempty"`
	AvatarDecorationData any    `json:"avatar_decoration_data,omitempty"`
	BannerColor          any    `json:"banner_color,omitempty"`
	MfaEnabled           bool   `json:"mfa_enabled,omitempty"`
	Locale               string `json:"locale,omitempty"`
	Email                string `json:"email,omitempty"`
	Verified             bool   `json:"verified,omitempty"`
	Phone                any    `json:"phone,omitempty"`
	NsfwAllowed          bool   `json:"nsfw_allowed,omitempty"`
	PremiumUsageFlags    int    `json:"premium_usage_flags,omitempty"`
	LinkedUsers          []any  `json:"linked_users,omitempty"`
	PurchasedFlags       int    `json:"purchased_flags,omitempty"`
	Bio                  string `json:"bio,omitempty"`
	AuthenticatorTypes   []any  `json:"authenticator_types,omitempty"`
}

type Stripe3DSource struct {
	ID            string `json:"id,omitempty"`
	Object        string `json:"object,omitempty"`
	Amount        int    `json:"amount,omitempty"`
	AmountDetails struct {
		Tip struct {
		} `json:"tip,omitempty"`
	} `json:"amount_details,omitempty"`
	AutomaticPaymentMethods any    `json:"automatic_payment_methods,omitempty"`
	CanceledAt              any    `json:"canceled_at,omitempty"`
	CancellationReason      any    `json:"cancellation_reason,omitempty"`
	CaptureMethod           string `json:"capture_method,omitempty"`
	ClientSecret            string `json:"client_secret,omitempty"`
	ConfirmationMethod      string `json:"confirmation_method,omitempty"`
	Created                 int    `json:"created,omitempty"`
	Currency                string `json:"currency,omitempty"`
	Description             string `json:"description,omitempty"`
	LastPaymentError        any    `json:"last_payment_error,omitempty"`
	Livemode                bool   `json:"livemode,omitempty"`
	NextAction              struct {
		Type         string `json:"type,omitempty"`
		UseStripeSdk struct {
			DirectoryServerEncryption struct {
				Algorithm                  string   `json:"algorithm,omitempty"`
				Certificate                string   `json:"certificate,omitempty"`
				DirectoryServerID          string   `json:"directory_server_id,omitempty"`
				KeyID                      string   `json:"key_id,omitempty"`
				RootCertificateAuthorities []string `json:"root_certificate_authorities,omitempty"`
			} `json:"directory_server_encryption,omitempty"`
			DirectoryServerName  string `json:"directory_server_name,omitempty"`
			Merchant             string `json:"merchant,omitempty"`
			OneClickAuthn        any    `json:"one_click_authn,omitempty"`
			ServerTransactionID  string `json:"server_transaction_id,omitempty"`
			ThreeDSecure2Source  string `json:"three_d_secure_2_source,omitempty"`
			ThreeDsMethodURL     string `json:"three_ds_method_url,omitempty"`
			ThreeDsOptimizations string `json:"three_ds_optimizations,omitempty"`
			Type                 string `json:"type,omitempty"`
		} `json:"use_stripe_sdk,omitempty"`
	} `json:"next_action,omitempty"`
	PaymentMethod                     string   `json:"payment_method,omitempty"`
	PaymentMethodConfigurationDetails any      `json:"payment_method_configuration_details,omitempty"`
	PaymentMethodTypes                []string `json:"payment_method_types,omitempty"`
	Processing                        any      `json:"processing,omitempty"`
	ReceiptEmail                      any      `json:"receipt_email,omitempty"`
	SetupFutureUsage                  any      `json:"setup_future_usage,omitempty"`
	Shipping                          any      `json:"shipping,omitempty"`
	Source                            any      `json:"source,omitempty"`
	Status                            string   `json:"status,omitempty"`
}
