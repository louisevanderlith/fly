package main

import (
	"flag"
	"log"

	"github.com/louisevanderlith/fly/patterns"
)

func main() {
	modePtr := flag.String("mode", "B", "[P]lay || [B]uild || [D]eploy (lab) || [F]ly Config")
	swaggerPtr := flag.Bool("swagger", false, "Updates Swagger docs and routers.")
	flag.Parse()

	log.Printf("FLYing Mode:%s Swagger:%t\n", *modePtr, *swaggerPtr)
	conf, err := patterns.DetectConfig(".", "DEV")

	if err != nil {
		panic(err)
	}

	switch *modePtr {
	case "B":
		conf.Build()
	case "P":
		conf.Play(*swaggerPtr)
	case "D":
		conf.Deploy()
	case "F":
		log.Print("Config Generated?")
	}

	log.Print("Thank you, FLY again soon!")
}
