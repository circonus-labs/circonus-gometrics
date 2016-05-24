package circonusgometrics

import (
	"log"
	"strconv"
	"strings"
)

func getCheck(id int) {
	url := strings.Join([]string{"/v2/check/", strconv.Itoa(id)}, "")
	checkDetails := apiCall(url)
	details, ok := checkDetails["_details"]
	if !ok {
		log.Printf("Cannot find submission URL at %s\n", url)
		return
	}
	dmap := details.(map[string]interface{})
	val := dmap["submission_url"]
	submissionUrl = val.(string)
}
