package main

import (
	"bytes"
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

//viewHandler is the basic http handler to display the form
func viewHandler(w http.ResponseWriter, r *http.Request) {
	filename := "blogIndex.html"
	body, _ := ioutil.ReadFile(filename)
	fmt.Fprintf(w, string(body))
}

//channel to send context to template
var ctxChan = make(chan (BlogCTX), 2)

//BlogCTX holds the context for the blog file
type BlogCTX struct {
	Title       string
	Description string
	Image       string
	Video       string
	Tags        string
	Categories  string
	Draft       string
	Date        string
	Body        string
}

// SaveHandler saves edited wiki data and reloads pages
func saveHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var blogData bytes.Buffer
	var ctx BlogCTX
	err := decoder.Decode(&ctx)
	if err != nil {
		fmt.Println(err)
	}
	//if the image != nil and starts with http
	//get the image, store it, and change the image name to match the created file, then send on the ctx channel
	if ctx.Image != "" {
		if ctx.Image[:4] == "http" {
			cmd := exec.Command("wget", "-O", "./static/images/"+ctx.Title+".jpg", ctx.Image)
			err = cmd.Start()
			if err != nil {
				fmt.Println(err)
			}
			ctx.Image = ctx.Title + ".jpg"
		}
	}
	ctx.Date = time.Now().Format("2006-01-02T15:04:05") //
	t, err := template.ParseFiles("./blogTemplate.t")
	if err != nil {
		fmt.Println("error parsing template: ", err)
	} //fill requestBody with the executed template and context
	err = t.Execute(&blogData, ctx)
	if err != nil {
		fmt.Println("Error executing template: ", err)
	}
	file, err := os.Create(`./content/project/` + ctx.Title + ".md")
	if err != nil {
		fmt.Println("error creating file: ", err)
	}
	file.Write(blogData.Bytes())
	file.Close()
	fmt.Println(blogData.String())
	fmt.Println("cat: ", ctx.Categories)

	time.Sleep(5000)
	cmd := exec.Command("hugo")
	err = cmd.Run()
	if err != nil {
		fmt.Println("error executing command: ", err)
	}
	blogData.Reset()
}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(listenPort, nil)
}
