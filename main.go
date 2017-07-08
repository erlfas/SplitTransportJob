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

type Reference struct {
	ReferenceNo   string `xml:"ReferenceNo"`
	ReferenceType string `xml:"ReferenceType"`
}

type TransportLegType struct {
	Value string
}

type Consignment struct {
	ConsignmentId string `xml:"consignmentId,attr"`

	DateAndTimes struct {
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
		GoodsValue                 string `xml:",chatdata"`
	} `xml:"GoodsValue"`

	NumberOfPackages struct {
		UnitCode string `xml:"unitCode,attr"`
		Value    string `xml:",chardata"`
	} `xml:"NumberOfPackages"`

	TotalGrossWeight struct {
		UnitCode string `xml:"unitCode,attr"`
		Value    string `xml:",chardata"`
	} `xml:"TotalGrossWeight"`

	TotalVolume struct {
		UnitCode string `xml:"unitCode,attr"`
		Value    string `xml:",chardata"`
	}

	References []Reference `xml:"Reference"`

	TransportLeg []TransportLegType `xml:"TransportLeg>TransportLegType"`
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
				//fmt.Println(tag)
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
				//fmt.Println(tag)
				isUnderConsignment = false
				buffer.WriteString(tag)
				bf := buffer.String()
				bfslice := make([]string, 1)
				bfslice[0] = bf
				consignments = append(consignments, bfslice...)
				buffer = bytes.NewBufferString("")
			} else if isUnderConsignment {
				tag := getEndTag(se)
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
