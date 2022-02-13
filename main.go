package main

import "concurrent-checkout/src/checkout"

func main() {

	for i := 0; i < 20; i++ {
		initialisedCheckout := checkout.InitCheckout()
		if i%2 == 0 {
			checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
			checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
		} else {
			checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
			checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
		}
	}
}
