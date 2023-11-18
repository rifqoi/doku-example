package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// DOKU Request Const
const CLIENT_ID = "Client-Id"
const REQUEST_ID = "Request-Id"
const REQUEST_TIMESTAMP = "Request-Timestamp"
const REQUEST_TARGET = "Request-Target"
const DIGEST = "Digest"
const SYMBOL_COLON = ":"

// Generate random string
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomString(length int) string {
	return StringWithCharset(length, charset)
}

// Generate Digest
func generateDigest(jsonBody string) string {
	converted := []byte(jsonBody)
	hasher := sha256.New()
	hasher.Write(converted)

	return (base64.StdEncoding.EncodeToString(hasher.Sum(nil)))
}

func generateSignature(clientId string, requestId string, requestTimestamp string, requestTarget string, digest string, secret string) string {
	// Prepare Signature Component
	fmt.Println("----- Component Signature -----")
	var componentSignature strings.Builder
	componentSignature.WriteString(CLIENT_ID + SYMBOL_COLON + clientId)
	componentSignature.WriteString("\n")
	componentSignature.WriteString(REQUEST_ID + SYMBOL_COLON + requestId)
	componentSignature.WriteString("\n")
	componentSignature.WriteString(REQUEST_TIMESTAMP + SYMBOL_COLON + requestTimestamp)
	componentSignature.WriteString("\n")
	componentSignature.WriteString(REQUEST_TARGET + SYMBOL_COLON + requestTarget)
	componentSignature.WriteString("\n")
	componentSignature.WriteString(DIGEST + SYMBOL_COLON + digest)

	fmt.Println(componentSignature.String())
	fmt.Println("")

	// Calculate HMAC-SHA256 base64 from all the components above
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(componentSignature.String()))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// Prepend encoded result with algorithm info HMACSHA256=
	return "HMACSHA256=" + signature
}

func generateInvoiceNumber() string {
	currentTime := time.Now().Format("20060102")

	randString := generateRandomString(4)

	invoiceNumber := fmt.Sprintf("INV-%s-%s", currentTime, randString)
	return invoiceNumber
}

func generateRequest(amount int) (string, error) {

	paymentLinkRequest := new(PaymentLinkRequest)
	paymentLinkRequest.Order.Amount = amount
	paymentLinkRequest.Payment.PaymentDueDate = 30

	// Invoice number ini random asalkan max length 128 karakter
	paymentLinkRequest.Order.InvoiceNumber = generateInvoiceNumber()

	jsonBody, err := paymentLinkRequest.ToJSON()
	if err != nil {
		return "", err
	}

	jsonString := string(jsonBody)

	return jsonString, nil
}
