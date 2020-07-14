package main

import (
	"fmt"
	"os"

	. "github.com/proullon/wikibookgen/api/model"

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
	{
		Name: "ListBook",
		Test: testListBook,
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
	lang := `fr`

	// test invalid model
	orderID, err := wikibookgen.Order(subject, `invalidmodel`, lang)
	if err == nil {
		return fmt.Errorf("Order(%s, invalidmodel): should have error", subject)
	}

	// test invalid subject
	orderID, err = wikibookgen.Order(`invalidsubject`, model, lang)
	if err == nil {
		return fmt.Errorf("Order(invalidsubject, %s): should have error", model)
	}

	// test nominal
	orderID, err = wikibookgen.Order(subject, model, lang)
	if err != nil {
		return fmt.Errorf("Order(%s): %s", subject, err)
	}
	if orderID == "" {
		return fmt.Errorf("Order(%s): invalid uuid", subject)
	}

	fakeOrderID := "invaliduuid"
	_, _, err = wikibookgen.OrderStatus(fakeOrderID)
	if err == nil {
		return fmt.Errorf("Expected error with OrderStatus(%s)", fakeOrderID)
	}

	status, uuid, err := wikibookgen.OrderStatus(orderID)
	if err != nil {
		return fmt.Errorf("OrderStatus(%s): %s", orderID, err)
	}
	fmt.Printf("Status:%s UUID:%s\n", status, uuid)

	return nil
}

func testListBook() error {
	var page int64 = 1
	var size int64 = 30
	lang := `fr`

	list, err := wikibookgen.ListWikibook(page, size, lang)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return fmt.Errorf("ListWikibook: empty list")
	}

	for _, l := range list {
		fmt.Printf("- %v\n", l)
	}

	return nil
}

func ClusterDepth(c *Cluster) int {
	return clusterDepth(c, 1)
}

func clusterDepth(c *Cluster, depth int) int {
	fmt.Printf("Cluster d%d : %v\n", depth, c.Members)

	var max int = depth
	for _, sub := range c.Subclusters {
		d := clusterDepth(sub, depth+1)
		if d > max {
			max = d
		}
	}

	return max
}
