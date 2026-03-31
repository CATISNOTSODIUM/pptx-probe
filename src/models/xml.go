package models

import (
	"bytes"
	"encoding/xml"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`

	// Optional fields
	Id    string `xml:"id,attr"`
	Level string `xml:"lvl,attr"`

	// Extra fields for convenience
	Parent *Node
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func Decode(data []byte) (*Node, error) {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
