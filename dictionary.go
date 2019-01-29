// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package dictionary

import (
	"bufio"
	"github.com/moov-io/base"
	"os"
)

type Dictionary struct {
	// File being parsed
	File File
	// Participants is a list of Participant structs
	Participants []*Participant
	//scanner provides a convenient interface for reading data
	scanner *bufio.Scanner
	// line being read
	line string
	// IndexACHRoutingNumber creates an index of ACHParticipants keyed by ACHParticipant.RoutingNumber
	IndexRoutingNumber map[string]*Participant
	// IndexACHCustomerName creates an index of ACHParticipants keyed by ACHParticipant.CustomerName
	IndexCustomerName map[string][]*Participant
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
}

// Participant holds a FedACH dir routing record as defined by Fed ACH Format
// https://www.frbservices.org/EPaymentsDirectory/achFormat.html
type Participant struct {
	// ToDo: Challenges of validation with the abstraction, if each is concrete, only one domain can be affected
	// Fed General
	// RoutingNumber The institution's routing number
	RoutingNumber string `json:"routingNumber"`
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string `json:"customerName"`
	// Location is the delivery address for ACH or WIRE
	Location `json:"achlocation"`

	// ACH
	// OfficeCode Main/Head Office or Branch. O=main B=branch
	OfficeCode string `json:"officeCode"`
	// ServicingFrbNumber Servicing Fed's main office routing number
	ServicingFrbNumber string `json:"servicingFrbNumber"`
	// RecordTypeCode The code indicating the ABA number to be used to route or send ACH items to the RFI
	// 0 = Institution is a Federal Reserve Bank
	// 1 = Send items to customer routing number
	// 2 = Send items to customer using new routing number field
	RecordTypeCode string `json:"recordTypeCod"`
	// Revised Date of last revision: YYYYMMDD, or blank
	ACHRevisedDate string `json:"achRevisedDate"`
	// NewRoutingNumber Institution's new routing number resulting from a merger or renumber
	NewRoutingNumber string `json:"newRoutingNumber"`
	// PhoneNumber The institution's phone number
	PhoneNumber string `json:"phoneNumber"`
	// StatusCode Code is based on the customers receiver code
	// 1=Receives Gov/Comm
	StatusCode string `json:"statusCode"`
	// ViewCode
	ViewCode string `json:"viewCode"`

	// WIRE
	// TelegraphicName is the short name of financial institution  Wells Fargo
	TelegraphicName string `json:"telegraphicName"`
	// FundsTransferStatus designates funds transfer status
	// Y - Eligible
	// N - Ineligible
	FundsTransferStatus string `json:"fundsTransferStatus"`
	// FundsSettlementOnlyStatus designates funds settlement only status
	// S - Settlement-Only
	FundsSettlementOnlyStatus string `json:"fundsSettlementOnlyStatus"`
	// BookEntrySecuritiesTransferStatus designates book entry securities transfer status
	BookEntrySecuritiesTransferStatus string `json:"bookEntrySecuritiesTransferStatus"`
	// Date of last revision: YYYYMMDD, or blank
	WIRERevisedDate string `json:"WIRERevisedDate"`
}

// Location City name and state code in the institution's delivery address
// ToDo: Challenges of validation with the abstraction, if each is concrete, only one domain can be affected
type Location struct {
	// Address
	ACHAddress string `json:"achAddress"`
	// City
	ACHCity string `json:"achCity"`
	// State
	ACHState string `json:"achState"`
	// PostalCode
	ACHPostalCode string `json:"achPostalCode"`
	// PostalCodeExtension
	ACHPostalCodeExtension string `json:"achPostalCodeExtension"`
	// City
	WIRECity string `json:"wireCity"`
	// State
	WIREState string `json:"wireState"`
}

// ToDo: Relate the Data?

// NewDictionary creates a Dictionary
func NewDictionary(file File) (Diction, error) {
	//f, err := os.Open("./data/FedACHdir.txt")
	f, err := os.Open(file.FilePath)
	if err != nil {
		file.errors.Add(err)
		return nil, file.errors
	}
	defer f.Close()

	switch file.FileName {
	case "FedACHdir.txt":
		achD := NewACHDictionary(f)
		err := achD.Read()
		if err != nil {
			file.errors.Add(err)
			return nil, file.errors
		}
		return achD, nil
	case "fpddir.txt":
		wireD := NewWIREDictionary(f)
		err := wireD.Read()
		if err != nil {
			file.errors.Add(err)
			return nil, file.errors
		}
		return wireD, nil
	}

	file.errors.Add(ErrCreateDictionary)
	return nil, file.errors
}
