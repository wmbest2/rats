package test

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

type NullFloat64 struct {
	sql.NullFloat64
}

type NullBool struct {
	sql.NullBool
}

type NullString struct {
	sql.NullString
}

func NewNullString(val string) NullString {
	return NullString{sql.NullString{val, true}}
}

func NewNullBool(val bool) NullBool {
	return NullBool{sql.NullBool{val, true}}
}

/*
 * JSON Marshalling
 */

//MarshalJSON is a wrapper around sql.NullString that satifies json.Marshaler
func (nf NullString) MarshalJSON() ([]byte, error) {
	var data interface{}

	if !nf.Valid {
		data = nil
	} else {
		data = nf.String
	}

	return json.Marshal(data)
}

func (nf NullString) UnmarshalJSON(data []byte) error {
	nf.String = string(data)
	nf.Valid = true
	return nil
}

//MarshalJSON is a wrapper around sql.NullBool that satifies json.Marshaler
func (nf NullBool) MarshalJSON() ([]byte, error) {
	var data interface{}

	if !nf.Valid {
		data = nil
	} else {
		data = nf.Bool
	}

	return json.Marshal(data)
}

func (nf NullBool) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseBool(string(data))
	nf.Bool = val
	nf.Valid = err != nil
	return nil
}

//MarshalJSON is a wrapper around sql.NullFloat64 that satifies json.Marshaler
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	var data interface{}

	if !nf.Valid {
		data = nil
	} else {
		data = nf.Float64
	}

	return json.Marshal(data)
}

func (nf NullFloat64) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseFloat(string(data), 64)
	nf.Float64 = val
	nf.Valid = err != nil
	return nil
}

/*
 * XML Marshalling
 */

func (nf *NullFloat64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var attr xml.Attr

	if !nf.Valid {
		attr = xml.Attr{name, "0"}
	} else {
		attr = xml.Attr{name, strconv.FormatFloat(nf.Float64, 'g', -1, 64)}
	}

	return attr, nil
}

func (nf *NullFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	data, err := strconv.ParseFloat(attr.Value, 64)
	nf.Float64 = data
	nf.Valid = true
	return err
}

func (nf *NullString) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if nf.Valid {
		return e.EncodeElement(nf.String, start)
	}
	return nil
}

func (nf *NullString) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var attr xml.Attr

	if !nf.Valid {
		//attr = nil //xml.Attr{name, nil}
	} else {
		attr = xml.Attr{name, nf.String}
	}

	return attr, nil
}

func (nf *NullString) UnmarshalXMLAttr(attr xml.Attr) error {
	nf.String = attr.Value
	nf.Valid = true
	return nil
}
