package models

import (
	"fmt"
	u "p3/utils"
)

type ECardinalOrient string

//Desc        string          `json:"description"`

type Site_Attributes struct {
	ID             int    `json:"id" gorm:"column:id"`
	Orientation    string `json:"orientation" gorm:"column:site_orientation"`
	UsableColor    string `json:"usableColor" gorm:"column:usable_color"`
	ReservedColor  string `json:"reservedColor" gorm:"column:reserved_color"`
	TechnicalColor string `json:"technicalColor" gorm:"column:technical_color"`
	Address        string `json:"address" gorm:"column:address"`
	Zipcode        string `json:"zipcode" gorm:"column:zipcode"`
	City           string `json:"city" gorm:"column:city"`
	Country        string `json:"country" gorm:"column:country"`
	Gps            string `json:"gps" gorm:"column:gps"`
}

type Site struct {
	//gorm.Model
	ID         int             `json:"id" gorm:"column:id"`
	Name       string          `json:"name" gorm:"column:site_name"`
	Category   string          `json:"category" gorm:"-"`
	Domain     string          `json:"domain" gorm:"column:site_domain"`
	ParentID   string          `json:"parentId" gorm:"column:site_parent_id"`
	Attributes Site_Attributes `json:"attributes"`

	Building []Building
}

func (site *Site) Validate() (map[string]interface{}, bool) {
	if site.Name == "" {
		return u.Message(false, "site Name should be on payload"), false
	}

	if site.Category == "" {
		return u.Message(false, "Category should be on the payload"), false
	}

	/*if site.Desc == "" {
		return u.Message(false, "Description should be on the payload"), false
	}*/

	if site.Domain == "" {
		return u.Message(false, "Domain should be on the payload"), false
	}

	if GetDB().Table("tenant").
		Where("id = ?", site.ParentID).First(&Tenant{}).Error != nil {

		return u.Message(false, "SiteParentID should be correspond to tenant ID"), false
	}

	/*if site.Color == "" {
		return u.Message(false, "Color should be on the payload"), false
	}*/

	switch site.Attributes.Orientation {
	case "NE", "NW", "SE", "SW":
	case "":
		return u.Message(false, "Orientation should be on the payload"), false

	default:
		return u.Message(false, "Orientation is invalid!"), false
	}

	//Successfully validated Site
	return u.Message(true, "success"), true
}

func (site *Site) Create() map[string]interface{} {
	if resp, ok := site.Validate(); !ok {
		return resp
	}

	//GetDB().Create(site)

	GetDB().Create(site)

	site.Attributes.ID = site.ID

	GetDB().Table("site_attributes").Create(&(site.Attributes))
	resp := u.Message(true, "success")
	resp["site"] = site
	return resp
}

//Would have to think about
//these functions more
//since I set it up
//to just obtain the first site
//The GORM command might be
//wrong too
func GetSites(id uint) []*Site {
	site := make([]*Site, 0)

	err := GetDB().Table("tenants").Where("id = ?", id).First(&Tenant{}).Error
	if err != nil {
		fmt.Println("yo the tenant wasnt found here")
		return nil
	}

	e := GetDB().Table("sites").Where("domain = ?", id).Find(&site).Error
	if e != nil {
		fmt.Println("yo the there isnt any site matching the foreign key")
		return nil
	}

	return site
}

func GetSite(id uint) *Site {
	site := &Site{}

	err := GetDB().Table("sites").Where("id = ?", id).First(site).Error
	if err != nil {
		fmt.Println("There was an error in getting site by ID")
		return nil
	}
	return site
}

func GetAllSites() []*Site {
	sites := make([]*Site, 0)

	err := GetDB().Table("sites").Find(&sites).Error
	if err != nil {
		fmt.Println("There was an error in getting site by ID")
		return nil
	}
	return sites
}

func DeleteSite(id uint) map[string]interface{} {

	//First check if the site exists
	if c := GetSite(id); c == nil {
		return u.Message(false, "There was an error in finding the site")
	}

	//This is a hard delete!
	e := GetDB().Unscoped().Table("sites").Delete(&Site{}, id).Error

	//The command below is a soft delete
	//Meaning that the 'deleted_at' field will be set
	//the record will remain but unsearchable
	//e := GetDB().Table("tenants").Delete(Tenant{}, id).Error
	if e != nil {
		return u.Message(false, "There was an error in deleting the site")
	}

	return u.Message(true, "success")
}

func DeleteSitesOfTenant(id uint) map[string]interface{} {

	//First check if the domain is valid
	if GetDB().Table("sites").Where("domain = ?", id).First(&Site{}).Error != nil {
		return u.Message(false, "The domain was not found")
	}

	//This is a hard delete!
	e := GetDB().Unscoped().Table("sites").
		Where("domain = ?", id).Delete(&Site{}).Error

	//The command below is a soft delete
	//Meaning that the 'deleted_at' field will be set
	//the record will remain but unsearchable
	//e := GetDB().Table("tenants").Delete(Tenant{}, id).Error
	if e != nil {
		return u.Message(false, "There was an error in deleting the site")
	}

	return u.Message(true, "success")
}

func UpdateSite(id uint, newSiteInfo *Site) map[string]interface{} {
	site := &Site{}

	err := GetDB().Table("sites").Where("id = ?", id).First(site).Error
	if err != nil {
		return u.Message(false, "Site was not found")
	}

	if newSiteInfo.Name != "" && newSiteInfo.Name != site.Name {
		site.Name = newSiteInfo.Name
	}

	if newSiteInfo.Category != "" && newSiteInfo.Category != site.Category {
		site.Category = newSiteInfo.Category
	}

	/*if newSiteInfo.Desc != "" && newSiteInfo.Desc != site.Desc {
		site.Desc = newSiteInfo.Desc
	}*/

	//Should it be possible to update domain
	//to new tenant? Will have to think about it more
	//if newSiteInfo.Domain

	/*if newSiteInfo.Color != "" && newSiteInfo.Color != site.Color {
		site.Color = newSiteInfo.Color
	}*/

	if newSiteInfo.Attributes.Orientation != "" {
		switch newSiteInfo.Attributes.Orientation {
		case "NE", "NW", "SE", "SW":
			site.Attributes.Orientation = newSiteInfo.Attributes.Orientation

		default:
		}
	}

	//Successfully validated the new data
	GetDB().Table("site").Save(site)
	return u.Message(true, "success")
}
