package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
			gin.H{},
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
			gin.H{},
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
			gin.H{},
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
			"users":           getAllObjectsByType("user"),
			"userPermissions": permissions,
		},
	)
}
func SingleUser(c *gin.Context) {
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

	userSql := "SELECT users.name, users.user_id, users.date_created, users.email FROM users WHERE users.user_id = $1"
	var name, id, created_at, email string
	err = db.QueryRow(userSql, customerid).Scan(&name, &id, &created_at, &email)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"user.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"user_name":       name,
			"user_id":         id,
			"created_at":      created_at,
			"email":           email,
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
			gin.H{},
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

	ticketSql := `SELECT tickets.name, tickets.description, products.name, customers.name, contacts.name
				FROM tickets
				INNER JOIN products ON tickets.product_id=products.product_id
				INNER JOIN customers ON tickets.customer_id=customers.customer_id
				INNER JOIN contacts ON tickets.contact_id=contacts.contact_id
				WHERE tickets.ticket_id = $1`
	var name, description, productName, customerName, contactName string
	err = db.QueryRow(ticketSql, ticketId).Scan(&name, &description, &productName, &customerName, &contactName)
	if err != nil {
		log.Println("Failed to execute query: ", err)
	}
	htmlDesc := template.HTML(description)
	// Unmarshal or Decode the JSON to the interface.
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"ticket.tmpl",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"name":            name,
			"description":     htmlDesc,
			"customer_name":   customerName,
			"product_name":    productName,
			"contact_name":    contactName,
			"ticket_id":       ticketId,
			"userPermissions": permissions,
		},
	)
}

func Tickets(c *gin.Context) {
	permissions := jwt.ExtractClaims(c)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ticketSql := `SELECT tickets.ticket_id, tickets.name, tickets.description, products.name, customers.name, contacts.name
				FROM tickets
				INNER JOIN products ON tickets.product_id=products.product_id
				INNER JOIN customers ON tickets.customer_id=customers.customer_id
				INNER JOIN contacts ON tickets.contact_id=contacts.contact_id`
	rows, err := db.Query(ticketSql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var id, name, desc, prod, cust, cont string
	listOfTickets := make(map[string]models.Ticket)
	for rows.Next() {
		fmt.Println("inside rows next")
		err = rows.Scan(&id, &name, &desc, &prod, &cust, &cont)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		var tempTicket models.Ticket
		tempTicket.Name = name
		tempTicket.Description = desc
		tempTicket.Product_name = prod
		tempTicket.Customer_name = cust
		tempTicket.Contact_name = cont
		listOfTickets[id] = tempTicket
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
		created_by := "3"

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
			fmt.Println(customerId)
			sqlStatement := `
			INSERT INTO contacts (name, date_created, created_by, customer_id)
			VALUES ($1, $2, $3, $4) RETURNING contact_id`
			err = db.QueryRow(sqlStatement, name, date_created, created_by, customerId).Scan(&newRecordId)
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
		}

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
	fmt.Println(result["name"] + "is the name")
	userSql := "SELECT name, user_id, email, password, can_create_contact FROM users WHERE email = $1"
	var name, id, email, pass, can_create_contact string
	err = db.QueryRow(userSql, result["name"]).Scan(&name, &id, &email, &pass, &can_create_contact)
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
		authUser.Id = id
		authUser.Can_create_contact = can_create_contact

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
