package models

import "time"

type Object struct {
	ID           int
	DateCreated  time.Time
	CreatedBy    int
	DateModified time.Time
	Notes        map[int]string
	IsActive     bool
	Name         string
	CustomFields map[string]string //custom field should be a separate struct
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
type Ticket struct {
	Object
	Product_id    string
	Customer_id   string
	Contact_id    string
	Description   string
	Product_name  string
	Customer_name string
	Contact_name  string
}

type ResponseObject struct {
	Id     int
	Status string
}

type User struct {
	Object
	Email              string `json:"email"`
	Password           string `json:"password"`
	Can_create_contact string
}

type AuthUser struct {
	Name               string
	Id                 string
	Can_create_contact string
	Email              string
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
