package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID            uint   `json:"ID" gorm:"primaryKey"`
	Title         string `json:"Title" gorm:"size:100;not null"`
	Text          string `json:"Text"`
	CategoryRefer uint   `json:"Category"`
	UserRefer     uint   `json:"User"`
	LikeCount     uint   `json:"LikeCount"`
	ImagePath     string `json:"ImagePath"`
}

type Category struct {
	gorm.Model
	ID    uint   `json:"ID" gorm:"primaryKey"`
	Name  string `json:"Name"`
	Posts []Post `gorm:"foreignKey:CategoryRefer" json:"Posts"`
}

type User struct {
	gorm.Model
	ID        uint   `json:"ID" gorm:"primaryKey"`
	Username  string `json:"Username"`
	ImagePath string `json:"ImagePath"`
	Posts     []Post `gorm:"foreignKey:UserRefer" json:"Posts"`
}

func (category Category) toString() string {
	return string(category.ID) + " " + category.Name
}

func (post Post) toString() string {
	return string(post.ID) + " Title:" + post.Title + " Text:" + post.Text + " Category:" + string(post.CategoryRefer)
}

func (user User) toString() string {
	return string(string(user.ID) + "Username:" + user.Username + " ImagePath:" + user.ImagePath)
}

// func connectToDatabase() {
//   dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
//   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//
//   if err!=nil{
//     panic("failed to connect database")
//   }
// }\

var dsn string = "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
var db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func createUser(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var newUser User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter user data")
	}

	fmt.Println("reqbody: ", reqBody)

	json.Unmarshal(reqBody, &newUser)
	fmt.Println(newUser.toString())

	db.Create(&User{Username: newUser.Username, ImagePath: newUser.ImagePath})

	fmt.Println("User with Username:" + newUser.Username + " added to database")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var users []User
	db.Find(&users)

	for _, v := range users {
		fmt.Println(v.toString())
	}

	//w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(users)
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	userID := mux.Vars(r)["id"]

	var user User

	if err := db.Where("ID = ?", userID).First(&user).Error; err != nil {
		fmt.Println("User with ID:" + userID + " not found")
		return
	}

	fmt.Println(user.toString())

	//  w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(user)
}

func getUserByUsername(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	fmt.Println("getUserByUsername")

	// q := r.URL.Query()

	// username := q["username"]

	username := mux.Vars(r)["username"]
	var user User

	fmt.Println("Username:" + username)

	if err := db.Where("Username = ?", username).First(&user).Error; err != nil {
		fmt.Println("User with username:" + username + " not found")
		return
	}

	fmt.Println(user.toString())

	//  w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	userID := mux.Vars(r)["id"]
	var updatedUser User

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the post")
	}
	json.Unmarshal(reqBody, &updatedUser)

	var user User

	if err := db.Where("ID = ?", userID).First(&user).Error; err != nil {
		fmt.Println("User with ID:" + userID + " not found")
		return
	}

	fmt.Println(user.toString())

	db.Model(&user).Updates(User{Username: updatedUser.Username, ImagePath: updatedUser.ImagePath})

	fmt.Println("Post with ID:" + userID + " updated")

	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	userID := mux.Vars(r)["id"]
	var user User

	if err := db.Where("ID = ?", userID).First(&user).Error; err != nil {
		fmt.Println("User with ID:" + userID + " not found")
		return
	}

	db.First(&user, userID)
	db.Delete(&user, userID)

	fmt.Println("User with ID:" + userID + " was deleted")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var newPost Post
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter post data")
	}

	fmt.Println("reqbody: ", reqBody)

	json.Unmarshal(reqBody, &newPost)
	fmt.Println("Title:" + newPost.Title + "Text:" + newPost.Text + "Category:" + string(newPost.CategoryRefer))

	db.Create(&Post{Title: newPost.Title, Text: newPost.Text, CategoryRefer: newPost.CategoryRefer, UserRefer: newPost.UserRefer, LikeCount: 0, ImagePath: newPost.ImagePath})

	fmt.Println("Post with Title:" + newPost.Title + " added to database")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var newCategory Category
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the category name")
	}

	json.Unmarshal(reqBody, &newCategory)
	w.WriteHeader(http.StatusCreated)

	db.Create(&Category{Name: newCategory.Name})
	fmt.Println(newCategory.Name + " added to database")

	json.NewEncoder(w).Encode(newCategory)
}

func getOnePost(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	postID := mux.Vars(r)["id"]

	var post Post

	if err := db.Where("ID = ?", postID).First(&post).Error; err != nil {
		fmt.Println("Category with ID:" + postID + " not found")
		return
	}

	fmt.Println(post.toString())

	//  w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(post)
}

func getOneCategory(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	categoryID := mux.Vars(r)["id"]

	var category Category

	db.First(&category, "ID = ?", categoryID)
	if err := db.Where("ID = ?", categoryID).First(&category).Error; err != nil {
		fmt.Println("Category with ID:" + categoryID + " not found")
		return
	}

	fmt.Println(category.toString())

	json.NewEncoder(w).Encode(category)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var posts []Post
	db.Find(&posts)

	for _, v := range posts {
		fmt.Println(v.toString())
	}

	//w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(posts)
}

func getPostsByRange(w http.ResponseWriter, r *http.Request) {
	limit := mux.Vars(r)["limit"]
	limit_int, _ := strconv.Atoi(limit)

	var posts []Post

	db.Limit(limit_int).Find(&posts)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(posts)
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	var categories []Category
	db.Find(&categories)

	for _, v := range categories {
		fmt.Println(v.toString())
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(categories)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	categoryID := mux.Vars(r)["id"]
	var updatedCategory Category

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the category name only in order to update")
	}
	json.Unmarshal(reqBody, &updatedCategory)

	var category Category

	if err := db.Where("ID = ?", categoryID).First(&category).Error; err != nil {
		fmt.Println("Category with ID:" + categoryID + " not found")
		return
	}

	fmt.Println(category.toString())

	db.Model(&category).Update("Name", updatedCategory.Name)
	fmt.Println("Category with ID:" + categoryID + " updated")

	json.NewEncoder(w).Encode(category)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	postID := mux.Vars(r)["id"]
	var updatedPost Post

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the post")
	}
	json.Unmarshal(reqBody, &updatedPost)

	var post Post

	if err := db.Where("ID = ?", postID).First(&post).Error; err != nil {
		fmt.Println("Post with ID:" + postID + " not found")
		return
	}

	fmt.Println(post.toString())

	db.Model(&post).Updates(Post{Title: updatedPost.Title, Text: updatedPost.Text, CategoryRefer: updatedPost.CategoryRefer, UserRefer: updatedPost.UserRefer, LikeCount: updatedPost.LikeCount, ImagePath: updatedPost.ImagePath})

	fmt.Println("Post with ID:" + postID + " updated")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	categoryID := mux.Vars(r)["id"]
	var category Category

	if err := db.Where("ID = ?", categoryID).First(&category).Error; err != nil {
		fmt.Println("Category with ID:" + categoryID + " not found")
		return
	}

	db.First(&category, categoryID)
	db.Delete(&category, categoryID)

	fmt.Println("Category with ID:" + categoryID + " was deleted")
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	// dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	panic("failed to connect database")
	// }

	postID := mux.Vars(r)["id"]
	var post Post

	if err := db.Where("ID = ?", postID).First(&post).Error; err != nil {
		fmt.Println("Post with ID:" + postID + " not found")
		return
	}

	db.First(&post, postID)
	db.Delete(&post, postID)

	fmt.Println("Post with ID:" + postID + " was deleted")
}

func getPostsByCategoryId(w http.ResponseWriter, r *http.Request) {
	categoryId := mux.Vars(r)["categoryId"]

	print("Category id:" + categoryId)

	var posts []Post

	if err := db.Where("category_refer = ?", categoryId).First(&posts).Error; err != nil {
		print(err)
		fmt.Println("Category with id:" + categoryId + " not found")
		json.NewEncoder(w).Encode(nil)
	}

	// if err := db.Where("category_refer = ?", category.ID).Find(&posts).Error; err != nil {
	// 	print(err)
	// 	fmt.Println("Post with categoryID:" + string(category.ID) + " not found")
	// 	return
	// }

	for _, v := range posts {
		fmt.Println(v.toString())
	}

	//w.Header().Add("Access-Control-Allow-Origin","*");
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(posts)
	// var post Post

	// if err := db.Where("ID = ?", postID).First(&post).Error; err != nil {
	// 	fmt.Println("Category with ID:" + postID + " not found")
	// 	return
	// }

	// fmt.Println(post.toString())

	// //  w.Header().Add("Access-Control-Allow-Origin","*");
	// w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// json.NewEncoder(w).Encode(post)
}

func main() {

	if err != nil {
		panic("failed to connect database")
	}

	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Category{})
	// db.AutoMigrate(&Post{})

	// db.AutoMigrate(
	// 	&User{},
	// 	&Category{},
	// 	&Post{},
	// )

	// db.Migrator().CreateTable(&User{})
	// db.Migrator().CreateTable(&Category{})
	// db.Migrator().CreateTable(&Post{})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/categories/{id}", getOneCategory).Methods("GET")
	router.HandleFunc("/categories/{id}", updateCategory).Methods("PATCH")
	router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")
	router.HandleFunc("/categories", getAllCategories).Methods("GET")
	router.HandleFunc("/categories", createCategory).Methods("POST")
	router.HandleFunc("/posts/{id}", getOnePost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PATCH")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/posts/categoryId/{categoryId}", getPostsByCategoryId).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts", getAllPosts).Methods("GET")
	router.HandleFunc("/posts/limit/{limit}", getPostsByRange).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getOneUser).Methods("GET")
	//	router.HandleFunc("/users?username={username}", getUserByUsername).Methods("GET")
	//router.HandleFunc("/users/?username={username}", getUserByUsername).Methods("GET")
	router.HandleFunc("/users/username/{username}", getUserByUsername).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))

	// start server listen
	// with error handling

	//log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))

	//
	// db.Create(&Category{ID:"1",Name:"Technology",Posts:nil})
	// var category Category
	//
	//  db.First(&category,"Name = ?","Technology")
	//
	// db.Create(&Post{ID:"4",Title:"Title3",Text:"Text3",CategoryRefer:"1"})
	// //  db.Create(&Post{ID:"2", Title:"Title2", Text:"Text2"})
	//
	//   var post Post
	//   db.First(&post,"CategoryRefer = ?","1")

	//  fmt.Println("Title:"+post.Title+" Text:"+post.Text)
}
