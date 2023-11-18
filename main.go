package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func generatePaymentLink(url string, body []byte, clientId, requestId, requestTimestamp, signature string) (*PaymentLinkResponse, error) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(CLIENT_ID, clientId)
	r.Header.Add(REQUEST_ID, requestId)
	r.Header.Add(REQUEST_TIMESTAMP, requestTimestamp)
	r.Header.Add("Signature", signature)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		respBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		jsonBytes, err := json.Marshal(respBytes)
		if err != nil {
			return nil, err
		}

		jsonString := string(jsonBytes)
		return nil, fmt.Errorf("Unsuccessful request to DOKU:", jsonString)
	}

	paymentResponse := PaymentLinkResponse{}
	err = json.NewDecoder(res.Body).Decode(&paymentResponse)
	if err != nil {
		return nil, fmt.Errorf("JSON decode failed:", err.Error())
	}

	return &paymentResponse, nil

}

func main() {
	jsonBody, err := generateRequest(20000)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonBody)

	// Time harus UTC berdasarkan dokumentasi
	now := time.Now().UTC()
	requestTimestamp := now.Format(time.RFC3339)

	reqURL := "https://api-sandbox.doku.com/checkout/v1/payment"

	// Genreate Digest from JSON Body
	digest := generateDigest(jsonBody)
	fmt.Println("----- Digest -----")
	fmt.Println(digest)
	fmt.Println("")

	// RequestID juga random berdasarkan dokumentasi
	requestID := generateRandomString(20)

	// Ganti ke client ID dan secretKey dari dashboard doku
	clientId := "PUT_CLIENT_ID_HERE"
	secretKey := "PUT_SECRET_KEY_HERE"

	// Generate Signature
	headerSignature := generateSignature(
		clientId,
		requestID,
		requestTimestamp,
		"/checkout/v1/payment", // For merchant request to DOKU, use DOKU path here. For HTTP Notification, use merchant path here
		digest,                 // Set empty string for this argumentes if HTTP Method is GET/DELETE
		secretKey)

	fmt.Println("----- Header Signature -----")
	fmt.Println(headerSignature)

	paymentResponse, err := generatePaymentLink(reqURL, []byte(jsonBody), clientId, requestID, requestTimestamp, headerSignature)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(paymentResponse)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)
	fmt.Println(jsonString)
}
