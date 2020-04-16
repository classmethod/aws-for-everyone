package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

type PersonRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Person struct {
	Id        string `dynamo:"Id"`
	FirstName string `dynamo:"FirstName"`
	LastName  string `dynamo:"LastName"`
}

type AwsSession struct {
	Sess *session.Session
	Err  error
}

var awsSess AwsSession

func init() {
	awsSess.Sess, awsSess.Err = session.NewSession()
}

func createResponse(code int, msg string) events.APIGatewayProxyResponse {
	header := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
	}
	res := events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    header,
		Body:       msg,
	}

	return res
}

func getPersons(table dynamo.Table) events.APIGatewayProxyResponse {
	var persons []Person

	err := table.Scan().All(&persons)
	if err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("scan error: %s", err.Error()))
	}

	json, err := json.Marshal(persons)
	if err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("create json error: %s", err.Error()))
	}

	return createResponse(http.StatusOK, string(json))
}

func addPerson(table dynamo.Table, reqBody string) events.APIGatewayProxyResponse {
	id := uuid.New()

	// Bodyを構造体に変換
	var personReq PersonRequest
	if err := json.Unmarshal([]byte(reqBody), &personReq); err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("decode json error: %s", err.Error()))
	}

	// 書き込むための構造体を作成
	person := Person{
		Id:        id.String(),
		LastName:  personReq.LastName,
		FirstName: personReq.FirstName,
	}

	err := table.Put(person).Run()
	if err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("add person error: %s", err.Error()))
	}

	res, err := json.Marshal(person)
	if err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("create json error: %s", err.Error()))
	}

	return createResponse(http.StatusCreated, string(res))
}

func deletePerson(table dynamo.Table, id string) events.APIGatewayProxyResponse {
	err := table.Delete("Id", id).Run()
	if err != nil {
		return createResponse(http.StatusInternalServerError,
			fmt.Sprintf("delete person error: %s", err.Error()))
	}

	return createResponse(http.StatusNoContent, "")
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AWS SDKのセッション作成でエラーが発生した場合の処理
	if awsSess.Err != nil {
		log.Printf("create aws session error: %s", awsSess.Err.Error())
		return createResponse(http.StatusInternalServerError,
			fmt.Sprint("internal server error")), nil
	}

	ddb := dynamo.New(awsSess.Sess)
	table := ddb.Table(os.Getenv("TABLE_NAME"))

	switch {
	// GET /persons
	case req.HTTPMethod == http.MethodGet && req.Path == "/persons":
		return getPersons(table), nil
	// POST /persons
	case req.HTTPMethod == http.MethodPost && req.Path == "/persons":
		return addPerson(table, req.Body), nil
	// DELETE /persons/{personId}
	case req.HTTPMethod == http.MethodDelete &&
		req.Path == fmt.Sprintf("/persons/%s", req.PathParameters["personId"]):
		return deletePerson(table, req.PathParameters["personId"]), nil
	}

	return createResponse(http.StatusNotFound, "not found path or not allowed method"), nil
}

func main() {
	lambda.Start(Handler)
}
