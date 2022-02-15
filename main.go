package main

import (
	"concurrent-checkout/src/checkout"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		initialisedCheckout := checkout.InitCheckout()
		wg.Add(1)
		if i%2 == 0 {
			go checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
			go checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
		} else {
			go checkout.FinaliseCheckoutWithNotificationRPC(initialisedCheckout)
			go checkout.FinaliseCheckoutWithReceiptRequest(initialisedCheckout)
		}
	}
	wg.Wait()
}
