// Lambda Invoker - using Go
// =========================
//
// I wanted to invoke a Lambda funciton remotely without using the AWS API Gateway
// or the full SDK. The code creates an AWS SigV4 header, and then invokes the
// Lambda through the HTTP Endpoint. It relies on AWS credentials already being
// configured on the local machine
//

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

// AWS Service and Region Details
const REGION string = "eu-west-2" // CHANGE THIS VALUE
const SERVICE string = "lambda"

// URI Path
const FUNCTION_NAME string = "## INSERT LAMBDA FUNCTION NAME ##" // CHANGE THIS VALUE

// PAYLOAD
const PAYLOAD string = "{\"example-key\":\"example-value\"}" // CHANGE THIS VALUE

func main() {
	// Update logger to include ms
	log.SetFlags(log.Lmicroseconds)

	// Create Context
	log.Printf("Initialising")
	ctx := context.TODO()

	// New Session
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(REGION))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Get Creds
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		log.Fatalf("expected no error, but received %v", err)
	}

	// Generate host name
	host := fmt.Sprintf("https://%s.%s.amazonaws.com", SERVICE, REGION)

	// Generate URI path
	path := fmt.Sprintf("/2015-03-31/functions/%s/invocations", FUNCTION_NAME)

	// Build the body
	body := strings.NewReader(PAYLOAD)

	// Create Signer
	signer := v4.NewSigner()

	// Build request
	log.Printf("Creating Request")
	req, err := http.NewRequest("POST", host, body)
	if err != nil {
		log.Fatalf("expected no error, but received %v", err)
	}

	// Set the request path
	req.URL.Path = path

	// Generate SHA256 hash of payload in a hex string format
	log.Printf("Generating Hash")
	payloadHash := sha256.New()
	payloadHash.Write([]byte(PAYLOAD))
	sha256_hash := hex.EncodeToString(payloadHash.Sum(nil))
	log.Printf("Hash: %s", sha256_hash)

	// Sign the Request
	err = signer.SignHTTP(ctx,
		creds,
		req,
		sha256_hash,
		SERVICE,
		REGION,
		time.Now())
	if err != nil {
		log.Fatalf("expected no error, but received %v", err)
	}

	// Output Signature
	log.Printf("Signature: %s", req.Header.Get("Authorization"))

	// Make Request
	log.Printf("Making web request")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("expected no error, but received %v", err)
	}
	defer resp.Body.Close()

	// Read Ouput
	outputBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("expected no error, but received %v", err)
	}

	// Output Result
	log.Printf("Return Status: %s", resp.Status)
	log.Printf("Body: %s", string(outputBody))

}
