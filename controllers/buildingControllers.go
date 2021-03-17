package controllers

import (
	"encoding/json"
	"net/http"
	"p3/models"
	u "p3/utils"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:operation POST /api/user/buildings buildings Create
// Creates a Building in the system.
// ---
// produces:
// - application/json
// parameters:
// - name: Name
//   in: query
//   description: Name of building
//   required: true
//   type: string
//   default: "Building A"
// - name: ParentID
//   description: 'ParentID of Building refers to Site'
//   required: true
//   type: string
//   default: "999"
// - name: Category
//   in: query
//   description: Category of Building (ex. Consumer Electronics, Medical)
//   required: true
//   type: string
//   default: "Research"
// - name: Description
//   in: query
//   description: Description of Building
//   required: false
//   type: string[]
//   default: ["Some abandoned building in Grenoble"]
// - name: Domain
//   description: 'Domain Of Building'
//   required: true
//   type: string
//   default: "Some Domain"
// - name: PosXY
//   in: query
//   description: 'Indicates the position in a XY coordinate format'
//   required: true
//   type: string
//   default: "{\"x\":-30.0,\"y\":0.0}"
// - name: PosXYU
//   in: query
//   description: 'Indicates the unit of the PosXY position. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: PosZ
//   in: query
//   description: 'Indicates the position in the Z axis'
//   required: true
//   type: string
//   default: "10"
// - name: PosZU
//   in: query
//   description: 'Indicates the unit of the Z coordinate position. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: Size
//   in: query
//   description: 'Size of Building in an XY coordinate format'
//   required: true
//   type: string
//   default: "{\"x\":25.0,\"y\":29.399999618530275}"
// - name: SizeU
//   in: query
//   description: 'The unit for Building Size. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: Height
//   in: query
//   description: 'Height of Building'
//   required: true
//   type: string
//   default: "5"
// - name: HeightU
//   in: query
//   description: 'The unit for Building Height. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: Floors
//   in: query
//   description: 'Number of floors'
//   required: true
//   type: string
//   default: "3"

// responses:
//     '200':
//         description: Created
//     '400':
//         description: Bad request

var CreateBuilding = func(w http.ResponseWriter, r *http.Request) {

	bldg := &models.Building{}
	err := json.NewDecoder(r.Body).Decode(bldg)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp, e := bldg.Create()
	switch e {
	case "validate":
		//w.WriteHeader(http.)
	case "internal":
		//
	}
	u.Respond(w, resp)
}

// swagger:operation GET /api/user/buildings/{id} buildings GetBuilding
// Gets Building using Building ID.
// ---
// produces:
// - application/json
// parameters:
// - name: ID
//   in: path
//   description: ID of Building
//   required: true
//   type: int
//   default: 999
// responses:
//     '200':
//         description: Found
//     '400':
//         description: Not Found

//Retrieve bldg using Bldg ID
var GetBuilding = func(w http.ResponseWriter, r *http.Request) {

	id, e := strconv.Atoi(mux.Vars(r)["id"])
	resp := u.Message(true, "success")

	if e != nil {
		u.Respond(w, u.Message(false, "Error while parsing path parameters"))
	}

	data, e1 := models.GetBuilding(uint(id))
	if data == nil {
		resp = u.Message(false, "Error while getting Building: "+e1)

		switch e1 {
		case "validate":
			//
		default:
		}

	} else {
		resp = u.Message(true, "success")
	}

	resp["data"] = data
	u.Respond(w, resp)
}

// swagger:operation GET /api/user/buildings buildings GetAllBuildings
// Gets All Buildings in the system.
// ---
// produces:
// - application/json
// parameters:
// responses:
//     '200':
//         description: Found
//     '400':
//         description: Not Found

var GetAllBuildings = func(w http.ResponseWriter, r *http.Request) {

	resp := u.Message(true, "success")

	data := models.GetAllBuildings()
	if data == nil {
		resp = u.Message(false, "unsuccessful")
	}

	resp["data"] = data
	u.Respond(w, resp)
}

// swagger:operation DELETE /api/user/buildings/{id} buildings DeleteBuilding
// Deletes a Building.
// ---
// produces:
// - application/json
// parameters:
// - name: ID
//   in: path
//   description: ID of desired building
//   required: true
//   type: int
//   default: 999
// responses:
//     '204':
//        description: Successful
//     '400':
//        description: Not found
var DeleteBuilding = func(w http.ResponseWriter, r *http.Request) {
	id, e := strconv.Atoi(mux.Vars(r)["id"])

	if e != nil {
		u.Respond(w, u.Message(false, "Error while parsing path parameters"))
	}

	v := models.DeleteBuilding(uint(id))
	u.Respond(w, v)
}

// swagger:operation PUT /api/user/buildings/{id} buildings UpdateBuilding
// Changes Building data in the system.
// If no new or any information is provided
// an OK will still be returned
// ---
// produces:
// - application/json
// parameters:
// - name: ID
//   in: path
//   description: ID of desired building
//   required: true
//   type: int
//   default: 999
// - name: Name
//   in: query
//   description: Name of building
//   required: false
//   type: string
//   default: "Building B"
// - name: Category
//   in: query
//   description: Category of Building (ex. Consumer Electronics, Medical)
//   required: false
//   type: string
//   default: "New Building"
// - name: Description
//   in: query
//   description: Description of Building
//   required: false
//   type: string[]
//   default: ["Derelict", "Building"]
// - name: Domain
//   description: 'Domain Of Building'
//   required: false
//   type: string
//   default: "Derelict Domain"
// - name: PosXY
//   in: query
//   description: 'Indicates the position in a XY coordinate format'
//   required: false
//   type: string
//   default: "{\"x\":999,\"y\":999}"
// - name: PosXYU
//   in: query
//   description: 'Indicates the unit of the PosXY position. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: false
//   type: string
//   default: "cm"
// - name: PosZ
//   in: query
//   description: 'Indicates the position in the Z axis'
//   required: false
//   type: string
//   default: "999"
// - name: PosZU
//   in: query
//   description: 'Indicates the unit of the Z coordinate position. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: false
//   type: string
//   default: "cm"
// - name: Size
//   in: query
//   description: 'Size of Building in an XY coordinate format'
//   required: false
//   type: string
//   default: "{\"x\":999,\"y\":999}"
// - name: SizeU
//   in: query
//   description: 'The unit for Building Size. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: false
//   type: string
//   default: "cm"
// - name: Height
//   in: query
//   description: 'Height of Building'
//   required: false
//   type: string
//   default: "999"
// - name: HeightU
//   in: query
//   description: 'The unit for Building Height. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: false
//   type: string
//   default: "cm"
// - name: Floors
//   in: query
//   description: 'Number of floors'
//   required: false
//   type: string
//   default: "999"

// responses:
//     '200':
//         description: Updated
//     '400':
//         description: Bad request
//Updates work by passing ID in path parameter
var UpdateBuilding = func(w http.ResponseWriter, r *http.Request) {

	bldg := &models.Building{}
	id, e := strconv.Atoi(mux.Vars(r)["id"])

	if e != nil {
		u.Respond(w, u.Message(false, "Error while parsing path parameters"))
	}

	err := json.NewDecoder(r.Body).Decode(bldg)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
	}

	v := models.UpdateBuilding(uint(id), bldg)
	u.Respond(w, v)
}
