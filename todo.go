package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// create new instance of echo

	//database init

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.DELETE},
	}))

	type Employee struct {
		Id     string `json:"id"`
		Name   string `json:"employee_name"`
		Salary string `json:"employee_salary"`
		Age    string `json:"employee_age"`
	}

	type Employees struct {
		Employees []Employee `json:"employee`
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go-todo")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db is connected")
	}
	defer db.Close()
	// make sure connnection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	//post method
	e.POST("/employee", func(c echo.Context) error {
		emp := new(Employee)
		if err := c.Bind(emp); err != nil {
			return err
		}

		sql := "INSERT INTO employee(employee_name,employee_age,employee_salary) VALUES(?,?,?)"
		stmt, err := db.Prepare(sql)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer stmt.Close()
		result, err2 := stmt.Exec(emp.Name, emp.Salary, emp.Age)
		// exit if we get an error
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(result.LastInsertId())
		return c.JSON(http.StatusCreated, emp.Name)
	})

	// delete
	e.DELETE("/employee/:id", func(c echo.Context) error {
		requested_id := c.Param("id")
		sql := "DELETE FROM employee WHERE id=?"
		stmt, err := db.Prepare(sql)
		if err != nil {
			fmt.Println(err)
		}
		result, err2 := stmt.Exec(requested_id)
		if err2 != nil {
			panic(err2)
		}
		fmt.Println(result.RowsAffected())
		return c.JSON(http.StatusOK, "Deleted")

	})

	// fetch by id
	e.GET("employee/:id", func(c echo.Context) error {
		requested_id := c.Param("id")
		fmt.Println(requested_id)
		var name string
		var id string
		var salary string
		var age string

		err = db.QueryRow("SELECT id,employee_name,employee_age,employee_salary FROM employee WHERE id=?", requested_id).Scan(&id, &name, &age, &salary)
		if err != nil {
			fmt.Println(err)
		}
		response := Employee{Id: id, Name: name, Salary: salary, Age: age}
		return c.JSON(http.StatusOK, response)
	})

	// fetch all data

	e.GET("employees", func(c echo.Context) error {
		sqlStatement := "SELECT id,employee_name,employee_salary,employee_age FROM employee order by id"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		result := Employees{}

		for rows.Next() {
			employee := Employee{}
			err2 := rows.Scan(&employee.Id, &employee.Name, &employee.Salary, &employee.Age)

			if err2 != nil {
				return err2
			}
			result.Employees = append(result.Employees, employee)
		}
		return c.JSON(http.StatusCreated, result)
	})

	// update

	e.PUT("/employee", func(c echo.Context) error {
		u := new(Employee)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "UPDATE employee SET employee_name=$1,employee_salary=$2,employee_age=$3 WHERE id=$4"
		res, err := db.Query(sqlStatement, u.Name, u.Salary, u.Age, u.Id)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, u.Id)
	})
	/*
		e.GET("/tasks", func(c echo.Context) error { return c.JSON(200, "GET Tasks") })
		e.PUT("/tasks", func(c echo.Context) error { return c.JSON(200, "PUT Tasks") })
		e.DELETE("/tasks/:id", func(c echo.Context) error { return c.JSON(200, "DELETE Task "+c.Param("id")) })
	*/
	// start as a web server
	e.Start(":8000")
}
