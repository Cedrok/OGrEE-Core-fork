package models

import (
	"fmt"
	u "p3/utils"
	"strings"
)

type Rack_Attributes struct {
	ID          int    `json:"id" gorm:"column:id"`
	PosXY       string `json:"posXY" gorm:"column:rack_pos_x_y"`
	PosXYU      string `json:"posXYUnit" gorm:"column:rack_pos_x_y_unit"`
	PosZ        string `json:"posZ" gorm:"column:rack_pos_z"`
	PosZU       string `json:"posZUnit" gorm:"column:rack_pos_z_unit"`
	Template    string `json:"template" gorm:"column:rack_template"`
	Orientation string `json:"orientation" gorm:"column:rack_orientation"`
	Size        string `json:"size" gorm:"column:rack_size"`
	SizeU       string `json:"sizeUnit" gorm:"column:rack_size_unit"`
	Height      string `json:"height" gorm:"column:rack_height"`
	HeightU     string `json:"heightUnit" gorm:"column:rack_height_unit"`
	Vendor      string `json:"vendor" gorm:"column:rack_vendor"`
	Type        string `json:"type" gorm:"column:rack_type"`
	Model       string `json:"model" gorm:"column:rack_model"`
	Serial      string `json:"serial" gorm:"column:rack_serial"`
}

type Rack struct {
	//gorm.Model
	ID              int             `json:"id" gorm:"column:id"`
	Name            string          `json:"name" gorm:"column:rack_name"`
	ParentID        string          `json:"parentId" gorm:"column:rack_parent_id"`
	Category        string          `json:"category" gorm:"-"`
	Domain          string          `json:"domain" gorm:"column:rack_domain"`
	DescriptionJSON []string        `json:"description" gorm:"-"`
	DescriptionDB   string          `json:"-" gorm:"column:rack_description"`
	Attributes      Rack_Attributes `json:"attributes"`

	//Site []Site
	//D is used to help the JSON marshalling
	//while Description will be used in
	//DB transactions
}

func (rack *Rack) Validate() (map[string]interface{}, bool) {
	if rack.Name == "" {
		return u.Message(false, "Rack Name should be on payload"), false
	}

	/*if rack.Category == "" {
		return u.Message(false, "Category should be on the payload"), false
	}

	if rack.Desc == "" {
		return u.Message(false, "Description should be on the payload"), false
	}*/

	if rack.Domain == "" {
		return u.Message(false, "Domain should should be on the payload"), false
	}

	if GetDB().Table("room").
		Where("id = ?", rack.ParentID).First(&Room{}).Error != nil {

		return u.Message(false, "ParentID should be correspond to Room ID"), false
	}

	if rack.Attributes.PosXY == "" {
		return u.Message(false, "XY coordinates should be on payload"), false
	}

	if rack.Attributes.PosXYU == "" {
		return u.Message(false, "PositionXYU string should be on the payload"), false
	}

	/*if rack.Attributes.PosZ == "" {
		return u.Message(false, "Z coordinates should be on payload"), false
	}

	if rack.Attributes.PosZU == "" {
		return u.Message(false, "PositionZU string should be on the payload"), false
	}*/

	/*if rack.Attributes.Template == "" {
		return u.Message(false, "Template should be on the payload"), false
	}*/

	switch rack.Attributes.Orientation {
	case "front", "rear", "left", "right":
	case "":
		return u.Message(false, "Orientation should be on the payload"), false

	default:
		return u.Message(false, "Orientation is invalid!"), false
	}

	if rack.Attributes.Size == "" {
		return u.Message(false, "Invalid size on the payload"), false
	}

	if rack.Attributes.SizeU == "" {
		return u.Message(false, "Rack size string should be on the payload"), false
	}

	if rack.Attributes.Height == "" {
		return u.Message(false, "Invalid Height on payload"), false
	}

	if rack.Attributes.HeightU == "" {
		return u.Message(false, "Rack Height string should be on the payload"), false
	}

	//Successfully validated Rack
	return u.Message(true, "success"), true
}

func (rack *Rack) Create() (map[string]interface{}, string) {
	if resp, ok := rack.Validate(); !ok {
		return resp, "validate"
	}

	rack.DescriptionDB = strings.Join(rack.DescriptionJSON, "XYZ")

	if e := GetDB().Create(rack).Error; e != nil {
		return u.Message(false, "Internal Error while creating Rack: "+e.Error()),
			"internal"
	}
	rack.Attributes.ID = rack.ID
	if e := GetDB().Create(&(rack.Attributes)).Error; e != nil {
		return u.Message(false, "Internal Error while creating Rack Attrs: "+e.Error()),
			"internal"
	}

	resp := u.Message(true, "success")
	resp["rack"] = rack
	return resp, ""
}

//Get the rack using ID
func GetRack(id uint) (*Rack, string) {
	rack := &Rack{}
	err := GetDB().Table("rack").Where("id = ?", id).First(rack).
		Table("rack_attributes").Where("id = ?", id).First(&(rack.Attributes)).Error
	if err != nil {
		fmt.Println(err)
		return nil, err.Error()
	}

	rack.DescriptionJSON = strings.Split(rack.DescriptionDB, "XYZ")
	return rack, ""
}

//Obtain all racks of a room
func GetRacks(room *Room) ([]*Rack, string) {
	racks := make([]*Rack, 0)

	err := GetDB().Table("racks").Where("foreignkey = ?", room.ID).Find(&racks).Error
	if err != nil {
		fmt.Println(err)
		return nil, err.Error()
	}

	return racks, ""
}

//Obtain all racks
func GetAllRacks() ([]*Rack, string) {
	racks := make([]*Rack, 0)
	attrs := make([]*Rack_Attributes, 0)
	err := GetDB().Find(&racks).Find(&attrs).Error
	if err != nil {
		fmt.Println(err)
		return nil, err.Error()
	}

	for i := range racks {
		racks[i].Attributes = *(attrs[i])
		racks[i].DescriptionJSON = strings.Split(racks[i].DescriptionDB, "XYZ")
	}

	return racks, ""
}

//More methods should be made to
//Meet CRUD capabilities
//Need Update and Delete
//These would be a bit more complicated
//So leave them out for now

func UpdateRack(id uint, newRackInfo *Rack) (map[string]interface{}, string) {
	rack := &Rack{}

	err := GetDB().Table("rack").Where("id = ?", id).First(rack).
		Table("rack_attributes").Where("id = ?", id).First(&(rack.Attributes)).Error
	if err != nil {
		return u.Message(false, "Error while checking Rack: "+err), err.Error()
	}

	if newRackInfo.Name != "" && newRackInfo.Name != rack.Name {
		rack.Name = newRackInfo.Name
	}

	if newRackInfo.Domain != "" && newRackInfo.Domain != rack.Domain {
		rack.Domain = newRackInfo.Domain
	}

	if dc := strings.Join(newRackInfo.DescriptionJSON, "XYZ"); dc != "" && strings.Compare(dc, rack.DescriptionDB) != 0 {
		rack.DescriptionDB = dc
	}

	if newRackInfo.Attributes.PosXY != "" && newRackInfo.Attributes.PosXY != rack.Attributes.PosXY {
		rack.Attributes.PosXY = newRackInfo.Attributes.PosXY
	}

	if newRackInfo.Attributes.PosXYU != "" && newRackInfo.Attributes.PosXYU != rack.Attributes.PosXYU {
		rack.Attributes.PosXYU = newRackInfo.Attributes.PosXYU
	}

	if newRackInfo.Attributes.PosZ != "" && newRackInfo.Attributes.PosZ != rack.Attributes.PosZ {
		rack.Attributes.PosZ = newRackInfo.Attributes.PosZ
	}

	if newRackInfo.Attributes.PosZU != "" && newRackInfo.Attributes.PosZU != rack.Attributes.PosZU {
		rack.Attributes.PosZU = newRackInfo.Attributes.PosZU
	}

	if newRackInfo.Attributes.Template != "" && newRackInfo.Attributes.Template != rack.Attributes.Template {
		rack.Attributes.Template = newRackInfo.Attributes.Template
	}

	if newRackInfo.Attributes.Orientation != "" {
		switch newRackInfo.Attributes.Orientation {
		case "NE", "NW", "SE", "SW":
			rack.Attributes.Orientation = newRackInfo.Attributes.Orientation

		default:
		}
	}

	if newRackInfo.Attributes.Size != "" && newRackInfo.Attributes.Size != rack.Attributes.Size {
		rack.Attributes.Size = newRackInfo.Attributes.Size
	}

	if newRackInfo.Attributes.SizeU != "" && newRackInfo.Attributes.SizeU != rack.Attributes.SizeU {
		rack.Attributes.SizeU = newRackInfo.Attributes.SizeU
	}

	if newRackInfo.Attributes.Height != "" && newRackInfo.Attributes.Height != rack.Attributes.Height {
		rack.Attributes.Height = newRackInfo.Attributes.Height
	}

	if newRackInfo.Attributes.HeightU != "" && newRackInfo.Attributes.HeightU != rack.Attributes.HeightU {
		rack.Attributes.HeightU = newRackInfo.Attributes.HeightU
	}

	if newRackInfo.Attributes.Vendor != "" && newRackInfo.Attributes.Vendor != rack.Attributes.Vendor {
		rack.Attributes.Vendor = newRackInfo.Attributes.Vendor
	}

	if newRackInfo.Attributes.Type != "" && newRackInfo.Attributes.Type != rack.Attributes.Type {
		rack.Attributes.Type = newRackInfo.Attributes.Type
	}

	if newRackInfo.Attributes.Model != "" && newRackInfo.Attributes.Model != rack.Attributes.Model {
		rack.Attributes.Model = newRackInfo.Attributes.Model
	}

	if newRackInfo.Attributes.Serial != "" && newRackInfo.Attributes.Serial != rack.Attributes.Serial {
		rack.Attributes.Serial = newRackInfo.Attributes.Serial
	}

	//Successfully validated the new data
	if e1 := GetDB().Table("rack").Save(rack).
		Table("rack_attributes").Save(&(rack.Attributes)).Error; e1 != nil {
		return u.Message(false, "Error while updating rack: "+e1), e1.Error()
	}
	return u.Message(true, "success"), ""
}

func DeleteRack(id uint) map[string]interface{} {

	//This is a hard delete!
	e := GetDB().Unscoped().Table("rack").Delete(&Rack{}, id).RowsAffected

	//The command below is a soft delete
	//Meaning that the 'deleted_at' field will be set
	//the record will remain but unsearchable
	//e := GetDB().Table("tenants").Delete(Tenant{}, id).Error
	if e == 0 {
		return u.Message(false, "There was an error in deleting the rack")
	}

	return u.Message(true, "success")
}
