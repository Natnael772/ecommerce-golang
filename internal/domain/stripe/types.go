package stripe

type StripeProvider struct {
	webhookSecret string
}

type ProviderWebhookEvent struct {
	Provider      string
	ProviderTxnID string
	OrderID       string
	Status        string
	FailureReason string
	RawEvent      interface{}
}
