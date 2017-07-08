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

type Job struct {
	ID          int
	Consignment string
}

func main() {

	path := os.Args[1]

	consignments := getConsignments(path)

	var tasks []*Task

	for id, consignment := range consignments {
		slice := make([]*Task, 1)
		jobID := id
		slice[0] = NewTask(jobID, func() string {
			resp, err := http.Post(
				"http://localhost:8080/TransportJobMapper/rest/transportjob/save",
				"application/xml",
				bytes.NewBuffer([]byte(consignment)))

			if err != nil {
				return "500"
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(body)

			return bodyString
		})
		tasks = append(tasks, slice...)
	}

	pool := NewPool(tasks, 8)
	pool.Run()

	var numErrors int
	for _, task := range pool.Tasks {
		if task.Result == "400" || task.Result == "500" {
			log.Fatal("Task #", task.ID, " failed.")
			numErrors++
		}
		if numErrors >= 10 {
			log.Fatal("Too many errors.")
			break
		}
	}
}

func worker(id int, jobs <-chan *Job, results chan<- string) {
	for {
		select {
		case job := <-jobs:
			fmt.Println("Worker #", id, " starting job #", job.ID)
			resp, err := http.Post(
				"http://localhost:8080/TransportJobMapper/rest/transportjob/save",
				"application/xml",
				bytes.NewBuffer([]byte(job.Consignment)))

			if err != nil {
				log.Fatal("Worker #", id, " failed to save job", err)
				results <- "500"
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(body)
			fmt.Println("Worker #", id, " finished job #", job.ID, " with status ", bodyString)
			results <- bodyString
		}
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
