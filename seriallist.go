/*

This file is part of serial-port-json-server. (https://github.com/chilipeppr/serial-port-json-server)

serial-port-json-server is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later version.

serial-port-json-server is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
serial-port-json-server. If not, see http://www.gnu.org/licenses/.

*/

// Supports Windows, Linux, Mac, and Raspberry Pi

package main

import (
	"encoding/xml"
	"strings"

	"github.com/facchinm/go-serial"
	//"fmt"

	//"io/ioutil"
	"log"
	//"os"
)

type OsSerialPort struct {
	Name         string
	FriendlyName string
	RelatedNames []string // for some devices there are 2 or more ports, i.e. TinyG v9 has 2 serial ports
	SerialNumber string
	DeviceClass  string
	Manufacturer string
	Product      string
	IdProduct    string
	IdVendor     string
}

func GetList() ([]OsSerialPort, error) {

	//log.Println("Doing GetList()")

	ports, err := serial.GetPortsList()

	arrPorts := []OsSerialPort{}
	for _, element := range ports {
		friendly := strings.Replace(element, "/dev/", "", -1)
		arrPorts = append(arrPorts, OsSerialPort{Name: element, FriendlyName: friendly})
	}

	//log.Printf("Done doing GetList(). arrPorts:%v\n", arrPorts)

	return arrPorts, err
}

func GetMetaList() ([]OsSerialPort, error) {
	metaportlist, err := getMetaList()
	if err.Err != nil {
		return nil, err.Err
	}
	return metaportlist, err.Err
}

func GetFriendlyName(portname string) string {
	log.Println("GetFriendlyName from base class")
	return ""
}

type Dict struct {
	Keys    []string `xml:"key"`
	Arrays  []Dict   `xml:"array"`
	Strings []string `xml:"string"`
	Dicts   []Dict   `xml:"dict"`
}

type Result struct {
	XMLName xml.Name `xml:"plist"`
	//Strings []string `xml:"dict>string"`
	Dict `xml:"dict"`
	//Phone   string
	//Groups  []string `xml:"Group>Value"`
}
