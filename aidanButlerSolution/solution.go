package main

import "net"
import "fmt"
import "bufio"
import "os"
import "bytes"
import "io/ioutil"
import "net/http"
import "github.com/xeipuuv/gojsonschema"
import "strconv"

func removeDuplicates(elements []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add
		} else {
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	return result
}

func convertEventIDsStringToArray(eventsString string) []int {
	intArray := make([]int, 0)
	// TODO: convert string contaiining sequence of comma separated numbers to array
	return intArray
}

func getJsonStrings() []string {
	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8000")

	iDBaseString := "/football/events"
	eventBaseString := "/football/events"
	marketBaseString := "/football/markets"

	// send to socket
	fmt.Fprintf(conn, iDBaseString)
	// listen for reply
	eventIDsString, _ := bufio.NewReader(conn).ReadString('\n')
	eventIDsArray := convertEventIDsStringToArray(eventIDsString)
	eventIDsArray = removeDuplicates(eventIDsArray)

	// Get event JSON strings
	eventsSlice := make([]string, 0)
	for _, id := range eventIDsArray {
		eventRequestString := eventBaseString + "/" + strconv.Itoa(id)
		fmt.Fprintf(conn, eventRequestString)
		eventJsonAsString, _ := bufio.NewReader(conn).ReadString('\n')
		eventsSlice = append(eventsSlice, eventJsonAsString)
	}

	// Get market JSON strings
	marketsSlice := make([]string, 0)
	for _, id := range eventIDsArray {
		marketRequestString := marketBaseString + "/" + strconv.Itoa(id)
		fmt.Fprintf(conn, marketRequestString)
		marketJsonAsString, _ := bufio.NewReader(conn).ReadString('\n')
		marketsSlice = append(marketsSlice, marketJsonAsString)
	}

	jsonStrings := buildJsonStrings(eventIDsArray, eventsSlice, marketsSlice)

	validatedJsonStrings := make([]string, 0)

	for _, jsonStr := range jsonStrings {
		if isValidJsonString(jsonStr) {
			validatedJsonStrings = append(validatedJsonStrings, jsonStr)
		}
	}

	return validatedJsonStrings
}

func buildJsonStrings(eventIDsArray []int, eventsSlice []string, marketsSlice []string) []string {
	combinedJsonStrings := make([]string, 0)
	// TODO: combine Json strings according to schema
	return combinedJsonStrings
}

func isValidJsonString(unvalidatedJsonString string) bool {
	schemaLoader := gojsonschema.NewReferenceLoader("schema.json")
	loader := gojsonschema.NewStringLoader(unvalidatedJsonString)

	result, err := gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		return true
	} else {
		return false
	}
}

func postJson(postToUrl string, jsonStrings []string) {
	url := postToUrl
	fmt.Println("URL = ", url)
	jsonString := ""
	for _, str := range jsonStrings {
		jsonString += str
	}
	var jsonStr = []byte(jsonString)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("Response status:", response.Status)
	fmt.Println("Response headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Response body:", string(body))
}

func getEnvironmentVariableValue() string {
	url := os.Getenv("STORE_ADDR")
	return url
}

func main() {
	url := getEnvironmentVariableValue() + "/event/"
	jsonStrings := getJsonStrings()
	postJson(url, jsonStrings)
}
