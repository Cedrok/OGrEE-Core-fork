package controllers

//This file contains code associated with initialising the Shell

import (
	"cli/com"
	l "cli/logger"
	"cli/models"
	"cli/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func InitConfigFilePath(configPath string) {
	State.ConfigPath = configPath
}

func InitHistoryFilePath(histPath string) {
	State.HistoryFilePath = histPath
}

func InitDebugLevel(verbose string) {
	var ok bool
	State.DebugLvl, ok = map[string]int{
		"NONE":    NONE,
		"ERROR":   ERROR,
		"WARNING": WARNING,
		"INFO":    INFO,
		"DEBUG":   DEBUG,
	}[verbose]
	if !ok {
		Println("Invalid Logging Mode detected. Resorting to default: ERROR")
		State.DebugLvl = 1
	}
}

func PingAPI() bool {
	_, err := models.Send("", State.APIURL, "", nil)
	return err == nil
}

// Intialise the ShellState
func InitState(conf *com.Config) error {
	State.Hierarchy = BuildBaseTree()
	if State.CurrPath == "" {
		State.CurrPath = "/Physical"
	}
	if State.PrevPath == "" {
		State.PrevPath = "/Physical"
	}
	//Set the filter attributes setting
	State.FilterDisplay = false

	//Set which objects Unity will be notified about
	State.ObjsForUnity = SetObjsForUnity(conf.Updates)
	State.DrawableObjs = SetObjsForUnity(conf.Drawable)
	State.DrawableJsons = make(map[string]map[string]interface{}, 16)

	for i := SITE; i < GROUP+1; i++ {
		ent := EntityToString(i)
		State.DrawableJsons[ent] = SetDrawableTemplate(ent, conf.DrawableJson)
	}

	//Set Draw Threshold
	SetDrawThreshold(int(conf.DrawLimit))

	resp, err := RequestAPI("GET", "/api/version", nil, http.StatusOK)
	if err != nil {
		return err
	}
	info, ok := resp.body["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid response from API on GET /api/version")
	}
	if cInf, ok := info["Customer"]; ok {
		if customer, ok := cInf.(string); ok {
			State.Customer = customer
		}
	}
	if State.Customer == "" {
		if State.DebugLvl > NONE {
			Println("Tenant Information not found!")
		}
		State.Customer = "UNKNOWN"
	}
	return nil
}

// Tries to establish a connection with OGrEE-3D and, if possible,
// starts a go routine for receiving messages from it
func InitOGrEE3DCommunication() error {
	errConnect := models.Ogree3D.Connect(State.Ogree3DURL, State.Timeout)
	if errConnect != nil {
		return ErrorWithInternalError{
			UserError:     errors.New("OGrEE-3D is not reachable"),
			InternalError: errConnect,
		}
	}

	errLogin := models.Ogree3D.Login(State.APIURL, GetKey(), State.DebugLvl)
	if errLogin != nil {
		return ErrorWithInternalError{
			UserError:     errors.New("OGrEE-3D login not possible"),
			InternalError: errLogin,
		}
	}

	Println("Established connection with OGrEE-3D!")

	go models.Ogree3D.ReceiveLoop()

	return nil
}

func InitTimeout(duration string) {
	var timeLen int
	var durationType string
	_, err := fmt.Sscanf(duration, "%d%s", &timeLen, &durationType)
	if err != nil {
		msg := "Invalid unity timeout format. Resorting to default of 10ms"
		l.GetWarningLogger().Println(msg)
		if State.DebugLvl > 1 {
			Println(msg)
		}
		State.Timeout = time.Millisecond * time.Duration(10)
	}
	switch durationType {
	case "ns":
		State.Timeout = time.Nanosecond * time.Duration(timeLen)
	case "us":
		State.Timeout = time.Microsecond * time.Duration(timeLen)
	case "ms":
		State.Timeout = time.Millisecond * time.Duration(timeLen)
	case "s":
		State.Timeout = time.Second * time.Duration(timeLen)
	default:
		msg := "Invalid duration unit found. Resorting to default of ms"
		l.GetWarningLogger().Println(msg)
		if State.DebugLvl > 1 {
			Println(msg)
		}
		State.Timeout = time.Millisecond * time.Duration(timeLen)
	}
}

func InitKey(apiKey string) {
	if apiKey != "" {
		State.APIKEY = apiKey
	} else {
		Println("Error: No API Key Found")
		if State.DebugLvl > 0 {
			l.GetErrorLogger().Println(
				"No API Key provided in env file nor as argument")
		}
	}
}

func InitURLs(apiURL string, ogree3DURL string) {
	apiURL = strings.TrimRight(apiURL, "/")
	_, err := url.ParseRequestURI(apiURL)
	if err != nil {
		msg := "apiURL is not valid!\n"
		msg += "Falling back to default API URL: http://localhost:3001"
		Println(msg)
		l.GetInfoLogger().Println(msg)
		State.APIURL = "http://localhost:3001"
	} else {
		State.APIURL = apiURL
	}

	err = State.SetOgree3DURL(ogree3DURL)
	if err != nil {
		Println(err.Error())
		State.SetDefaultOgree3DURL()
	}
}

// Helper for InitState will insert objs
// and init DrawThreshold
func SetObjsForUnity(objs []string) []int {
	res := []int{}
	allDetected := false
	for _, obj := range objs {
		obj = strings.ToLower(obj)
		if val := EntityStrToInt(obj); val != -1 {
			res = append(res, val)
		} else if obj == "all" {
			//Exit the loop and use default code @ end of function
			allDetected = true
			break
		}
	}
	//Set the array to all and exit
	//GROUP is the greatest value int enum type
	//So we use that for the cond guard
	if allDetected {
		for idx := 0; idx < GROUP+1; idx++ {
			res = append(res, idx)
		}
		res = append(res, DOMAIN)
	}
	return res
}

func SetDrawThreshold(limit int) {
	//Set Draw Threshold
	if limit < 0 {
		if State.DebugLvl > 0 {
			Println("Setting Draw Limit to default")
		}
		State.DrawThreshold = 50 //50 is default value
	} else {
		State.DrawThreshold = limit
	}
}

func SetDrawableTemplate(entity string, DrawableJson map[string]string) map[string]interface{} {
	if objStr, ok := DrawableJson[entity]; ok && objStr != "" {
		objStr = strings.Trim(objStr, "'\"")
		//Now retrieve file
		ans := map[string]interface{}{}
		if !filepath.IsAbs(objStr) {
			objStr = utils.ExeDir() + "/" + objStr
		}
		f, e := os.ReadFile(objStr)
		if e == nil {
			json.Unmarshal(f, &ans)
			return ans
		}
	}
	l.GetWarningLogger().Println("Specified template for " + entity + " not found")
	if State.DebugLvl > 1 {
		Println("Specified template for " + entity +
			" not found, resorting to defaults")
	}
	return nil
}

func Login(user string, password string) (*User, string, error) {
	data := map[string]any{"email": user, "password": password}
	resp, err := RequestAPI("POST", "/api/login", data, http.StatusOK)
	if err != nil {
		return nil, "", err
	}
	account, accountOk := (resp.body["account"].(map[string]interface{}))
	token, tokenOk := account["token"].(string)
	userID, userIDOk := account["_id"].(string)
	if !accountOk || !tokenOk || !userIDOk {
		return nil, "", fmt.Errorf("invalid response from API")
	}
	return &User{user, userID}, token, nil
}
