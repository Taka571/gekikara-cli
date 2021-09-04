/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"googlemaps.github.io/maps"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List up to 20 restaurants from the desired address.",
	Long: `List up to 20 restaurants from the desired address.
For example: $ gekikara ls 渋谷駅`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify address.")
		} else if len(args) > 1 {
			fmt.Println("Too many args.")
		} else {
			listUpGekikara(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listUpGekikara(address string) {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GCP_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	geocodingReq := maps.GeocodingRequest{
		Address:  address,
		Language: "ja",
	}

	geocodingRes, err := client.Geocode(context.TODO(), &geocodingReq)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	nearbySearchReq := maps.NearbySearchRequest{
		Location: &geocodingRes[0].Geometry.Location,
		Radius:   1000,
		Keyword:  "激辛",
		Language: "ja",
		Type:     "restaurant",
	}

	nearbySearchRes, err := client.NearbySearch(context.TODO(), &nearbySearchReq)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for i, res := range nearbySearchRes.Results {
		humanizeIndex := fmt.Sprintf("%3s", strconv.Itoa(i+1)+".")
		rating := fmt.Sprintf("%6s", "☆"+fmt.Sprintf("%.1f", res.Rating))
		ratingsTotal := fmt.Sprintf("%7s", "("+strconv.Itoa(res.UserRatingsTotal)+"件"+")")
		name := res.Name

		fmt.Println(humanizeIndex + rating + ratingsTotal + "  " + name)
	}
}
