package api

import (
	"io/ioutil"
	"log"
	"net/http"
)
import "strings"

type HttpHeader struct {
	name  string
	value string
}

type AuthorizedApiRequest struct {
	headers []HttpHeader
	baseUrl string
	token   string
}

func newAuthorizedApiRequest(token string) AuthorizedApiRequest {
	req := AuthorizedApiRequest{}
	authorisationHeader := HttpHeader{
		"Authorization",
		token,
	}
	contentTypeHeader := HttpHeader{
		"Content-Type",
		"application/json",
	}
	req.headers = append(req.headers, authorisationHeader, contentTypeHeader)

	return req
}

type apiRequest interface {
	getBody() string
	getUrl() string
	getMethod() string
	addCustomHeader()
	getHeaders() []HttpHeader
}

type GetSelectionRequest struct {
	authedReq AuthorizedApiRequest
	path      string
	body      string
	method    string
}

func (req *GetSelectionRequest) getBody() string {
	return ""
}

func (req *GetSelectionRequest) getUrl() string {
	return req.authedReq.baseUrl + "/selection"
}

func (req *GetSelectionRequest) getMethod() string {
	return "GET"
}

func (req *GetSelectionRequest) addCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *GetSelectionRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func newGetSelectionRequest(authedReq AuthorizedApiRequest) GetSelectionRequest {
	req := GetSelectionRequest{}
	req.authedReq = authedReq

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

func (req *AddItemRequest) addCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *AddItemRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func newAddItemRequest(authedReq AuthorizedApiRequest, item string) AddItemRequest {
	req := AddItemRequest{}
	req.authedReq = authedReq
	req.item = item

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

func (req *SetPaymentMethodRequest) addCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *SetPaymentMethodRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func newSetPaymentMethodRequest(authedReq AuthorizedApiRequest, paymentMethodUri string) SetPaymentMethodRequest {
	req := SetPaymentMethodRequest{}
	req.authedReq = authedReq
	req.paymentMethodUri = paymentMethodUri

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

func (req *InitCheckoutRequest) addCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *InitCheckoutRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func newInitCheckoutRequest(authedReq AuthorizedApiRequest, body string) InitCheckoutRequest {
	req := InitCheckoutRequest{}
	req.authedReq = authedReq
	req.body = body

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

func (req *ReceiptRequest) addCustomHeader(header HttpHeader) {
	req.authedReq.headers = append(req.authedReq.headers, header)
}

func (req *ReceiptRequest) getHeaders() []HttpHeader {
	return req.authedReq.headers
}

func newReceiptRequest(authedReq AuthorizedApiRequest, body string) ReceiptRequest {
	req := ReceiptRequest{}
	req.authedReq = authedReq
	req.body = body

	return req
}

func sendApiRequest(apiRequest apiRequest) {
	client := &http.Client{}
	reader := strings.NewReader(apiRequest.getBody())
	request, _ := http.NewRequest(apiRequest.getMethod(), apiRequest.getUrl(), reader)

	for _, header := range apiRequest.getHeaders() {
		request.Header.Set(header.name, header.value)
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
}

func finalizeCheckout() {

	authedReq := AuthorizedApiRequest{}
	getSelection := newGetSelectionRequest(authedReq)

	//init selection

	//get selection
	//add item
	//set payment method
	//inject payment details

	//receipt + notification calls
}
