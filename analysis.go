package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//MyCandidate struct as part of the build up for json object
type MyCandidate struct {
	Candidate string `json:"candidate"`
	Votes     string `json:"votes"`
}

//MyParty struct as part of the build up for json object
type MyParty struct {
	Party   string              `json:"party"`
	Results []map[string]string `json:"results"`
}

//MyElections struct as part of the build up for json object
type MyElections struct {
	Name      string    `json:"name"`
	Fips      string    `json:"fips"`
	Elections []MyParty `json:"elections"`
}

//MyCounty struct json object to be returned by the API call
type MyCounty struct {
	Counties []MyElections
}

//LoadFile is a standalone function to return a read csv file.
func LoadFile(csvfile string) (newData [][]string) {
	newFile, err := os.Open(csvfile)
	if err != nil {
		fmt.Println("Failed to load file csv file", err)
	}
	defer newFile.Close()

	myReader := csv.NewReader(newFile)
	myReader.Read()
	// myReader.Comma = ','
	myReader.FieldsPerRecord = -1
	newData, err = myReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

/*Contains function checks whether an item is in the slice, this function is used to compute unique
items in the slice*/
func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

//GetUniqueFips fuction will return a set of fips from the loaded csv file.
func GetUniqueFips() []string {
	fipSet := make([]string, 0)
	for _, data := range LoadFile("demResults.csv") {
		fipData := data[1]
		if ok := Contains(fipSet, fipData); !ok {
			fipSet = append(fipSet, fipData)
		}

	}
	return fipSet
}

/*GetUniqueFipsCounties fuction will return a set of counties from the loaded csv file,
using the set of unique fips as reference.*/
func GetUniqueFipsCounties() []map[string]string {
	//returns a Slice of unique fips county
	countySet := make([]map[string]string, 0)
	myList := GetUniqueFips()
	for _, num := range myList {
		for _, data := range LoadFile("demResults.csv") {
			if data[1] == num {
				countyMap := make(map[string]string)
				countyMap["county"] = data[0]
				countyMap["fips"] = data[1]
				countySet = append(countySet, countyMap)
				break
			}
		}
	}
	return countySet
}

//RepCounties returns list of maps(MyCandidate) republican candidates for a each fips.
func RepCounties(tFips string) []map[string]string {
	repResultList := make([]map[string]string, 0)
	for _, data := range LoadFile("repResults.csv") {
		if tFips == data[1] {
			repMap := make(map[string]string)
			repMap["candidate"] = data[2]
			repMap["votes"] = data[3]
			repResultList = append(repResultList, repMap)
		}
	}
	return repResultList
}

//DemCounties function returns list of maps(MyCandidate) democratic candidates for a each fips.
func DemCounties(tFips string) []map[string]string {
	demResultList := make([]map[string]string, 0)
	for _, data := range LoadFile("demResults.csv") {
		if tFips == data[1] {
			demMap := make(map[string]string)
			demMap["candidate"] = data[2]
			demMap["votes"] = data[3]
			demResultList = append(demResultList, demMap)
		}
	}
	return demResultList
}

//AllCounties function is the handler for api with route(/counties)
func AllCounties(w http.ResponseWriter, r *http.Request) {
	uniquefipscounties := GetUniqueFipsCounties()
	var output MyCounty
	outputList := make([]MyElections, 0)

	for _, row := range uniquefipscounties {
		demResults := DemCounties(row["fips"])
		repResults := RepCounties(row["fips"])
		demParty := MyParty{Party: "Democratic", Results: demResults}
		repParty := MyParty{Party: "Republican", Results: repResults}
		mergeList := make([]MyParty, 0)
		mergeList = append(mergeList, demParty, repParty)
		countiesMap := MyElections{Name: row["county"], Fips: row["fips"], Elections: mergeList}
		outputList = append(outputList, countiesMap)
	}
	output = MyCounty{outputList}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
	// jsonString, err := json.Marshal(output)
	// if err != nil{
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(jsonString))
}

//GetFromFip function is the handler for api with route(/counties/<fips>)
func GetFromFip(w http.ResponseWriter, r *http.Request) {
	uniquefipscounties := GetUniqueFipsCounties()
	var output MyCounty
	outputList := make([]MyElections, 0)
	params := mux.Vars(r)

	for _, row := range uniquefipscounties {
		if row["fips"] == params["fips"] {
			demResults := DemCounties(row["fips"])
			repResults := RepCounties(row["fips"])
			demParty := MyParty{Party: "Democratic", Results: demResults}
			repParty := MyParty{Party: "Republican", Results: repResults}
			mergeList := make([]MyParty, 0)
			mergeList = append(mergeList, demParty, repParty)
			countiesMap := MyElections{Name: row["county"], Fips: row["fips"], Elections: mergeList}
			outputList = append(outputList, countiesMap)
		}
	}
	output = MyCounty{outputList}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

//The main function that return all the routes.
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/counties", AllCounties).Methods("GET")
	router.HandleFunc("/counties/{fips}", GetFromFip).Methods("GET")
	// AllCounties()

	log.Fatal(http.ListenAndServe(":8080", router))
}
