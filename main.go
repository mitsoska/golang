package main

import (
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"log"
	jwtware "github.com/gofiber/jwt/v3"
	"time"
)

type Post struct {
	Id uint
	Name string
	Url string
	Time string
	Content string
}

var Database *sql.DB
var err error

func main() {
	log.Println("Creating database")
	filename := "database.db"
	
	// Create the database file if it does not exist	
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Create(filename)
	}

	Database, err = sql.Open("sqlite3", filename)	

	if err != nil {
		log.Fatal("failed to open database file")
	}

	createTable(Database)

	defer Database.Close()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// The static files are all going to be in the public directory
	app.Static("/", "./public")

	app.Get("/chemistry_posts", func(c *fiber.Ctx) error {
		posts := SQL_get_data(Database, "post/chemistry.html")
		name := get_cookie_name(c)
		
		return c.Render("forums", fiber.Map {
			"topic": "Θεωρία Χημείας",
			"Url": "post/chemistry.html",
			"posts": posts,
			"Logged": name,
		}, "layouts/main")
	})

	app.Get("/mathematics_posts", func(c *fiber.Ctx) error {
		posts := SQL_get_data(Database, "post/mathematics.html")
		name := get_cookie_name(c)
		
		return c.Render("forums", fiber.Map {
			"topic": "Θεωρία Μαθηματικών",
				"Url": "post/mathematics.html",
				"posts": posts,
				"Logged": name,
			}, "layouts/main")
	})

	app.Get("/physics_posts", func(c *fiber.Ctx) error {
		posts := SQL_get_data(Database, "post/physics.html")
		name := get_cookie_name(c)
		
		return c.Render("forums", fiber.Map {
			"topic": "Θεωρία Φυσικής",
				"Url": "post/physics.html",
				"posts": posts,
				"Logged": name,
			}, "layouts/main")
	})

	app.Get("/material_posts", func(c *fiber.Ctx) error {
		posts := SQL_get_data(Database, "post/material.html")
		name := get_cookie_name(c)
		
		return c.Render("forums", fiber.Map {
			"topic": "Υλικό",
				"Url": "post/material.html",
				"posts": posts,
				"Logged": name,
			}, "layouts/main")
	})

	app.Get("/cs_posts", func(c *fiber.Ctx) error {
		posts := SQL_get_data(Database, "post/cs.html")
		name := get_cookie_name(c)
		
		return c.Render("forums", fiber.Map {
			"topic": "Χρήσιμα προγράμματα",
				"Url": "post/cs.html",
				"posts": posts,
				"Logged": name,
			}, "layouts/main")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		chemistry := SQL_get_most_recent(Database, "post/chemistry.html")
		math := SQL_get_most_recent(Database, "post/mathematics.html")
		physics := SQL_get_most_recent(Database, "post/physics.html")
		material := SQL_get_most_recent(Database, "post/material.html")
		cs := SQL_get_most_recent(Database, "post/cs.html")

		name := get_cookie_name(c)
		
		return c.Render("index", fiber.Map {
			"Chemistry": chemistry,
			"Mathematics": math,
			"Physics": physics,
			"Material": material,
			"CS": cs,
			"Logged": name,
			}, "layouts/main")
	})

	app.Get("/register.html", func(c *fiber.Ctx) error {
		name := get_cookie_name(c)
		
		return c.Render("register", fiber.Map {
			"Logged": name,
			}, "layouts/main")
	})

	app.Get("/login.html", func(c *fiber.Ctx) error {
		name := get_cookie_name(c)
		
		return c.Render("login", fiber.Map {
			"Logged": name,
			}, "layouts/main")
	})


	app.Post("/register", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		log.Println("Registering. ... ", username, password)
		register_user(username, password)

		// Redirect back to home
		return c.Redirect("/login.html")
	})

	app.Post("/Login", login)

	app.Post("/post/chemistry.html", func(c *fiber.Ctx) error {
		content := c.FormValue("content")

		username := get_cookie_name(c)
		
		// Save content into the database
		insertComment(Database, username, "post/chemistry.html", content)
		log.Println("Inserted new comment", username, "post/chemistry.html", content)
		return c.Redirect("/chemistry_posts")
	})

	app.Post("/post/mathematics.html", func(c *fiber.Ctx) error {
		content := c.FormValue("content")

		username := get_cookie_name(c)
		
		// Save content into the database
		insertComment(Database, username, "post/mathematics.html", content)
		log.Println("Inserted new comment", username, "post/mathematics.html", content)
		return c.Redirect("/mathematics_posts")
	})

	app.Post("/post/physics.html", func(c *fiber.Ctx) error {
		content := c.FormValue("content")

		username := get_cookie_name(c)
		
		// Save content into the database
		insertComment(Database, username, "post/physics.html", content)
		log.Println("Inserted new comment", username, "post/physics.html", content)
		return c.Redirect("/physics_posts")
	})

	app.Post("/post/material.html", func(c *fiber.Ctx) error {
		content := c.FormValue("content")

		username := get_cookie_name(c)
		
		// Save content into the database
		insertComment(Database, username, "post/material.html", content)
		log.Println("Inserted new comment", username, "post/material.html", content)
		return c.Redirect("/material_posts")
	})

	app.Post("/post/cs.html", func(c *fiber.Ctx) error {
		content := c.FormValue("content")

		username := get_cookie_name(c)
		
		// Save content into the database
		insertComment(Database, username, "post/cs.html", content)
		log.Println("Inserted new comment", username, "post/cs.html", content)
		return c.Redirect("/cs_posts")
	})


	// Creating a JWT middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))
	
	//app.Post("/protected", Protected)
	
	log.Fatal(app.Listen(":3838"))
}

func createTable(database *sql.DB) {

	table := `
CREATE TABLE IF NOT EXISTS user (
           id INTEGER PRIMARY KEY,
           username TEXT,
           password TEXT
);
`
	table2 := `	
CREATE TABLE IF NOT EXISTS comment (
           id INTEGER PRIMARY KEY,
           name TEXT,
           url TEXT,
           date TEXT,
           comment TEXT
        );
`

	statement, err := database.Prepare(table)
	statement2, err := database.Prepare(table2)
	
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
	statement2.Exec()
}

func insertComment(database *sql.DB, name string, url string, comment string) {
	insertCommentSQL := `INSERT  INTO comment(name, url, date, comment) VALUES (?, ?, ?, ?)`

	statement, err := Database.Prepare(insertCommentSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(name, url, time.Now().UTC().Format("02/01/2006"), comment)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func register_user(name string, password string) {
	insertUserSQL := `INSERT  INTO user(username, password) VALUES (?, ?)`
	statement, err := Database.Prepare(insertUserSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(name, password)

	if err != nil {
		log.Fatalln(err.Error())
	}


}

func displayComments(database *sql.DB) {
	row, err := database.Query("SELECT * FROM comment ORDER BY name")

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer row.Close()

	for row.Next() {
		var id int
		var name string
		var comment string
		row.Scan(&id, &name, &comment)
		log.Println("Student: ", " ", name, " ", comment)
	}
}

func SQL_get_data(database *sql.DB, url string) []Post {
	row, err := database.Query("SELECT * FROM comment WHERE url=?", url)

	if err != nil {
		log.Fatalln(err.Error())
	}
	
	defer row.Close()
	i := 0

	// Load all comments into a dynamic array
	posts := []Post{}
	
	for row.Next() {
		var post Post
		row.Scan(&post.Id, &post.Name, &post.Url, &post.Time, &post.Content)

		posts = append(posts, post);
		i += 1
	}

	return posts
}

func SQL_get_most_recent(database *sql.DB, url string) Post {
	row, err := database.Query("SELECT * FROM comment WHERE url=? ORDER BY id DESC LIMIT 1", url)


	if err != nil {
		log.Fatalln(err.Error())
	}
	
	defer row.Close()

	var post Post

	for row.Next() {
		row.Scan(&post.Id, &post.Name, &post.Url, &post.Time, &post.Content)
	}
	
	return post
}

func SQL_get_count(database *sql.DB) (count int ) {
	rows, _ := database.Query("SELECT * FROM comment ORDER BY name")

	count = 0

	for rows.Next() {
		count += 1
	}

	return count
}

