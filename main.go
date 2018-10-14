package main

import (
	"flag"
	"log"
)

func main() {
	modePtr := flag.String("mode", "B", "[P]lay || [B]uild || [D]eploy (lab)")
	swaggerPtr := flag.Bool("swagger", false, "Updates Swagger docs and routers.")
	flag.Parse()

	log.Printf("/*FLYing Mode:%s Swagger:%t*/\n", *modePtr, *swaggerPtr)
	conf, err := loadConfig()

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
	}

	log.Print("Thank you, FLY again soon!")
}
