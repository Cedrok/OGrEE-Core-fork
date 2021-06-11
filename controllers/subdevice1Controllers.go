package controllers

import (
	"encoding/json"
	"net/http"
	"p3/models"
	u "p3/utils"
)

// swagger:operation POST /api/user/subdevices1 subdevices1 CreateSubdevice1
// Creates a Subdevice1 in the system.
// ---
// produces:
// - application/json
// parameters:
// - name: Name
//   in: query
//   description: Name of subdevice1
//   required: true
//   type: string
//   default: "Subdevice1A"
// - name: Category
//   in: query
//   description: Category of Subdevice1 (ex. Consumer Electronics, Medical)
//   required: true
//   type: string
//   default: "internal"
// - name: Description
//   in: query
//   description: Description of Subdevice1
//   required: false
//   type: string[]
//   default: ["Some abandoned subdevice1 in Grenoble"]
// - name: Domain
//   description: 'Domain of Subdevice1'
//   required: true
//   type: string
//   default: "Some Domain"
// - name: ParentID
//   description: 'Parent of Subdevice1 refers to Rack ID'
//   required: true
//   type: int
//   default: 999
// - name: Orientation
//   in: query
//   description: 'Indicates the location. Only values of
//   "front", "rear", "frontflipped", "rearflipped" are acceptable'
//   required: true
//   type: string
//   default: "front"
// - name: Template
//   in: query
//   description: 'Subdevice1 template'
//   required: true
//   type: string
//   default: "Some Template"
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
//   description: 'Size of Subdevice in an XY coordinate format'
//   required: true
//   type: string
//   default: "{\"x\":25.0,\"y\":29.399999618530275}"
// - name: SizeUnit
//   in: query
//   description: 'Extraneous Size Unit Attribute'
//   required: false
//   type: string
//   default: "{\"x\":999,\"y\":999}"
// - name: SizeU
//   in: query
//   description: 'The unit for Subdevice1 Size. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: Slot
//   in: query
//   description: 'Subdevice1 Slot (if any)'
//   required: false
//   type: string
//   default: "01"
// - name: PosU
//   in: query
//   description: 'Extraneous Position Unit Attribute'
//   required: false
//   type: string
//   default: "???"
// - name: Height
//   in: query
//   description: 'Height of Subdevice1'
//   required: true
//   type: string
//   default: "5"
// - name: HeightU
//   in: query
//   description: 'The unit for Subdevice1 Height. Only values of
//   "mm", "cm", "m", "U", "OU", "tile" are acceptable'
//   required: true
//   type: string
//   default: "m"
// - name: Vendor
//   in: query
//   description: 'Vendor of Subdevice1'
//   required: false
//   type: string
//   default: "Some Vendor"
// - name: Model
//   in: query
//   description: 'Model of Subdevice1'
//   required: false
//   type: string
//   default: "Some Model"
// - name: Type
//   in: query
//   description: 'Type of Subdevice1'
//   required: false
//   type: string
//   default: "Some Type"
// - name: Serial
//   in: query
//   description: 'Serial of Subdevice1'
//   required: false
//   type: string
//   default: "Some Serial"

// responses:
//     '201':
//         description: Created
//     '400':
//         description: Bad request
var CreateSubdevice1 = func(w http.ResponseWriter, r *http.Request) {

	subdevice1 := &models.Subdevice1{}
	err := json.NewDecoder(r.Body).Decode(subdevice1)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		u.ErrLog("Error while decoding request body", "CREATE SUBDEVICE1", "", r)
		return
	}

	resp, e := subdevice1.Create()

	switch e {
	case "validate":
		w.WriteHeader(http.StatusBadRequest)
		u.ErrLog("Error while creating Subdevice1", "CREATE SUBDEVICE1", e, r)
	case "internal":
		//
	default:
		w.WriteHeader(http.StatusCreated)
	}
	u.Respond(w, resp)
}
