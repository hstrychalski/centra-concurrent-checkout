package checkout

import (
	"concurrent-checkout/src/api"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type InitialisedCheckout struct {
	orderNumber       int
	token             string
	paymentMethodType PaymentMethodType
	authedReq         api.AuthorizedApiRequest
}

func InitCheckout() InitialisedCheckout {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	baseUrl := os.Getenv("API_BASE_URL")
	itemId := os.Getenv("ITEM_ID")
	paymentMethodURI := os.Getenv("PAYMENT_METHOD_URI")
	paymentMethod := os.Getenv("PAYMENT_METHOD")

	authedReq := api.NewAuthorizedApiRequest(baseUrl)
	getSelectionReq := api.NewGetSelectionRequest(authedReq)

	getSelectionResult := api.SendApiRequest(&getSelectionReq)
	authedReq.AuthorizeWithToken(getSelectionResult.Token)

	addItemReq := api.NewAddItemRequest(authedReq, itemId)
	api.SendApiRequest(&addItemReq)

	setPaymentMethodReq := api.NewSetPaymentMethodRequest(authedReq, paymentMethodURI)
	api.SendApiRequest(&setPaymentMethodReq)

	initCheckoutJSON := getInitCheckoutBody()
	initCheckoutReq := api.NewInitCheckoutRequest(authedReq, initCheckoutJSON)
	initCheckoutReq.AddCustomHeader(api.HttpHeader{Name: "Test-Case", Value: "1"})

	initCheckoutResult := api.SendApiRequest(&initCheckoutReq)

	paymentMethodType, err := NewPaymentMethodType(paymentMethod)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	orderNumber := ParseOrderNumberFromHtmlSnippet(initCheckoutResult.FormHTML, paymentMethodType)
	fmt.Printf("Initialised checkout for order number:%d\n", orderNumber)

	res := InitialisedCheckout{}
	res.token = getSelectionResult.Token
	res.orderNumber = orderNumber
	res.paymentMethodType = paymentMethodType
	res.authedReq = authedReq

	return res
}

func FinaliseCheckoutWithReceiptRequest(initedCheckout InitialisedCheckout) {
	receiptReq, _ := CreateReceiptRequest(initedCheckout.paymentMethodType, initedCheckout.authedReq)
	receiptReq.AddCustomHeader(api.HttpHeader{Name: "Test-Case", Value: "3"})
	time.Sleep(5 * time.Second)
	res := api.SendApiRequest(receiptReq)

	fmt.Printf("Finalised order with number:%s\n", res.OrderJSON.OrderNumber)
}

func FinaliseCheckoutWithNotificationRPC(initedCheckout InitialisedCheckout) {

	notificationRPC := api.NewQliroNotificationRPC(
		os.Getenv("NOTIFICATION_BASE_URL"),
		os.Getenv("PAYMENT_METHOD_URI"),
		os.Getenv("PAYMENT_PLUGIN_ID"),
		os.Getenv("SERVER_SIDE_SECRET"),
		initedCheckout.orderNumber)
	notificationRPC.AddHeader(api.HttpHeader{Name: "Test-Case", Value: "3"})

	responseBody := api.SendNotificationRPC(notificationRPC)
	fmt.Printf("Notification response: %s\n", responseBody)
}

type PaymentMethodType interface {
	getPaymentMethodType() string
}

type PaymentMethodQliro struct{}

func (PaymentMethodQliro) getPaymentMethodType() string {
	return "qliro"
}

func NewPaymentMethodType(paymentMethodType string) (PaymentMethodType, error) {
	if paymentMethodType == "qliro" {
		return PaymentMethodQliro{}, nil
	}

	return nil, errors.New("unsupported payment method type")
}

func CreateReceiptRequest(paymentMethodType PaymentMethodType, authedReq api.AuthorizedApiRequest) (*api.ReceiptRequest, error) {
	if paymentMethodType.getPaymentMethodType() == (PaymentMethodQliro{}).getPaymentMethodType() {
		req := api.NewReceiptRequest(authedReq, "")
		return &req, nil
	}

	return nil, errors.New("unsupported payment method type")
}

func ParseOrderNumberFromHtmlSnippet(formHtml string, paymentMethod PaymentMethodType) int {
	if paymentMethod.getPaymentMethodType() == (PaymentMethodQliro{}).getPaymentMethodType() {
		regex := regexp.MustCompile("\\d+")
		res := regex.FindAllString(formHtml, -1)
		log.Printf(res[0])
		parseRes, _ := strconv.Atoi(res[0])
		return parseRes
	}

	log.Fatalln("Unsupported payment method")
	return 0
}

func getInitCheckoutBody() string {
	return `{
			"paymentReturnPage": "http://example.com/success",
			"paymentFailedPage": "http://example.com/failure",
			"termsAndConditions": "true",
			"address": {
			"newsletter": false,
				"email": "abc123@example.com",
				"phoneNumber": "123456789",
				"firstName": "Test Billing",
				"lastName": "Testson Billing",
				"address1": "Address One",
				"address2": "Address Two",
				"zipCode": "12345",
				"city": "Malmo",
				"country": "SE"
		},
		"shippingAddress": {
			"phoneNumber": "123456789",
				"firstName": "Test Shipping",
				"lastName": "Testson Shipping",
				"address1": "ShipAddress One",
				"address2": "ShipAddress Two",
				"zipCode": "12345",
				"city": "Stockholm",
				"country": "SE"
		}
	}`
}
