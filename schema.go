package memphis

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

const (
	schemaObjectName = "Schema"
)

type Schema struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	CreatedByUsername string `json:"created_by_username"`
	SchemaContent     string `json:"schema_content"`
	MessageStructName string `json:"message_struct_name"`
}

type createSchemaReq struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	CreatedByUsername string `json:"created_by_username"`
	SchemaContent     string `json:"schema_content"`
	MessageStructName string `json:"message_struct_name"`
}

type createSchemaResp struct {
	Err string `json:"error"`
}

type removeSchemaReq struct {
	Name string `json:"name"`
}

func (s *Schema) getCreationSubject() string {
	return "$memphis_schema_creations"
}

func (s *Schema) getDestructionSubject() string {
	return ""
}

func (s *Schema) getCreationReq() any {
	return createSchemaReq{
		Name:              s.Name,
		Type:              s.Type,
		CreatedByUsername: s.CreatedByUsername,
		SchemaContent:     s.SchemaContent,
	}
}

func (s *Schema) handleCreationResp(resp []byte) error {
	cr := &createSchemaResp{}
	err := json.Unmarshal(resp, cr)
	if err != nil {
		return defaultHandleCreationResp(resp)
	}

	if cr.Err != "" {
		return memphisError(errors.New(cr.Err))
	}
	return nil
}

func (s *Schema) getDestructionReq() any {
	return nil
}

// CreateSchema - validates and uploads a new schema to the Broker
func (c *Conn) CreateSchema(name, schemaType, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return memphisError(err)
	}

	schemaContent := string(data)

	err = validateSchemaName(name)
	if err != nil {
		return memphisError(err)
	}

	err = validateSchemaType(schemaType)
	if err != nil {
		return memphisError(err)
	}

	s := Schema{
		Name:              name,
		Type:              schemaType,
		CreatedByUsername: c.username,
		SchemaContent:     schemaContent,
		MessageStructName: "",
	}

	if err = c.create(&s); err != nil {
		return memphisError(err)
	}

	return nil
}

func validateSchemaName(schemaName string) error {
	return validateName(schemaName, schemaObjectName)
}

func validateName(name, objectType string) error {
	emptyErrStr := fmt.Sprintf("%v name can not be empty", objectType)
	tooLongErrStr := fmt.Sprintf("%v should be under 128 characters", objectType)
	invalidCharErrStr := fmt.Sprintf("Only alphanumeric and the '_', '-', '.' characters are allowed in %v", objectType)
	firstLetterErrStr := fmt.Sprintf("%v name can not start or end with non alphanumeric character", objectType)

	emptyErr := errors.New(emptyErrStr)
	tooLongErr := errors.New(tooLongErrStr)
	invalidCharErr := errors.New(invalidCharErrStr)
	firstLetterErr := errors.New(firstLetterErrStr)

	if len(name) == 0 {
		return emptyErr
	}

	if len(name) > 128 {
		return tooLongErr
	}

	re := regexp.MustCompile("^[a-z0-9_.-]*$")

	validName := re.MatchString(name)
	if !validName {
		return invalidCharErr
	}

	if name[0:1] == "." || name[0:1] == "-" || name[0:1] == "_" || name[len(name)-1:] == "." || name[len(name)-1:] == "-" || name[len(name)-1:] == "_" {
		return firstLetterErr
	}

	return nil
}

func validateSchemaType(schemaType string) error {
	invalidTypeErrStr := "unsupported schema type"
	invalidTypeErr := errors.New(invalidTypeErrStr)
	invalidSupportTypeErrStr := "avro is not supported at this time"
	invalidSupportTypeErr := errors.New(invalidSupportTypeErrStr)

	if schemaType == "protobuf" || schemaType == "json" || schemaType == "graphql" {
		return nil
	} else if schemaType == "avro" {
		return invalidSupportTypeErr
	} else {
		return invalidTypeErr
	}
}

// TODO: autoRegSchema()
func autoRegSchema(msg any, stationName string) error {
	// 1. TODO: read flag ScanMsgChema. Where to store that flag?  if flag ScanMsgSchema is true then
	
	// 2. TODO: scan message to return schema
	if opts.scanMsgSchema {
		err := scanMsgSchema(opt.Message)
	}
	// 3. Save new schema
	func (c *Conn) CreateSchema(name, schemaType, path string) error {}
	// 4. Enforce schema - attache schema to station
	func (c *Conn) EnforceSchema(name string, stationName string) error {}	
}

// TODO: func scanMsgChema(msg any)
func scanMsgChema(msg any) {
// TODO: getMsgFormat() - validate msg type to check if msg type is supported or not and to get the type
// TODO: readSchema() - once type is known, read schema of that specific type
}

// TODO: extractSchema takes an object and generates the schema
func extractSchema(message, schemaType string) (schemaContent, schemaType string) {
	var schemaContent string
	switch schemaType {
	case "protobuf":
		// Protocol Buffer Descriptors
		// https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect
	case "json":
		// genJsonSchema generates json schema in the format described on json-schema.org
		// https://github.com/xeipuuv/gojsonschema
		// https://github.com/invopop/jsonschema
		// https://git.sr.ht/~emersion/go-jsonschema
		// https://github.com/mcuadros/go-jsonschema-generator
		// https://github.com/fybrik/json-schema-generator
	case "graphql":
		// https://gqlgen.com/feature-comparison/
	default:
		return schemaContent, schemaType
	}
	//these are on the broker schema handler:	
	"github.com/graph-gophers/graphql-go"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/santhosh-tekuri/jsonschema/v5"
}