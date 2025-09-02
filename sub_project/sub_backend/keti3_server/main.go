package main

import (
	"log"

	"keti3/pkg/tal_content_arch_lib/pontos"
	"keti3/router"
)

func main() {
	chain := &pontos.ChainOptions{
		IsInitDefault: true,
		ConfigFile:    "./app.toml",
	}

	if err := pontos.NewApp(chain); err != nil {
		panic(err)
	}

	// register router
	router.RegisterRouter(pontos.Server.GinEngine())

	err := pontos.Server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server stop err:%v\n", err)
	} else {
		log.Println("Server exit")
	}
}
