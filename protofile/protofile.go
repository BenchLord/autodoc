package protofile

import (
	"io"
	"os"
	"strconv"
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
	messages []Message
}

// GetServices returns the services defined in the file
func (p *ProtoFile) GetServices() []Service {
	return p.services
}

// GetMessages returns the messages defined in the file
func (p *ProtoFile) GetMessages() []Message {
	return p.messages
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
		if strings.HasPrefix(statement, "service") {
			defer func() {
				service := parseServiceStatement(statement, &pf)
				pf.services = append(pf.services, service)
			}()
		}
		if strings.HasPrefix(statement, "message") {
			pf.messages = append(pf.messages, parseMessageStatement(statement))
		}
	}
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

func parseServiceStatement(s string, pf *ProtoFile) Service {
	s = strings.Trim(s, "service ")
	service := Service{}

	service.Name = strings.Split(s, " ")[0]

	methods := strings.Split(strings.Trim(s[strings.Index(s, "{"):], "{ }"), ";")
	for _, method := range methods {
		if method == "" {
			continue
		}
		service.Methods = append(service.Methods, parseMethodStatement(method, pf))
	}
	return service
}

func parseMethodStatement(s string, pf *ProtoFile) Method {
	s = strings.Trim(s, "prc ")
	method := Method{}

	// get name

	// check if req is repeated
	// use req name to get method info

	// check if res is repeated
	// use res name to get method info

	return method
}

func parseMessageStatement(s string) Message {
	s = strings.Trim(s, "message ")

	name := strings.Split(s, " ")[0]
	message := Message{Name: name, Attributes: []Attribute{}}

	first := strings.Index(s, "{") + 1
	last := strings.Index(s, "}")
	s = s[first:last]
	attributes := strings.Split(s, ";")
	for _, attribute := range attributes {
		attribute = strings.Trim(attribute, " ")
		if len(attribute) == 0 {
			continue
		}
		a := parseAttributeStatement(attribute)

		message.Attributes = append(message.Attributes, a)
	}

	return message
}

func parseAttributeStatement(s string) Attribute {
	attribute := Attribute{}
	words := strings.Split(s, " ")

	if strings.ToLower(words[0]) == "repeated" {
		attribute.Repeated = true
		words = words[1:]
	}

	attribute.AttributeType = words[0]
	attribute.Name = words[1]
	num, _ := strconv.Atoi(words[3])
	attribute.Number = num
	return attribute
}
