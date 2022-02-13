package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HttpHeader struct {
	Name  string
	Value string
}

type AuthorizedApiRequest struct {
	headers []HttpHeader
	baseUrl string
}

func NewAuthorizedApiRequest(baseUrl string) AuthorizedApiRequest {
	req := AuthorizedApiRequest{}

	contentTypeHeader := HttpHeader{
		"Content-Type",
		"application/json",
	}
	req.headers = append(req.headers, contentTypeHeader)
	req.baseUrl = baseUrl

	return req
}

func (req *AuthorizedApiRequest) AuthorizeWithToken(token string) {
	authorisationHeader := HttpHeader{
		"api-token",
		token,
	}

	req.headers = append(req.headers, authorisationHeader)
}

type CheckoutApiRequest interface {
	getBody() string
	getUrl() string
	getMethod() string
	AddCustomHeader(header HttpHeader)
	getHeaders() []HttpHeader
}

type GetSelectionRequest struct {
	authedReq AuthorizedApiRequest
	path      string
	body      string
	method    string
}

func (req GetSelectionRequest) getBody() string {
	return ""
}

func (req GetSelectionRequest) getUrl() string {
	return req.authedReq.baseUrl + "/selection"
}

func (req GetSelectionRequest) getMethod() string {
	return "GET"
}

func (req *GetSelectionRequest) AddCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req GetSelectionRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func NewGetSelectionRequest(authedReq AuthorizedApiRequest) GetSelectionRequest {
	req := GetSelectionRequest{authedReq: authedReq}

	return req
}

type AddItemRequest struct {
	authedReq AuthorizedApiRequest
	item      string
	path      string
	body      string
	method    string
}

func (req *AddItemRequest) getBody() string {
	return ""
}

func (req *AddItemRequest) getUrl() string {
	return req.authedReq.baseUrl + "/items/" + req.item
}

func (req *AddItemRequest) getMethod() string {
	return "POST"
}

func (req *AddItemRequest) AddCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *AddItemRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func NewAddItemRequest(authedReq AuthorizedApiRequest, item string) AddItemRequest {
	req := AddItemRequest{authedReq: authedReq, item: item}

	return req
}

type SetPaymentMethodRequest struct {
	authedReq        AuthorizedApiRequest
	paymentMethodUri string
	path             string
	body             string
	method           string
}

func (req *SetPaymentMethodRequest) getBody() string {
	return ""
}

func (req *SetPaymentMethodRequest) getUrl() string {
	return req.authedReq.baseUrl + "/payment-methods/" + req.paymentMethodUri
}

func (req *SetPaymentMethodRequest) getMethod() string {
	return "PUT"
}

func (req *SetPaymentMethodRequest) AddCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *SetPaymentMethodRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func NewSetPaymentMethodRequest(authedReq AuthorizedApiRequest, paymentMethodUri string) SetPaymentMethodRequest {
	req := SetPaymentMethodRequest{authedReq: authedReq, paymentMethodUri: paymentMethodUri}

	return req
}

type InitCheckoutRequest struct {
	authedReq AuthorizedApiRequest
	path      string
	body      string
	method    string
}

func (req *InitCheckoutRequest) getBody() string {
	return req.body
}

func (req *InitCheckoutRequest) getUrl() string {
	return req.authedReq.baseUrl + "/payment"
}

func (req *InitCheckoutRequest) getMethod() string {
	return "POST"
}

func (req *InitCheckoutRequest) AddCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *InitCheckoutRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func NewInitCheckoutRequest(authedReq AuthorizedApiRequest, body string) InitCheckoutRequest {
	req := InitCheckoutRequest{authedReq: authedReq, body: body}

	return req
}

type ReceiptRequest struct {
	authedReq AuthorizedApiRequest
	path      string
	body      string
	method    string
}

func (req *ReceiptRequest) getBody() string {
	return req.body
}

func (req *ReceiptRequest) getUrl() string {
	return req.authedReq.baseUrl + "/payment-result"
}

func (req *ReceiptRequest) getMethod() string {
	return "POST"
}

func (req *ReceiptRequest) AddCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *ReceiptRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func NewReceiptRequest(authedReq AuthorizedApiRequest, body string) ReceiptRequest {
	req := ReceiptRequest{authedReq: authedReq, body: body}

	return req
}

type NotificationRPC interface {
	getUrl() string
	getBody() string
	getHeaders() []HttpHeader
}

type QliroNotificationRPC struct {
	baseUrl          string
	paymentMethodURI string
	pluginId         string
	serverSideSecret string
	body             string
	headers          []HttpHeader
}

func NewQliroNotificationRPC(baseUrl string, paymentMethodURI string, pluginId string, serverSideSecret string, orderNumber int) QliroNotificationRPC {
	body := fmt.Sprintf(`{
  		"OrderId": 12345,
  		"MerchantReference": "%d",
  		"Status": "Completed",
  		"Timestamp": "2016-03-03T11:43:05.567",
  		"NotificationType": "CustomerCheckoutStatus",
  		"PaymentTransactionId": 1234
	}`, orderNumber)

	return QliroNotificationRPC{baseUrl: baseUrl, paymentMethodURI: paymentMethodURI, pluginId: pluginId, serverSideSecret: serverSideSecret, body: body}
}

func (notificationRPC QliroNotificationRPC) getUrl() string {
	return notificationRPC.baseUrl + "/" + notificationRPC.paymentMethodURI + "/" + notificationRPC.pluginId + "/checkout-status-update/?access-key=" + notificationRPC.serverSideSecret
}

func (notificationRPC QliroNotificationRPC) getBody() string {
	return notificationRPC.body
}

func (notificationRPC QliroNotificationRPC) getHeaders() []HttpHeader {
	return notificationRPC.headers
}

func (notificationRPC *QliroNotificationRPC) AddHeader(header HttpHeader) {
	notificationRPC.headers = append(notificationRPC.headers, header)
}

type ApiResponse struct {
	Token     string    `json:"token"`
	FormHTML  string    `json:"formHtml"`
	OrderJSON OrderJSON `json:"order"`
}

type OrderJSON struct {
	OrderNumber string `json:"order"`
}

func SendApiRequest(apiRequest CheckoutApiRequest) ApiResponse {
	client := &http.Client{}
	reader := strings.NewReader(apiRequest.getBody())
	request, _ := http.NewRequest(apiRequest.getMethod(), apiRequest.getUrl(), reader)

	for _, header := range apiRequest.getHeaders() {
		request.Header.Set(header.Name, header.Value)
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var result ApiResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Fatalln(err)
	}

	return result
}

func SendNotificationRPC(notification NotificationRPC) string {
	client := &http.Client{}
	reader := strings.NewReader(notification.getBody())
	request, _ := http.NewRequest("POST", notification.getUrl(), reader)

	for _, header := range notification.getHeaders() {
		request.Header.Set(header.Name, header.Value)
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
