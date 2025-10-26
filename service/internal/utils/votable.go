package utils

import (
	"encoding/xml"
)

// VOTable represents the VOTable structure
type VOTable struct {
	XMLName  xml.Name `xml:"VOTABLE"`
	Version  string   `xml:"version,attr"`
	Xmlns    string   `xml:"xmlns,attr"`
	Resource Resource `xml:"RESOURCE"`
}

// Resource represents a RESOURCE element in VOTable
type Resource struct {
	Type   string   `xml:"type,attr"`
	Infos  []Info   `xml:"INFO"`
	Params []Param  `xml:"PARAM"`
	Tables []Table  `xml:"TABLE"`
	Coosys []Coosys `xml:"COOSYS"`
}

// Info represents an INFO element in VOTable
type Info struct {
	Name        string `xml:"name,attr"`
	Value       string `xml:"value,attr"`
	Description string `xml:"DESCRIPTION,omitempty"`
}

// Param represents a PARAM element in VOTable
type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Unit  string `xml:"unit,attr,omitempty"`
	Ucd   string `xml:"ucd,attr,omitempty"`
}

// Table represents a TABLE element in VOTable
type Table struct {
	Name        string  `xml:"name,attr"`
	Description string  `xml:"DESCRIPTION,omitempty"`
	Fields      []Field `xml:"FIELD"`
	Params      []Param `xml:"PARAM"`
	Groups      []Group `xml:"GROUP"`
	Data        Data    `xml:"DATA"`
}

// Coosys represents a COOSYS element in VOTable
type Coosys struct {
	ID      string `xml:"ID,attr,omitempty"`
	Equinox string `xml:"equinox,attr,omitempty"`
	System  string `xml:"system,attr,omitempty"`
	Epoch   string `xml:"epoch,attr,omitempty"`
}

// Field represents a FIELD element in VOTable
type Field struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"DESCRIPTION,omitempty"`
	ID          string `xml:"ID,attr,omitempty"`
	Datatype    string `xml:"datatype,attr"`
	Unit        string `xml:"unit,attr,omitempty"`
	Ucd         string `xml:"ucd,attr,omitempty"`
	ArraySize   string `xml:"arraysize,attr,omitempty"`
}

// Group represents a GROUP element in VOTable
type Group struct {
	ID     string     `xml:"ID,attr,omitempty"`
	Name   string     `xml:"name,attr,omitempty"`
	Fields []FieldRef `xml:"FIELDref"`
	Params []Param    `xml:"PARAM"`
}

// FieldRef represents a FIELDref element in VOTable
type FieldRef struct {
	Ref string `xml:"ref,attr"`
}

// Data represents a DATA element in VOTable
type Data struct {
	TableData TableData `xml:"TABLEDATA"`
}

// TableData represents a TABLEDATA element in VOTable
type TableData struct {
	Rows []Row `xml:"TR"`
}

// Row represents a TR element in VOTable
type Row struct {
	Columns []Column `xml:"TD"`
}

// Column represents a TD element in VOTable
type Column struct {
	Value string `xml:",chardata"`
}

// NewVOTable creates a new VOTable from string
func NewVOTableFromString(xmlRepr string) (*VOTable, error) {
	var votable VOTable
	err := xml.Unmarshal([]byte(xmlRepr), &votable)
	if err != nil {
		return nil, err
	}
	return &votable, nil
}

// NewVOTableFromBytes creates a new VOTable from bytes
func NewVOTableFromBytes(xmlRepr []byte) (*VOTable, error) {
	var votable VOTable
	err := xml.Unmarshal(xmlRepr, &votable)
	if err != nil {
		return nil, err
	}
	return &votable, nil
}
