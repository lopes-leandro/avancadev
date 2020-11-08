package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var coupons Coupons

// Coupon guarda um código
type Coupon struct {
	Code string
}

// Coupons guarda uma lista de códigos
type Coupons struct {
	Coupon []Coupon
}

// Check verifica se o código é valido
func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "not_applied"
		}
	}
	coupon := Coupon{
		Code: code,
	}
	coupons.Coupon = append(coupons.Coupon, coupon)
	return "applied"
}

// Result é usado para realizar o retorno
type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	appied := coupons.Check(coupon)

	result := Result{Status: appied}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}