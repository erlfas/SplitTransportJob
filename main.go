package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type AdditionalServiceCode struct {
	Code string `xml:"AdditionalServiceCode"`
}

type TransportJob struct {
}

type Consignment struct {
	ConsignmentId string `xml:"consignmentId,attr"`
	DateAndTimes  struct {
		LoadingDate struct {
			Date string `xml:"Date"`
		} `xml:"LoadingDate"`
	} `xml:"DateAndTimes"`

	Service struct {
		BasicServiceCode       string                  `xml:"BasicServiceCode"`
		AdditionalServiceCodes []AdditionalServiceCode `xml:"AdditionalServiceCode"`
	} `xml:"Service"`
	GoodsValue struct {
		CurrencyIdentificationCode string `xml:"currencyIdentificationCode,attr"`
		GoodsValue                 string
	} `xml:"GoodsValue"`
}

func main() {
	xmlFile, err := os.Open("/home/erlend/Downloads/z14_730825601_21062017110235403.xml")
	if err != nil {
		log.Fatal(err)
	}
	decoder := xml.NewDecoder(xmlFile)

	var consignments []string
	isUnderConsignment := false
	buffer := bytes.NewBufferString("")

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "Consignment" {
				tag := getStartTag(se.Name.Local)
				//fmt.Println(tag)
				isUnderConsignment = true
				buffer.WriteString(tag)
			} else if isUnderConsignment {
				tag := getStartTag(se.Name.Local)
				//fmt.Println(tag)
				buffer.WriteString(tag)
			}
			break
		case xml.EndElement:
			if se.Name.Local == "Consignment" {
				tag := getEndTag(se.Name.Local)
				//fmt.Println(tag)
				isUnderConsignment = false
				buffer.WriteString(tag)
				bf := buffer.String()
				bfslice := make([]string, 1)
				bfslice[0] = bf
				consignments = append(consignments, bfslice...)
				buffer = bytes.NewBufferString("")
			} else if isUnderConsignment {
				tag := getEndTag(se.Name.Local)
				//fmt.Println(tag)
				buffer.WriteString(tag)
			}
			break
		case xml.CharData:
			if isUnderConsignment {
				s := strings.TrimSpace(string(se.Copy()))
				if len(s) > 0 {
					//fmt.Println(s)
					buffer.WriteString(s)
				}
			}
			break
		}
	}

	fmt.Println(len(consignments))

	for _, s := range consignments {
		fmt.Println("C: ", s)
	}
}

func getEndTag(s string) string {
	return "</" + s + ">"
}

func getStartTag(s string) string {
	return "<" + s + ">"
}
