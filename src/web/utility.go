package main

import (
	"net/http"
	"os"
	"io"
	"log"
	"strings"
	"strconv"
    "html/template"
	"github.com/joho/godotenv"
)


var (
	port string
	formfile string
	comment string
	file_info string
	upload_url string
	browse_url string
	data_url string
	comment_url string
	delete_url string
	upload_folder string
	template_folder string
	max_upload_size int64
	webpage webpage_template
)

type webpage_template struct {
	homepage *template.Template
	invalid_url *template.Template
	upload_page *template.Template
	upload_success *template.Template
	upload_failed *template.Template
	browse_page *template.Template
	data_page *template.Template
}


func init() {
	var ok bool
	
	if err := godotenv.Load(); err != nil {
		log.Print("load .env error")
	}
	
	if port, ok = os.LookupEnv("PORT"); !ok {
		log.Print(".env error: port")
		port = "8080"
	}
	if formfile, ok = os.LookupEnv("FORMFILE"); !ok {
		log.Print(".env error: formfile")
		formfile = "html_formfile"
	}
	if comment, ok = os.LookupEnv("COMMENT"); !ok {
		log.Print(".env error: comment")
		comment = "html_comment"
	}
	if file_info, ok = os.LookupEnv("FILE_INFO"); !ok {
		log.Print(".env error: file info")
		file_info = "html_file_info"
	}
	if upload_url, ok = os.LookupEnv("UPLOAD_URL"); !ok {
		log.Print(".env error: upload url")
		upload_url = "/upload_url"
	}
	if browse_url, ok = os.LookupEnv("BROWSE_URL"); !ok {
		log.Print(".env error: browse url")
		browse_url = "/browse_url/"
	}
	if data_url, ok = os.LookupEnv("DATA_URL"); !ok {
		log.Print(".env error: data url")
		data_url = "/data_url/"
	}
	if comment_url, ok = os.LookupEnv("COMMENT_URL"); !ok {
		log.Print(".env error: comment url")
		comment_url = "/comment_url"
	}
	if delete_url, ok = os.LookupEnv("DELETE_URL"); !ok {
		log.Print(".env error: delete url")
		delete_url = "/delete_url"
	}
	if upload_folder, ok = os.LookupEnv("UPLOAD_FOLDER_NAME"); !ok {
		log.Print(".env error: upload folder")
		upload_folder = "upload_folder"
	}
	if template_folder, ok = os.LookupEnv("TEMPLATE_FOLDER_NAME"); !ok {
		log.Print(".env error: template folder")
		template_folder = "templates"
	}
	
	if mb, err :=  strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE")); err != nil {
		log.Print(".env error: max upload size")
		max_upload_size = 1024 * 1024
	} else {
		max_upload_size = int64(mb * 1024 * 1024)
	}
	
	preload_template()
	
	if err := os.MkdirAll(upload_folder, os.ModePerm); err != nil {
		log.Fatal("mkdir err: upload folder")
	}
}

func preload_template() {
	webpage.homepage = template.Must(template.ParseFiles(template_folder + "/homepage.html"))
	webpage.invalid_url = template.Must(template.ParseFiles(template_folder + "/invalid_url.html"))
	webpage.upload_page = template.Must(template.ParseFiles(template_folder + "/upload_page.html"))
	webpage.upload_success = template.Must(template.ParseFiles(template_folder + "/upload_success.html"))
	webpage.upload_failed = template.Must(template.ParseFiles(template_folder + "/upload_failed.html"))
	webpage.browse_page = template.Must(template.ParseFiles(template_folder + "/browse_page.html"))
	webpage.data_page = template.Must(template.ParseFiles(template_folder + "/data_page.html"))
}

func upload_utility(w http.ResponseWriter, r *http.Request) (bool, string) {
	var err error
	
	// limit file size
	r.Body = http.MaxBytesReader(w, r.Body, max_upload_size)
	if err = r.ParseMultipartForm(max_upload_size); err != nil {
		return false, err.Error()
	}
	
	file, fileHeader, err := r.FormFile(formfile)
	if err != nil {
		return false, err.Error()
	}
	
	defer file.Close()
	
	// store file
	var filename = strings.ToLower(fileHeader.Filename)
	dst, err := os.Create(upload_folder  + "/" + filename)
	if err != nil {
		return false, err.Error()
	}
	
	defer dst.Close()
	
	if _, err = io.Copy(dst, file); err != nil {
		return false, err.Error()
	}
	
	return true, filename
}
