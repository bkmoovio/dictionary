package dictionary

import (
	"testing"
)

func TestACHDirectoryRead(t *testing.T) {

	fedACH := File{
		FilePath: "./data/FedACHdir.txt",
		FileName: "FedACHdir.txt",
		FileType: ".txt",
	}

	dictionary, err := NewDictionary(fedACH)

	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}

	dictP := dictionary.GetParticipants()

	if len(dictP) != 18198 {
		t.Errorf("Expected '18198' got: %v", len(dictionary.GetParticipants()))
	}

	if dictP[0].CustomerName != "FEDERAL RESERVE BANK" {
		t.Errorf("Expected `FEDERAL RESERVE BANK` got : %v", dictP[0].CustomerName)
	}

	if fi, ok := dictionary.GetIndexRoutingNumber()["073905527"]; ok {
		if fi.CustomerName != "LINCOLN SAVINGS BANK" {
			t.Errorf("Expected `LINCOLN SAVINGS BANK` got : %v", fi.CustomerName)
		}
	} else {
		t.Errorf("ach routing number `073905527` not found")
	}

	if fi, ok := dictionary.GetIndexCustomerName()["LOWER VALLEY CU"]; ok {
		for _, f := range fi {
			if f.CustomerName != "LOWER VALLEY CU" {
				t.Errorf("Expected `LOWER VALLEY CU` got : %v", f.CustomerName)
			}
		}
	} else {
		t.Errorf("Customer Name `LOWER VALLEY CU` not found")
	}

}

func TestWIREDirectoryRead(t *testing.T) {

	fedWIRE := File{
		FilePath: "./data/fpddir.txt",
		FileName: "fpddir.txt",
		FileType: ".txt",
	}

	dictionary, err := NewDictionary(fedWIRE)

	if err != nil {
		t.Fatalf("%T: %s", err, err)
	}

	dictP := dictionary.GetParticipants()

	if len(dictP) != 7693 {
		t.Errorf("Expected '7693' got: %v", len(dictionary.GetParticipants()))
	}

	if dictP[0].CustomerName != "FEDERAL RESERVE BANK OF BOSTON" {
		t.Errorf("Expected `FEDERAL RESERVE BANK OF BOSTON` got : %v", dictP[0].CustomerName)
	}

	if fi, ok := dictionary.GetIndexRoutingNumber()["325280039"]; ok {
		if fi.TelegraphicName != "MAC FCU" {
			t.Errorf("Expected `MAC FCU` got : %v", fi.CustomerName)
		}
	} else {
		t.Errorf("ach routing number `325280039` not found")
	}

	if fi, ok := dictionary.GetIndexCustomerName()["TRUGROCER FEDERAL CREDIT UNION"]; ok {
		for _, f := range fi {
			if f.CustomerName != "TRUGROCER FEDERAL CREDIT UNION" {
				t.Errorf("Expected `TRUGROCER FEDERAL CREDIT UNION` got : %v", f.CustomerName)
			}
		}
	} else {
		t.Errorf("Customer Name `TRUGROCER FEDERAL CREDIT UNION` not found")
	}
}
