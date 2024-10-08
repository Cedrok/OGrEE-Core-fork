package utils

//Builds json messages and
//returns json response

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/elliotchance/pie/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BuildHash string
var BuildTree string
var BuildTime string
var GitCommitDate string

const AsteriskLine = "******************************************************"

const (
	DOMAIN = iota
	// hierarchical root entities
	STRAYOBJ
	SITE
	// hierarchical entities with mandatory parent
	BLDG
	ROOM
	RACK
	DEVICE
	AC
	CABINET
	CORRIDOR
	GENERIC
	PWRPNL
	GROUP
	// logical non hierarchical entities
	ROOMTMPL
	OBJTMPL
	BLDGTMPL
	TAG
	LAYER
	VIRTUALOBJ
)

type Namespace string

const (
	Any            Namespace = ""
	Physical       Namespace = "physical"
	PStray         Namespace = "physical.stray"
	PHierarchy     Namespace = "physical.hierarchy"
	Organisational Namespace = "organisational"
	Logical        Namespace = "logical"
	LObjTemplate   Namespace = "logical.objtemplate"
	LBldgTemplate  Namespace = "logical.bldgtemplate"
	LRoomTemplate  Namespace = "logical.roomtemplate"
	LTags          Namespace = "logical.tag"
	LLayers        Namespace = "logical.layer"
)

const HN_DELIMETER = "."  // hierarchyName path delimiter
const RESET_TAG = "RESET" // used as email to identify a reset token
const HIERARCHYOBJS_ENT = "hierarchy_object"

type RequestFilters struct {
	FieldsToShow []string  `schema:"fieldOnly"`
	StartDate    string    `schema:"startDate"`
	EndDate      string    `schema:"endDate"`
	Limit        string    `schema:"limit"`
	Namespace    Namespace `schema:"namespace"`
	Id           string    `schema:"id"`
}

type LayerObjsFilters struct {
	Root        string `schema:"root"`
	IsRecursive bool   `schema:"recursive"`
}

type HierarchyFilters struct {
	Namespace      Namespace `schema:"namespace"`
	StartDate      string    `schema:"startDate"`
	EndDate        string    `schema:"endDate"`
	Limit          string    `schema:"limit"`
	WithCategories bool      `schema:"withcategories"`
}

// Error definitions
type ErrType int

const (
	ErrUnauthorized ErrType = iota
	ErrForbidden
	ErrDuplicate
	ErrBadFormat
	ErrInvalidValue
	ErrDBError
	ErrInternal
	ErrNotFound
	WarnShouldChangePass
)

type Error struct {
	Type    ErrType
	Message string
	Details []string
}

func (err Error) Error() string {
	return err.Message
}

func GetBuildDate() string {
	return BuildTime
}

func GetCommitDate() string {
	return GitCommitDate
}

func GetBuildHash() string {
	return BuildHash
}

func GetBuildTree() string {
	return BuildTree
}

func Connect() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

const HttpResponseContentType = "application/json"

func WriteOptionsHeader(w http.ResponseWriter, methods string) {
	w.Header().Add("Content-Type", HttpResponseContentType)
	w.Header().Add("Allow", methods+", OPTIONS")
}

func RespDataWrapper(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"message": message, "data": data}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
	w.Header().Add("Content-Type", HttpResponseContentType)
}

func RespondWithError(w http.ResponseWriter, err *Error) {
	errMap := map[string]interface{}{"message": err.Message}
	if len(err.Details) > 0 {
		errMap["errors"] = err.Details
	}
	w.WriteHeader(ErrTypeToStatusCode(err.Type))
	json.NewEncoder(w).Encode(errMap)
	w.Header().Add("Content-Type", HttpResponseContentType)
}

func ErrLog(message, funcname, details string, r *http.Request) {
	if flag.Lookup("test.v") != nil {
		return
	}
	f, err := os.OpenFile("resources/debug.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	ip := r.RemoteAddr

	log.SetOutput(f)
	log.Println(message + " FOR FUNCTION: " + funcname)
	log.Println("FROM IP: " + ip)
	log.Println(details)
}

func FilteredReqFromQueryParams(link *url.URL) bson.M {
	queryValues, _ := url.ParseQuery(link.RawQuery)
	bsonMap := bson.M{}

	for key := range queryValues {
		if key != "fieldOnly" && key != "startDate" && key != "endDate" &&
			key != "limit" && key != "namespace" {
			keyValue := ConvertString(queryValues.Get(key))
			AddFilterToReq(bsonMap, key, keyValue)
		}
	}
	return bsonMap
}

func AddFilterToReq(bsonMap primitive.M, key string, keyValue any) {
	if key == "parentId" {
		regex := applyWildcards(keyValue.(string)) + `\.(` + NAME_REGEX + ")"
		bsonMap["id"] = regexToMongoFilter(regex)
		return
	} else if key == "tag" {
		// tag is in tags list
		bsonMap["tags"] = bson.M{"$eq": keyValue}
		return
	} else if reflect.TypeOf(keyValue).Kind() == reflect.String &&
		strings.Contains(keyValue.(string), "*") {
		regex := applyWildcards(keyValue.(string))
		keyValue = regexToMongoFilter(regex)
	}

	switch key {
	case "id", "name", "category",
		"description", "domain",
		"createdDate", "lastUpdated", "slug":
		bsonMap[key] = keyValue
	default:
		bsonMap["attributes."+key] = keyValue
	}
}

func ErrTypeToStatusCode(errType ErrType) int {
	switch errType {
	case ErrForbidden:
		return http.StatusForbidden
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrDuplicate, ErrBadFormat:
		return http.StatusBadRequest
	case ErrDBError, ErrInternal:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

var Entities = []int{
	DOMAIN,
	STRAYOBJ, SITE,
	BLDG, ROOM, RACK, DEVICE, AC, CABINET, CORRIDOR, GENERIC, PWRPNL, GROUP,
	ROOMTMPL, OBJTMPL, BLDGTMPL, TAG, LAYER, VIRTUALOBJ,
}

var EntitiesWithTags = []int{
	STRAYOBJ, SITE, BLDG, ROOM, RACK, DEVICE, AC, CABINET, CORRIDOR, GENERIC, PWRPNL, GROUP, VIRTUALOBJ,
}

var EntitiesWithAttributeCheck = []int{
	CORRIDOR, GROUP, DEVICE, VIRTUALOBJ,
}

var RoomChildren = []int{RACK, CORRIDOR, GENERIC}

func EntityHasTags(entity int) bool {
	return pie.Contains(EntitiesWithTags, entity)
}

func IsEntityHierarchical(entity int) bool {
	return !IsEntityNonHierarchical(entity)
}

func IsEntityNonHierarchical(entity int) bool {
	return entity >= ROOMTMPL && entity < VIRTUALOBJ
}

func EntityToString(entity int) string {
	switch entity {
	case SITE:
		return "site"
	case BLDG:
		return "building"
	case ROOM:
		return "room"
	case RACK:
		return "rack"
	case DEVICE:
		return "device"
	case AC:
		return "ac"
	case PWRPNL:
		return "panel"
	case DOMAIN:
		return "domain"
	case STRAYOBJ:
		return "stray_object"
	case ROOMTMPL:
		return "room_template"
	case OBJTMPL:
		return "obj_template"
	case BLDGTMPL:
		return "bldg_template"
	case CABINET:
		return "cabinet"
	case GROUP:
		return "group"
	case CORRIDOR:
		return "corridor"
	case GENERIC:
		return "generic"
	case TAG:
		return "tag"
	case LAYER:
		return "layer"
	case VIRTUALOBJ:
		return "virtual_obj"
	default:
		return "INVALID"
	}
}

func EntitiesStrToInt(entities []string) []int {
	return pie.Map[string, int](entities, EntityStrToInt)
}

func EntityStrToInt(entity string) int {
	switch entity {
	case "site":
		return SITE
	case "building", "bldg":
		return BLDG
	case "room":
		return ROOM
	case "rack":
		return RACK
	case "device":
		return DEVICE
	case "ac":
		return AC
	case "panel":
		return PWRPNL
	case "domain":
		return DOMAIN
	case "stray_object":
		return STRAYOBJ
	case "room_template":
		return ROOMTMPL
	case "obj_template":
		return OBJTMPL
	case "bldg_template":
		return BLDGTMPL
	case "cabinet":
		return CABINET
	case "group":
		return GROUP
	case "corridor":
		return CORRIDOR
	case "generic":
		return GENERIC
	case "tag":
		return TAG
	case "layer":
		return LAYER
	case "virtual_obj":
		return VIRTUALOBJ
	default:
		return -1
	}
}

func NamespaceToString(namespace Namespace) string {
	ref := reflect.ValueOf(namespace)
	return ref.String()
}

func GetEntitiesById(namespace Namespace, hierarchyId string) []string {
	var entNames []string
	switch namespace {
	case Organisational:
		entNames = append(entNames, EntityToString(DOMAIN))
	case Logical:
		for entity := GROUP; entity <= VIRTUALOBJ; entity++ {
			entNames = append(entNames, EntityToString(entity))
		}
	case LObjTemplate:
		entNames = append(entNames, EntityToString(OBJTMPL))
	case LBldgTemplate:
		entNames = append(entNames, EntityToString(BLDGTMPL))
	case LRoomTemplate:
		entNames = append(entNames, EntityToString(ROOMTMPL))
	case LTags:
		entNames = append(entNames, EntityToString(TAG))
	case LLayers:
		entNames = append(entNames, EntityToString(LAYER))
	case PStray:
		entNames = append(entNames, EntityToString(STRAYOBJ))
	case Physical, PHierarchy, Any:
		entities := []int{VIRTUALOBJ}

		if hierarchyId == "" || hierarchyId == "**" {
			// All entities of each namespace
			switch namespace {
			case Physical:
				for entity := STRAYOBJ; entity <= GROUP; entity++ {
					entities = append(entities, entity)
				}
			case PHierarchy:
				for entity := SITE; entity <= GROUP; entity++ {
					entities = append(entities, entity)
				}
			case Any:
				entities = Entities
			}
		} else {
			if namespace == Any {
				entities = append(entities, DOMAIN)
			}

			// Add entities according to hierarchyName possibilities
			if strings.Contains(hierarchyId, ".**") {
				var initialEntity int
				finalEntity := GROUP

				switch strings.Count(hierarchyId, HN_DELIMETER) {
				case 1, 2:
					initialEntity = BLDG
				case 3:
					initialEntity = ROOM
				case 4:
					initialEntity = RACK
				case 5:
					initialEntity = DEVICE
				default:
					// only devices
					initialEntity = DEVICE
					finalEntity = DEVICE
				}

				for entity := initialEntity; entity <= finalEntity; entity++ {
					entities = append(entities, entity)
				}
			} else {
				switch strings.Count(hierarchyId, HN_DELIMETER) {
				case 0:
					entities = append(entities, SITE)
					if namespace == Any {
						entities = append(entities, OBJTMPL, ROOMTMPL, BLDGTMPL, TAG, LAYER)
					}
					if namespace == Any || namespace == Physical {
						entities = append(entities, STRAYOBJ)
					}
				case 1:
					entities = append(entities, BLDG)
				case 2:
					entities = append(entities, ROOM)
				case 3:
					entities = append(entities, RACK, AC, CORRIDOR, PWRPNL, CABINET, GROUP, GENERIC)
				case 4:
					entities = append(entities, DEVICE, GROUP)
				default:
					entities = append(entities, DEVICE)
				}
			}
		}

		// Convert entities to string
		for _, entInt := range entities {
			entNames = append(entNames, EntityToString(entInt))
		}
	}

	return entNames
}

func GetParentOfEntityByInt(entity int) int {
	switch entity {
	case DOMAIN:
		return DOMAIN
	case AC, PWRPNL, CABINET, CORRIDOR, GENERIC:
		return ROOM
	case ROOMTMPL, OBJTMPL, BLDGTMPL, TAG, GROUP, STRAYOBJ, LAYER:
		return -1
	default:
		return entity - 1
	}
}

func FormatNotifyData(msgType, entityStr string, data any) string {
	if entityStr == "tag" {
		msgType = msgType + "-tag"
	} else if entityStr == "layer" {
		msgType = msgType + "-layer"
	}
	//convert to json then string
	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(map[string]any{"type": msgType, "data": data})
	if err != nil {
		println("Error notifying 3D client: unable to encode json data")
	}
	return buff.String()
}

// Helper functions
func StrSliceContains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

var floatType = reflect.TypeOf(float64(0))

func GetFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func ConvertString(strValue string) any {
	if num, err := strconv.ParseFloat(strValue, 64); err == nil {
		// is number
		return num
	} else if nums, err := StringToFloatSlice(strValue); err == nil {
		// is array of numbers
		return nums
	} else if strs, err := StringToStrSlice(strValue); err == nil {
		// is array of strings
		return strs
	} else if boolean, err := strconv.ParseBool(strValue); err == nil {
		// is boolean
		return boolean
	}
	// is string
	return strValue
}

func StringToFloatSlice(strValue string) ([]float64, error) {
	numbers := []float64{}
	if len(strValue) < 2 || strValue[0] != '[' || strValue[len(strValue)-1:] != "]" {
		return numbers, fmt.Errorf("not a vector")
	}
	strSplit := strings.Split(strValue[1:len(strValue)-1], ",")
	for _, val := range strSplit {
		if n, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
			numbers = append(numbers, n)
		} else {
			return numbers, fmt.Errorf("invalid vector format")
		}
	}
	return numbers, nil
}

func StringToStrSlice(strValue string) ([]string, error) {
	if len(strValue) < 2 || strValue[0] != '[' || strValue[len(strValue)-1:] != "]" {
		return []string{}, fmt.Errorf("not a vector")
	}
	strs := strings.Split(strings.ReplaceAll(strValue[1:len(strValue)-1], "\"", ""), ",")
	if len(strs[0]) <= 0 {
		return strs, fmt.Errorf("invalid vector format")
	}
	for i := range strs {
		strs[i] = strings.TrimSpace(strs[i])
	}
	return strs, nil
}
