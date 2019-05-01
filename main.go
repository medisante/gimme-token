package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"golang.org/x/crypto/ssh/terminal"
)

// These two values are static and can be used by any API client
// of the Medisanté Eliot API.
const (
	cognitoUserPoolID = "eu-central-1_P8l0OEy9K"
	oAuth2ClientID    = "7g6sangq93eruio247tlsp883n"
)

func main() {
	ctx := context.Background()
	// Get username from the flags the program was called with
	var username string
	flag.StringVar(&username, "username", "", "Your username")
	flag.Parse()
	if len(username) == 0 {
		log.Fatal("Please provide a valid username")
	}

	// Get password from the user's input so it isn't stored in their shell history
	fmt.Print("Enter your password: ")
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatal(fmt.Sprintf("error reading password: %s", err))
	}
	password := string(bytePassword)
	if len(password) == 0 {
		log.Fatal("Please provide a valid password")
	}

	// Configure Cognito client from AWS SDK
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error loading default AWS config: %s", err))
	}
	cfg.Region = endpoints.EuCentral1RegionID
	cfg.Credentials = aws.AnonymousCredentials
	cognitoClient := cognitoidentityprovider.New(cfg)

	// Configure SRP client
	csrp, err := cognitosrp.NewCognitoSRP(username, password, cognitoUserPoolID, oAuth2ClientID, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating SRP client: %s", err))
	}

	// Create InitiateAuth request and send it to AWS Cognito
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       cognitoidentityprovider.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	}
	err = input.Validate()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error validating InitiateAuth input: %s", err))
	}
	req := cognitoClient.InitiateAuthRequest(input)
	resp, err := req.Send(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error sending InitiateAuth request: %s", err))
	}
	if resp.ChallengeName != cognitoidentityprovider.ChallengeNameTypePasswordVerifier {
		log.Fatal(fmt.Sprintf("Unknown auth challenge: %s", resp.ChallengeName))
	}

	// Create response to the password verifier challenge and send it to AWS Cognito
	challengeInput, err := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating password verifier challenge: %s", err))
	}
	err = challengeInput.Validate()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error validating challenge input: %s", err))
	}
	chalReq := cognitoClient.RespondToAuthChallengeRequest(challengeInput)
	chalResp, err := chalReq.Send(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error sending challenge request: %s", err))
	}

	// Print the received token
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Your token:")
	fmt.Println(*chalResp.AuthenticationResult.IdToken)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("Your tokens expires in %v seconds", *chalResp.AuthenticationResult.ExpiresIn))
}
