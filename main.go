
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

type Post struct {
  gorm.Model
  ID uint `json:"ID" gorm:"primaryKey"`
  Title string `json:"Title"`
  Text string `json:"Text"`
  CategoryRefer uint `json:"Category"`
}

type Category struct {
  gorm.Model
  ID uint `json:"ID" gorm:"primaryKey"`
  Name string `json:"Name"`
  Posts []Post `gorm:"foreignKey:CategoryRefer" json:"Posts"`
}

func (category Category) toString() string {
    return string(category.ID)+" "+category.Name
}

func (post Post) toString() string {
    return string(post.ID)+" Title:"+post.Title+" Text:"+post.Text+" Category:"+string(post.CategoryRefer)
}


// func connectToDatabase() {
//   dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
//   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//
//   if err!=nil{
//     panic("failed to connect database")
//   }
// }

func createPost(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  var newPost Post
  reqBody, err:=ioutil.ReadAll(r.Body)
  if err!=nil {
    fmt.Fprintf(w,"Kindly enter post data")
  }

  json.Unmarshal(reqBody,&newPost)
  w.WriteHeader(http.StatusCreated)
  db.Create(&Post{Title:newPost.Title,Text:newPost.Text,CategoryRefer:newPost.CategoryRefer})


  fmt.Println("Post with Title:"+newPost.Title+" added to database")
  json.NewEncoder(w).Encode(newPost)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  var newCategory Category
  reqBody, err := ioutil.ReadAll(r.Body)
  if err!=nil {
    fmt.Fprintf(w,"Kindly enter data with the category name")
  }

  json.Unmarshal(reqBody, &newCategory)
  w.WriteHeader(http.StatusCreated)

  db.Create(&Category{Name: newCategory.Name})
  fmt.Println(newCategory.Name+" added to database")

  json.NewEncoder(w).Encode(newCategory)
}

func getOnePost(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  postID := mux.Vars(r)["id"]

  var post Post

  if err := db.Where("ID = ?",postID).First(&post).Error; err!=nil{
    fmt.Println("Category with ID:"+postID+" not found")
    return
  }

  fmt.Println(post.toString())

//  w.Header().Add("Access-Control-Allow-Origin","*");
  w.Header().Add("Content-Type","application/json; charset=utf-8")

  json.NewEncoder(w).Encode(post)
}

func getOneCategory(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  categoryID := mux.Vars(r)["id"]

  var category Category

  db.First(&category,"ID = ?", categoryID)
  if err := db.Where("ID = ?",categoryID).First(&category).Error; err!=nil{
    fmt.Println("Category with ID:"+categoryID+" not found")
    return
  }

  fmt.Println(category.toString())

  json.NewEncoder(w).Encode(category)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  var posts []Post
  db.Find(&posts)

  for _,v := range posts {
    fmt.Println(v.toString())
  }

  //w.Header().Add("Access-Control-Allow-Origin","*");
  w.Header().Add("Content-Type","application/json; charset=utf-8")

  json.NewEncoder(w).Encode(posts)
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  var categories []Category
  db.Find(&categories)

  for _,v := range categories {
      fmt.Println(v.toString())
  }

  w.Header().Add("Content-Type","application/json; charset=utf-8")

  json.NewEncoder(w).Encode(categories)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  categoryID := mux.Vars(r)["id"]
  var updatedCategory Category

  reqBody,err := ioutil.ReadAll(r.Body)
  if err!=nil {
    fmt.Fprintf(w,"Kindly enter data with the category name only in order to update")
  }
  json.Unmarshal(reqBody, &updatedCategory)

  var category Category

  if err := db.Where("ID = ?",categoryID).First(&category).Error; err!=nil{
    fmt.Println("Category with ID:"+categoryID+" not found")
    return
  }

  fmt.Println(category.toString())

  db.Model(&category).Update("Name",updatedCategory.Name)
  fmt.Println("Category with ID:"+categoryID+" updated")

  json.NewEncoder(w).Encode(category)
}

func updatePost(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  postID := mux.Vars(r)["id"]
  var updatedPost Post

  reqBody,err := ioutil.ReadAll(r.Body)
  if err!=nil {
    fmt.Fprintf(w,"Kindly enter data with the post")
  }
  json.Unmarshal(reqBody, &updatedPost)

  var post Post

  if err := db.Where("ID = ?",postID).First(&post).Error; err!=nil{
    fmt.Println("Post with ID:"+postID+" not found")
    return
  }

  fmt.Println(post.toString())

  db.Model(&post).Updates(Post{Title:updatedPost.Title, Text:updatedPost.Text, CategoryRefer:updatedPost.CategoryRefer})

  fmt.Println("Post with ID:"+postID+" updated")

  json.NewEncoder(w).Encode(post)
}

func deleteCategory(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  categoryID:=mux.Vars(r)["id"]
  var category Category

  if err := db.Where("ID = ?",categoryID).First(&category).Error; err!=nil{
    fmt.Println("Category with ID:"+categoryID+" not found")
    return
  }


  db.First(&category,categoryID)
  db.Delete(&category,categoryID)

  fmt.Println("Category with ID:"+categoryID+" was deleted")
}

func deletePost(w http.ResponseWriter,r *http.Request) {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  postID:=mux.Vars(r)["id"]
  var post Post

  if err := db.Where("ID = ?",postID).First(&post).Error; err!=nil{
    fmt.Println("Post with ID:"+postID+" not found")
    return
  }

  db.First(&post,postID)
  db.Delete(&post,postID)

  fmt.Println("Post with ID:"+postID+" was deleted")
}

func main() {
  dsn := "sonat:Es@184720158971@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err!=nil{
    panic("failed to connect database")
  }

  db.AutoMigrate(&Post{})
  db.AutoMigrate(&Category{})


  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/categories/{id}",getOneCategory).Methods("GET")
  router.HandleFunc("/categories/{id}",updateCategory).Methods("PATCH")
  router.HandleFunc("/categories/{id}",deleteCategory).Methods("DELETE")
  router.HandleFunc("/categories",getAllCategories).Methods("GET")
  router.HandleFunc("/categories", createCategory).Methods("POST")
  router.HandleFunc("/posts/{id}",getOnePost).Methods("GET")
  router.HandleFunc("/posts/{id}",updatePost).Methods("PATCH")
  router.HandleFunc("/posts/{id}",deletePost).Methods("DELETE")
  router.HandleFunc("/posts",createPost).Methods("POST")
  router.HandleFunc("/posts",getAllPosts).Methods("GET")
  log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
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
