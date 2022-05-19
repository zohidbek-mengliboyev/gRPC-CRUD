package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	pbu "grpccrud/api_gateway/genproto/book_service"
	pb "grpccrud/api_gateway/genproto/user_service"

	"github.com/gin-gonic/gin"
	
)

var client pb.UserManagementClient
var bookClient pbu.BookManagementClient


func CreateNewUser(ct *gin.Context) {
	n := pb.NewUser{}
	err := ct.BindJSON(&n)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := client.CreateNewUser(ct, &n)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct.JSON(http.StatusCreated, res)
}

// GetUsers godoc
// @Summary List of users
// @Description get users.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} models.Users
// @Failure 400  {object}  httputil.HTTPError
// @Failure 404  {object}  httputil.HTTPError
// @Failure 500  {object}  httputil.HTTPError
// @Router / [get]
func GetUsers(ct *gin.Context) {
	param := pb.GetUsersParams{}

	res, err := client.GetUsers(ct, &param)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func GetById(ct *gin.Context) {

	id, err := strconv.ParseInt(ct.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := pb.UserRequest{Id: id}

	res, err := client.GetById(ct, &obj)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func UpdateUser(ct *gin.Context) {
	u := pb.User{}
	err := ct.BindJSON(&u)
	if err != nil {
		fmt.Println("here", err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("RESULT: ", u.Age, u.Name)
	id, err := strconv.ParseInt(ct.Param("id"), 10, 64)
	if err != nil {
		fmt.Println("not here", err)
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := client.UpdateUser(ct, &pb.UpdateUserRequest{Id: id,
		User: &u,
	})
	if err != nil {
		log.Println("33333333333", err.Error())
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func DeleteUser(ct *gin.Context) {
	id, err := strconv.ParseInt(ct.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := pb.UserRequest{Id: id}

	res, err := client.DeleteUser(ct, &obj)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.Header("Entity", fmt.Sprintf("%v", res.GetSuccess()))
	ct.JSON(http.StatusNoContent, res)
}

func Create(ct *gin.Context) {

	b := pbu.NewBook{}
	err := ct.BindJSON(&b)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := bookClient.Create(ct, &b)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ct.JSON(http.StatusCreated, res)
}

func GetAll(ct *gin.Context) {

	bookParams := pbu.BookParams{}

	res, err := bookClient.GetAll(ct, &bookParams)
	if err != nil {
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func GetUserBooks(ct *gin.Context) {

	user_id, err := strconv.ParseInt(ct.Param("user_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := pbu.UserBookRequest{UserId: user_id}

	res, err := bookClient.GetUserBooks(ct, &obj)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	object := pb.UserRequest{Id: user_id}

	rev, err := client.GetById(ct, &object)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res.Age = rev.User.Age
	res.Name = rev.User.Name
	res.Id = rev.User.Id

	ct.JSON(http.StatusOK, res)
}

func Get(ct *gin.Context) {

	book_id, err := strconv.ParseInt(ct.Param("book_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := pbu.BookRequest{BookId: book_id}

	res, err := bookClient.Get(ct, &obj)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func Update(ct *gin.Context) {
	b := pbu.Book{}
	err := ct.BindJSON(&b)
	if err != nil {
		fmt.Println("here", err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("RESULT: ", b.BookId, b.BookName, b.Author)
	book_id, err := strconv.ParseInt(ct.Param("book_id"), 10, 64)
	if err != nil {
		fmt.Println("not here", err)
		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := bookClient.Update(ct, &pbu.UpdateBookRequest{BookId: book_id,
		Book: &b,
	})
	if err != nil {
		log.Println("33333333333", err.Error())
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.JSON(http.StatusOK, res)
}

func Delete(ct *gin.Context) {

	book_id, err := strconv.ParseInt(ct.Param("book_id"), 10, 64)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	obj := pbu.BookRequest{BookId: book_id}

	res, err := bookClient.Delete(ct, &obj)
	if err != nil {
		fmt.Println(err)

		ct.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.Header("Entity", fmt.Sprintf("%v", res.GetSuccess()))
	ct.JSON(http.StatusNoContent, res)
}
