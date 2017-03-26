package main

import (
	//"io"
	"net/http"
	"fmt"
	//"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"html/template"
	//"strconv"
	"net/url"
	"strings"
	"strconv"
)

const webSiteName  = "my first go blog"
const LIMIT = 5;
var db *sql.DB
var count int

type assignData struct {
	Data 	[LIMIT]homeList
	Name	string
	Page	int
	NextPage	int
	LastPage	int
}

type assignArticleData struct {
	Data 	articleData
	Name	string
}

type homeList struct {
	Id	string
	Title	string
}
type articleData struct {
	Id	string
	Title	string
	Content template.HTML
}

func DbConnect() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@/qipajun")	//对应数据库的用户名和密码
	//defer db.Close()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DbConnect success")
	}
	return db
}

func queryHomeList(db *sql.DB,sql string) (list [LIMIT]homeList, err error)  {
	res, err := db.Query(sql)
	defer db.Close()

	if err != nil {
		log.Println(err)
		return list,err
	}
	var i int = 0;
	for res.Next() {
		var title string
		var id string
		if err := res.Scan(&title,&id); err != nil {
			log.Fatal(err)
			return list,err
		}
		//fmt.Println(id,title)
		list[i] = homeList{id,title}
		i ++
	}
	count = i
	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	return list,err
}
func queryArticle(db *sql.DB,sql string) (data articleData, err error)  {
	res, err := db.Query(sql)
	defer db.Close()

	if err != nil {
		log.Println(err)
		return data,err
	}
	//var i int = 0;
	for res.Next() {
		var title string
		var id string
		var content string
		if err := res.Scan(&title,&id,&content); err != nil {
			log.Fatal(err)
			return data,err
		}
		data = articleData{id,title,template.HTML(content)}
		//i ++
	}
	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	return data,err
}


func Article(w http.ResponseWriter, request *http.Request)  {
	u, err := url.Parse(request.RequestURI)
	if err != nil {
		panic(err)
	}
	path := strings.Split(u.Path,"/")
	//fmt.Println(path)
	var aid int = 0
	if path[1] != "article" || path[1] != "" {
		Aid := strings.Replace(string(path[2]),".html","",-1)
		Aid = strings.Replace(Aid," ","",-1)
		aid,err = strconv.Atoi(Aid)

		if err != nil {
			aid = 0
		}
	}

	fmt.Println("out =>",aid)
	//request.ParseForm()
	//id := request.Form.Get("id")
	//aid, err := strconv.Atoi(id)
	//if err != nil {
	//	aid = 0
	//}
	// := fmt.Sprintf("%s-%d","article",aid)


	//fmt.Println(output)
	//io.WriteString(w, output)
	renderArticle(aid, w)
}

func renderArticle(aid int ,w http.ResponseWriter)  {
	db = DbConnect()
	defer db.Close()
	sql := fmt.Sprintf("select post_title as title,ID as id ,post_content as content from duosutewp_posts where id = %d",aid);
	data, err := queryArticle(db,sql)
	if err != nil {
		log.Fatal(err)
	}

	assignData := assignArticleData{data,webSiteName}
	//fmt.Println(data)
	t,_ := template.ParseFiles("article.html")
	t.Execute(w, assignData)
}



func getPage(request *http.Request)  (page int){

	u, err := url.Parse(request.RequestURI)
	if err != nil {
		panic(err)
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		page = 0
	}

	if len(m["page"]) > 0 {
		page,err = strconv.Atoi(m["page"][0])
		if err != nil {
			page = 0
		}
	}

	return page
}
func Home(w http.ResponseWriter, r *http.Request)  {

	db = DbConnect()
	defer db.Close()

	page := getPage(r)

	startNum := page * LIMIT;

	sql := fmt.Sprintf("select post_title as title,ID as id from duosutewp_posts order by id desc limit %d,%d ",startNum,LIMIT);
	list, err := queryHomeList(db,sql)
	if err != nil {
		log.Fatal(err)
	}

	nextPage := page
	if count >= LIMIT {
		nextPage = nextPage + 1
	}
	lastPage := page

	if lastPage <= 0 {
		lastPage = 0
	} else {
		lastPage = lastPage -1
	}

	assignData := assignData{list, webSiteName, page, nextPage, lastPage}

	t,_ := template.ParseFiles("index.html")
	t.Execute(w, assignData)

}

func main() {
	http.HandleFunc("/index", Home)
	http.HandleFunc("/article/", Article)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}