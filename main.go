package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func validateEmail(email string) bool {

	if strings.Contains(strings.TrimSpace(email), " ") {
		return false
	}

	if !strings.Contains(email, "@") {
		return false
	}

	if match, _ := regexp.MatchString(`[a-z0-9]+@[a-z0-9]+\.[a-z]`, strings.TrimSpace(email)); !match {
		return false
	}

	return true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(" Go Email Verifier:\n")
	fmt.Print("=====================\n\n")

	fmt.Print("Enter an email address to verify : \t")
	for scanner.Scan() {
		email := scanner.Text()
		if validateEmail(email) {
			a := strings.Split(email, "@")
			checkDomain(email, a[1])
		} else {
			log.Fatal("Invalid Email")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v \n", err)

	}
}

func checkDomain(email, domain string) {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	spfRecord = "none"
	dmarcRecord = "none"

	mxRecord, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecord) > 0 {
		hasMX = true
	}

	txtRecord, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error : %v\n", err)
	}

	for _, record := range txtRecord {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		log.Printf("Error : %v\n", err)

	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	if hasDMARC && hasMX && hasSPF {
		fmt.Printf("%v is an email address \n", email)
	}
	fmt.Println("Results")
	fmt.Printf("domain : %v \n hasMx : %v \n hasSPF : %v \n SPF Records : %v \n hasDMARC : %v \n DMARC Records: %v \n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	fmt.Printf("\nYou may try another email address or CTRL-C to close the program:\t")
}
