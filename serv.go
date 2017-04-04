package main

import (
	"os/exec"
		"time"
	"os"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
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
	ctx.Date = time.Now().Format("2017-03-18T22:18:08-05:00")


		if ctx.Image != ""{
			if ctx.Image[4:] == "http"{
				cmd := exec.Command("wget", "-O", "./static/images/"+ctx.Title+".jpg", ctx.Image)
				err := cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Waiting for command to finish...")
				err = cmd.Wait()
				if err != nil{
					fmt.Printf("Command finished with error: %v", err)
				}
				ctx.Image = ctx.Title+".jpg"
			}
		}

	t, err := template.ParseFiles("./blogTemplate.t")
	if err != nil {
		fmt.Println("error parsing template: ", err)
	} //fill requestBody with the executed template and context
	err = t.Execute(&blogData, ctx)
	if err != nil {
		fmt.Println("Error executing template: ", err)
	}
	file, err := os.Create("./content/projects/"+ctx.Title[1:len(ctx.Title)-1]+".md")
	if err != nil{
		fmt.Println("error creating file: ", err)
	}
	 file.Write(blogData.Bytes())
	 file.Close()
	fmt.Println(blogData.String())
	fmt.Println("cat: ", ctx.Categories)
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
