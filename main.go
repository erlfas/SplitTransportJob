package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	path := os.Args[1]

	consignments := getConsignments(path)

	jobs := make(chan string, 100)
	results := make(chan string, 100)

	for w := 1; w <= 8; w++ {
		go worker(w, jobs, results)
	}

	numConsignments := len(consignments)

	fmt.Println("# consignments: ", numConsignments)

	for _, consignment := range consignments {
		jobs <- consignment
	}

	close(jobs)

	for r := 1; r <= numConsignments; r++ {
		fmt.Println(<-results)
	}
}

func worker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {
		fmt.Println("Worker #", id, " processing job.")
		resp, err := http.Post("http://localhost:8080/TransportJobMapper/rest/transportjob/save", "application/xml", bytes.NewBuffer([]byte(job)))
		if err != nil {
			log.Fatal(id, ": Failed to save job", err)
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		results <- bodyString
	}
}

func getHeader(path string) string {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	decoder := xml.NewDecoder(xmlFile)

	buffer := bytes.NewBufferString("")

MainLoop:
	for {
		token, err := decoder.RawToken()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Consignment" {
				break MainLoop
			} else {
				buffer.WriteString(getStartTag(se))
			}
			break
		case xml.Directive:
			buffer.WriteString(string(se))
			break
		case xml.EndElement:
			buffer.WriteString(getEndTag(se))
			break
		case xml.CharData:
			buffer.WriteString(string(se))
			break
		case xml.ProcInst:
			buffer.WriteString(getProcInst(se))
			break
		}
	}

	return buffer.String()
}

func getConsignments(path string) []string {
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	decoder := xml.NewDecoder(xmlFile)

	var consignments []string
	isUnderConsignment := false
	buffer := bytes.NewBufferString("")

	for {
		token, err := decoder.RawToken()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Consignment" {
				tag := getStartTag(se)
				isUnderConsignment = true
				buffer.WriteString(tag)
			} else if isUnderConsignment {
				tag := getStartTag(se)
				buffer.WriteString(tag)
			}
			break
		case xml.EndElement:
			if se.Name.Local == "Consignment" {
				tag := getEndTag(se)
				isUnderConsignment = false
				buffer.WriteString(tag)
				bf := buffer.String()
				bfslice := make([]string, 1)
				bfslice[0] = bf
				consignments = append(consignments, bfslice...)
				buffer = bytes.NewBufferString("")
			} else if isUnderConsignment {
				tag := getEndTag(se)
				buffer.WriteString(tag)
			}
			break
		case xml.CharData:
			if isUnderConsignment {
				buffer.WriteString(string(se))
			}
			break
		}
	}

	header := getHeader(path)
	footer := "\n</TransportJob>"

	for i, c := range consignments {
		consignments[i] = header + c + footer
	}

	return consignments
}

func getProcInst(elm xml.ProcInst) string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString("<?")
	buffer.WriteString(elm.Target)
	buffer.WriteString(" ")
	buffer.WriteString(string(elm.Inst))
	buffer.WriteString("?>")

	return buffer.String()
}

func getStartTag(elm xml.StartElement) string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString("<")
	buffer.WriteString(elm.Name.Local)

	for _, atr := range elm.Attr {
		buffer.WriteString(" ")
		buffer.WriteString(atr.Name.Local)
		buffer.WriteString("=")
		buffer.WriteString("'")
		buffer.WriteString(atr.Value)
		buffer.WriteString("'")
	}

	buffer.WriteString(">")

	return buffer.String()
}

func getEndTag(elm xml.EndElement) string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString("</")
	buffer.WriteString(elm.Name.Local)
	buffer.WriteString(">")

	return buffer.String()
}
