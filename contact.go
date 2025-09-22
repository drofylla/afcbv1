package main

import (
	"errors"
	"fmt"
	"strings"

	gonanoid "github.com/matoous/go-nanoid"
)

// struct for contact details
type Contact struct {
	ID          string
	ContactType string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
}

// slice of Contact structs
type Contacts []Contact

// generate unique 6-character ID using custom alphabet & numbers
func genID() (string, error) {
	id, err := gonanoid.Generate("drofylla12301993", 6)
	if err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}
	return id, nil
}

// create new Contact
func (c *Contacts) New(contactType, firstName, lastName, email, phone string) (Contact, error) {
	//gen new ID for new Contact
	id, err := genID()
	if err != nil {
		return Contact{}, errors.New("unable to generate ID: " + err.Error())
	}

	//create new contact
	contact := Contact{
		ID:          id,
		ContactType: contactType,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Phone:       phone,
	}

	//append contacts slice
	*c = append(*c, contact)

	//return new contact
	return contact, nil
}

func (c *Contacts) Save(id, contactType, firstName, lastName, email, phone string) error {
	if id == "" {
		//gen new id
		newID, err := genID()
		if err != nil {
			return errors.New("fail to generate ID for new contact")
		}

		contact := Contact{
			ID:          newID,
			ContactType: contactType,
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			Phone:       phone,
		}
		*c = append(*c, contact)
		return nil
	}

	//iterate through slice to find existing contact
	for i := range *c {
		if (*c)[i].ID == id {
			//update field of existing contact
			(*c)[i].ContactType = contactType
			(*c)[i].FirstName = firstName
			(*c)[i].LastName = lastName
			(*c)[i].Email = email
			(*c)[i].Phone = phone
			return nil
		}
	}
	return errors.New("contact not found")
}

func (c *Contacts) Update(id string, updates map[string]string) error {
	fmt.Printf("Searching for contact with ID: %s\n", id)
	fmt.Printf("Available contacts: %+v\n", *c)

	for i := range *c {
		if (*c)[i].ID == id {
			fmt.Printf("Found contact %+v\n", (*c)[i])
			for field, value := range updates {
				key := strings.ToLower(strings.ReplaceAll(field, " ", ""))
				switch key {
				case "contacttype":
					(*c)[i].ContactType = value
				case "firstname":
					(*c)[i].FirstName = value
				case "lastname":
					(*c)[i].LastName = value
				case "email":
					(*c)[i].Email = value
				case "phone":
					(*c)[i].Phone = value
				default:
					return fmt.Errorf("Invalid field: %s\n", field)
				}
			}
			fmt.Printf("Contact info updated: %+v\n", (*c)[i])
			return nil
		}
	}
	return errors.New("Unable to update contact info due to no ID found")
}

func (c *Contacts) Delete(id string) error {
	for i, contact := range *c {
		if contact.ID == id {
			*c = append((*c)[:i], (*c)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("contact id %s not found", id)
}
