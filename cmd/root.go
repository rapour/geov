package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/profile"
	"github.com/rapour/geov"
	"github.com/spf13/cobra"
)

var in string
var out string
var ratio float64

func init() {
	rootCmd.PersistentFlags().StringVar(&in, "input", "", "input file path")
	rootCmd.PersistentFlags().Float64Var(&ratio, "ratio", 1, "ratio of simplification")
	rootCmd.PersistentFlags().StringVar(&out, "output", "", "output file path")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")

}

var rootCmd = &cobra.Command{
	Use: "map",
	Run: func(cmd *cobra.Command, args []string) {

		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()

		bin, err := os.ReadFile(in)
		if err != nil {
			log.Fatal(err)
		}

		out, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}

		mp, err := geov.Unmarshal(bin)
		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()
		smp := geov.Simplify(mp, geov.Visvalingam, ratio)

		fmt.Printf("processing time %f\n", time.Since(t).Seconds())

		err = smp.SVG(out)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {

	if err := Execute(); err != nil {
		log.Fatal(err)
	}

}
