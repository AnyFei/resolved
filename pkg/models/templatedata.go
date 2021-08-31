package models

// type Ticket struct {
// 	ID           string
// 	CustomFields map[string]string `json:"Ticket,omitempty"`
// }

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Customers []Customer
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Tickets   []Ticket
	Headers   map[int]string
}
