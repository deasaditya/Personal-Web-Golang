package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Template struct {
	templates *template.Template
}

type Project struct {
	Id           int
	Name         string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Description  string
	Technologies []string
	Nodejs       bool
	Reactjs      bool
	Nextjs       bool
	Typescript   bool
	Image        string
}
type User struct {
	ID       int
	Title    string
	Email    string
	Password string
}

var dataProject = []Project{
	{
		Name:        "Deas Aditya",
		Duration:    "2 Months",
		Description: "belajar, berlatih",
		Nodejs:      true,
		Reactjs:     false,
		Nextjs:      true,
		Typescript:  false,
	},
	{
		Name:        "Deas Aditya",
		Duration:    "2 Months",
		Description: "Bismillah",
		Nodejs:      true,
		Reactjs:     true,
		Nextjs:      false,
		Typescript:  false,
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public

	e.Static("/upload", "upload")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/hello", helloWorld)                              //localhost:5000/hello
	e.GET("/", home)                                         //localhost:5000
	e.GET("/contact-me", contact)                            //localhost:5000/contact
	e.GET("/detail/:id", projectDetail)                      //localhost:5000/blog-detail/0 | :id = url params
	e.GET("/formaddproject", formaddproject)                 //localhost:5000/form-blog
	e.POST("/addproject", middleware.UploadFile(addproject)) //localhost:5000/add-blog
	e.GET("/delete-project/:id", deleteproject)
	e.GET("/edit-project/:id", editProject)
	e.POST("/formedit/:id", middleware.UploadFile(edit))
	e.GET("/form-login", formlogin)
	e.POST("/login", login)
	e.GET("/formregister", formregister)
	e.POST("/register", register)
	e.GET("/logout", logout)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func home(c echo.Context) error {

	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["title"],
	}

	delete(sess.Values, "message")

	sess.Save(c.Request(), c.Response())

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM public.tb_projects;")
	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		duration := each.EndDate.Sub(each.StartDate)
		var Durations string
		if duration.Hours()/24 < 7 {
			Durations = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + "Days"
		} else if duration.Hours()/24/7 < 4 {
			Durations = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
		} else if duration.Hours()/24/30 < 12 {
			Durations = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
		} else {
			Durations = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
		}
		each.Duration = Durations
		result = append(result, each)

	}

	Project := map[string]interface{}{
		"Projects": result,
		"Flash":    flash,
	}
	return c.Render(http.StatusOK, "index.html", Project)
}

func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact-me.html", nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	startDateFormat := ProjectDetail.StartDate.Format("02 January 2006")
	endDateFormat := ProjectDetail.EndDate.Format("02 January 2006")
	duration := ProjectDetail.EndDate.Sub(ProjectDetail.StartDate)

	var Durations string
	if duration.Hours()/24 < 7 {
		Durations = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + "Days"
	} else if duration.Hours()/24/7 < 4 {
		Durations = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if duration.Hours()/24/30 < 12 {
		Durations = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		Durations = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	detailProject := map[string]interface{}{
		"Project":  ProjectDetail,
		"Duration": Durations,
		"Start":    startDateFormat,
		"End":      endDateFormat,
	}

	return c.Render(http.StatusOK, "detail.html", detailProject)
}

func formaddproject(c echo.Context) error {
	return c.Render(http.StatusOK, "Add My Project.html", nil)
}

func addproject(c echo.Context) error {
	name := c.FormValue("Project-Name")
	startdate := c.FormValue("Start-Date")
	enddate := c.FormValue("End-Date")
	description := c.FormValue("Description")
	techno := c.Request().Form["check"]
	image := c.Get("dataFile").(string) // => image-982349187nfjka.png
	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects( name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6)", name, startdate, enddate, description, techno, image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	edit := Project{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1;", id).Scan(&edit.Id, &edit.Name, &edit.StartDate, &edit.EndDate, &edit.Description, &edit.Technologies, &edit.Image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var Node, React, Golang, Python bool
	for _, tech := range edit.Technologies {
		if tech == "node" {
			Node = true
		}
		if tech == "react" {
			React = true
		}
		if tech == "golang" {
			Golang = true
		}
		if tech == "python" {
			Python = true
		}
	}

	StartDateFormat := edit.StartDate.Format("2006-01-02")
	EndDateFormat := edit.EndDate.Format("2006-01-02")

	editResult := map[string]interface{}{
		"Edit":   edit,
		"Id":     id,
		"Start":  StartDateFormat,
		"End":    EndDateFormat,
		"Nodejs": Node,
		"React":  React,
		"Golang": Golang,
		"Python": Python,
	}
	return c.Render(http.StatusOK, "Update My Project.html", editResult)
}

func edit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	name := c.FormValue("Project-Name")
	StartDate := c.FormValue("Start-Date")
	EndDate := c.FormValue("End-Date")
	description := c.FormValue("Description")
	techno := c.Request().Form["check"]
	image := c.Get("dataFile").(string) // => image-982349187nfjka.png
	SD, _ := time.Parse("2006-01-02", StartDate)

	ED, _ := time.Parse("2006-01-02", EndDate)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7;", name, SD, ED, description, techno, image, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func formregister(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (title, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		fmt.Println(err)
		redirectWithMessage(c, "Register failed, try again", false, "/formregister")
	}
	return redirectWithMessage(c, "Register success", true, "/form-login")
}

func formlogin(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Title, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email Salah !", false, "/form-login")
	}
	fmt.Println(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/form-login")
	}
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 //3 jam
	sess.Values["message"] = "Login Success !"
	sess.Values["status"] = true //show alert
	sess.Values["title"] = user.Title
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true //akses login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}
