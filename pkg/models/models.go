package models

import (
	"database/sql"
	"database/sql/driver"
	"html/template"
	"time"
)

type NullString string

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		*s = ""
	}
	*s = NullString(strVal)
	return nil
}
func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

type Object struct {
	ID           int
	DateCreated  time.Time
	CreatedBy    int
	DateModified time.Time
	Notes        map[int]string
	IsActive     bool
	Name         string
}

type Customer struct {
	Object
	Domains map[int]string
}

type Contact struct {
	Object
	Email         string
	Phone         string
	Customer_id   string
	Customer_name string
}

type UserGroup struct {
	Name string
	ID   int
}

type Ticket struct {
	Object
	Product_id       sql.NullString
	Customer_id      sql.NullString
	Contact_id       sql.NullString
	Description      string
	HTML_description template.HTML
	Product_name     sql.NullString
	Customer_name    sql.NullString
	Contact_name     sql.NullString
	Group_id         sql.NullString
	Group_name       sql.NullString
	CustomFields     []CustomField
}

type ResponseObject struct {
	Id     int
	Status string
}
type ActiveField struct {
	FieldID  string `json:"CustomFieldID"`
	OptionID string `json:"OptionID"`
	Value    string `json:"OptionValue"`
}
type CustomField struct {
	ID       string
	Name     string
	IsActive bool
	Type     int
	Options  []CustomFieldOption
}
type CustomFieldOption struct {
	ID       string `json:"id"`
	Value    string `json:"value"`
	IsActive bool   `json:"active"`
}
type User struct {
	Object
	Email                NullString `json:"email"`
	Password             NullString `json:"password"`
	Can_create_contacts  bool
	Can_edit_contacts    bool
	Can_create_customers bool
	Can_edit_customers   bool
	Can_create_products  bool
	Can_edit_products    bool
	Can_create_users     bool
	Can_edit_users       bool
	Groups               []UserGroup
	Avatar               string
}

type TicketAction struct {
	Action_id       string
	Ticket_id       string
	Action_text     template.HTML
	Action_type     string
	Date_created    string
	Created_by      string
	Created_by_name string
	Source          string
}
type AuthUser struct {
	Name                 string
	ID                   string
	Can_create_contacts  string
	Can_edit_contacts    string
	Can_create_customers string
	Can_edit_customers   string
	Can_create_products  string
	Can_edit_products    string
	Can_create_users     string
	Can_edit_users       string
	Email                string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}
