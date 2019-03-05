package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"golang.org/x/crypto/ssh/terminal"
)

const oAuth2lientID = "3m25f2rmefr47lq00rlh3j1p17"

func main() {
	var username string
	flag.StringVar(&username, "username", "", "Your username")
	flag.Parse()

	if len(username) == 0 {
		log.Fatal("Please provide a valid username")
	}

	fmt.Print("Enter your password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(fmt.Errorf("error reading password: %s", err))
	}
	password := string(bytePassword)
	if len(password) == 0 {
		log.Fatal("Please provide a valid password")
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading default AWS config: %s", err))
	}
	cognitoClient := cognitoidentityprovider.New(cfg)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: cognitoidentityprovider.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
		ClientId: aws.String(oAuth2lientID),
	}
	err = input.Validate()
	if err != nil {
		log.Fatal(fmt.Errorf("error validating InitiateAuth input: %s", err))
	}
	req := cognitoClient.InitiateAuthRequest(input)
	resp, err := req.Send()
	if err != nil {
		log.Fatal(fmt.Errorf("error sending InitiateAuth request: %s", err))
	}

	if resp.AuthenticationResult.IdToken == nil {
		log.Fatal("No ID token found in response. Please use an existing Medisanté application to check if your account is set up correctly.")
	}

	fmt.Println("")
	fmt.Println(*resp.AuthenticationResult.IdToken)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("Tokens expires in %v seconds", *resp.AuthenticationResult.ExpiresIn))
}
