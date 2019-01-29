// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package dictionary

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

// ACHDictionary of Participant records
type ACHDictionary struct {
	ACHDictionary Dictionary
}

// NewACHDictionary creates a ACHDictionary
func NewACHDictionary(r io.Reader) *ACHDictionary {
	achD := new(ACHDictionary)
	achD.ACHDictionary.IndexRoutingNumber = make(map[string]*Participant)
	achD.ACHDictionary.IndexCustomerName = make(map[string][]*Participant)
	achD.ACHDictionary.scanner = bufio.NewScanner(r)
	return achD
}

// Read parses a single line or multiple lines of FedACHdir text
func (achD *ACHDictionary) Read() error {
	// read through the entire file
	for achD.ACHDictionary.scanner.Scan() {
		achD.ACHDictionary.line = achD.ACHDictionary.scanner.Text()

		if utf8.RuneCountInString(achD.ACHDictionary.line) != 155 {
			achD.ACHDictionary.errors.Add(NewRecordWrongLengthErr(155, len(achD.ACHDictionary.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return achD.ACHDictionary.errors
		}
		if err := achD.parseParticipant(); err != nil {
			achD.ACHDictionary.errors.Add(err)
			return achD.ACHDictionary.errors
		}
	}
	if err := achD.createIndexCustomerName(); err != nil {
		achD.ACHDictionary.errors.Add(err)
		return achD.ACHDictionary.errors
	}
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (achD *ACHDictionary) parseParticipant() error {
	p := new(Participant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = achD.ACHDictionary.line[:9]
	// OfficeCode (1): O
	p.OfficeCode = achD.ACHDictionary.line[9:10]
	// ServicingFrbNumber (9): 011000015
	p.ServicingFrbNumber = achD.ACHDictionary.line[10:19]
	// RecordTypeCode (1): 0
	p.RecordTypeCode = achD.ACHDictionary.line[19:20]
	// ChangeDate (6): 122415
	p.ACHRevisedDate = achD.ACHDictionary.line[20:26]
	// NewRoutingNumber (9): 000000000
	p.NewRoutingNumber = achD.ACHDictionary.line[26:35]
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(achD.ACHDictionary.line[35:71], " ")
	// Address (36): 1000 PEACHTREE ST N.E.
	p.ACHAddress = strings.Trim(achD.ACHDictionary.line[71:107], " ")
	// City (20): ATLANTA
	p.ACHCity = strings.Trim(achD.ACHDictionary.line[107:127], " ")
	// State (2): GA
	p.ACHState = achD.ACHDictionary.line[127:129]
	// PostalCode (5): 30309
	p.ACHPostalCode = achD.ACHDictionary.line[129:134]
	// PostalCodeExtension (4): 4470
	p.ACHPostalCodeExtension = achD.ACHDictionary.line[134:138]
	// PhoneNumber(10): 8773722457
	p.PhoneNumber = achD.ACHDictionary.line[138:148]
	// StatusCode (1): 1
	p.StatusCode = achD.ACHDictionary.line[148:149]
	// ViewCode (1): 1
	p.ViewCode = achD.ACHDictionary.line[149:150]

	achD.ACHDictionary.Participants = append(achD.ACHDictionary.Participants, p)
	achD.ACHDictionary.IndexRoutingNumber[p.RoutingNumber] = p
	return nil
}

// createIndexACHCustomerName creates an index of Financial Institutions keyed by ACHParticipant.CustomerName
func (achD *ACHDictionary) createIndexCustomerName() error {
	for _, achP := range achD.ACHDictionary.Participants {
		achD.ACHDictionary.IndexCustomerName[achP.CustomerName] = append(achD.ACHDictionary.IndexCustomerName[achP.CustomerName], achP)
	}
	return nil
}

// RoutingNumberSearch returns a FEDACH participant based on a ACHParticipant.RoutingNumber.  Routing Number validation
// is only that it exists in IndexParticipant.  Expecting 9 digits, checksum needs to be included.
func (achD *ACHDictionary) RoutingNumberSearch(s string) *Participant {
	if _, ok := achD.ACHDictionary.IndexRoutingNumber[s]; ok {
		return achD.ACHDictionary.IndexRoutingNumber[s]
	}
	return nil
}

// FinancialInstitutionSearch returns a FEDACH participant based on a ACHParticipant.CustomerName
func (achD *ACHDictionary) FinancialInstitutionSearch(s string) []*Participant {
	if _, ok := achD.ACHDictionary.IndexCustomerName[s]; ok {
		return achD.ACHDictionary.IndexCustomerName[s]
	}
	return nil
}

// GetParticipants returns a slice of participants for the dictionary
func (achD *ACHDictionary) GetParticipants() []*Participant {
	return achD.ACHDictionary.Participants
}

func (achD *ACHDictionary) GetIndexRoutingNumber() map[string]*Participant {
	return achD.ACHDictionary.IndexRoutingNumber
}

func (achD *ACHDictionary) GetIndexCustomerName() map[string][]*Participant {
	return achD.ACHDictionary.IndexCustomerName
}
