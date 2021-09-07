package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

var (
	ErrBadRequestResponse = errors.New(http.StatusText(http.StatusBadRequest))
	ErrInternalServerResponse = errors.New(http.StatusText(http.StatusInternalServerError))

	DestRegistry = os.Getenv("DestRegistry")
	DestRegistryUser = os.Getenv("DestRegistryUser")
	DestRegistryPassword = os.Getenv("DestRegistryPassword")
)

type ImagePushInfo struct {
	registry string `json:"registry"`
	image string `json:"image"`
	user string  `json:"user"`
	password string  `json:"password"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := "success"
	statusCode := 200

	imagePushInfo, err := getImagePushInfo(request.Body)
	if err != nil {
		logError(err)
		return events.APIGatewayProxyResponse{}, ErrBadRequestResponse
	}

	skopeoHelper := skopeo{}
	skopeoHelper.login(imagePushInfo.user, imagePushInfo.password, imagePushInfo.registry, true)
	skopeoHelper.login(DestRegistryUser, DestRegistryPassword, DestRegistry, false)

	_, err = skopeoHelper.copy(
		imagePushInfo.image,
		getDestImageUri(imagePushInfo.image),
		true,
		false)
	if err != nil {
		logError(err)
		return events.APIGatewayProxyResponse{}, ErrInternalServerResponse
	}

	return events.APIGatewayProxyResponse{
		Body: body,
		StatusCode: statusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getImagePushInfo(input string) (ImagePushInfo, error) {
	var pushInfo ImagePushInfo
	err := json.Unmarshal([]byte(input), &pushInfo)
	return pushInfo, err
}

func getDestImageUri(image string) string {
	return fmt.Sprintf("%v/lab/%v", DestRegistry, image)
}

func logError(err error) {
	fmt.Printf("Error: %v\n", err)
}