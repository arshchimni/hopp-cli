package methods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

//Colls hold the structure of the basic `postwoman-collection.json`
type Colls struct {
	Name    string    `json:"name"`
	Folders []string  `json:"folders"`
	Request []Reqdata `json:"requests"`
}

//Reqdata hold the format of the request part in `postwoman-collection.json`
type Reqdata struct {
	URL     string     `json:"url"`
	Path    string     `json:"path"`
	Method  string     `json:"method"`
	Auth    string     `json:"auth"`
	User    string     `json:"httpUser"`
	Pass    string     `json:"httpPassword"`
	Token   string     `json:"bearerToken"`
	Ctype   string     `json:"contentType"`
	Name    string     `json:"name"`
	Heads   []string   `json:"headers"`
	Params  []string   `json:"params"`
	Bparams []Bpardata `json:"bodyParams"`
}

//Bpardata hold the format of the bodyParams of `postwoman-collection.json`
type Bpardata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//ReadCollection reads the PostWoman Collection Json File and does the Magic Stuff
func ReadCollection(c *cli.Context) {
	data, err := ioutil.ReadFile(c.String("pt"))
	if err != nil {
		fmt.Print(err)
	}
	var jsondat []Colls
	err = json.Unmarshal([]byte(data), &jsondat)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Name:\t" + color.HiMagentaString(jsondat[0].Name))
	for i := 0; i < len(jsondat[0].Request); i++ {
		request(jsondat, i)
	}

}
func request(c []Colls, i int) {
	colors := color.New(color.FgHiRed, color.Bold)
	if c[0].Request[i].Method == "GET" {
		out, err := getsend(c, i, "GET")
		if err != nil {
			fmt.Print(err)
		}
		methods := color.HiYellowString(c[0].Request[i].Method)
		fURL := colors.Sprintf(c[0].Request[i].URL + c[0].Request[i].Path)
		fmt.Printf("%s \t%s\t%s", fURL, methods, out)
	} else {
		out, err := sendpopa(c, i, c[0].Request[i].Method)
		if err != nil {
			fmt.Print(err)
		}
		methods := color.HiYellowString(c[0].Request[i].Method)
		fURL := colors.Sprintf(c[0].Request[i].URL + c[0].Request[i].Path)
		fmt.Printf("%s |\t%s |\t%s |\t%s", color.HiGreenString(c[0].Request[i].Name), fURL, methods, out)
	}

}
func getsend(c []Colls, ind int, method string) (string, error) {
	color := color.New(color.FgCyan, color.Bold)
	var url = c[0].Request[ind].URL + c[0].Request[ind].Path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	if c[0].Request[ind].Token != "" {
		var bearer = "Bearer " + c[0].Request[ind].Token
		req.Header.Add("Authorization", bearer)
	}
	if c[0].Request[ind].User != "" && c[0].Request[ind].Pass != "" {
		un := c[0].Request[ind].User
		pw := c[0].Request[ind].Pass
		req.Header.Add("Authorization", "Basic "+basicAuth(un, pw))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	defer resp.Body.Close()
	//fmt.Print(resp.Header)
	s := color.Sprintf("Status: %s\tStatusCode:\t%d\n", resp.Status, resp.StatusCode)
	return s, nil
}

func sendpopa(c []Colls, ind int, method string) (string, error) {
	color := color.New(color.FgCyan, color.Bold)
	var jsonStr []byte
	var url = c[0].Request[ind].URL + c[0].Request[ind].Path
	if len(c[0].Request[ind].Bparams) > 0 {
		jsonStr = []byte(string(c[0].Request[ind].Bparams[0].Key[0] + c[0].Request[ind].Bparams[0].Value[0]))
	} else {
		jsonStr = nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", c[0].Request[ind].Ctype)
	if err != nil {
		return "", err
	}
	if c[0].Request[ind].Token != "" {
		var bearer = "Bearer " + c[0].Request[ind].Token
		req.Header.Add("Authorization", bearer)
	}
	if c[0].Request[ind].User != "" && c[0].Request[ind].Pass != "" {
		un := c[0].Request[ind].User
		pw := c[0].Request[ind].Pass
		req.Header.Add("Authorization", "Basic "+basicAuth(un, pw))
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	defer resp.Body.Close()
	//fmt.Print(resp.Header)
	s := color.Sprintf("Status: %s\tStatusCode:\t%d\n", resp.Status, resp.StatusCode)
	return s, nil

}
