package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"strings"
)

var (
	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New(http.StatusText(http.StatusUnauthorized))
)

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := getToken(request)
	err := authToken(token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return events.APIGatewayCustomAuthorizerResponse{}, ErrNon200Response
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func main() {
	lambda.Start(handler)
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func getToken(request events.APIGatewayCustomAuthorizerRequest) string {
	// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-input.html
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	if len(tokenSlice) > 1 {
		return tokenSlice[len(tokenSlice)-1]
	}
	return ""
}

func authToken(token string) error {
	if token == "" {
		return ErrNon200Response
	}

	req, err := http.NewRequest("GET", "https://api.github.com/orgs/104corp/repos?per_page=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return ErrNon200Response
	}

	return nil
}