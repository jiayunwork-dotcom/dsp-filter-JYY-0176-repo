package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"dsp-filter/internal/api"
	"dsp-filter/internal/design"
	"dsp-filter/internal/response"
	"dsp-filter/internal/zplane"
)

//go:embed web/*
var webFS embed.FS

//go:embed example/*
var exampleFS embed.FS

func main() {
	httpAddr := flag.String("http", "", "serve web UI and API on this address, e.g. :8080")
	flag.Parse()
	if *httpAddr != "" {
		sub, err := fs.Sub(webFS, "web")
		if err != nil {
			log.Fatal(err)
		}
		examples, err := fs.Sub(exampleFS, "example")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("dsp-filter listening on %s", *httpAddr)
		log.Fatal(http.ListenAndServe(*httpAddr, api.New(sub, examples)))
	}
	args := flag.Args()
	if len(args) == 0 {
		usage()
	}
	switch args[0] {
	case "design":
		runDesign(args[1:])
	case "response":
		runResponse(args[1:])
	case "zplane":
		runZPlane(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "unknown command "+args[0])
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: dsp-filter [-http :8080] design -kind fir -order 30 -cutoff 0.2 -window hamming | response <b.json> | zplane <b.json>")
	os.Exit(2)
}

func runDesign(args []string) {
	fs := flag.NewFlagSet("design", flag.ExitOnError)
	kind := fs.String("kind", "fir", "fir or iir")
	order := fs.Int("order", 30, "filter order")
	cutoff := fs.Float64("cutoff", 0.2, "normalized cutoff in (0, 0.5)")
	window := fs.String("window", "hamming", "fir window: rect/hann/hamming")
	fs.Parse(args)
	f, err := design.Design(&design.DesignSpec{
		Kind: design.Kind(*kind), Order: *order, Cutoff: *cutoff, Window: *window,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("kind=%s order=%d\n", f.Kind, f.Order())
	fmt.Printf("b = %v\n", f.B)
	fmt.Printf("a = %v\n", f.A)
}

func runResponse(args []string) {
	if len(args) < 1 {
		usage()
	}
	f, err := loadFilter(args[0])
	if err != nil {
		log.Fatal(err)
	}
	res, err := response.Compute(f.B, f.A, response.Grid(128))
	if err != nil {
		log.Fatal(err)
	}
	pass := passbandCheck(res.MagnitudeDB, 0.2, 3)
	fmt.Printf("points=%d passband@0.2dB=%.4g 3dB@fc=%v\n", len(res.Freq), res.MagnitudeDB[51], pass)
	fmt.Printf("group delay avg = %.4g samples\n", avg(res.GroupDelay))
}

func runZPlane(args []string) {
	if len(args) < 1 {
		usage()
	}
	f, err := loadFilter(args[0])
	if err != nil {
		log.Fatal(err)
	}
	zp := zplane.ZeroPoles(f.B, f.A)
	fmt.Printf("zeros=%d poles=%d stable=%v\n", len(zp.Zeros), len(zp.Poles), zp.Stable)
}

func loadFilter(path string) (*design.Filter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f design.Filter
	if err := json.Unmarshal(data, &f); err == nil && len(f.B) > 0 {
		return &f, nil
	}
	var spec design.DesignSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return design.Design(&spec)
}

func avg(v []float64) float64 {
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

func passbandCheck(mag []float64, cutoff float64, tol float64) bool {
	idx := int(cutoff * float64(len(mag)-1) * 2)
	if idx >= len(mag) {
		idx = len(mag) - 1
	}
	return mag[idx] > -tol
}
