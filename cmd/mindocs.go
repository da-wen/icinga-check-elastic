// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const qs string = `{
	"query": {
		"bool": {
			"must": [
				{
					"range": {
						"%s": {
							"gte":"%s"
						}
					}
				}
			]
		}
	  
	}
  }
`

var url, index, field, auth string
var ecode, period, cTimeout int
var treshold int64

// countCmd represents the mincount command
var mindocsCmd = &cobra.Command{
	Use:   "mindocs",
	Short: "Counts documents for a given index patern.",
	Long: `Fails with exit code 2 , if doc count for a period is lower then defined treshold
`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(mindocsCmd)
	// required flags
	mindocsCmd.Flags().StringVar(&url, "url", "", "string of url like http://localhost:9200")
	mindocsCmd.Flags().StringVar(&index, "index", "", "string of idex pattern or full index name like my-awesome-index-*")
	mindocsCmd.Flags().StringVar(&field, "field", "", "string of filed name like internal.created")
	mindocsCmd.MarkFlagRequired("url")
	mindocsCmd.MarkFlagRequired("index")
	mindocsCmd.MarkFlagRequired("field")

	mindocsCmd.Flags().IntVarP(&period, "period", "p", 30, "sets mintues for period now - x minutes")
	mindocsCmd.Flags().Int64VarP(&treshold, "min", "m", 100, "defines minimum amount of docs that are required when command fails")
	mindocsCmd.Flags().IntVarP(&ecode, "exit", "e", 2, "exit code to be used for fail")
	mindocsCmd.Flags().StringVarP(&auth, "auth", "a", "", "basic auth for header authentication. format=username:password")
}

func logExit(m string, err error) {
	if err != nil {
		fmt.Println(m, err)
	} else {
		fmt.Println(m)
	}

	os.Exit(5)
}

func run() {
	p := fmt.Sprintf("%s/%s/_count", url, index)
	now := time.Now()
	timePast := now.Add(time.Duration(-period) * time.Minute)
	s := fmt.Sprintf(qs, field, timePast.Format("2006-01-02T15:04:05"))
	qb := []byte(s)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: false,
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(http.MethodGet, p, bytes.NewBuffer(qb))
	if err != nil {
		logExit("Request can not be created", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if len(auth) != 0 {
		bAuth := []byte(auth)
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(bAuth))
	}

	res, err := client.Do(req)
	if err != nil {
		logExit("Es request error", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logExit("Status response not readable", err)
	}
	if res.StatusCode >= 400 {
		var resErr map[string]interface{}
		json.Unmarshal(body, &resErr)
		logExit(fmt.Sprintf("Response code invalid: %v", resErr["error"]), nil)
	}

	var count struct {
		Count int64 `json:"count"`
	}
	json.Unmarshal(body, &count)

	if count.Count < treshold {
		fmt.Printf("Critical Doc count is under threshold [%d < %d] | docs=%d", count.Count, treshold, count.Count)
		os.Exit(ecode)
	}

	fmt.Printf("OK Doc count is over threshold  [%d > %d] | docs=%d", count.Count, treshold, count.Count)
}
