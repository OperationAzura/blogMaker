package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	//fmt.Println(ctx.Image)
	fmt.Println("is image blank?")
	if ctx.Image != "" {
		fmt.Println(ctx.Image[:4])
		if ctx.Image[:4] == "http" {
			fmt.Println("http was found on image")
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

			ctx.Image = strings.Replace(ctx.Title, " ", "", -1) + ".jpg"
		} else {
			fmt.Println("http not detected")
			fmt.Println("splitting string of base")
			//imgBase64 := strings.SplitAfter(ctx.Image, "base64")
			fmt.Println("Standard base64 decode")
			//fmt.Println("len of split base64 slice: ", len(imgBase64))
			//fmt.Println("len of base64 string to be used: ", len(imgBase64[1]))
			//fmt.Println("pice cut off: ", imgBase64[1][1:16])
			var imgBytes []byte
			//ctx.Image = strings.Replace(ctx.Image, `/`, "", -1)
			var i = 0
			var success = false
			for i < 20 && !success {
				fmt.Println("Starting Iteration: ", i)
				imgBytes, err = base64.RawURLEncoding.DecodeString(ctx.Image[i:]) //base64.StdEncoding.DecodeString(imgBase64[1])
				if err != nil {
					fmt.Println("RawURL decode error", err)
					imgBytes, err = base64.URLEncoding.DecodeString(ctx.Image[i:]) //base64.StdEncoding.DecodeString(imgBase64[1])
					if err != nil {
						fmt.Println("URL decode error", err)
						imgBytes, err = base64.StdEncoding.DecodeString(ctx.Image[i:])
						if err != nil {
							fmt.Println("Std decode error", err)
							imgBytes, err = base64.RawStdEncoding.DecodeString(ctx.Image[i:])
							if err != nil {
								fmt.Println("RawStd decode error", err)
							} else {
								fmt.Println("success on RawStd")
								success = true
							}
						} else {
							fmt.Println("success on Std")
							success = true
						}
					} else {
						fmt.Println("success on URL")
						success = true
					}
				} else {
					fmt.Println("success on RawURL")
					success = true
				}
				i++
			}
			fmt.Println("Attempging to build new reader")
			imgReader := bytes.NewReader(imgBytes)
			fmt.Println("decoding into image object")
			imgJpg, err := jpeg.Decode(imgReader)
			if err != nil {
				fmt.Println("Error decoding to jpeg: ", err)
			}
			fmt.Println("opening file")
			f, err := os.OpenFile("./static/images/"+strings.Replace(ctx.Title, " ", "", -1)+".jpg", os.O_WRONLY|os.O_CREATE, 0777)
			if err != nil {
				fmt.Println("Error opening file for jpeg: ", err)
			}
			fmt.Println("Encoding to file")
			_ = imgBytes
			jpeg.Encode(f, imgJpg, nil)
			ctx.Image = strings.Replace(ctx.Title, " ", "", -1) + ".jpg"
		}
	}
	t, err := template.New("blogTemplate.t").Funcs(NextIndex).ParseFiles("./blogTemplate.t")
	if err != nil {
		fmt.Println("error parsing template: ", err)
	} //fill requestBody with the executed template and context
	err = t.Execute(&blogData, ctx)
	if err != nil {
		fmt.Println("Error executing template: ", err)
	}
	file, err := os.Create(`./content/project/` + strings.Replace(ctx.Title, " ", "", -1) + ".md")
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
