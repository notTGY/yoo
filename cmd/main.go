package main

import (
	"github.com/nottgy/yoo"
	"log"
	"os"
)

func main() {
	id := os.Getenv("ID")
	c, _ := yoo.NewClient(os.Getenv("YOOKASSA_SHOP_ID"), os.Getenv("YOOKASSA_SECRET_KEY"))

	if id != "" {
		r, _ := c.GetPayment(id)
		log.Println(r)
		return
	}

	id, crt, err := c.CreatePayment(9900, "example@example.com", "тестовый платеж", "hello-world")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ID: %s\nCRT: %s\n", id, crt)
}
