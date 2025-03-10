package main

import "net/http"

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Not Found"
	
	// rendering the template
	app.render(w, r, http.StatusOK, "error.tmpl", tmplData)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Oooops"
	
	// setting the error title and message
	tmplData.Error.Title = "Error 405"
	tmplData.Error.Message = "Something went wrong!"
	
	// rendering the template
	app.render(w, r, http.StatusOK, "error.tmpl", tmplData)
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	
	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Home"
	
	// rendering the template
	app.render(w, r, http.StatusOK, "home.tmpl", tmplData)
}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {

}
func (app *application) commandDevice(w http.ResponseWriter, r *http.Request) {

}
func (app *application) getDeviceInfo(w http.ResponseWriter, r *http.Request) {

}
