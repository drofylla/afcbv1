package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

var contacts Contacts

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

const dataFile = "AFcb.json"

var conCard = template.Must(template.New("card").Parse(`
    <div class="card bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition-all duration-300" id="contact-{{.ID}}">
    <div class="details">
        <span class="id text-xs font-semibold text-gray-500">ID: {{.ID}}</span>
        <strong class="name block text-xl font-bold text-gray-800 mt-1">{{.FirstName}} {{.LastName}}</strong>
        <span class="type inline-block mt-2 px-3 py-1 rounded-full text-sm font-medium
            {{if eq .ContactType "Personal"}}bg-blue-100 text-blue-800
            {{else if eq .ContactType "Work"}}bg-green-100 text-green-800
            {{else if eq .ContactType "Family"}}bg-purple-100 text-purple-800
            {{else}}bg-gray-100 text-gray-800{{end}}">
            {{.ContactType}}
        </span>
        <div class="details mt-3 text-gray-600">
            <div class="flex items-center mb-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                {{.Email}}
            </div>
            <div class="flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
                {{.Phone}}
            </div>
        </div>
    </div>
    <div class="actions flex justify-end mt-4 space-x-2">
    		<button class="edit-btn p-2 rounded-lg border border-gray-300 hover:border-blue-500 hover:bg-blue-50 transition-colors"
                hx-get="/modal/edit/{{.ID}}"
                hx-target="#modal-container"
                hx-swap="innerHTML"
                title="Edit">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 20h9"/>
                    <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/>
                </svg>
            </button>
        <button class="delete-btn p-2 rounded-lg border border-gray-300 hover:border-red-500 hover:bg-red-50 transition-colors"
                hx-delete="/contacts/{{.ID}}"
                hx-target="#contact-{{.ID}}"
                hx-swap="outerHTML"
                hx-confirm="Are you sure you want to delete this contact?"
                title="Delete">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6m5 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/>
                <line x1="10" y1="11" x2="10" y2="17"/>
                <line x1="14" y1="11" x2="14" y2="17"/>
            </svg>
        </button>
    </div>
</div>
`))

var addModalHTML = `
<div id="contact-modal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="flex justify-end">
            <button hx-target="#contact-modal" hx-swap="outerHTML" hx-get="/modal/close" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <h3 class="text-xl font-bold mb-4">Add New Contact</h3>
        <form id="contactForm"
              hx-post="/contacts"
              hx-target="#contact-list"
              hx-swap="afterbegin"
              hx-on::after-request="if(event.detail.successful) htmx.remove(htmx.find('#contact-modal'))">
            <input type="hidden" id="contact-id" name="id">
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="contactType">Contact Type</label>
                <select id="contactType" name="ContactType" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
                    <option value="Personal">Personal</option>
                    <option value="Work">Work</option>
                    <option value="Family">Family</option>
                </select>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="firstName">First Name</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="firstName" name="FirstName" type="text" placeholder="First Name" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="lastName">Last Name</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="lastName" name="LastName" type="text" placeholder="Last Name" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="email">Email</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" name="Email" type="email" placeholder="Email" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="phone">Phone</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="phone" name="Phone" type="tel" placeholder="Phone" required>
            </div>
            <div class="flex items-center justify-end">
                <button type="button" hx-target="#contact-modal" hx-swap="outerHTML" hx-get="/modal/close" class="bg-gray-500 text-white font-bold py-2 px-4 rounded-lg shadow-md hover:bg-gray-600 transition-colors duration-300 mr-2">Cancel</button>
                <button type="submit" class="bg-blue-600 text-white font-bold py-2 px-4 rounded-lg shadow-md hover:bg-blue-700 transition-colors duration-300">Save Contact</button>
            </div>
        </form>
    </div>
</div>
`

var editModalHTML = `
<div id="contact-modal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="flex justify-end">
            <button hx-target="#contact-modal" hx-swap="outerHTML" hx-get="/modal/close" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <h3 class="text-xl font-bold mb-4">Edit Contact</h3>
        <form id="contactForm"
              hx-put="/contacts/{{.ID}}"
              hx-target="#contact-{{.ID}}"
              hx-swap="outerHTML"
              hx-on::after-request="if(event.detail.successful) htmx.remove(htmx.find('#contact-modal'))">
            <input type="hidden" name="id" value="{{.ID}}">
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="contactType">Contact Type</label>
                <select id="contactType" name="ContactType" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
                    <option value="Personal" {{if eq .ContactType "Personal"}}selected{{end}}>Personal</option>
                    <option value="Work" {{if eq .ContactType "Work"}}selected{{end}}>Work</option>
                    <option value="Family" {{if eq .ContactType "Family"}}selected{{end}}>Family</option>
                </select>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="firstName">First Name</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="firstName" name="FirstName" type="text" value="{{.FirstName}}" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="lastName">Last Name</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="lastName" name="LastName" type="text" value="{{.LastName}}" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="email">Email</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" name="Email" type="email" value="{{.Email}}" required>
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="phone">Phone</label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="phone" name="Phone" type="tel" value="{{.Phone}}" required>
            </div>
            <div class="flex items-center justify-end">
                <button type="button" hx-target="#contact-modal" hx-swap="outerHTML" hx-get="/modal/close" class="bg-gray-500 text-white font-bold py-2 px-4 rounded-lg shadow-md hover:bg-gray-600 transition-colors duration-300 mr-2">Cancel</button>
                <button type="submit" class="bg-blue-600 text-white font-bold py-2 px-4 rounded-lg shadow-md hover:bg-blue-700 transition-colors duration-300">Save Changes</button>
            </div>
        </form>
    </div>
</div>
`

func renderCard(w http.ResponseWriter, c Contact) {
	w.Header().Set("Content-Type", "text/html")
	conCard.Execute(w, c)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Printf("=== GET /contacts called ===\n")
	fmt.Printf("Returning %d contacts to client\n", len(contacts))

	// Check if the contacts slice is empty
	if len(contacts) == 0 {
		fmt.Printf("No contacts found, returning empty message\n")
		fmt.Fprintf(w, `<div class="flex items-center justify-center p-8 bg-gray-100 text-gray-500 rounded-lg shadow-md">
            No contacts found. Add your first contact!
        </div>`)
		return
	}

	cardRendered := 0
	for _, c := range contacts {
		fmt.Printf("Rendering contact: %s %s (ID: %s)\n", c.FirstName, c.LastName, c.ID)
		if err := conCard.Execute(w, c); err != nil {
			fmt.Printf("Error rendering contact %s: %v\n", c.ID, err)
			http.Error(w, "Error rendering card: "+err.Error(), http.StatusInternalServerError)
			return
		}
		cardRendered++
	}
	fmt.Printf("Successfully rendered %d contacts\n", cardRendered)
	fmt.Printf("=== END GET /contacts ===\n")
}

func addContact(w http.ResponseWriter, r *http.Request) {
	// make sure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	contactType := r.FormValue("ContactType")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")
	phone := r.FormValue("Phone")

	fmt.Printf("Received form date - ID: '%s', Type: '%s', Name: '%s %s', Email: '%s', Phone: '%s'\n",
		id, contactType, firstName, lastName, email, phone)

	if id != "" {
		//update existing contact
		updates := map[string]string{
			"ContactType": contactType,
			"FirstName":   firstName,
			"LastName":    lastName,
			"Email":       email,
			"Phone":       phone,
		}

		if err := contacts.Update(id, updates); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		//find and return updated contact
		for _, c := range contacts {
			if c.ID == id {
				if err := contacts.SaveToFile(dataFile); err != nil {
					http.Error(w, "Fail to save contacts: "+err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Printf("Contact updated and save to file: %s\n", dataFile)
				renderCard(w, c)
				return
			}
		}
		http.Error(w, "Contact not found after update", http.StatusNotFound)
		return
	}

	//email validation
	if !emailRegex.MatchString(email) {
		http.Error(w, "Invalid email address input", http.StatusBadRequest)
		return
	}

	//use New method
	newContact, err := contacts.New(contactType, firstName, lastName, email, phone)
	if err != nil {
		http.Error(w, "Fail to create contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("New contact created with ID: %s. Total contact: %d\n", newContact.ID, len(contacts))

	//save to file
	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "unable to save contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Succesfully saved contact to file: %s\n", dataFile)

	renderCard(w, newContact)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	//make sure it's PUT or PATCH request
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("UPDATE request received for id: %s, form values: %+v\n", id, r.Form)

	// Validate required fields
	contactType := r.FormValue("ContactType")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")
	phone := r.FormValue("Phone")

	if contactType == "" || firstName == "" || lastName == "" || email == "" || phone == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Email validation
	if !emailRegex.MatchString(email) {
		http.Error(w, "Invalid email address format", http.StatusBadRequest)
		return
	}

	updates := map[string]string{
		"ContactType": r.FormValue("ContactType"),
		"FirstName":   r.FormValue("FirstName"),
		"LastName":    r.FormValue("LastName"),
		"Email":       r.FormValue("Email"),
		"Phone":       r.FormValue("Phone"),
	}

	fmt.Printf("Attempting to update contact %s with: %+v\n", id, updates)

	// find contact to ensure it exists
	var found bool
	for i := range contacts {
		if contacts[i].ID == id {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	//update contact
	if err := contacts.Update(id, updates); err != nil {
		fmt.Println("Update error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//save to file
	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "Fail to save contact to file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//return updated contact
	for _, c := range contacts {
		if c.ID == id {
			fmt.Printf("Successfully update contact: %+v\n", c)
			renderCard(w, c)
			return
		}
	}

	http.Error(w, "contact not found after update", http.StatusNotFound)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println("DELETE request received for id:", id)

	if err := contacts.Delete(id); err != nil {
		fmt.Println("Delete error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "failed to save contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return empty content - HTMX remove element
	w.WriteHeader(http.StatusOK)
}

func searchContacts(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	fmt.Printf("Search request received for keyword: '%s'\n", keyword) //log

	w.Header().Set("Content-type", "text/html")

	if keyword == "" {
		fmt.Println("No keyword provided, returning all contacts") //log

		for _, c := range contacts {
			if err := conCard.Execute(w, c); err != nil {
				http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		return
	}

	results := contacts.Search(keyword)
	fmt.Printf("Found %d results for keyword '%s'\n", len(results), keyword) //log

	if len(results) == 0 {
		fmt.Fprintf(w, `<div class="no-results">No contacts found for "%s"</div>`, template.HTMLEscapeString(keyword))
		return
	}

	for _, c := range results {
		if err := conCard.Execute(w, c); err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ALL MODAL RELATED //
// add modal render the add contact form modal
func addModal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("modal").Parse(addModalHTML))
	tmpl.Execute(w, nil)
}

func editModal(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	contact, err := contacts.Find(id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("edit-modal").Parse(editModalHTML))
	tmpl.Execute(w, contact)
}

func closeForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func main() {
	//initialize contacts
	contacts = Contacts{}

	//load contacts from file
	if err := contacts.LoadContacts(dataFile); err != nil {
		fmt.Printf("Error loading contacts: %v\n", err)
		fmt.Println("Starting with empty contacts list")
		contacts = Contacts{}
	}

	router := mux.NewRouter()

	//serve main html file
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	//add contact API endpoints
	router.HandleFunc("/contacts", getContacts).Methods("GET")
	router.HandleFunc("/contacts", addContact).Methods("POST")
	router.HandleFunc("/modal/add", addModal).Methods("GET")
	router.HandleFunc("/modal/edit/{id}", editModal).Methods("GET")
	router.HandleFunc("modal/close", closeForm).Methods("GET")
	router.HandleFunc("/contacts/{id}", updateContact).Methods("PUT", "PATCH")
	router.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")
	router.HandleFunc("/search", searchContacts).Methods("GET")

	//server start
	fmt.Println("AFcb started at http://localhost:1330")
	http.ListenAndServe(":1330", router)
}
