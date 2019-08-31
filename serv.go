package main

import (
	"bytes"
	"strings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"time"
)

const listenPort = ":8080"

func viewHandler(w http.ResponseWriter, r *http.Request) {
	filename := "blogIndex.html"
	body, _ := ioutil.ReadFile(filename)
	fmt.Fprintf(w, string(body))
}

//BlogCTX holds the context for the blog file
type BlogCTX struct {
	Title       string
	Description string
	Image       string
	Video       string
	Tags        string
	TagSlice    []string
	Categories  string
	CatSlice    []string
	Draft       string
	Date        string
	Body        string
}

// SaveHandler saves edited wiki data and reloads pages
func saveHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var blogData = bytes.Buffer{}
	var ctx BlogCTX
	ctx.TagSlice = make([]string, 0)
	ctx.CatSlice = make([]string, 0)
	err := decoder.Decode(&ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strings.Split(ctx.Tags, " "))
	ctx.TagSlice = strings.Split(ctx.Tags, " ")
	ctx.CatSlice = strings.Split(ctx.Categories, " ")
	ctx.Date = time.Now().Format("2006-01-02T15:04:05") //
	ctx.Body = ctx.Body
	fmt.Println(ctx.Image)
	if ctx.Image != "" {
		
			fmt.Println(ctx.Image[:4])
			if ctx.Image[:4] == "http" {
				cmd := exec.Command("wget", "-O", "./static/images/"+strings.Replace(ctx.Title, " ", "", -1)+".jpg", ctx.Image)
				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Waiting for image command to finish...")
				err = cmd.Wait()
				if err != nil {
					fmt.Printf("Command finished with error: %v", err)
				}
				fmt.Println("done waiting for image")
			}
		
		ctx.Image = strings.Replace(ctx.Title, " ", "", -1) + ".jpg"

	}

	t, err := template.New("blogTemplate.t").Funcs(NextIndex).ParseFiles("./blogTemplate.t")
	if err != nil {
		fmt.Println("error parsing template: ", err)
	} //fill requestBody with the executed template and context
	err = t.Execute(&blogData, ctx)
	if err != nil {
		fmt.Println("Error executing template: ", err)
	}
	file, err := os.Create(`./content/project/` + strings.Replace(ctx.Title," ", "", -1 ) + ".md")
	if err != nil {
		fmt.Println("error creating file: ", err)
	}
	file.Write(blogData.Bytes())
	file.Close()
	go func() {
		cmd := exec.Command("hugo")
		err = cmd.Run()
		if err != nil {
			fmt.Println("error executing command: ", err)
		}
	}()
	/*fmt.Println("Title: ", t.Title)
	fmt.Println("Description: ", t.Description)
	fmt.Println("Image: ", t.Image)
	fmt.Println("Video: ", t.Video)
	fmt.Println("Tags: ", t.Tags)
	fmt.Println("Categories: ", t.Categories)
	fmt.Println("Draft: ", t.Draft)
	fmt.Println("Body: ", t.Body)*/
	blogData.Reset()
}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(listenPort, nil)
}

//NextIndex is a function for the go templates to do basic math
var NextIndex = template.FuncMap{
	"NextIndex": func(i int) int {
		return i + 1
	},
}
