package src

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	vegeta "github.com/vegeta/lib"
)

type Account struct {
	Type          string `json:"type"`
	RoutingNumber int    `json:"routingNumber"`
	AccountNumber int    `json:"accountNumber"`
	Balance       int    `json:"balance"`
	Interest      int    `json:"interest"`
}

type AttackTargets struct {
	targets []vegeta.Target
}

func init() {

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func GenerateLoadData(count int, duration int, api string) {

	numberOfTargets := count * duration

	fmt.Println("Generating LoadData for number of requests", count, api)
	var targets []vegeta.Target
	// targets := make([]vegeta.Target, numberOfTargets)

	rand.Seed(time.Now().UnixNano())
	ftype := []func(){generateAccounts(api, &targets), generateGETRequests(api, &targets)}

	for index := 0; index < numberOfTargets; index++ {
		generator := ftype[rand.Intn(len(ftype))]
		generator()
	}

	log.WithFields(log.Fields{"Number of targets generated": len(targets)}).Debug()

	rate := uint64(count) // per second
	du := time.Duration(duration) * time.Second

	for index := 0; index < len(targets); index++ {
		fmt.Println(string(targets[index].Body), targets[index].URL)
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(vegeta.NewStaticTargeter(targets...), rate, du, "abc") {
		fmt.Println(res.Error)
		metrics.Add(res)
	}
	metrics.Close()

	log.WithFields(log.Fields{"99th percentile": metrics.Latencies.P99, "rate": metrics.Rate, "requests": metrics.Requests, "duration": metrics.Duration}).Info()

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func generateAccounts(api string, targets *[]vegeta.Target) func() {

	return func() {
		account := generateRandomAccount()
		log.WithFields(log.Fields{"Account: ": account}).Debug()
		body, err := json.Marshal(account)
		check(err)
		var header = make(http.Header)
		header.Add("content-type", "application/json")

		target := vegeta.Target{
			Method: http.MethodPost,
			URL:    api + "/account",
			Body:   body,
			Header: header,
		}

		addValue(targets, target)
	}
}

func generateGETRequests(api string, targets *[]vegeta.Target) func() {

	return func() {
		appURLs := []string{
			"/customers",
			"/account",
			"/patients",
			"/saveSettings",
			"/loadSettings",
			"/customers/3",
			"/account/2",
			"/customers/1",
			"/error",
			"/debugEscaped?firstName=%22%22",
			"/account/1",
			"/account/3",
			"/search/user?foo=new%20java.lang.ProcessBuilder(%7B%27%2Fbin%2Fbash%27%2C%27-c%27%2C%27echo%203vilhax0r%3E%2Ftmp%2Fhacked%27%7D).start()",
			"/debug?customerId=ID-4242&clientId=1&firstName=%22%22&lastName=%22%22&dateOfBirth=10-11-17&ssn=%22%22&socialSecurityNum=%22%22&tin=%22%22&phoneNumber=%22%22",
			"/debugEscaped?firstName=%22%22",
			"/admin/login"}

		target := vegeta.Target{
			Method: http.MethodGet,
			URL:    api + appURLs[rand.Intn(len(appURLs))],
		}
		addValue(targets, target)
	}

}

func addValue(s *[]vegeta.Target, target vegeta.Target) {
	*s = append(*s, target)
	// fmt.Printf("In addValue: s is %v\n", s)
}

func generateRandomAccount() Account {
	var account Account
	accountType := []string{"SAVING", "CHECKING"}
	account.RoutingNumber = rand.Intn(50000)
	account.Balance = rand.Intn(50000)
	account.Interest = rand.Intn(15)
	account.Type = accountType[rand.Intn(len(accountType))]
	return account
}
