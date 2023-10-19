package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/chrusty/go-tableprinter"
)

type Config struct {
	Token        string `json:"token"`
	JfrogBaseURL string `json:"jfrogbaseurl"`
	ODays        string `json:"odays"`
}

var AccessToken = ""
var JfURL = ""
var lastdays = ""

type resultData struct {
	Name            string `header:"User Name"`
	TotalDownSize   int    `header:"Download Size"`
	TotalDownloads  int    `header:"Downloaded files"`
	TotalUploadSize int    `header:"Upload Size"`
	TotalUploads    int    `header:"Uploaded files"`
}
type resultDataArr struct {
	Items []resultData
}

type resultStruct struct {
	Results []struct {
		Repo       string    `json:"repo"`
		Path       string    `json:"path"`
		Name       string    `json:"name"`
		Type       string    `json:"type"`
		Size       int       `json:"size"`
		Created    time.Time `json:"created"`
		CreatedBy  string    `json:"created_by"`
		Modified   time.Time `json:"modified"`
		ModifiedBy string    `json:"modified_by"`
		Updated    time.Time `json:"updated"`
		Stats      []struct {
			Downloaded      time.Time `json:"downloaded"`
			DownloadedBy    string    `json:"downloaded_by"`
			Downloads       int       `json:"downloads"`
			RemoteDownloads int       `json:"remote_downloads"`
		} `json:"stats"`
	} `json:"results"`
}

// type resultStructArr struct {
// 	Items []repoStruct
// }

// func (r *resultDataArr) AddItem(item resultData) {
// 	r.Items = append(r.Items, item)

// }
func (r *resultDataArr) AddItem(item resultData) {
	r.Items = append(r.Items, item)

}
func jsonFileToStruct(fileName string) resultStruct {
	file, _ := ioutil.ReadFile(fileName)

	data := resultStruct{}

	_ = json.Unmarshal([]byte(file), &data)

	return data

}

func searchByKey(d resultStruct, keyName string, keyValue string) resultStruct {
	dsearched := resultStruct{}
	for i := 0; i < len(d.Results); i++ {
		r1 := reflect.ValueOf(d.Results[i])
		f := reflect.Indirect(r1).FieldByName(keyName)
		// fmt.Println("Product Id: ", f)
		v := fmt.Sprintf("%v", f)
		if v == keyValue {
			// fmt.Println(d.Results[i])
			dsearched.Results = append(dsearched.Results, d.Results[i])
		}
		// r = append(r, v)
	}
	return dsearched
}
func getUniqValueList(d resultStruct, u resultStruct, keyName string) []string {
	r := []string{}
	for i := 0; i < len(d.Results); i++ {
		r1 := reflect.ValueOf(d.Results[i])
		f := reflect.Indirect(r1).FieldByName(keyName)
		// fmt.Println("Product Id: ", f)
		v := fmt.Sprintf("%v", f)
		if v == "_system_" {
			continue
		}
		r = append(r, v)
	}
	for i := 0; i < len(u.Results); i++ {
		r1 := reflect.ValueOf(u.Results[i])
		f := reflect.Indirect(r1).FieldByName(keyName)
		// fmt.Println("Product Id: ", f)
		v := fmt.Sprintf("%v", f)
		if v == "_system_" {
			continue
		}
		r = append(r, v)
	}
	return Unique(r)
}

func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqSlice = append(uniqSlice, v)
	}
	return uniqSlice
}
func PrintTable(items resultDataArr) {
	// fmt.Println(rdata)

	Alldownloades := 0
	Alluploads := 0
	AllDownloadSize := 0
	AllUploadSize := 0
	for _, myd1 := range items.Items {
		AllDownloadSize = AllDownloadSize + myd1.TotalDownSize
		Alldownloades = Alldownloades + myd1.TotalDownloads
		Alluploads = Alluploads + myd1.TotalUploads
		AllUploadSize = AllUploadSize + myd1.TotalUploadSize
	}
	k := resultData{Name: "Total", TotalDownSize: AllDownloadSize, TotalDownloads: Alldownloades, TotalUploadSize: AllUploadSize, TotalUploads: Alluploads}
	items.AddItem(k)
	sort.Slice(items.Items, func(i, j int) bool {
		return items.Items[j].TotalDownSize+items.Items[j].TotalUploadSize < items.Items[i].TotalDownSize+items.Items[i].TotalUploadSize
	})
	// Print the slice of structs as table, as shown above.
	// fmt.Println(items.Items)
	// printer.Print(items.Items)

	// printer := tableprinter.New().WithBorders(true)
	tableprinter.SetSortedHeaders(false)
	tableprinter.SetBorder(true)
	tableprinter.Print(items.Items)
}
func summaryData(d resultStruct, u resultStruct, name string) resultData {
	totalUsize := 0
	totalSize := 0
	for _, myd := range d.Results {
		if myd.CreatedBy == "_system_" {
			continue
		}
		sizea := myd.Size
		var dcount int
		for _, mydow := range myd.Stats {
			dcount = mydow.Downloads
		}
		totalSize = totalSize + sizea*dcount
	}
	for _, myu := range u.Results {
		if myu.CreatedBy == "_system_" {
			continue
		}
		totalUsize = totalUsize + myu.Size
	}

	rd := resultData{Name: name, TotalDownSize: totalSize, TotalUploadSize: totalUsize, TotalDownloads: len(d.Results), TotalUploads: len(u.Results)}
	return rd
}
func getDownloadStruct() resultStruct {
	mydtmp := `items.find(
		{"stat.downloaded" : {"$last" : "{{.}}"}}
	).include("stat")`
	// fmt.Println(mydtmp)
	t, err := template.New("Download").Parse(mydtmp)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, lastdays)
	if err != nil {
		panic(err)
	}
	result := tpl.String()
	// fmt.Println(result)
	repoData := restpostcall(result)
	var repData resultStruct
	json.Unmarshal(repoData, &repData)
	return repData
}

func getUploadStruct() resultStruct {
	myutmp := `items.find(
		{"updated" : {"$last" : "{{.}}"}}
	).include("stat")`
	// fmt.Println(mydtmp)
	t, err := template.New("Download").Parse(myutmp)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, lastdays)
	if err != nil {
		panic(err)
	}
	result := tpl.String()
	// fmt.Println(result)
	repoData := restpostcall(result)
	var repData resultStruct
	json.Unmarshal(repoData, &repData)
	return repData
}

func restpostcall(jsondata string) []byte {
	// var repData AutoGenerated
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	body := strings.NewReader(jsondata)
	req, err := http.NewRequest("POST", JfURL+"/artifactory/api/search/aql", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	return body1

}

func main() {
	// ddata := jsonFileToStruct("d.txt")
	// udata := jsonFileToStruct("u.txt")
	jsonFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var config Config
	json.Unmarshal(byteValue, &config)
	JfURL = config.JfrogBaseURL
	AccessToken = config.Token
	lastdays = config.ODays
	ddata := getDownloadStruct()
	udata := getUploadStruct()
	// fmt.Println(data)
	r := getUniqValueList(ddata, udata, "Repo")
	u := getUniqValueList(ddata, udata, "CreatedBy")

	repoSummary := resultDataArr{}
	UserSummary := resultDataArr{}

	for i := 0; i < len(r); i++ {
		u1 := searchByKey(ddata, "Repo", r[i])
		d1 := searchByKey(udata, "Repo", r[i])
		r1 := summaryData(u1, d1, r[i])
		repoSummary.AddItem(r1)
	}
	for i := 0; i < len(u); i++ {

		u2 := searchByKey(ddata, "CreatedBy", u[i])
		d2 := searchByKey(udata, "CreatedBy", u[i])
		r2 := summaryData(u2, d2, u[i])
		UserSummary.AddItem(r2)
	}
	// fmt.Println(res)
	fmt.Println("Report Based on Repo")
	PrintTable(repoSummary)
	fmt.Println("Report Based on Users")
	PrintTable(UserSummary)
}
