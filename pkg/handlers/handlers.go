package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anyfei/resolved/pkg/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	// "github.com/AnyFei/resolved/pkg/config"
	// "github.com/AnyFei/resolved/pkg/models"
	// "github.com/AnyFei/resolved/pkg/render"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qferty"
	dbname   = "postgres"
)

// Home is the handler for the home page
func Home(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	//m.App.Session.Start(w, r).Set("test", "test")
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"home.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title":           "Home Page",
			"userPermissions": permissions,
		},
	)
}

func Products(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(permissions["CanCreateContacts"])
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"products.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"products":        getAllObjectsByType("product"),
			"userPermissions": permissions,
			"logged":          true,
		},
	)
}
func SingleProduct(c *gin.Context) {
	customerid := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSql := "SELECT products.name, products.product_id, products.date_created FROM products WHERE products.product_id = $1"
	var product_name, product_id, created_at string
	err = db.QueryRow(userSql, customerid).Scan(&product_name, &product_id, &created_at)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"product.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"product_name":    product_name,
			"product_id":      product_id,
			"created_at":      created_at,
			"userPermissions": permissions,
		},
	)
}
func DisplayCreateProduct(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateContacts"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"create_product.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"no_permission.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "product",
			},
		)
	}
}
func Contacts(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(permissions["CanCreateContacts"])
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"contacts.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"contacts":        getAllObjectsByType("contact"),
			"userPermissions": permissions,
		},
	)
}
func SingleContact(c *gin.Context) {
	userid := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSql := "SELECT contacts.name, contacts.contact_id, contacts.date_created, customers.name, customers.customer_id, contacts.email FROM contacts INNER JOIN customers ON contacts.customer_id=customers.customer_id WHERE contacts.contact_id = $1"
	var name, contact_id, customer_name, customer_id, created_at, email string
	err = db.QueryRow(userSql, userid).Scan(&name, &contact_id, &created_at, &customer_name, &customer_id, &email)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"contact.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"name":            name,
			"contact_id":      contact_id,
			"customer_name":   customer_name,
			"customer_id":     customer_id,
			"created_at":      created_at,
			"email":           email,
			"userPermissions": permissions,
		},
	)
}
func DisplayCreateContact(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateContacts"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"create_contact.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"no_permission.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "contact",
			},
		)
	}
}
func Customers(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(permissions["CanCreateContacts"])
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"customers.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"customers":       getAllObjectsByType("customer"),
			"userPermissions": permissions,
		},
	)
}
func SingleCustomer(c *gin.Context) {
	customerid := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSql := "SELECT customers.name, customers.customer_id, customers.date_created FROM customers WHERE customers.customer_id = $1"
	var customer_name, customer_id, created_at string
	err = db.QueryRow(userSql, customerid).Scan(&customer_name, &customer_id, &created_at)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"customer.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"customer_name":   customer_name,
			"customer_id":     customer_id,
			"created_at":      created_at,
			"userPermissions": permissions,
		},
	)
}
func DisplayCreateCustomer(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateContacts"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"create_customer.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"no_permission.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "customer",
			},
		)
	}
}
func Users(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(permissions["CanCreateContacts"])
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"users.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"Users":           getAllUsers(),
			"userPermissions": permissions,
		},
	)
}
func SingleUser(c *gin.Context) {
	userId := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSql := `SELECT users.name, users.user_id, users.date_created, users.email, users.is_active, can_edit_contacts, users.can_create_contacts,  
				users.can_create_customers, users.can_edit_customers, users.can_create_products, users.can_edit_products, users.can_create_users, users.can_edit_users
				FROM users WHERE users.user_id = $1`
	var user models.User
	err = db.QueryRow(userSql, userId).Scan(&user.Name, &user.ID, &user.DateCreated, &user.Email, &user.IsActive, &user.Can_edit_contacts, &user.Can_create_contacts,
		&user.Can_create_customers, &user.Can_edit_customers, &user.Can_create_products, &user.Can_edit_products, &user.Can_create_users, &user.Can_edit_users)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}

	userGroupsSql := `SELECT groups.group_id, groups.name FROM user_groups INNER JOIN groups ON user_groups.group_id = groups.group_id WHERE user_groups.user_id = $1 
					  ORDER BY user_groups.added_on DESC`
	groupsRows, err := db.Query(userGroupsSql, userId)
	if err != nil {
		fmt.Println(err)
	}
	defer groupsRows.Close()

	user.Groups = make([]models.UserGroup, 0)
	var tempUserGroup models.UserGroup
	for groupsRows.Next() {
		err = groupsRows.Scan(&tempUserGroup.ID, &tempUserGroup.Name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		user.Groups = append(user.Groups, tempUserGroup)
	}
	fmt.Println(user.Can_create_contacts)
	fmt.Println(user.Can_edit_contacts)

	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"user.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"User":            user,
			"userPermissions": permissions,
		},
	)
}
func DisplayCreateUser(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateContacts"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"create_user.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"no_permission.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "user",
			},
		)
	}
}
func DisplayCreateTicket(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	//m.App.Session.Start(w, r).Set("test", "test")
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"create_ticket.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title":           "Home Page",
			"userPermissions": permissions,
		},
	)
}

func SingleTicket(c *gin.Context) {
	ticketId := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ticketSql := `SELECT tickets.name, tickets.description, products.name, products.product_id, customers.name, customers.customer_id, contacts.name, contacts.contact_id, groups.name, groups.group_id, tickets.date_created
				FROM tickets
				LEFT JOIN products ON tickets.product_id=products.product_id
				LEFT JOIN customers ON tickets.customer_id=customers.customer_id
				LEFT JOIN contacts ON tickets.contact_id=contacts.contact_id
				LEFT JOIN groups ON tickets.group_id=groups.group_id
				WHERE tickets.ticket_id = $1`
	var ticket models.Ticket
	err = db.QueryRow(ticketSql, ticketId).Scan(&ticket.Name, &ticket.Description, &ticket.Product_name, &ticket.Product_id,
		&ticket.Customer_name, &ticket.Customer_id, &ticket.Contact_name, &ticket.Contact_id, &ticket.Group_name, &ticket.Group_id, &ticket.DateCreated)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	ticket.ID, _ = strconv.Atoi(ticketId)
	objectSql := `SELECT ticket_actions.action_id, ticket_actions.action_text, ticket_actions.action_type, ticket_actions.date_created, ticket_actions.created_by, users.name 
					FROM ticket_actions INNER JOIN users ON ticket_actions.created_by=users.user_id 
					WHERE ticket_id = $1 
					ORDER BY ticket_actions.date_created DESC`
	rows, err := db.Query(objectSql, ticketId)
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
	}
	defer rows.Close()
	listOfActions := make([]models.TicketAction, 0) //this part can be faster, but it's better to make it slow now and imporve it later on
	for rows.Next() {
		var tempAction models.TicketAction
		err = rows.Scan(&tempAction.Action_id, &tempAction.Action_text, &tempAction.Action_type, &tempAction.Date_created, &tempAction.Created_by, &tempAction.Created_by_name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		listOfActions = append(listOfActions, tempAction)
	}
	ticket.HTML_description = template.HTML(ticket.Description)
	ticket.CustomFields = getCustomFields()
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		http.StatusOK,
		"ticket.tmpl",
		gin.H{
			"Ticket":          ticket,
			"userPermissions": permissions,
			"Actions":         listOfActions,
		},
	)
}

func getCustomFields() []models.CustomField {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cfQuery := `SELECT custom_fields.name, custom_fields.customfield_id, custom_fields.type 
			FROM custom_fields;`

	rows, err := db.Query(cfQuery)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	customFields := make([]models.CustomField, 0)
	//var cfOption models.CustomFieldOption
	for rows.Next() {
		var tempCF models.CustomField
		tempCF.Options = make([]models.CustomFieldOption, 0)
		err = rows.Scan(&tempCF.Name, &tempCF.ID, &tempCF.Type)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		if tempCF.Type == 3 {
			optionsQuery := `SELECT custom_field_options.option_id, custom_field_options.value 
			FROM custom_field_options
			WHERE custom_field_options.customfield_id=` + tempCF.ID + `;`
			fmt.Println(optionsQuery)

			options, err := db.Query(optionsQuery)
			if err != nil {
				fmt.Println(err)
			}
			defer options.Close()
			for options.Next() {
				var tempOption models.CustomFieldOption
				err = options.Scan(&tempOption.ID, &tempOption.Value)
				tempCF.Options = append(tempCF.Options, tempOption)
				fmt.Println(tempOption.Value)

			}
			fmt.Println(tempCF.Options)

		}
		customFields = append(customFields, tempCF)

	}
	return customFields
}
func groupQueryCreator(groups []int) string {
	query := "tickets.group_id = " + strconv.Itoa(groups[0])
	for i := 1; i < len(groups); i++ {
		groupString := strconv.Itoa(groups[i])
		query = query + " OR tickets.group_id = " + groupString
	}
	return query
}
func TicketsByGroup(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userGroups := make([]int, 0)
	userGroupsSql := `SELECT group_id FROM user_groups WHERE user_id = $1`
	groupsRows, err := db.Query(userGroupsSql, 5)
	if err != nil {
		fmt.Println(err)
	}
	var tempGroupId int
	for groupsRows.Next() {
		err = groupsRows.Scan(&tempGroupId)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		userGroups = append(userGroups, tempGroupId)
	}
	defer groupsRows.Close()

	fmt.Println(groupQueryCreator(userGroups))
	ticketSql := `SELECT tickets.ticket_id, tickets.name, products.product_id, products.name, customers.customer_id, customers.name, contacts.contact_id, contacts.name, groups.group_id, groups.name
				FROM tickets
				INNER JOIN products ON tickets.product_id=products.product_id
				INNER JOIN customers ON tickets.customer_id=customers.customer_id
				INNER JOIN contacts ON tickets.contact_id=contacts.contact_id
				INNER JOIN user_groups ON tickets.group_id=user_groups.group_id
				INNER JOIN groups ON tickets.group_id=groups.group_id
				WHERE ` + groupQueryCreator(userGroups)
	rows, err := db.Query(ticketSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var tempTicket models.Ticket
	listOfTickets := make(map[string]models.Ticket)
	for rows.Next() {
		fmt.Println("inside rows next")
		err = rows.Scan(&tempTicket.ID, &tempTicket.Name, &tempTicket.Product_id, &tempTicket.Product_name, &tempTicket.Customer_id, &tempTicket.Customer_name,
			&tempTicket.Contact_id, &tempTicket.Contact_name, &tempTicket.Group_id, &tempTicket.Group_name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		ticketid := strconv.Itoa(tempTicket.ID)
		listOfTickets[ticketid] = tempTicket
	}
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"tickets.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"tickets":         listOfTickets,
			"userPermissions": permissions,
		},
	)
}
func TicketsByGroupJson(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userGroups := make([]int, 0)
	userGroupsSql := `SELECT group_id FROM user_groups WHERE user_id = $1`
	groupsRows, err := db.Query(userGroupsSql, 5)
	if err != nil {
		fmt.Println(err)
	}
	var tempGroupId int
	for groupsRows.Next() {
		err = groupsRows.Scan(&tempGroupId)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		userGroups = append(userGroups, tempGroupId)
	}
	defer groupsRows.Close()

	fmt.Println(groupQueryCreator(userGroups))
	ticketSql := `SELECT tickets.ticket_id, tickets.name, products.product_id, products.name, customers.customer_id, customers.name, contacts.contact_id, contacts.name, groups.group_id, groups.name
				FROM tickets
				INNER JOIN products ON tickets.product_id=products.product_id
				INNER JOIN customers ON tickets.customer_id=customers.customer_id
				INNER JOIN contacts ON tickets.contact_id=contacts.contact_id
				INNER JOIN user_groups ON tickets.group_id=user_groups.group_id
				INNER JOIN groups ON tickets.group_id=groups.group_id
				WHERE ` + groupQueryCreator(userGroups)
	rows, err := db.Query(ticketSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var tempTicket models.Ticket
	listOfTickets := make(map[string]models.Ticket)
	for rows.Next() {
		fmt.Println("inside rows next")
		err = rows.Scan(&tempTicket.ID, &tempTicket.Name, &tempTicket.Product_id, &tempTicket.Product_name, &tempTicket.Customer_id, &tempTicket.Customer_name,
			&tempTicket.Contact_id, &tempTicket.Contact_name, &tempTicket.Group_id, &tempTicket.Group_name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		ticketid := strconv.Itoa(tempTicket.ID)
		listOfTickets[ticketid] = tempTicket
	}
	c.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"tickets": listOfTickets,
		},
	)
}
func Groups(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(permissions["CanCreateContacts"])
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"products.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"products":        getAllObjectsByType("product"),
			"userPermissions": permissions,
			"logged":          true,
		},
	)
}
func SingleGroup(c *gin.Context) {
	groupId := c.Param("id")
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userSql := "SELECT groups.name, groups.group_id, groups.date_created FROM groups WHERE groups.group_id = $1"
	var group_name, group_id, created_at string
	err = db.QueryRow(userSql, groupId).Scan(&group_name, &group_id, &created_at)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"group.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"group_name":      group_name,
			"group_id":        group_id,
			"created_at":      created_at,
			"userPermissions": permissions,
		},
	)
}
func TicketCustomValues(c *gin.Context) {
	ticketId := c.Param("id")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fieldValues := make([]models.ActiveField, 0)
	cfSql := `Select tableProps.fieldid, tableProps.optionId, cfo.value from (select jsonb_object_keys(tickets.custom_values) as fieldid, tickets.custom_values->jsonb_object_keys(tickets.custom_values)->>'value' as optionId 
	from tickets WHERE ticket_id=` + ticketId + `) as tableProps
	inner join custom_field_options cfo on (tableProps.optionId)::int = cfo.option_id ;`
	cfValues, err := db.Query(cfSql)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	defer cfValues.Close()

	for cfValues.Next() {
		var activeValue models.ActiveField
		cfValues.Scan(&activeValue.FieldID, &activeValue.OptionID, &activeValue.Value)
		fieldValues = append(fieldValues, activeValue)
	}
	fmt.Println(fieldValues)
	c.JSON(200, gin.H{
		"ActiveValues": fieldValues,
	})

}

func DisplayCreateGroup(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateContacts"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"create_group.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"no_permission.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "group",
			},
		)
	}
}
func CreateObject(c *gin.Context) {
	var result map[string]string
	c.BindJSON(&result)
	fmt.Println(result["type"])
	if checkIfHavePermissions(c, result["type"]) {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		var newRecordId string
		// Unmarshal or Decode the JSON to the interface.
		name := result["name"]
		date_created := time.Now()
		created_by := jwt.ExtractClaims(c)["UserID"]

		switch result["type"] {
		case "Customer":
			sqlStatement := `
			INSERT INTO customers (name, date_created, created_by)
			VALUES ($1, $2, $3) RETURNING customer_id`
			err = db.QueryRow(sqlStatement, name, date_created, created_by).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "Product":
			sqlStatement := `
			INSERT INTO products (name, date_created, created_by)
			VALUES ($1, $2, $3) RETURNING product_id`
			err = db.QueryRow(sqlStatement, name, date_created, created_by).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "User":
			email := result["email"]
			fmt.Println(email)
			sqlStatement := `
			INSERT INTO users (name, email, date_created, created_by)
			VALUES ($1, $2, $3, $4) RETURNING user_id`
			err = db.QueryRow(sqlStatement, name, email, date_created, created_by).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "Contact":
			fmt.Println("inside contact switch")
			customerId := result["id"]
			contactEmail := result["contactEmail"]
			fmt.Println(customerId)
			sqlStatement := `
			INSERT INTO contacts (name, email, date_created, created_by, customer_id)
			VALUES ($1, $2, $3, $4, $5) RETURNING contact_id`
			err = db.QueryRow(sqlStatement, name, contactEmail, date_created, created_by, customerId).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "Ticket":
			fmt.Println("inside ticket switch")
			customerId := result["customer"]
			productId := result["product"]
			title := result["title"]
			description := result["description"]
			contactId := result["contact"]
			fmt.Println(customerId)
			sqlStatement := `
			INSERT INTO tickets (name, description, product_id, contact_id, customer_id, date_created, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ticket_id`
			err = db.QueryRow(sqlStatement, title, description, productId, contactId, customerId, date_created, created_by).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "Group":
			sqlStatement := `
			INSERT INTO groups (name, date_created, created_by)
			VALUES ($1, $2, $3) RETURNING group_id`
			err = db.QueryRow(sqlStatement, name, date_created, created_by).Scan(&newRecordId)
			if err != nil {
				fmt.Println(err)
			}
		case "UserGroup":
			userId, _ := strconv.Atoi(result["UserID"])
			groupId, _ := strconv.Atoi(result["GroupID"])
			sqlStatement := `
			INSERT INTO user_groups (user_id, group_id, added_on)
			VALUES ($1, $2, $3)`
			_, err = db.Exec(sqlStatement, userId, groupId, date_created)
			if err != nil {
				fmt.Println(err)
				fmt.Println("here")

			}
		}

		c.JSON(200,
			gin.H{"Id": newRecordId},
		)
	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}

func CreateCustomValue(c *gin.Context) {
	var result map[string]string
	c.BindJSON(&result)
	name := result["name"]
	cfType := 3
	cfValues := result["customValues"]
	var slice []string
	_ = json.Unmarshal([]byte(cfValues), &slice)
	var customField models.CustomField
	var cfOption models.CustomFieldOption
	customField.Name = name
	customField.Type = cfType
	customField.Options = make([]models.CustomFieldOption, 0)
	fmt.Println(slice)
	var optionsString string

	if checkIfHavePermissions(c, "Contact") {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		var newRecordId string
		// Unmarshal or Decode the JSON to the interface.
		sqlStatement := `
			INSERT INTO custom_fields (name, type)
			VALUES ($1, $2) RETURNING customfield_id`
		err = db.QueryRow(sqlStatement, customField.Name, customField.Type).Scan(&newRecordId)
		if err != nil {
			fmt.Println(err)
			fmt.Println("here")

		}
		for _, v := range slice {
			fmt.Println(string(v))
			cfOption.Value = string(v)
			cfOption.IsActive = true
			optionsString += "('" + SQLStringEscape(cfOption.Value) + "', '" + newRecordId + "'), "
			customField.Options = append(customField.Options, cfOption)

		}
		optionsString = strings.TrimRight(optionsString, ", ")
		optionSt := `
			INSERT INTO custom_field_options (value, customfield_id)
			VALUES ` + optionsString
		fmt.Println(optionSt)
		_, err = db.Exec(optionSt)
		if err != nil {
			fmt.Println(err)
			fmt.Println("here")

		}
		c.JSON(200,
			gin.H{"Id": newRecordId},
		)
	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}
func DisplayAdmin(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	if permissions["CanCreateUsers"] == "true" {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"admin.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"userPermissions": permissions,
			},
		)
	} else {
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"admin.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"name": "admin section",
			},
		)
	}
}
func DeleteObject(c *gin.Context) {
	var result map[string]string
	c.BindJSON(&result)
	fmt.Println(result["type"])
	if checkIfHavePermissions(c, result["type"]) {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		switch result["type"] {
		case "Customer":
			customerID := result["ID"]
			sqlStatement := `
			DELETE FROM customers WHERE customer_id = $1`
			_, err = db.Exec(sqlStatement, customerID)
			if err != nil {
				fmt.Println(err)
			}
		case "Product":
			productID := result["ID"]
			sqlStatement := `
			DELETE FROM products WHERE product_id = $1`
			_, err = db.Exec(sqlStatement, productID)
			if err != nil {
				fmt.Println(err)
			}
		case "User":
			userID := result["ID"]
			sqlStatement := `
			DELETE FROM users WHERE user_id = $1`
			_, err = db.Exec(sqlStatement, userID)
			if err != nil {
				fmt.Println(err)
			}
		case "Contact":
			contactID := result["ID"]
			sqlStatement := `
			DELETE FROM contacts WHERE contact_id = $1`
			_, err = db.Exec(sqlStatement, contactID)
			if err != nil {
				fmt.Println(err)
			}
		}

	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}
func DeleteUserGroup(c *gin.Context) {
	var result map[string]string
	c.BindJSON(&result)
	userID := jwt.ExtractClaims(c)["UserID"].(string)
	groupID := result["GroupID"]
	fmt.Println("UserID: " + userID + ", GroupID: " + groupID)
	if checkIfHavePermissions(c, result["type"]) {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		sqlStatement := `
			DELETE FROM user_groups WHERE user_id = $1 AND group_id=$2`
		_, err = db.Exec(sqlStatement, userID, groupID)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}
func UpdateObject(c *gin.Context) {
	var result map[string]string
	c.BindJSON(&result)
	fmt.Println(result)
	if checkIfHavePermissions(c, result["type"]) {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		switch result["type"] {
		case "Customer":
			itemType := "customer"
			id := result["id"]
			delete(result, "id")
			delete(result, "type")
			fmt.Println(createCustomUpdateQuery(result, itemType, id))
			_, err = db.Exec(createCustomUpdateQuery(result, itemType, id))
			if err != nil {
				fmt.Println(err)
			}
		case "Product":
			itemType := "product"
			id := result["id"]
			delete(result, "id")
			delete(result, "type")
			fmt.Println(createCustomUpdateQuery(result, itemType, id))
			_, err = db.Exec(createCustomUpdateQuery(result, itemType, id))
			if err != nil {
				fmt.Println(err)
			}
		case "User":
			itemType := "user"
			id := result["id"]
			delete(result, "id")
			delete(result, "type")
			fmt.Println(createCustomUpdateQuery(result, itemType, id))
			_, err = db.Exec(createCustomUpdateQuery(result, itemType, id))
			if err != nil {
				fmt.Println(err)
			}
		case "Contact":
			itemType := "contact"
			id := result["id"]
			delete(result, "id")
			delete(result, "type")
			fmt.Println(createCustomUpdateQuery(result, itemType, id))
			_, err = db.Exec(createCustomUpdateQuery(result, itemType, id))
			if err != nil {
				fmt.Println(err)
			}
		case "UserPermissions":
			permissionType := result["permission"]
			value := result["newValue"]
			id := result["userId"]
			permissionsSql := `UPDATE users SET ` + permissionType + `=` + value + ` WHERE user_id=` + id
			fmt.Println(permissionsSql)
			_, err = db.Exec(permissionsSql)
			if err != nil {
				fmt.Println(err)
			}
		case "TicketCustomValue":
			cfId := result["field_id"]
			optionId := result["option_id"]
			ticketId := result["ticket_id"]
			permissionsSql := `UPDATE tickets
			SET custom_values = custom_values ||'{"` + cfId + `": {"value": "` + optionId + `"}}'::jsonb
			WHERE ticket_id = ` + ticketId + `;`
			fmt.Println(permissionsSql)
			_, err = db.Exec(permissionsSql)
			if err != nil {
				fmt.Println(err)
			}
		}

	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}

func AddAction(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	var result map[string]string
	c.BindJSON(&result)
	if checkIfHavePermissions(c, result["type"]) {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		var newRecordId string
		date_created := time.Now()
		created_by := claims["UserID"]
		description := result["description"]
		ticket_id := result["ticket_id"]
		fmt.Println(result["action_type"])
		actionType, _ := strconv.Atoi(result["action_type"])
		sqlStatement := `
			INSERT INTO ticket_actions (ticket_id, action_text, action_type, date_created, created_by)
			VALUES ($1, $2, $3, $4, $5) RETURNING action_id`
		fmt.Println(actionType)
		err = db.QueryRow(sqlStatement, ticket_id, description, actionType, date_created, created_by).Scan(&newRecordId)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200,
			gin.H{
				"Id": newRecordId,
			},
		)
	} else {
		c.JSON(401, gin.H{
			"Message": "Not authorized",
		})
	}
}
func DisplayLogin(c *gin.Context) {
	//m.App.Session.Start(w, r).Set("test", "test")
	var result map[string]string
	c.BindJSON(&result)
	prevPath := result["fullPath"]
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"login.page.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title":    "Home Page",
			"prevPath": prevPath,
		},
	)
}

func Login(c *gin.Context) (models.AuthUser, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var result map[string]string

	// Unmarshal or Decode the JSON to the interface.
	c.BindJSON(&result)
	fmt.Println(psqlInfo)
	userSql := "SELECT name, user_id, email, password, can_create_contacts, can_edit_contacts,  can_create_customers, can_edit_customers, can_create_products, can_edit_products, can_create_users, can_edit_users FROM users WHERE email = $1"
	var name, id, email, pass, can_create_contacts, can_edit_contacts, can_create_customers, can_edit_customers, can_create_products, can_edit_products, can_create_users, can_edit_users string
	err = db.QueryRow(userSql, result["name"]).Scan(&name, &id, &email, &pass, &can_create_contacts, &can_edit_contacts, &can_create_customers, &can_edit_customers, &can_create_products, &can_edit_products, &can_create_users, &can_edit_users)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}

	//compare the user from the request, with the one we defined:
	if result["password"] != "" && result["name"] != "" && result["password"] == pass && result["name"] == email {
		fmt.Println("Password from form: " + result["password"])
		fmt.Println("Password from DB: " + pass)
		var authUser models.AuthUser
		authUser.Name = name
		authUser.Email = email
		authUser.ID = id
		authUser.Can_create_contacts = can_create_contacts
		authUser.Can_edit_contacts = can_edit_contacts
		authUser.Can_create_customers = can_create_customers
		authUser.Can_edit_customers = can_edit_customers
		authUser.Can_create_products = can_create_products
		authUser.Can_edit_products = can_edit_products
		authUser.Can_create_users = can_create_users
		authUser.Can_edit_users = can_edit_users

		return authUser, nil
	} else {
		fmt.Println("Password from form: " + result["password"])
		fmt.Println("Password from DB: " + pass)
		fmt.Println("WRONG PASSWORD")
		c.JSON(http.StatusUnauthorized, gin.H{"Status": "Wrong credentials"})
		return models.AuthUser{}, errors.New("Wrong credentials")

	}

}
func CustomerContacts(c *gin.Context) {
	customerid := c.Param("id")
	c.JSON(
		200,
		gin.H{
			"contacts": getAllContactsByCustomer(customerid),
		},
	)
}
func checkIfHavePermissions(c *gin.Context, typeOfPermission string) bool {
	permissions := jwt.ExtractClaims(c)

	// if permissions["CanCreate"+typeOfPermission+"s"] == "true" { -- should be this, DB needs to be updated
	if permissions["CanCreateContacts"] == "true" {

		return true
	}
	return false
}
func AllCustomers(c *gin.Context) {
	// [] json: unsupported type: map[int]main.Foo
	c.JSON(200, gin.H{
		"customers": getAllObjectsByType("customer"),
	})
}
func AllTicketsJson(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ticketSql := `SELECT tickets.ticket_id, tickets.name, products.product_id, products.name, customers.customer_id, customers.name, contacts.contact_id, contacts.name, groups.group_id, groups.name
				FROM tickets
				LEFT JOIN products ON tickets.product_id=products.product_id
				LEFT JOIN customers ON tickets.customer_id=customers.customer_id
				LEFT JOIN contacts ON tickets.contact_id=contacts.contact_id
				LEFT JOIN groups ON tickets.group_id=groups.group_id`
	rows, err := db.Query(ticketSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var tempTicket models.Ticket
	listOfTickets := make(map[string]models.Ticket)
	for rows.Next() {
		fmt.Println("inside rows next")
		err = rows.Scan(&tempTicket.ID, &tempTicket.Name, &tempTicket.Product_id, &tempTicket.Product_name, &tempTicket.Customer_id, &tempTicket.Customer_name,
			&tempTicket.Contact_id, &tempTicket.Contact_name, &tempTicket.Group_id, &tempTicket.Group_name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		ticketid := strconv.Itoa(tempTicket.ID)
		listOfTickets[ticketid] = tempTicket
	}
	c.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"tickets": listOfTickets,
		},
	)
}
func AllObjectsByType(c *gin.Context) {
	// [] json: unsupported type: map[int]main.Foo
	objectType := c.Param("type")

	switch objectType {
	case "customers":
		c.JSON(200, gin.H{
			"customers": getAllObjectsByType("customer"),
		})
	case "contacts":
		c.JSON(200, gin.H{
			"contacts": getAllObjectsByType("contact"),
		})
	case "products":
		c.JSON(200, gin.H{
			"products": getAllObjectsByType("product"),
		})
	case "users":
		c.JSON(200, gin.H{
			"users": getAllObjectsByType("user"),
		})
	case "tickets":
		c.JSON(200, gin.H{
			"tickets": getAllObjectsByType("ticket"),
		})
	case "groups":
		c.JSON(200, gin.H{
			"groups": getAllObjectsByType("group"),
		})
	}
}
func getAllObjectsByType(objectType string) map[string]string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Unmarshal or Decode the JSON to the interface.
	objectSql := "SELECT " + objectType + "_id, name FROM " + objectType + "s"
	rows, err := db.Query(objectSql)
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
	}
	defer rows.Close()
	listOfObjects := make(map[string]string)
	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		listOfObjects[id] = name
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return listOfObjects
}
func getAllContactsByCustomer(customerId string) map[string]string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Unmarshal or Decode the JSON to the interface.
	objectSql := "SELECT contacts.name, contacts.contact_id FROM contacts INNER JOIN customers ON contacts.customer_id=customers.customer_id WHERE customers.customer_id=$1"
	rows, err := db.Query(objectSql, customerId)
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
	}
	defer rows.Close()
	listOfContacts := make(map[string]string)
	for rows.Next() {
		var name string
		var contactId string
		err = rows.Scan(&name, &contactId)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		listOfContacts[contactId] = name
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return listOfContacts
}

func getAllUsers() []models.User {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Unmarshal or Decode the JSON to the interface.
	objectSql := "SELECT user_id, name, date_created, email FROM users"
	rows, err := db.Query(objectSql)
	if err != nil {
		// handle this error better than this
		fmt.Println(err)
	}
	defer rows.Close()
	listOfUsers := make([]models.User, 0)
	for rows.Next() {
		var tempUser models.User
		err = rows.Scan(&tempUser.ID, &tempUser.Name, &tempUser.DateCreated, &tempUser.Email)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		tempUser.Avatar = strconv.Itoa(tempUser.ID % 10)
		listOfUsers = append(listOfUsers, tempUser)
	}
	return listOfUsers
}

//function used to create custom sql query
func createCustomUpdateQuery(updateValues map[string]string, itemType string, id string) string {
	var updateString string
	for k, v := range updateValues {
		updateString = updateString + k + " = '" + v + "', "
	}
	updateString = strings.TrimSuffix(updateString, ", ")
	query := ("UPDATE " + itemType + "s SET " + updateString + " WHERE " + itemType + "_id = " + id)
	fmt.Println(query)
	return query

}

func SQLStringEscape(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\'':
			sb.WriteByte('\'')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
