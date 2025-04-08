package server

import (
	"net/http"

	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func RunServer() {
	log.SetFormatter(&logrus.TextFormatter{})

	RegisterRoutes()

	// port, err := strconv.Atoi(os.Getenv("PORT"))
	// var portAddr string
	// if err != nil {
	// 	portAddr = ":8080"
	// 	log.Info("Using default port 8080")
	// } else {
	// 	portAddr = ":" + strconv.Itoa(port)
	// 	log.Infof("Starting server on  %v", portAddr)
	// }

	// if err := http.ListenAndServe(portAddr, nil); err != nil {
	// 	log.Fatal("Server failed:", err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Using default port 8080")
	} else {
		log.Infof("Starting server on :%s", port)
	}

	// Start the server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}

}
