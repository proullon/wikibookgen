package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	wikibookgen "github.com/proullon/wikibookgen/api/sdk"
)

type test struct {
	Name string
	Test func() error
}

var tt = []test{
	{
		Name: "Order",
		Test: testOrder,
	},
}

func main() {
	var endpoint string

	endpoint = os.Getenv("WIKIBOOKGEN_API")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8090"
	}

	wikibookgen.Initialize(endpoint, "")

	for _, t := range tt {
		err := t.Test()
		if err != nil {
			log.Errorf("Test %s: Failed (%s)\n", t.Name, err)
		} else {
			log.Infof("Test %s: OK\n", t.Name)
		}
	}

}

func testOrder() error {
	subject := `Mathématiques`
	model := `tour`

	// test invalid model
	orderID, err := wikibookgen.Order(subject, `invalidmodel`)
	if err == nil {
		return fmt.Errorf("Order(%s, invalidmodel): should have error", subject)
	}

	// test invalid subject
	orderID, err = wikibookgen.Order(`invalidsubject`, model)
	if err == nil {
		return fmt.Errorf("Order(invalidsubject, %s): should have error", model)
	}

	// test nominal
	orderID, err = wikibookgen.Order(subject, model)
	if err != nil {
		return fmt.Errorf("Order(%s): %s", subject, err)
	}
	if orderID == "" {
		return fmt.Errorf("Order(%s): invalid uuid", subject)
	}

	return nil
}
