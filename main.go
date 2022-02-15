package main

import "concurrent-checkout/src/checkout"

func main() {

	for i := 0; i < 20; i++ {
		initialisedCheckout := checkout.InitCheckout()
		if i%2 == 0 {
			go checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
			go checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
		} else {
			go checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
			go checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
		}
	}
}
