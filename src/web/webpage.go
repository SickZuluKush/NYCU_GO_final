package main 

import (
	"log"
	"os"
	"strings"
	"net/http"
    "path/filepath"
    "go_final/db"
)


func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) != 1 {
		webpage.invalid_url.Execute(w, nil)
	} else {
		webpage.homepage.Execute(w, struct {
			BROWSE string
			UPLOAD string
		}{
			browse_url,
			upload_url,
		})
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		webpage.upload_page.Execute(w, struct {
			UPLOAD string
			FILENAME string
		}{
			upload_url,
			formfile,
		})
	} else if r.Method == "POST" {
		if ok, msg := upload_utility(w, r); ok {
			db.Reset_record(msg)
			webpage.upload_success.Execute(w, browse_url + msg)
		} else {
			webpage.upload_failed.Execute(w, struct {
				ERR_MSG string
				UPLOAD string
			}{
				msg,
				upload_url,
			})
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	p, err := filepath.Rel(browse_url, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if p == "." {
		files, err := os.ReadDir("./" + upload_folder)
		if err != nil {
			log.Print(err)
			return
		}
		
		var filenames []string
		
		for _, file := range(files) {
			filenames = append(filenames, file.Name())
		}
		
		webpage.browse_page.Execute(w, struct {
			BROWSE string
			FILES []string
		}{
			browse_url,
			filenames,
		})
	} else if _, err := os.Stat("./" + upload_folder + "/" + p); err == nil  {
		comments := db.Get_record(strings.ToLower(p))
		
		webpage.data_page.Execute(w, struct {
				DATA string
				BROWSE string
				COMMENT_URL string
				COMMENT string
				DELETE_URL string
				FILE_INFO string
				FILE string
				COMMENTS []string
		}{
			data_url + p,
			browse_url,
			comment_url,
			comment,
			delete_url,
			file_info,
			p,
			comments,
		})
	} else {
		webpage.invalid_url.Execute(w, nil)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == data_url {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	http.StripPrefix(data_url, http.FileServer(http.Dir("./" + upload_folder))).ServeHTTP(w, r)
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		if info, err := os.Stat("./" + upload_folder + "/" + r.FormValue(file_info)); err != nil {
			log.Print(err)
		} else if !info.IsDir() && r.FormValue(comment) != "" {
			db.Add_record(strings.ToLower(r.FormValue(file_info)), r.FormValue(comment))			
		}
		
		http.Redirect(w, r, browse_url + r.FormValue(file_info), http.StatusSeeOther)		
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		if err := os.Remove("./" + upload_folder + "/" + r.FormValue(file_info)); err != nil {
			log.Print("file remove err")
		}
		
		db.Reset_record(strings.ToLower(r.FormValue(file_info)))
		http.Redirect(w, r, browse_url, http.StatusSeeOther)
	}
}

func setup_web() {
	http.HandleFunc("/", homepageHandler)
    http.HandleFunc(upload_url, uploadHandler)
    http.HandleFunc(browse_url, browseHandler)
    http.HandleFunc(data_url, dataHandler)
    http.HandleFunc(comment_url, commentHandler)
    http.HandleFunc(delete_url, deleteHandler)
    
    log.Fatal(http.ListenAndServe(":" + port, nil))
}
