package main

import "fmt"

func main() {
	fmt.Println("AFCB v1 started")

	//initialize empty Contacts slice
	var contacts Contacts

	fmt.Println("\nTest Add: Create new contact")
	newContact, err := contacts.New("Family", "Orm", "Korn", "ok@kshhh.co", "270-5200-227")
	if err != nil {
		fmt.Printf("Error to create new contact: %v\n", err)
	} else {
		fmt.Printf("New contact successfully added: %+v\n", newContact)
	}
}
