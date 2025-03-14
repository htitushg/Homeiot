package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Not Found"

	// rendering the template
	app.render(w, r, http.StatusNotFound, "error.tmpl", tmplData)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Oooops"

	// setting the error title and message
	tmplData.Error.Title = "Error 405"
	tmplData.Error.Message = "Something went wrong!"

	// rendering the template
	app.render(w, r, http.StatusMethodNotAllowed, "error.tmpl", tmplData)
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Home"

	// rendering the template
	app.render(w, r, http.StatusOK, "home.tmpl", tmplData)
}

// Dashboard handler - renders the IoT dashboard page
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	// Retrieve basic template data
	tmplData := app.newTemplateData(r)
	tmplData.Title = "Home IoT - Dashboard"

	tmplData.Devices = []struct {
		DeviceID string
		Name     string
		Status   string
		Battery  string
	}{
		{"device123", "Smart Light", "Online", "85%"},
		{"device456", "Thermostat", "Offline", "N/A"},
	}

	// Render the dashboard template
	app.render(w, r, http.StatusOK, "dashboard.tmpl", tmplData)
}

// CommandDevice handler - allows sending a command to a specific IoT device
func (app *application) commandDevice(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST requests are allowed
	if r.Method != http.MethodPost {
		app.methodNotAllowed(w, r)
		return
	}

	// Parse device command from request body
	type CommandRequest struct {
		DeviceID string `json:"device_id"`
		Command  string `json:"command"`
	}

	var cmdReq CommandRequest
	err := json.NewDecoder(r.Body).Decode(&cmdReq)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Simulate sending the command to the device (e.g., via MQTT, API call, etc.)
	// For now, we just log it
	app.logger.Debug(fmt.Sprintf("Sending command '%s' to device '%s'", cmdReq.Command, cmdReq.DeviceID))

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Command sent"})
	if err != nil {
		return
	}
}

// GetDeviceInfo handler - retrieves device details
func (app *application) getDeviceInfo(w http.ResponseWriter, r *http.Request) {
	// Ensure only GET requests are allowed
	if r.Method != http.MethodGet {
		app.methodNotAllowed(w, r)
		return
	}

	// Get device ID from query parameters
	deviceID := r.URL.Query().Get("device_id")
	if deviceID == "" {
		http.Error(w, "Missing device_id parameter", http.StatusBadRequest)
		return
	}

	// Simulated device info
	deviceInfo := map[string]interface{}{
		"device_id": deviceID,
		"name":      "Smart Light",
		"status":    "Online",
		"battery":   "85%",
	}

	// Send device info as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(deviceInfo)
	if err != nil {
		return
	}
}
