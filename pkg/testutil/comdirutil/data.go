package comdirutil

import (
	"fmt"
	"math/rand"
	"time"
)

var registry map[string]*Employee

func init() {
	registry = make(map[string]*Employee)
	n := 100
	for i := 0; i < n; i++ {
		registry[fmt.Sprintf("id-%d", i)] = &Employee{
			Person: &Person{
				ID:   fmt.Sprintf("id-%d", i),
				Name: fmt.Sprintf("name-%d", i),
			},
			Meta: &Meta{
				JoinDate:   time.Now(),
				EndDate:    time.Now().Add(time.Hour * time.Duration(rand.Intn(100))),
				Department: fmt.Sprintf("department-%d", i),
			},
			Manager: fmt.Sprintf("id-%d", rand.Intn(n)),
			Manages: []string{
				fmt.Sprintf("id-%d", rand.Intn(n)),
				fmt.Sprintf("id-%d", rand.Intn(n)),
				fmt.Sprintf("id-%d", rand.Intn(n)),
			},
		}
	}
}
