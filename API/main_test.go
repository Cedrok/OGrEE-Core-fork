package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"p3/app"
	"p3/models"
	u "p3/utils"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {
	// teardown()
	getAdminToken()
	exitCode := m.Run()
	os.Exit(exitCode)
}

var adminId primitive.ObjectID
var adminToken string

func getAdminToken() {
	// Create admin account
	admin := models.Account{}
	admin.Email = "admin@admin.com"
	admin.Password = "admin123"
	admin.Roles = map[string]models.Role{"*": "manager"}
	newAcc, _ := admin.Create(map[string]models.Role{"*": "manager"})
	if newAcc != nil {
		adminId = newAcc.ID
		adminToken = newAcc.Token
	}
}

func teardown() {
	ctx, _ := u.Connect()
	models.GetDB().Drop(ctx)
}

func makeRequest(method, url string, requestBody []byte) *httptest.ResponseRecorder {
	router := Router(app.JwtAuthentication)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateLoginAccount(t *testing.T) {
	// Test create new account
	requestBody := []byte(`{
		"email": "test@test.com",
    	"password": "pass123secret",
		"roles":{"*":"manager"}
	}`)
	recorder := makeRequest("POST", "/api/users", requestBody)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	_, exists := response["account"].(map[string]interface{})["token"].(string)
	assert.Equal(t, true, exists)

	// Test recreate existing account
	recorder = makeRequest("POST", "/api/users", requestBody)
	println(recorder.Body.String())
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// Test login
	recorder = makeRequest("POST", "/api/login", requestBody)
	assert.Equal(t, http.StatusOK, recorder.Code)

	json.Unmarshal(recorder.Body.Bytes(), &response)
	_, exists = response["account"].(map[string]interface{})["token"].(string)
	assert.Equal(t, true, exists)
}

func TestObjects(t *testing.T) {
	var response map[string]interface{}
	var parentId string

	// Create objects from schema examples
	for _, entInt := range []int{u.DOMAIN, u.SITE, u.BLDG, u.ROOM, u.RACK, u.DEVICE} {
		// Get object from schema
		entStr := u.EntityToString(entInt)
		data, _ := ioutil.ReadFile("models/schemas/" + entStr + "_schema.json")
		var obj map[string]interface{}
		json.Unmarshal(data, &obj)
		obj = obj["examples"].([]interface{})[0].(map[string]interface{})
		if entInt != u.SITE && entInt != u.DOMAIN {
			// Add parentId
			obj["parentId"] = parentId
		}
		data, _ = json.Marshal(obj)

		// Create (POST)
		recorder := makeRequest("POST", "/api/"+entStr+"s", data)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		json.Unmarshal(recorder.Body.Bytes(), &response)
		fmt.Println(response)
		id, exists := response["data"].(map[string]interface{})["id"].(string)
		assert.Equal(t, true, exists)

		// Verify create with GET
		recorder = makeRequest("GET", "/api/"+entStr+"s/"+id, nil)
		assert.Equal(t, http.StatusOK, recorder.Code)
		var responseGET map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &responseGET)
		delete(responseGET, "message")
		delete(response, "message")
		assert.Equal(t, true, reflect.DeepEqual(response, responseGET))

		// Update with PUT
		oldName := obj["name"].(string)
		obj["name"] = entStr + "Test"
		data, _ = json.Marshal(obj)
		recorder = makeRequest("PUT", "/api/"+entStr+"s/"+id, data)
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Verify it
		id = strings.Replace(id, oldName, obj["name"].(string), 1)
		recorder = makeRequest("GET", "/api/"+entStr+"s/"+id, nil)
		assert.Equal(t, http.StatusOK, recorder.Code)
		parentId = id
	}

	// Try to patch site name
	hierarchyName := "TESTPATCH"
	requestBody := []byte(`{
		"name": "` + hierarchyName + `"
	}`)
	recorder := makeRequest("PATCH", "/api/sites/siteTest", requestBody)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify the whole tree
	recorder = makeRequest("GET", "/api/sites/"+hierarchyName+"/all", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	json.Unmarshal(recorder.Body.Bytes(), &response)
	response = response["data"].(map[string]interface{})
	for _, entInt := range []int{u.BLDG, u.ROOM, u.RACK, u.DEVICE} {
		entStr := u.EntityToString(entInt)
		println(entStr)
		child := response["children"].([]interface{})[0].(map[string]interface{})
		hierarchyName = hierarchyName + u.HN_DELIMETER + entStr + "Test"
		assert.Equal(t, hierarchyName, child["id"].(string))
		if entInt != u.DEVICE {
			response = child
		}
	}

	// Delete everything
	recorder = makeRequest("DELETE", "/api/sites/TESTPATCH", nil)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	recorder = makeRequest("GET", "/api/hierarchy", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, 0,
		len(response["data"].(map[string]interface{})["tree"].(map[string]interface{})["physical"].(map[string]interface{})))
}
