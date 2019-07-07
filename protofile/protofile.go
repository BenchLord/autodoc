package protofile

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name        string
	Req         Message
	IsReqStream bool
	Res         Message
	IsResStream bool
}

type Message struct {
	Name       string
	Attributes []Attribute
}

type Attribute struct {
	Name          string
	AttributeType string
	Repeated      bool
	Number        int
}

type ProtoFile struct {
	services []Service
	methods  []Method
	messages []Message
}

/*
125: }
123: {
10: newline
40: (
41: )
32: space
59: ;
*/

func NewFile(name string) (*ProtoFile, error) {

	//methodMap := make(map[string][]Attribute)

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	pf := ProtoFile{}

	statements, err := makeStatementList(f)
	if err != nil {
		return nil, err
	}

	for _, statement := range statements {
		if strings.HasPrefix(statement, "message") {
			pf.messages = append(pf.messages, parseMessageStatement(statement))
		}

	}
	fmt.Println(pf.messages)
	return &pf, nil
}

func makeStatementList(f *os.File) ([]string, error) {
	statements := make([]string, 0)
	statement := make([]byte, 0)
	open := false
	for {
		buffer := make([]byte, 1)
		_, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		char := buffer[0]
		if char == 123 {
			open = true
		}
		if char == 125 {
			open = false
		}
		if char == 10 {
			if !open {
				if len(statement) > 0 {
					statements = append(statements, string(statement))
				}
				statement = make([]byte, 0)
			}
			continue
		}
		statement = append(statement, buffer[0])
	}
	return statements, nil
}

func parseMessageStatement(s string) Message {
	s = strings.Trim(s, "message ")

	name := strings.Split(s, " ")[0]

	attributes := make([]Attribute, 0)
	first := strings.Index(s, "{") + 1
	last := strings.Index(s, "}")
	attributesStatement := s[first:last]
	attributeStatements := strings.Split(attributesStatement, ";")
	for _, attribute := range attributeStatements {
		attribute = strings.Trim(attribute, " ")
		if len(attribute) == 0 {
			continue
		}
		// get name, type, and num from string
		attributes = append(attributes, Attribute{})
	}

	return Message{
		Name:       name,
		Attributes: attributes,
	}
}
