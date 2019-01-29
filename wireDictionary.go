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

// WIREDictionary of Participant records
type WIREDictionary struct {
	WIREDictionary Dictionary
}

// NewWIREDictionary creates a WIREDictionary
func NewWIREDictionary(r io.Reader) *WIREDictionary {
	wireD := new(WIREDictionary)
	wireD.WIREDictionary.IndexRoutingNumber = make(map[string]*Participant)
	wireD.WIREDictionary.IndexCustomerName = make(map[string][]*Participant)
	wireD.WIREDictionary.scanner = bufio.NewScanner(r)
	return wireD

}

// Read parses a single line or multiple lines of FedWIREdir text
func (wireD *WIREDictionary) Read() error {
	// read through the entire file
	for wireD.WIREDictionary.scanner.Scan() {
		wireD.WIREDictionary.line = wireD.WIREDictionary.scanner.Text()

		if utf8.RuneCountInString(wireD.WIREDictionary.line) != 101 {
			wireD.WIREDictionary.errors.Add(NewRecordWrongLengthErr(101, len(wireD.WIREDictionary.line)))
			// Return with error if the record length is incorrect as this file is a FED file
			return wireD.WIREDictionary.errors
		}
		if err := wireD.parseParticipant(); err != nil {
			wireD.WIREDictionary.errors.Add(err)
			return wireD.WIREDictionary.errors
		}
	}
	if err := wireD.createIndexCustomerName(); err != nil {
		wireD.WIREDictionary.errors.Add(err)
		return wireD.WIREDictionary.errors
	}
	return nil
}

// TODO return a parsing error if the format or file is wrong.
func (wireD *WIREDictionary) parseParticipant() error {
	p := new(Participant)

	//RoutingNumber (9): 011000015
	p.RoutingNumber = wireD.WIREDictionary.line[:9]
	// TelegraphicName (18): FED
	p.TelegraphicName = strings.Trim(wireD.WIREDictionary.line[9:27], " ")
	// CustomerName (36): FEDERAL RESERVE BANK
	p.CustomerName = strings.Trim(wireD.WIREDictionary.line[27:63], " ")
	// State (2): GA
	p.WIREState = wireD.WIREDictionary.line[63:65]
	// City (25): ATLANTA
	p.WIRECity = strings.Trim(wireD.WIREDictionary.line[65:90], " ")
	// FundsTransferStatus (1): Y or N
	p.FundsTransferStatus = wireD.WIREDictionary.line[90:91]
	// FundsSettlementOnlyStatus (1): " " or S - Settlement-Only
	p.FundsSettlementOnlyStatus = wireD.WIREDictionary.line[91:92]
	// BookEntrySecuritiesTransferStatus (1): Y or N
	p.BookEntrySecuritiesTransferStatus = wireD.WIREDictionary.line[92:93]
	// Date YYYYMMDD (8): 122415
	p.WIRERevisedDate = wireD.WIREDictionary.line[93:101]
	wireD.WIREDictionary.Participants = append(wireD.WIREDictionary.Participants, p)
	wireD.WIREDictionary.IndexRoutingNumber[p.RoutingNumber] = p
	return nil
}

// createIndexWIRECustomerName creates an index of Financial Institutions keyed by ACHParticipant.CustomerName
func (wireD *WIREDictionary) createIndexCustomerName() error {
	for _, wireP := range wireD.WIREDictionary.Participants {
		wireD.WIREDictionary.IndexCustomerName[wireP.CustomerName] = append(wireD.WIREDictionary.IndexCustomerName[wireP.CustomerName], wireP)
	}
	return nil
}

// RoutingNumberSearch returns a FEDWIRE participant based on a Participant.RoutingNumber
func (wireD *WIREDictionary) RoutingNumberSearch(s string) *Participant {
	if _, ok := wireD.WIREDictionary.IndexRoutingNumber[s]; ok {
		return wireD.WIREDictionary.IndexRoutingNumber[s]
	}
	return nil
}

// FinancialInstitutionSearch returns a FEDACH participant based on a WIREParticipant.CustomerName
func (wireD *WIREDictionary) FinancialInstitutionSearch(s string) []*Participant {
	if _, ok := wireD.WIREDictionary.IndexCustomerName[s]; ok {
		return wireD.WIREDictionary.IndexCustomerName[s]
	}
	return nil
}

// GetParticipants returns a slice of participants for the dictionary
func (wireD *WIREDictionary) GetParticipants() []*Participant {
	return wireD.WIREDictionary.Participants
}

func (wireD *WIREDictionary) GetIndexRoutingNumber() map[string]*Participant {
	return wireD.WIREDictionary.IndexRoutingNumber
}

func (wireD *WIREDictionary) GetIndexCustomerName() map[string][]*Participant {
	return wireD.WIREDictionary.IndexCustomerName
}
