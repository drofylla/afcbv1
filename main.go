package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

var contacts Contacts

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

const dataFile = "AFcb.json"

var conCard = template.Must(template.New("card").Parse(`
	<!-- contact-card.html -->
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
            onclick="openEditModal('{{.ID}}', '{{.ContactType}}', '{{.FirstName}}', '{{.LastName}}', '{{.Email}}', '{{.Phone}}')"
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

func renderCard(w http.ResponseWriter, c Contact) {
	w.Header().Set("Content-Type", "text/html")
	conCard.Execute(w, c)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Printf("=== GET /contacts called ===\n")
	fmt.Printf("Returning %d contacts to client\n", len(contacts))

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

func main() {
	fmt.Println("AFcb started at http://localhost:1330")
}
