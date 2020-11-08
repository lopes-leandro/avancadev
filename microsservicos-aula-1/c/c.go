package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/hashicorp/go-retryablehttp"
)

// Coupon mantem um único coupon
type Coupon struct {
	Code string
}

// Coupons guarda todos os coupons
type Coupons struct {
	Coupon []Coupon
}

// Check verifica se o cupon é valido
func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

// Result formata o resultado
type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	resultAppliedCoupon := makeHTTPCall("http://localhost:9093", coupon)

	result := Result{Status: "not_applied"}

	if resultAppliedCoupon.Status == "applied" {
		result.Status =  valid
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}

func makeHTTPCall(urlMicroservice string, coupon string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}