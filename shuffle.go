// Command shuffle -yamlFile=<path_to_file> reads the data from an input YAML file, shuffles the array of provided names and sends out an email with a random quote of the day.
// Can be used for the daily meeting order.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// the details needed to run the script (read from the input YAML file)
type emailDetails struct {
	SMTPAddress string   // the SMTP Server:Port to use to send the email
	Sender      string   // the email used as the sender
	Names       []string // the names to be shuffled and included in the email
	Receivers   []string // the email addresses to which to send the email
}

// struct for the random quote of the day
type quoteData struct {
	Quote  string
	Author string
}

// parseArgs parses the CLI arguement that provides the path to the yaml file with the details for the script
func parseArgs() *string {
	yamlFile := flag.String("yamlPath", "email_details.yaml", "Path to the YAML file with the details for the email.")
	flag.Parse()
	return yamlFile
}

// shuffleNames takes a slice of strings (i.e. the names) and shuffles the elements in-place
func shuffleNames(names []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })
}

// getRandomQuote returns a random quote of the day from https://quotes.rest/qod
// returned data is a map[string]string with keys "quote" and "author"
func getRandomQuote() quoteData {

	log.Println("Getting random quote")

	// fallback function to return in case of anything goes wrong (probably not the best way of doing this)
	failedQuote := func() quoteData {
		log.Println("Get quote request failed")
		return quoteData{
			Quote:  "API request failed, no quote today, guys...",
			Author: "The try except clause",
		}
	}

	// query the random quotes site
	resp, err := http.Get("https://quotes.rest/qod")
	if err != nil {
		log.Println(err)
		return failedQuote()
	}
	defer resp.Body.Close()

	// read all packets from the response's body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return failedQuote()
	}

	// make the container that will hold the decoded data
	decodedData := make(map[string]map[string][]quoteData) //just using map[string]interface{} also works, but don't know which is best
	if err := json.Unmarshal(data, &decodedData); err != nil {
		// BUG(me): This returns an error everytime and I still can't understand why
		log.Println(err)
		// return failedQuote()
	}

	quote := decodedData["contents"]["quotes"][0]

	log.Println("Get quote request completed successfully")

	return quote
}

// buildEmail creates an email message (as a slice of bytes) in the format expected by smtp.SendMail.
// subject is used in the Subject header and names are used to build the body
func buildEmail(subject string, names []string, quote quoteData) []byte {
	var b strings.Builder

	fmt.Fprintf(&b, "Subject: %s\r\n\r\nToday's daily order will be the following:\r\n\r\n", subject)

	// list out each name on a new line
	for _, name := range names {
		fmt.Fprintf(&b, "%s \r\n", name)
	}

	// add the random quote to the ending
	fmt.Fprintf(&b, "\r\n\r\nRandom quote of the day:\r\n\"%s\"\r\n- %s\r\n\r\n", quote.Quote, quote.Author)
	fmt.Fprintf(&b, "Regards,\r\nThe random generator")

	msg := []byte(b.String())

	return msg
}

// sendEmail wraps smtp.SendMail to send an email
func sendEmail(smtpAddress string, senderEmail string, receiverEmails []string, msg []byte, retry int) error {

	log.Println("Sending email")

	if err := smtp.SendMail(smtpAddress, nil, senderEmail, receiverEmails, msg); err != nil {
		log.Println(err)
		// retry a couple of times just in case
		if retry == 3 {
			return err
		}
		retry++
		return sendEmail(smtpAddress, senderEmail, receiverEmails, msg, retry)

	}

	log.Println("Email sent successfully")
	return nil

}

// readEmailDetailsYaml opens the given YAML file and unmarshals its contents to a emailDetails struct which it returns
func readEmailDetailsYaml(yamlFile *string) *emailDetails {

	log.Printf("Reading YAML file: %s \r\n", *yamlFile)

	yamlData, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		log.Fatalf("File not found: %s \r\n", *yamlFile)
	}

	decodedData := emailDetails{}
	if err := yaml.Unmarshal(yamlData, &decodedData); err != nil {
		log.Fatalf("Couldn't decode the YAML file: %s\r\n", err)
	}

	log.Println("YAML file read successfully")
	return &decodedData
}

func main() {

	// get path to yaml from args
	yamlFile := parseArgs()

	// process the yaml file
	emailData := readEmailDetailsYaml(yamlFile)

	// shuffle the names in place
	shuffleNames(emailData.Names)

	// get the random quote
	randomQuote := getRandomQuote()

	// build the email message
	emailMessage := buildEmail("Daily Name Generator", emailData.Names, randomQuote)

	// send the email
	sendEmail(emailData.SMTPAddress, emailData.Sender, emailData.Receivers, emailMessage, 0)

}
