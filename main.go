package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Vars struct {
	Hostname          string
	StackName         string
	ETCDURL           string
	PublicIPV4        string
	PrivateIPV4       string
	SSHAuthorizedKeys string
}

func main() {

	metadataService := "ops:8080/metadata/v1"

	templateFilename := "cloud-config.template.yaml"
	t := template.New(templateFilename).Delims("<<", ">>")
	t, err := t.ParseFiles(templateFilename)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/metadata/v1/user-data", func(c *gin.Context) {

		filename := "custom.json"

		// Open custom.json
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			errorMessage := fmt.Sprintf("[%s] You must create a custom.json file. custom.json.example should get you started.", err)
			c.String(404, errorMessage)
			return
		}

		// Parse custom.json
		decoder := json.NewDecoder(f)
		vars := Vars{}
		err = decoder.Decode(&vars)
		if err != nil {
			errorMessage := fmt.Sprintf("Invalid custom.json file. [%s]", err)
			c.String(404, errorMessage)
			return
		}

		// Decorate with a few variables of our own
		vars.Hostname = hostname(c)
		vars.PublicIPV4 = ip(c)
		vars.PrivateIPV4 = ip(c)

		cloudConfig := generateCloudConfig(t, &vars)

		c.String(http.StatusOK, cloudConfig)
	})

	router.GET("/metadata/v1/ip", func(c *gin.Context) {
		// s := "ip route get 8.8.8.8 | awk '{print $NF; exit}'"
		c.String(http.StatusOK, ip(c))
	})

	router.GET("/metadata/v1/hostname", func(c *gin.Context) {
		c.String(http.StatusOK, hostname(c))
	})

	router.GET("/metadata/v1/install", func(c *gin.Context) {
		s := "" +
			"curl -s \"" + metadataService + "/user-data\" > user-data\n" +
			"sudo coreos-install -d /dev/sda -C beta -c ./user-data\n"
		c.String(http.StatusOK, s)
	})

	router.Run(":8080")
}

func generateCloudConfig(t *template.Template, vars *Vars) string {

	var cloudConfigBuffer bytes.Buffer
	err := t.Execute(&cloudConfigBuffer, vars)
	if err != nil {
		log.Println(err)
	}
	return cloudConfigBuffer.String()
}

func ip(c *gin.Context) string {
	clientIP, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return clientIP
}

func hostname(c *gin.Context) string {
	hostname := "sn-" + strings.Replace(ip(c), ".", "-", -1)
	return hostname
}
