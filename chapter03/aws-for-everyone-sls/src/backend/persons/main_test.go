package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	dynamodbEndpoint = "http://localhost:4569"
	region           = "ap-northeast-1"
)

func createSession(t *testing.T) AwsSession {
	t.Helper()

	awsSess = AwsSession{}

	awsSess.Sess, awsSess.Err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("DUMMY", "DUMMY", "DUMMY"),
		Endpoint:    aws.String(dynamodbEndpoint),
		Region:      aws.String(region),
	})
	if awsSess.Err != nil {
		t.Fatal(awsSess.Err)
	}

	return awsSess
}

func createTable(t *testing.T, ddb *dynamo.DB, tableName string) {
	t.Helper()

	type PersonsTable struct {
		Id string `dynamo:"Id,hash"`
	}

	res := ddb.CreateTable(tableName, PersonsTable{})
	if err := res.Run(); err != nil {
		t.Fatal(err)
	}
}

func createRecord(t *testing.T, ddb *dynamo.DB, tableName string) []Person {
	t.Helper()

	table := ddb.Table(tableName)
	item := Person{Id: "556350d2-e993-4fb9-8242-c496a0664bb3", LastName: "Taro", FirstName: "Yamada"}

	err := table.Put(item).Run()
	if err != nil {
		t.Fatal(err)
	}

	return []Person{item}
}

func deleteTable(t *testing.T, ddb *dynamo.DB, tableName string) {
	t.Helper()

	table := ddb.Table(tableName)
	if err := table.DeleteTable().Run(); err != nil {
		t.Fatal(err)
	}
}

func TestGetPerson(t *testing.T) {
	awsSess = createSession(t)
	ddb := dynamo.New(awsSess.Sess)
	tableName := "test_table_for_get_persons"

	createTable(t, ddb, tableName)
	defer deleteTable(t, ddb, tableName)

	wantBody := createRecord(t, ddb, tableName)

	if err := os.Setenv("TABLE_NAME", tableName); err != nil {
		t.Fatal(err)
	}

	// Handlerは必ずerrを返さない
	req := events.APIGatewayProxyRequest{
		Path:       "/persons",
		HTTPMethod: http.MethodGet,
	}
	res, _ := Handler(req)

	var gotBody []Person
	if err := json.Unmarshal([]byte(res.Body), &gotBody); err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("got: %v, want: %v", res.StatusCode, http.StatusOK)
	}

	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Errorf("got: %v, want: %v", res.Body, wantBody)
	}
}

func TestAddPerson(t *testing.T) {
	awsSess = createSession(t)
	ddb := dynamo.New(awsSess.Sess)
	tableName := "test_table_for_add_person"

	createTable(t, ddb, tableName)
	defer deleteTable(t, ddb, tableName)

	if err := os.Setenv("TABLE_NAME", tableName); err != nil {
		t.Fatal(err)
	}

	personReq := PersonRequest{
		FirstName: "Taro",
		LastName:  "Yamada",
	}

	reqBody, err := json.Marshal(personReq)
	if err != nil {
		t.Fatal(err)
	}

	// Handlerは必ずerrを返さない
	req := events.APIGatewayProxyRequest{
		Path:       "/persons",
		HTTPMethod: http.MethodPost,
		Body:       string(reqBody),
	}
	res, _ := Handler(req)

	var gotBody Person
	if err := json.Unmarshal([]byte(res.Body), &gotBody); err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("got: %v, want: %v", res.StatusCode, http.StatusCreated)
	}

	if gotBody.Id == "" {
		t.Errorf("got: %v, want: [UUID]", gotBody.Id)
	}

	if gotBody.FirstName != personReq.FirstName {
		t.Errorf("got: %v, want: %v", gotBody.FirstName, personReq.FirstName)
	}

	if gotBody.LastName != personReq.LastName {
		t.Errorf("got: %v, want: %v", gotBody.LastName, personReq.LastName)
	}

	table := ddb.Table(tableName)

	var persons []Person

	err = table.Scan().All(&persons)
	if err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(persons)
	if err != nil {
		t.Fatal(err)
	}

	var gotData []Person

	err = json.Unmarshal(jsonBytes, &gotData)
	if err != nil {
		t.Fatal(err)
	}

	wantData := Person{
		Id:        "[UUID]",
		FirstName: "Taro",
		LastName:  "Yamada",
	}

	if gotData[0].Id == "" {
		t.Errorf("got: %v, want: [UUID]", gotData[0].Id)
	}

	if gotData[0].FirstName != wantData.FirstName {
		t.Errorf("got: %v, want: %v", gotData[0].FirstName, wantData.FirstName)
	}

	if gotData[0].LastName != wantData.LastName {
		t.Errorf("got: %v, want: %v", gotData[0].LastName, wantData.LastName)
	}
}

func TestDeletePerson(t *testing.T) {
	awsSess = createSession(t)
	ddb := dynamo.New(awsSess.Sess)
	tableName := "test_table_for_delete_person"

	createTable(t, ddb, tableName)
	defer deleteTable(t, ddb, tableName)

	reqBody := createRecord(t, ddb, tableName)

	if err := os.Setenv("TABLE_NAME", tableName); err != nil {
		t.Fatal(err)
	}

	pathParameters := map[string]string{
		"personId": reqBody[0].Id,
	}

	// Handlerは必ずerrを返さない
	req := events.APIGatewayProxyRequest{
		Path:           fmt.Sprintf("/persons/%s", reqBody[0].Id),
		HTTPMethod:     http.MethodDelete,
		PathParameters: pathParameters,
	}
	res, _ := Handler(req)

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("got: %v, want: %v", res.StatusCode, http.StatusNoContent)
	}

	if res.Body != "" {
		t.Errorf("got: %v, want: [Blank]", res.Body)
	}
}
