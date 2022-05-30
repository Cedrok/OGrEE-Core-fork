package models

import (
	"context"
	"fmt"
	u "p3/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TENANT = iota
	SITE
	BLDG
	ROOM
	RACK
	DEVICE
	AC
	ROW
	CABINET
	CORIDOR
	PWRPNL
	SENSOR
	SEPARATOR
	TILE
	GROUP
	ROOMTMPL
	OBJTMPL
	STRAYDEV
	STRAYSENSOR
)

//Function will recursively iterate through nested obj
//and accumulate whatever is found into category arrays
func parseDataForNonStdResult(ent string, eNum, end int, data map[string]interface{}) map[string][]map[string]interface{} {
	var nxt string
	ans := map[string][]map[string]interface{}{}
	add := data[u.EntityToString(eNum+1)+"s"].([]map[string]interface{})

	//NEW REWRITE
	for i := eNum; i+2 < end; i++ {
		idx := u.EntityToString(i + 1)
		//println("trying IDX: ", idx)
		firstArr := add

		ans[idx+"s"] = firstArr

		for q := range firstArr {
			nxt = u.EntityToString(i + 2)
			println("NXT: ", nxt)
			ans[nxt+"s"] = append(ans[nxt+"s"],
				ans[idx+"s"][q][nxt+"s"].([]map[string]interface{})...)
		}
		add = ans[nxt+"s"]

	}

	return ans
}

//Mongo returns '_id' instead of id
func fixID(data map[string]interface{}) map[string]interface{} {
	if v, ok := data["_id"]; ok {
		data["id"] = v
		delete(data, "_id")
	}
	return data
}

func validateParent(ent string, entNum int, t map[string]interface{}) (map[string]interface{}, bool) {

	//Check ParentID is valid
	if t["parentId"] == nil {
		return u.Message(false, "ParentID is not valid"), false
	}

	objID, err := primitive.ObjectIDFromHex(t["parentId"].(string))
	if err != nil {
		return u.Message(false, "ParentID is not valid"), false
	}

	switch entNum {
	case DEVICE:
		x, _ := GetEntity(bson.M{"_id": objID}, "rack")
		y, _ := GetEntity(bson.M{"_id": objID}, "device")
		if x == nil && y == nil {
			return u.Message(false,
				"ParentID should be correspond to Existing ID"), false
		}
	case SENSOR, GROUP:
		w, _ := GetEntity(bson.M{"_id": objID}, "device")
		x, _ := GetEntity(bson.M{"_id": objID}, "rack")
		y, _ := GetEntity(bson.M{"_id": objID}, "room")
		z, _ := GetEntity(bson.M{"_id": objID}, "building")
		if w == nil && x == nil && y == nil && z == nil {
			return u.Message(false,
				"ParentID should be correspond to Existing ID"), false
		}
	default:
		parentInt := u.GetParentOfEntityByInt(entNum)
		parent := u.EntityToString(parentInt)

		ctx, cancel := u.Connect()
		if GetDB().Collection(parent).
			FindOne(ctx, bson.M{"_id": objID}).Err() != nil {
			println("ENTITY VALUE: ", ent)
			println("We got Parent: ", parent, " with ID:", t["parentId"].(string))
			return u.Message(false,
				"ParentID should correspond to Existing ID"), false

		}
		defer cancel()
	}
	return nil, true
}

func ValidatePatch(ent int, t map[string]interface{}) (map[string]interface{}, bool) {
	for k := range t {
		switch k {
		case "name", "category", "domain":
			//Only for Entities until GROUP
			//And OBJTMPL
			if ent < GROUP+1 || ent == OBJTMPL {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot be nullified!"), false
				}
			}

		case "parentId":
			if ent < ROOMTMPL && ent > TENANT {
				x, ok := validateParent(u.EntityToString(ent), ent, t)
				if !ok {
					return x, ok
				}
			}
			//STRAYDEV's schema is very loose
			//thus we can safely invoke validateEntity
			if ent == STRAYDEV {
				x, ok := ValidateEntity(ent, t)
				if !ok {
					return x, ok
				}
			}

		case "attributes.color": // TENANT
			if ent == TENANT {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes.orientation": //SITE, ROOM, RACK, DEVICE
			if ent >= SITE && ent <= DEVICE {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes.usableColor",
			"attributes.reservedColor",
			"attributes.technicalColor": //SITE
			if ent == SITE {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes.posXY", "attributes.posXYUnit": // BLDG, ROOM, RACK
			if ent >= BLDG && ent <= RACK {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes": //TENANT ... SENSOR, OBJTMPL
			if (ent >= TENANT && ent < ROOMTMPL) || ent == OBJTMPL {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes.size", "attributes.sizeUnit",
			"attributes.height", "attributes.heightUnit":
			//BLDG ... DEVICE
			if ent >= BLDG && ent <= DEVICE {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "attributes.floorUnit": //ROOM
			if ent == ROOM {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "slug", "colors": //TEMPLATES
			if ent == OBJTMPL || ent == ROOMTMPL {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "orientation", "sizeWDHm", "reservedArea",
			"technicalArea", "separators", "tiles": //ROOMTMPL
			if ent == ROOMTMPL {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

		case "description", "slots",
			"sizeWDHmm", "fbxModel": //OBJTMPL
			if ent == OBJTMPL {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}
			}

			/*case "type":
			if ent == SENSOR {
				if v, _ := t[k]; v == nil {
					return u.Message(false,
						"Field: "+k+" cannot nullified!"), false
				}

				if t[k] != "rack" &&
					t[k] != "device" && t[k] != "room" {
					return u.Message(false,
						"Incorrect values given for: "+k+"!"+
							"Please provide rack or device or room"), false
				}
			}*/

		}
	}
	return nil, true

}

func ValidateEntity(entity int, t map[string]interface{}) (map[string]interface{}, bool) {
	var objID primitive.ObjectID
	var err error
	//parentObj := nil
	/*
		TODO:
		Need to capture device if it is a parent
		and check that the device parent has a slot
		attribute
	*/
	switch entity {
	case TENANT, SITE, BLDG, ROOM, RACK, DEVICE, AC,
		PWRPNL, SEPARATOR, CABINET, ROW,
		TILE, CORIDOR, SENSOR:
		if t["name"] == nil || t["name"] == "" {
			return u.Message(false, "Name should be on payload"), false
		}

		/*if t["category"] == nil || t["category"] == "" {
			return u.Message(false, "Category should be on the payload"), false
		}*/

		if t["domain"] == nil || t["domain"] == "" {
			return u.Message(false, "Domain should be on the payload"), false
		}

		//Check if Parent ID is valid
		//Tenants do not have Parents
		if entity == DEVICE {
			objID, err = primitive.ObjectIDFromHex(t["parentId"].(string))
			if err != nil {
				return u.Message(false, "ParentID is not valid"), false
			}

			ctx, cancel := u.Connect()
			if GetDB().Collection("rack").FindOne(ctx,
				bson.M{"_id": objID}).Err() != nil &&
				GetDB().Collection("device").FindOne(ctx,
					bson.M{"_id": objID}).Err() != nil {
				return u.Message(false, "ParentID should be correspond to Existing ID"), false
			}
			defer cancel()

		} else if entity > TENANT && entity <= SENSOR {
			_, ok := t["parentId"].(string)
			if !ok {
				return u.Message(false, "ParentID is not valid"), false
			}
			objID, err = primitive.ObjectIDFromHex(t["parentId"].(string))
			if err != nil {
				return u.Message(false, "ParentID is not valid"), false
			}
			parentInt := u.GetParentOfEntityByInt(entity)
			if parentInt == -2 { //Sensor
				parentSet := []string{"room", "rack", "device"}
				found := false

				//First check if sensor type is present
				/*if t["type"] == nil || t["type"] == "" {
					return u.Message(false, "Sensor type must be on payload"), false
				}*/

				for i := range parentSet {
					ctx, cancel := u.Connect()
					if GetDB().Collection(parentSet[i]).
						FindOne(ctx, bson.M{"_id": objID}).Err() == nil {
						found = true
						//Ensure sensor type and parent entity
						//are consistent
						if t["type"] != parentSet[i] {
							return u.Message(false, "Sensor type must match parent entity"), false
						}

						i = len(parentSet)

					}
					defer cancel()
				}
				if found == false {
					return u.Message(false, "Sensor ParentID should correspond to Existing ID"), false
				}
			} else {
				parent := u.EntityToString(parentInt)

				ctx, cancel := u.Connect()
				if GetDB().Collection(parent).
					FindOne(ctx, bson.M{"_id": objID}).Err() != nil {
					println("ENTITY VALUE: ", entity)
					println("We got Parent: ", parent, " with ID:", t["parentId"].(string))
					return u.Message(false, "ParentID should correspond to Existing ID"), false

				}
				defer cancel()
			}

		}

		if entity < AC || entity == PWRPNL ||
			entity == SEPARATOR || entity == GROUP ||
			entity == ROOMTMPL || entity == OBJTMPL {
			if _, ok := t["attributes"]; !ok {
				return u.Message(false, "Attributes should be on the payload"), false
			} else {
				if v, ok := t["attributes"].(map[string]interface{}); !ok {
					return u.Message(false, "Attributes should be on the payload"), false
				} else {
					switch entity {
					case TENANT:
						if _, ok := v["color"]; !ok || v["color"] == "" {
							return u.Message(false,
								"Color Attribute must be specified on the payload"), false
						}

					case SITE:
						switch v["orientation"] {
						case "EN", "NW", "WS", "SE", "NE", "SW":
						case "", nil:
							return u.Message(false, "Orientation should be on the payload"), false

						default:
							return u.Message(false, "Orientation is invalid!"), false
						}

					case BLDG:
						if v["posXY"] == "" || v["posXY"] == nil {
							return u.Message(false, "XY coordinates should be on payload"), false
						}

						if v["posXYUnit"] == "" || v["posXYUnit"] == nil {
							return u.Message(false, "PositionXYUnit should be on the payload"), false
						}

						if v["size"] == "" || v["size"] == nil {
							return u.Message(false, "Invalid building size on the payload"), false
						}

						if v["sizeUnit"] == "" || v["sizeUnit"] == nil {
							return u.Message(false, "Building size unit should be on the payload"), false
						}

						if v["height"] == "" || v["height"] == nil {
							return u.Message(false, "Invalid Height on payload"), false
						}

						if v["heightUnit"] == "" || v["heightUnit"] == nil {
							return u.Message(false, "Building Height unit should be on the payload"), false
						}

					case ROOM:
						if v["posXY"] == "" || v["posXY"] == nil {
							return u.Message(false, "XY coordinates should be on payload"), false
						}

						if v["posXYUnit"] == "" || v["posXYUnit"] == nil {
							return u.Message(false, "PositionXYUnit should be on the payload"), false
						}

						switch v["floorUnit"] {
						case "f", "m", "t":
						case "", nil:
							return u.Message(false, "floorUnit should be on the payload"), false
						default:
							return u.Message(false, "floorUnit is invalid!"), false

						}

						switch v["orientation"] {
						case "-E-N", "-E+N", "+E-N", "+E+N", "+N+E":
						case "-N-W", "-N+W", "+N-W", "+N+W":
						case "-W-S", "-W+S", "+W-S", "+W+S":
						case "-S-E", "-S+E", "+S-E", "+S+E":
						case "", nil:
							return u.Message(false, "Orientation should be on the payload"), false

						default:
							return u.Message(false, "Orientation is invalid!"), false
						}

						if v["size"] == "" || v["size"] == nil {
							return u.Message(false, "Invalid size on the payload"), false
						}

						if v["sizeUnit"] == "" || v["sizeUnit"] == nil {
							return u.Message(false, "Room size unit should be on the payload"), false
						}

						if v["height"] == "" || v["height"] == nil {
							return u.Message(false, "Invalid Height on payload"), false
						}

						if v["heightUnit"] == "" || v["heightUnit"] == nil {
							return u.Message(false, "Room Height unit should be on the payload"), false
						}

					case RACK:
						if v["posXY"] == "" || v["posXY"] == nil {
							return u.Message(false, "XY coordinates should be on payload"), false
						}

						if v["posXYUnit"] == "" || v["posXYUnit"] == nil {
							return u.Message(false, "PositionXYUnit should be on the payload"), false
						}

						switch v["orientation"] {
						case "front", "rear", "left", "right":
						case "", nil:
							return u.Message(false, "Orientation should be on the payload"), false

						default:
							return u.Message(false, "Orientation is invalid!"), false
						}

						if v["size"] == "" || v["size"] == nil {
							return u.Message(false, "Invalid size on the payload"), false
						}

						if v["sizeUnit"] == "" || v["sizeUnit"] == nil {
							return u.Message(false, "Rack size unit should be on the payload"), false
						}

						if v["height"] == "" || v["height"] == nil {
							return u.Message(false, "Invalid Height on payload"), false
						}

						if v["heightUnit"] == "" || v["heightUnit"] == nil {
							return u.Message(false, "Rack Height unit should be on the payload"), false
						}

					case DEVICE:
						switch v["orientation"] {
						case "front", "rear", "frontflipped", "rearflipped":
						case "", nil:
							return u.Message(false, "Orientation should be on the payload"), false

						default:
							return u.Message(false, "Orientation is invalid!"), false
						}

						if v["size"] == "" || v["size"] == nil {
							return u.Message(false, "Invalid size on the payload"), false
						}

						if v["sizeUnit"] == "" || v["sizeUnit"] == nil {
							return u.Message(false, "Device size unit should be on the payload"), false
						}

						if v["height"] == "" || v["height"] == nil {
							return u.Message(false, "Invalid Height on payload"), false
						}

						if v["heightUnit"] == "" || v["heightUnit"] == nil {
							return u.Message(false, "Device Height unit should be on the payload"), false
						}

						if side, ok := v["side"]; ok {
							switch side {
							case "front", "rear", "frontflipped", "rearflipped":
							default:
								msg := "The 'Side' value (if given) must be one of" +
									"the given values: front, rear, frontflipped, rearflipped"
								return u.Message(false, msg), false
							}
						}

					}
				}
			}
		}
	case ROOMTMPL, OBJTMPL:
		if t["slug"] == "" || t["slug"] == nil {
			return u.Message(false, "Slug should be on payload"), false
		}

		if _, ok := t["colors"]; !ok {
			return u.Message(false,
				"Colors should be on payload"), false
		}

		if entity == OBJTMPL {
			if _, ok := t["description"]; !ok {
				return u.Message(false,
					"Description should be on payload"), false
			}

			/*if _, ok := t["category"]; !ok {
				return u.Message(false,
					"Category should be on payload"), false
			}*/

			if _, ok := t["sizeWDHmm"]; !ok {
				return u.Message(false,
					"Size,Width,Depth (mm) Array should be on payload"), false
			}

			if _, ok := t["fbxModel"]; !ok {
				return u.Message(false,
					"fbxModel should be on payload"), false
			}

			if _, ok := t["attributes"]; !ok {
				return u.Message(false,
					"Attributes should be on payload"), false
			}

			if _, ok := t["slots"]; !ok {
				return u.Message(false,
					"fbxModel should be on payload"), false
			}

		} else { //ROOMTMPL
			if _, ok := t["orientation"]; !ok {
				return u.Message(false,
					"Orientation should be on payload"), false
			}

			if _, ok := t["sizeWDHm"]; !ok {
				return u.Message(false,
					"Size,Width,Depth Array should be on payload"), false
			}

			if _, ok := t["technicalArea"]; !ok {
				return u.Message(false,
					"TechnicalArea should be on payload"), false
			}

			if _, ok := t["reservedArea"]; !ok {
				return u.Message(false,
					"ReservedArea should be on payload"), false
			}

			if _, ok := t["separators"]; !ok {
				return u.Message(false,
					"Separators should be on payload"), false
			}

			if _, ok := t["tiles"]; !ok {
				return u.Message(false,
					"Tiles should be on payload"), false
			}
		}
	case GROUP:
		if t["name"] == "" || t["name"] == nil {
			return u.Message(false, "Name should be on payload"), false
		}

		if r, ok := validateParent("group", entity, t); !ok {
			return r, ok
		}

	case STRAYDEV, STRAYSENSOR:
		//Check for parent if PID provided
		if t["parentId"] != nil && t["parentId"] != "" {
			if pid, ok := t["parentId"].(string); ok {
				ID, _ := primitive.ObjectIDFromHex(pid)

				ctx, cancel := u.Connect()
				if GetDB().Collection("stray_device").FindOne(ctx,
					bson.M{"_id": ID}).Err() != nil {
					return u.Message(false,
						"ParentID should be an Existing ID or null"), false
				}
				defer cancel()
			} else {
				return u.Message(false,
					"ParentID should be an Existing ID or null"), false
			}

		}

		if t["name"] == nil || t["name"] == "" {
			return u.Message(false, "Please provide a valid name"), false
		}

		//Need to check for uniqueness before inserting
		//this is helpful for the validation endpoints
		ctx, cancel := u.Connect()
		entStr := u.EntityToString(entity)

		if c, _ := GetDB().Collection(entStr).CountDocuments(ctx,
			bson.M{"name": t["name"]}); c != 0 {
			msg := "Error a " + entStr + " with the name provided already exists." +
				"Please provide a unique name"
			return u.Message(false, msg), false
		}
		defer cancel()

	}

	//Successfully validated the Object
	return u.Message(true, "success"), true
}

func CreateEntity(entity int, t map[string]interface{}) (map[string]interface{}, string) {
	message := ""
	if resp, ok := ValidateEntity(entity, t); !ok {
		return resp, "validate"
	}

	//Set timestamp
	t["createdDate"] = primitive.NewDateTimeFromTime(time.Now())
	t["lastUpdated"] = t["createdDate"]

	ctx, cancel := u.Connect()
	entStr := u.EntityToString(entity)
	res, e := GetDB().Collection(entStr).InsertOne(ctx, t)
	if e != nil {
		return u.Message(false,
				"Internal error while creating "+entStr+": "+e.Error()),
			e.Error()
	}
	defer cancel()

	//Remove _id
	t["id"] = res.InsertedID
	//t = fixID(t)

	switch entity {
	case ROOMTMPL:
		message = "successfully created room_template"
	case OBJTMPL:
		message = "successfully created obj_template"
	default:
		message = "successfully created object"
	}

	resp := u.Message(true, message)
	resp["data"] = t
	return resp, ""
}

func GetEntity(req bson.M, ent string) (map[string]interface{}, string) {
	t := map[string]interface{}{}

	ctx, cancel := u.Connect()
	e := GetDB().Collection(ent).FindOne(ctx, req).Decode(&t)
	if e != nil {
		return nil, e.Error()
	}
	defer cancel()
	//Remove _id
	t = fixID(t)

	//If entity has '_' remove it
	if strings.Contains(ent, "_") {
		FixUnderScore(t)
	}
	return t, ""
}

func GetManyEntities(ent string, req bson.M, opts *options.FindOptions) ([]map[string]interface{}, string) {
	data := make([]map[string]interface{}, 0)
	ctx, cancel := u.Connect()
	c, err := GetDB().Collection(ent).Find(ctx, req, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err.Error()
	}
	defer cancel()

	data, e1 := ExtractCursor(c, ctx)
	if e1 != "" {
		fmt.Println(e1)
		return nil, e1
	}

	//Remove underscore If the entity has '_'
	if strings.Contains(ent, "_") == true {
		for i := range data {
			FixUnderScore(data[i])
		}
	}

	return data, ""
}

func DeleteEntityManual(entity string, req bson.M) (map[string]interface{}, string) {
	//Finally delete the Entity
	ctx, cancel := u.Connect()
	c, _ := GetDB().Collection(entity).DeleteOne(ctx, req)
	if c.DeletedCount == 0 {
		return u.Message(false, "There was an error in deleting the entity"), "not found"
	}
	defer cancel()

	return u.Message(true, "success"), ""
}

func DeleteEntity(entity string, id primitive.ObjectID) (map[string]interface{}, string) {
	var t map[string]interface{}
	var e string
	eNum := u.EntityStrToInt(entity)
	if eNum > DEVICE {
		//Delete the non hierarchal objects
		t, e = GetEntityHierarchy(id, entity, eNum, eNum+eNum)
	} else {
		t, e = GetEntityHierarchy(id, entity, eNum, AC)
	}

	if e != "" {
		return u.Message(false,
			"There was an error in deleting the entity: "+e), "not found"
	}

	return deleteHelper(t, eNum)
}

func deleteHelper(t map[string]interface{}, ent int) (map[string]interface{}, string) {
	if t != nil {

		if v, ok := t["children"]; ok {
			if x, ok := v.([]map[string]interface{}); ok {
				for i := range x {
					if ent == STRAYDEV {
						deleteHelper(x[i], ent)
					} else {
						deleteHelper(x[i], ent+1)
					}

				}
			}
		}

		println("So we got: ", ent)

		if ent == RACK {
			ctx, cancel := u.Connect()
			GetDB().Collection("sensor").DeleteMany(ctx,
				bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})

			GetDB().Collection("group").DeleteMany(ctx,
				bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})
			defer cancel()
		}

		//Delete associated non hierarchal objs
		if ent == ROOM {
			//ITER Through all nonhierarchal objs
			ctx, cancel := u.Connect()
			for i := AC; i < GROUP+1; i++ {
				ent := u.EntityToString(i)
				GetDB().Collection(ent).DeleteMany(ctx, bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})
			}
			defer cancel()
		}

		//Delete hierarchy under stray-device
		if ent == STRAYDEV {
			ctx, cancel := u.Connect()
			entity := u.EntityToString(u.STRAYSENSOR)
			GetDB().Collection(entity).DeleteMany(ctx, bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})
			defer cancel()
		}

		if ent == DEVICE {
			DeleteDeviceF(t["id"].(primitive.ObjectID))
		} else {
			ctx, cancel := u.Connect()
			entity := u.EntityToString(ent)
			println(entity)
			c, _ := GetDB().Collection(entity).DeleteOne(ctx, bson.M{"_id": t["id"].(primitive.ObjectID)})
			if c.DeletedCount == 0 {
				return u.Message(false, "There was an error in deleting the entity"), "not found"
			}
			defer cancel()

		}
	}
	return nil, ""
}

func UpdateEntity(ent string, req bson.M, t *map[string]interface{}, isPatch bool) (map[string]interface{}, string) {
	var e *mongo.SingleResult
	updatedDoc := bson.M{}
	retDoc := options.ReturnDocument(options.After)

	//Update timestamp requires first obj retrieval
	//there isn't any way for mongoDB to make a field
	//immutable in a document
	oldObj, e1 := GetEntity(req, ent)
	if e1 != "" {
		return u.Message(false, "Error: "+e1), e1
	}
	(*t)["lastUpdated"] = primitive.NewDateTimeFromTime(time.Now())
	(*t)["createdDate"] = oldObj["createdDate"]

	ctx, cancel := u.Connect()
	if isPatch == true {

		msg, ok := ValidatePatch(u.EntityStrToInt(ent), *t)
		if !ok {
			return msg, "invalid"
		}
		e = GetDB().Collection(ent).FindOneAndUpdate(ctx,
			req, bson.M{"$set": *t},
			&options.FindOneAndUpdateOptions{ReturnDocument: &retDoc})
		if e.Err() != nil {
			return u.Message(false, "failure: "+e.Err().Error()), e.Err().Error()
		}
	} else {

		//Ensure that the update will be valid
		msg, ok := ValidateEntity(u.EntityStrToInt(ent), *t)
		if !ok {
			return msg, "invalid"
		}

		e = GetDB().Collection(ent).FindOneAndReplace(ctx,
			req, *t,
			&options.FindOneAndReplaceOptions{ReturnDocument: &retDoc})

		if e.Err() != nil {
			return u.Message(false, "failure: "+e.Err().Error()), e.Err().Error()
		}
	}

	//Obtain new document then
	//Fix the _id / id discrepancy
	e.Decode(&updatedDoc)
	updatedDoc = fixID(updatedDoc)

	//Response Message
	message := ""
	switch u.EntityStrToInt(ent) {
	case ROOMTMPL:
		message = "successfully updated room_template"
	case OBJTMPL:
		message = "successfully updated obj_template"
	default:
		message = "successfully updated object"
	}

	defer cancel()
	resp := u.Message(true, message)
	resp["data"] = updatedDoc
	return resp, ""
}

func GetEntityHierarchy(ID primitive.ObjectID, ent string, start, end int) (map[string]interface{}, string) {
	var childEnt string
	if start < end {
		if ent == "stray_sensor" {
			println("DEBUG SENSOR IN H CALL")
		}
		top, e := GetEntity(bson.M{"_id": ID}, ent)
		if top == nil {
			return nil, e
		}
		top = fixID(top)

		children := []map[string]interface{}{}
		pid := ID.Hex()
		//Get sensors & groups
		if ent == "rack" || ent == "device" {
			//Get sensors
			sensors, _ := GetManyEntities("sensor", bson.M{"parentId": pid}, nil)

			//Get groups
			groups, _ := GetManyEntities("group", bson.M{"parentId": pid}, nil)

			if sensors != nil {
				children = append(children, sensors...)
			}
			if groups != nil {
				children = append(children, groups...)
			}
		}

		if ent == "device" || ent == "stray_device" {
			childEnt = ent
		} else {
			childEnt = u.EntityToString(start + 1)
		}

		subEnts, _ := GetManyEntities(childEnt, bson.M{"parentId": pid}, nil)

		for idx := range subEnts {
			tmp, _ := GetEntityHierarchy(subEnts[idx]["id"].(primitive.ObjectID), childEnt, start+1, end)
			if tmp != nil {
				subEnts[idx] = tmp
			}
		}

		if subEnts != nil {
			children = append(children, subEnts...)
		}

		if ent == "room" {
			for i := AC; i < CABINET+1; i++ {
				roomEnts, _ := GetManyEntities(u.EntityToString(i), bson.M{"parentId": pid}, nil)
				if roomEnts != nil {
					children = append(children, roomEnts...)
				}
			}
			for i := PWRPNL; i < TILE+1; i++ {
				roomEnts, _ := GetManyEntities(u.EntityToString(i), bson.M{"parentId": pid}, nil)
				if roomEnts != nil {
					children = append(children, roomEnts...)
				}
			}
			roomEnts, _ := GetManyEntities(u.EntityToString(CORIDOR), bson.M{"parentId": pid}, nil)
			if roomEnts != nil {
				children = append(children, roomEnts...)
			}
			roomEnts, _ = GetManyEntities(u.EntityToString(GROUP), bson.M{"parentId": pid}, nil)
			if roomEnts != nil {
				children = append(children, roomEnts...)
			}
		}

		if ent == "stray_device" {
			sSensors, _ := GetManyEntities("stray_sensor", bson.M{"parentId": pid}, nil)
			if sSensors != nil {
				children = append(children, sSensors...)
			}
		}

		if children != nil && len(children) > 0 {
			top["children"] = children
		}

		return top, ""
	}
	return nil, ""
}

func GetEntitiesUsingAncestorNames(ent string, id primitive.ObjectID, ancestry []map[string]string) ([]map[string]interface{}, string) {
	top, e := GetEntity(bson.M{"_id": id}, ent)
	if e != "" {
		return nil, e
	}

	//Remove _id
	top = fixID(top)

	pid := (top["id"].(primitive.ObjectID)).Hex()

	var x map[string]interface{}
	var e1 string
	for i := range ancestry {
		for k, v := range ancestry[i] {

			println("KEY:", k, " VAL:", v)

			if v == "all" {
				println("K:", k)
				//println("ID", x["_id"].(primitive.ObjectID).String())
				/*if k == "device" {
					return GetDeviceFByParentID(pid) nil, ""
				}*/
				return GetManyEntities(k, bson.M{"parentId": pid}, nil)
			}

			x, e1 = GetEntity(bson.M{"parentId": pid, "name": v}, k)
			if e1 != "" {
				println("Failing here")
				return nil, ""
			}
			pid = (x["id"].(primitive.ObjectID)).Hex()
		}
	}

	return nil, ""
}

func GetEntityUsingAncestorNames(ent string, id primitive.ObjectID, ancestry []map[string]string) (map[string]interface{}, string) {
	top, e := GetEntity(bson.M{"_id": id}, ent)
	if e != "" {
		return nil, e
	}

	//Remove _id
	top = fixID(top)

	pid := (top["id"].(primitive.ObjectID)).Hex()

	var x map[string]interface{}
	var e1 string
	for i := range ancestry {
		for k, v := range ancestry[i] {

			println("KEY:", k, " VAL:", v)

			x, e1 = GetEntity(bson.M{"parentId": pid, "name": v}, k)
			if e1 != "" {
				println("Failing here")
				return nil, ""
			}
			pid = (x["id"].(primitive.ObjectID)).Hex()
		}
	}

	//Remove _id
	x = fixID(x)
	return x, ""
}

func GetHierarchyByName(entity, name string, entnum, end int) (map[string]interface{}, string) {

	t, e := GetEntity(bson.M{"name": name}, entity)
	if e != "" {
		fmt.Println(e)
		return nil, e
	}

	//Remove _id
	t = fixID(t)

	var subEnt string
	if entnum == STRAYDEV {
		subEnt = "stray_device"
	} else {
		subEnt = u.EntityToString(entnum + 1)
	}

	tid := t["id"].(primitive.ObjectID).Hex()

	//Get immediate children
	children, e1 := GetManyEntities(subEnt, bson.M{"parentId": tid}, nil)
	if e1 != "" {
		println("Are we here")
		println("SUBENT: ", subEnt)
		println("PID: ", tid)
		return nil, e1
	}
	t["children"] = children

	//Get the rest of hierarchy for children
	for i := range children {
		var x map[string]interface{}
		var subIdx string
		if subEnt == "stray_device" { //only set entnum+1 for tenants
			subIdx = subEnt
		} else {
			subIdx = u.EntityToString(entnum + 1)
		}
		subID := (children[i]["id"].(primitive.ObjectID))
		x, _ =
			GetEntityHierarchy(subID, subIdx, entnum+1, end)
		if x != nil {
			children[i] = x
		}
	}

	return t, ""

}

func GetEntitiesUsingTenantAsAncestor(ent, id string, ancestry []map[string]string) ([]map[string]interface{}, string) {
	top, e := GetEntity(bson.M{"name": id}, ent)
	if e != "" {
		return nil, e
	}

	//Remove _id
	top = fixID(top)

	pid := (top["id"].(primitive.ObjectID)).Hex()

	var x map[string]interface{}
	var e1 string
	println("ANCS-LEN: ", len(ancestry))
	for i := range ancestry {
		for k, v := range ancestry[i] {

			println("KEY:", k, " VAL:", v)

			if v == "all" {
				println("K:", k)
				return GetManyEntities(k, bson.M{"parentId": pid}, nil)
			}

			x, e1 = GetEntity(bson.M{"parentId": pid, "name": v}, k)
			if e1 != "" {
				println("Failing here")
				println("E1: ", e1)
				return nil, ""
			}
			pid = (x["id"].(primitive.ObjectID)).Hex()
		}
	}

	return nil, ""
}

func GetEntityUsingTenantAsAncestor(ent, id string, ancestry []map[string]string) (map[string]interface{}, string) {
	top, e := GetEntity(bson.M{"name": id}, ent)
	if e != "" {
		return nil, e
	}

	pid := (top["id"].(primitive.ObjectID)).Hex()

	var x map[string]interface{}
	var e1 string
	for i := range ancestry {
		for k, v := range ancestry[i] {

			println("KEY:", k, " VAL:", v)

			x, e1 = GetEntity(bson.M{"parentId": pid, "name": v}, k)
			if e1 != "" {
				println("Failing here")
				return nil, ""
			}
			pid = (x["id"].(primitive.ObjectID)).Hex()
		}
	}

	return x, ""
}

func GetEntitiesOfAncestor(id interface{}, ent int, entStr, wantedEnt string) ([]map[string]interface{}, string) {
	var ans []map[string]interface{}
	var t map[string]interface{}
	var e, e1 string
	if ent == TENANT {

		t, e = GetEntity(bson.M{"name": id}, "tenant")
		if e != "" {
			return nil, e
		}

	} else {
		ID, _ := primitive.ObjectIDFromHex(id.(string))
		t, e = GetEntity(bson.M{"_id": ID}, entStr)
		if e != "" {
			return nil, e
		}
	}

	sub, e1 := GetManyEntities(u.EntityToString(ent+1),
		bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()}, nil)
	if e1 != "" {
		return nil, e1
	}

	if wantedEnt == "" {
		wantedEnt = u.EntityToString(ent + 2)
	}

	for i := range sub {
		x, _ := GetManyEntities(wantedEnt,
			bson.M{"parentId": sub[i]["id"].(primitive.ObjectID).Hex()}, nil)
		ans = append(ans, x...)
	}
	return ans, ""
}

//DEV FAMILY FUNCS

func DeleteDeviceF(entityID primitive.ObjectID) (map[string]interface{}, string) {
	//var deviceType string

	t, e := GetEntityHierarchy(entityID, "device", 0, 999)
	if e != "" {
		return u.Message(false,
			"There was an error in deleting the entity"), "not found"
	}

	return deleteDeviceHelper(t)
}

func deleteDeviceHelper(t map[string]interface{}) (map[string]interface{}, string) {
	println("entered ddH")
	if t != nil {

		if v, ok := t["children"]; ok {
			if x, ok := v.([]map[string]interface{}); ok {
				for i := range x {
					deleteDeviceHelper(x[i])
				}
			}
		}

		ctx, cancel := u.Connect()
		//Delete relevant non hierarchal objects
		GetDB().Collection("sensor").DeleteMany(ctx,
			bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})

		GetDB().Collection("group").DeleteMany(ctx,
			bson.M{"parentId": t["id"].(primitive.ObjectID).Hex()})

		c, _ := GetDB().Collection("device").DeleteOne(ctx, bson.M{"_id": t["id"].(primitive.ObjectID)})
		if c.DeletedCount == 0 {
			return u.Message(false, "There was an error in deleting the entity"), "not found"
		}
		defer cancel()

	}
	return nil, ""
}

func ExtractCursor(c *mongo.Cursor, ctx context.Context) ([]map[string]interface{}, string) {
	ans := []map[string]interface{}{}
	for c.Next(ctx) {
		x := map[string]interface{}{}
		err := c.Decode(x)
		if err != nil {
			fmt.Println(err)
			return nil, err.Error()
		}
		//Remove _id
		x = fixID(x)
		ans = append(ans, x)
	}
	return ans, ""
}

//Removes underscore in object category if present
func FixUnderScore(x map[string]interface{}) {
	if catInf, ok := x["category"]; ok {
		if cat, _ := catInf.(string); strings.Contains(cat, "_") == true {
			x["category"] = strings.Replace(cat, "_", "-", 1)
		}
	}
}
