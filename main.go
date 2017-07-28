package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {

	if _, ok := os.LookupEnv("PRIVATE_IPV4"); ok == false {
		log.Fatal("ERROR: PRIVATE_IPV4 must be defined.")
	}

	if _, ok := os.LookupEnv("MESOS_SLAVE_COMMON_PATH"); ok == false {
		log.Fatal("ERROR: MESOS_SLAVE_COMMON_PATH must be defined.")
	}

	if _, ok := os.LookupEnv("MESOS_EXPORTER_PORT"); ok == false {
		log.Fatal("ERROR: MESOS_EXPORTER_PORT must be defined.")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "19001"
	}

	log.Infof("mesos-slave-attribute-appender listening on port %s", port)
	http.HandleFunc("/metrics", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	// get metrics from mesos-exporter
	mesosExporterAddr := fmt.Sprintf("http://%s:%s/metrics", os.Getenv("PRIVATE_IPV4"), os.Getenv("MESOS_EXPORTER_PORT"))
	res, err := http.Get(mesosExporterAddr)
	check(err)
	exporterMetrics, err := ioutil.ReadAll(res.Body)
	check(err)

	// append mesos slave attributes to the metrics
	metricsWithAttributes, err := appendAttributes(exporterMetrics)
	check(err)

	// serve the updated metrics
	w.Write(metricsWithAttributes)
}

func appendAttributes(metrics []byte) ([]byte, error) {

	// make attributes prometheus label-friendly format
	attributes, err := attributesToLabels()
	check(err)

	// replace all } with attributes
	metrics = bytes.Replace(metrics, []byte("}"), []byte(attributes), -1)

	return metrics, nil
}

// get the slave attributes in prometheus label format
func attributesToLabels() (string, error) {

	// contents of mesos-slave-common looks like this:  MESOS_ATTRIBUTES=cloud:vmware;rack:c12;subnet-cidr:10.226.128.0/21

	attributeBytes, err := ioutil.ReadFile(os.Getenv("MESOS_SLAVE_COMMON_PATH"))
	check(err)

	attributes := string(attributeBytes)

	// TODO: change all these string things to bytes so we don't need to convert types unnecessarily

	attributes = strings.TrimPrefix(attributes, "MESOS_ATTRIBUTES=")
	attributes = strings.TrimSuffix(attributes, "\n")

	attributesSplit := strings.Split(attributes, ";")

	attributesAsLabels := ""

	for _, attr := range attributesSplit {
		// key will be stuff before the colon, value will be stuff after
		attrSplit := strings.Split(attr, ":")

		// format like prometheus-- key must have underscores, not hyphens, to be a valid prometheus label
		r := strings.NewReplacer("-", "_")
		attrSplit[0] = fmt.Sprintf(r.Replace(attrSplit[0]))
		attributesAsLabels = fmt.Sprintf("%s,%s=\"%s\"", attributesAsLabels, attrSplit[0], attrSplit[1])
	}

	// add } to the end
	attributesAsLabels = fmt.Sprintf("%s}", attributesAsLabels)

	return attributesAsLabels, nil
}

func check(err error) {
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}
