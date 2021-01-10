package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
)

var conf Config

// Config struct for app config
type Config struct {
	// Email is the struct for Sender and Receiver Data
	Email struct {
		Receiver struct {
			// Receiver Address is a list/map of emails to receive the notification
			Address []string `json:"address"`
		} `json:"receiver"`
		Sender struct {
			// Address is the email you are sending data from
			Address string `json:"address"`
			// Password is the password for the senders email account
			Password string `json:"password"`
		} `json:"sender"`
	} `json:"email"`

	// Items is the list/map of all items you want to watch
	Items  []string `json:"items"`
	Limits struct {
		Price struct {
			// Max is the maximum price you want to pay
			Max int `json:"max"`
			// Min is the minimum price you want to pay
			Min int `json:"min"`
		} `json:"price"`
	} `json:"limits"`
}

func init() {
	// load config.json file
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatal("Can't load config.json file with item numbers and email addresses.")
	}

	// unmarshal configs to Config struct
	json.Unmarshal(file, &conf)

	log.Println("Configs have been successfully loaded.")
}

// Payload for web request data
type Payload struct {
	MainItem struct {
		Description struct {
			Title string `json:"Title"`
		} `json:"Description"`
		// InStock is the boolean value of it the item is in stock
		Instock bool `json:"Instock"`

		// FinalPrice shows the price after all discounts are applied
		FinalPrice    float64 `json:"FinalPrice"`
		StockCount    int     `json:"Stock"`
		ItemNumber    string  `json:"ItemNumber"`
		AddToCartType int     `json:"AddToCartType"`
	} `json:"MainItem"`
	Additional struct {
		LimitQuantity int `json:"LimitQuantity"`
	} `json:"Additional"`
}

func main() {
	log.Println("Starting inventory search...")
	// loop for items in config to build and execute http requests

	for _, item := range conf.Items {

		weburl := fmt.Sprintf("https://www.newegg.com/product/api/ProductRealtime?ItemNumber=%v", item)
		client := &http.Client{}

		req, err := http.NewRequest("GET", weburl, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		data := &Payload{}
		json.NewDecoder(resp.Body).Decode(&data)
		defer resp.Body.Close()

		// extra request error checking because newegg doesn't return anything other than 200s -_-
		if string(data.MainItem.Description.Title) == "null" {
			log.Println("- " + item + " request error. Does " + weburl + " exist?")
			break
		}

		priceint, _ := math.Modf(data.MainItem.FinalPrice)

		// make sure items meet price limits requirements
		if int(priceint) > conf.Limits.Price.Max || int(priceint) < conf.Limits.Price.Min {
			log.Println(fmt.Sprintf("- $%0.2f does not meet the price requirements.  %v in stock. %v", data.MainItem.FinalPrice, data.MainItem.StockCount, weburl))
			break
		}

		// if its in stock then send email, AddToCartType 0 is available for cart
		if data.MainItem.Instock && data.MainItem.AddToCartType == 0 {
			log.Println(fmt.Sprintf("- [IN STOCK] - %v Total.  %v", data.Additional.LimitQuantity, weburl))
			sendMail(data.MainItem.Description.Title, weburl, fmt.Sprintf("%0.2f", data.MainItem.FinalPrice), int(data.MainItem.StockCount), data.Additional.LimitQuantity)
		} else {
			log.Println("- [NOT IN STOCK] - " + weburl)
		}

	}

	log.Println("Search complete.")
}

// sendMail sends the notification email using Gmail.
func sendMail(title, url, price string, total, limit int) {
	from := conf.Email.Sender.Address
	pass := conf.Email.Sender.Password
	to := conf.Email.Receiver.Address

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ", ") + "\n" +
		"Subject: NEWEGG-WATCHER | IN STOCK!\n\n" +
		"Url: " + url + "\n\n" +
		"Title: " + title + "\n" +
		"Price: " + price + "\n" +
		"Limit: " + strconv.Itoa(limit) + "\n" +
		"Total: " + strconv.Itoa(total) + "\n\n\n\n\n" +
		"- Sent using EggBtr"

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, []byte(msg))

	if err != nil {
		log.Printf("Email smtp error: %s", err)
	} else {
		log.Println("Email successfully sent to " + strings.Join(to, ", "))
	}

	return
}
