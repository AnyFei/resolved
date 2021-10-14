package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anyfei/resolved/pkg/handlers"
	"github.com/anyfei/resolved/pkg/models"
	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
)

// "github.com/AnyFei/resolved/pkg/config"
// "github.com/AnyFei/resolved/pkg/handlers"
// "github.com/AnyFei/resolved/pkg/models"
// "github.com/AnyFei/resolved/pkg/render"
// "github.com/AnyFei/resolved/pkg/config"

//"github.com/alexedwards/scs/v2"
//"github.com/alexedwards/scs/v2"

const portNumber = ":8080"

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var tempPath string

const (
	addr     = "imap.gmail.com:993"
	username = "resolved.contact.test@gmail.com"
	password = "ecuoywbaxeescyby"
)

// main is the main function
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.tmpl")

	// the jwt middleware
	AuthMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "test zone",
		Key:            []byte("secret key"),
		Timeout:        time.Minute * 30,
		MaxRefresh:     time.Minute * 30,
		IdentityKey:    "Email",
		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		CookieDomain:   "localhost",
		CookieName:     "token", // default jwt
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.AuthUser); ok {
				return jwt.MapClaims{
					"Email":              v.Email,
					"Name":               v.Name,
					"CanCreateContacts":  v.Can_create_contacts,
					"CanEditContacts":    v.Can_edit_contacts,
					"CanCreateCustomers": v.Can_create_customers,
					"CanEditCustomers":   v.Can_edit_customers,
					"CanCreateProducts":  v.Can_create_products,
					"CanEditProducts":    v.Can_edit_products,
					"CanCreateUsers":     v.Can_create_users,
					"CanEditUsers":       v.Can_edit_users,
					"UserID":             v.ID,
					"Logged":             true,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.AuthUser{
				ID:                   claims["UserID"].(string),
				Email:                claims["Email"].(string),
				Name:                 claims["Name"].(string),
				Can_create_contacts:  claims["CanCreateContacts"].(string),
				Can_edit_contacts:    claims["CanEditContacts"].(string),
				Can_create_customers: claims["CanCreateCustomers"].(string),
				Can_edit_customers:   claims["CanEditCustomers"].(string),
				Can_create_products:  claims["CanCreateProducts"].(string),
				Can_edit_products:    claims["CanEditProducts"].(string),
				Can_create_users:     claims["CanCreateUsers"].(string),
				Can_edit_users:       claims["CanEditUsers"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			//c.JSON(200, gin.H{"logged": true})
			authUser, err := handlers.Login(c)
			if err != nil {
				c.JSON(200, gin.H{"Status": "error"})
				return nil, jwt.ErrFailedAuthentication
			}
			fmt.Println(authUser)
			return &authUser, nil

		},
		LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
			fmt.Println(tempPath)
			c.JSON(200,
				gin.H{
					"Status":   "ok",
					"prevPath": tempPath,
				})
			tempPath = ""
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.AuthUser); ok && v.Name != "" {
				// c.HTML(http.StatusOK, "home.tmpl", gin.H{"loggedIn": true})
				return true
			}

			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println(c.FullPath())
			if c.FullPath() != "/login" && c.FullPath() != "" && c.FullPath() != "/refresh_token" {
				tempPath = strings.ReplaceAll(c.FullPath(), ":id", c.Param("id"))
				fmt.Println(tempPath)
			}
			c.HTML( // Set the HTTP status to 200 (OK)
				http.StatusOK,
				// Use the index.html template
				"unauthorized.tmpl",
				// Pass the data that the page uses (in this case, 'title')
				gin.H{
					"title": "Unauthorized",
				})
		},
		LogoutResponse: func(c *gin.Context, i int) {
			c.HTML(200, "home.tmpl", gin.H{})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := AuthMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("AuthMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	router.POST("/login", AuthMiddleware.LoginHandler)

	router.NoRoute(AuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.HTML( // Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the index.html template
			"404.tmpl",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "404",
			})
	})

	auth := router.Group("/")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", AuthMiddleware.RefreshHandler)
	//auth.Use(AuthMiddleware.RefreshHandler)
	auth.Use(AuthMiddleware.MiddlewareFunc())
	{

		auth.GET("/customers", handlers.Customers)
		auth.GET("/customers/:id", handlers.SingleCustomer)
		auth.GET("/customers/:id/contacts", handlers.CustomerContacts)
		auth.GET("/customers/new", handlers.DisplayCreateCustomer)

		auth.GET("/contacts", handlers.Contacts)
		auth.GET("/contacts/new", handlers.DisplayCreateContact)
		auth.GET("/contacts/:id", handlers.SingleContact)

		auth.GET("/products", handlers.Products)
		auth.GET("/products/new", handlers.DisplayCreateProduct)
		auth.GET("/products/:id", handlers.SingleProduct)

		auth.GET("/users", handlers.Users)
		auth.GET("/users/new", handlers.DisplayCreateUser)
		auth.GET("/users/:id", handlers.SingleUser)

		auth.GET("/tickets/new", handlers.DisplayCreateTicket)
		auth.GET("/tickets/:id", handlers.SingleTicket)
		auth.GET("/tickets", handlers.TicketsByGroup)
		auth.GET("/tickets/mygroups", handlers.TicketsByGroupJson)
		auth.GET("/tickets/all", handlers.AllTicketsJson)
		auth.GET("/tickets/:id/customValues", handlers.TicketCustomValues)

		auth.POST("/new", handlers.CreateObject)
		auth.POST("/delete", handlers.DeleteObject)
		auth.POST("/delete/usergroup", handlers.DeleteUserGroup)
		auth.POST("/update", handlers.UpdateObject)
		auth.POST("/allCustomers", handlers.AllCustomers)

		auth.GET("/groups/new", handlers.DisplayCreateGroup)
		auth.GET("/groups/:id", handlers.SingleGroup)
		auth.GET("/groups", handlers.Groups)

		auth.POST("/add_action", handlers.AddAction)

		auth.GET("/admin", handlers.DisplayAdmin)
		auth.POST("/admin/add-custom-field", handlers.CreateCustomValue)

		auth.GET("/all/:type", handlers.AllObjectsByType)

		auth.GET("/", handlers.Home)

	}
	router.GET("/login", handlers.DisplayLogin)
	router.GET("/logout", AuthMiddleware.LogoutHandler)

	router.Run("0.0.0.0:8080")

}
