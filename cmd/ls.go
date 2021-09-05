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
	Example: "$ gekikara ls 渋谷駅",
	Short:   "List up to 20 restaurants from the desired address.",
	Long:    "List up to 20 restaurants from the desired address.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify address.")
		} else if len(args) > 1 {
			fmt.Println("Too many args.")
		} else {
			listUpGekikara(cmd, args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringP("filter", "f", "", "Additional keyword filter.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listUpGekikara(cmd *cobra.Command, address string) {
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

	var keyword string = "激辛"
	filter, _ := cmd.Flags().GetString("filter")
	if filter != "" {
		keyword = "\"" + keyword + filter + "\""
	}
	nearbySearchReq := maps.NearbySearchRequest{
		Location: &geocodingRes[0].Geometry.Location,
		Radius:   1000,
		Keyword:  keyword,
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
