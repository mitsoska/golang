package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/mattn/go-sqlite3"
	jtoken "github.com/golang-jwt/jwt/v4"
	"time"
	"log"
	"github.com/golang-jwt/jwt/v4"
)

type login_request struct {
	// Using tags
	Username string `json:"name"` 
	Password string `json:"password"`
}

type login_response struct {
	Token string `json:"token"`
}

type User struct {
	Id int
	Username string
	Password string
}

const Secret = "secret"

func new_auth_middleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config {
		SigningKey: []byte(secret),
	})
}

func find_by_credentials(username string, password string) (*User, error) {
	log.Println("Trying to find a user", username, password)
	
	query := `SELECT * FROM user WHERE username=? AND password=?`

	statement, err := Database.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}

	
	defer statement.Close()
	
	var user User

	err = statement.QueryRow(username, password).Scan(&user.Id, &user.Username, &user.Password)
	
	if err != nil {
		return nil, err 
	}

	return &user, nil
}

func login(c *fiber.Ctx) error {

	log.Println("LOGGING IN")
	
	new_login_request := new(login_request)

	// BodyParser binds the request body to the struct
	if err := c.BodyParser(new_login_request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Find user
	user, error := find_by_credentials(new_login_request.Username, new_login_request.Password)

	if error != nil {
		return c.Redirect("/login.html")
		//return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		//"error": error.Error(),
	}
	

	// Creating the JWT claims. The token lasts for a day
	claims := jtoken.MapClaims {
		"Id": user.Id,
		"Username": user.Username,
		"Exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Creating the token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// Generate encoded token
	tok, err := token.SignedString([]byte(Secret))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie {
		Name: "tok",
		Value: tok,
	})

	log.Println("Ok logged in correctly")
	return c.Redirect("/")

}

func Protected(c *fiber.Ctx) error {

	// Get the user from the context and return it

	// Locals: A method that stores variables scoped to the request, and therefore are available only to the routes that match the request
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	name := claims["username"].(string)
	return c.SendString("Welcome " + name)
}

func get_cookie_name(c *fiber.Ctx) (string) {
	cookie := c.Cookies("tok")

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	
	if err != nil {
		return "Ανώνυμος: "
	}

	payload := token.Claims.(jwt.MapClaims)

	return payload["Username"].(string)
}

