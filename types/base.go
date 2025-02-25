package types

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/m0ssc0de/scale.go/utiles"
)

type ScaleDecoderOption struct {
	Spec             int
	SubType          string
	Module           string
	ValueList        []string
	Metadata         *MetadataStruct
	FixedLength      int
	SignedExtensions []SignedExtension `json:"signed_extensions"`
	TypeName         string
}

type TypeMapping struct {
	Names []string
	Types []string
}

type IScaleDecoder interface {
	Init(data ScaleBytes, option *ScaleDecoderOption)
	Process()
	Encode(interface{}) string
	buildStruct()
	NextBytes(int) []byte
	GetNextU8() int
	reset()
}

type ScaleDecoder struct {
	Data        ScaleBytes      `json:"-"`
	TypeString  string          `json:"-"`
	SubType     string          `json:"-"`
	Value       interface{}     `json:"-"`
	RawValue    string          `json:"-"`
	TypeMapping *TypeMapping    `json:"-"`
	Metadata    *MetadataStruct `json:"-"`
	Spec        int             `json:"-"`
	Module      string          `json:"-"`
	TypeName    string
}

func (s *ScaleDecoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option != nil {
		if option.Metadata != nil {
			s.Metadata = option.Metadata
		}
		if option.SubType != "" {
			s.SubType = option.SubType
		}
		if option.Spec != 0 {
			s.Spec = option.Spec
		}
		if option.Module != "" {
			s.Module = option.Module
		}
		if option.TypeName != "" {
			s.TypeName = option.TypeName
		}
	}
	s.Data = data
	s.RawValue = ""
	s.Value = nil
	if s.TypeMapping == nil && s.TypeString != "" {
		s.buildStruct()
	}
}

func (s *ScaleDecoder) Process() {}

func (s *ScaleDecoder) Encode(interface{}) string { return "" }

func (s *ScaleDecoder) NextBytes(length int) []byte {
	data := s.Data.GetNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) GetNextU8() int {
	b := s.NextBytes(1)
	if len(b) > 0 {
		return int(b[0])
	}
	return 0
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.NextBytes(1)
	return utiles.BytesToHex(data) == "01"
}

func (s *ScaleDecoder) reset() {
	s.Data.Data = []byte{}
	s.Data.Offset = 0
}

func (s *ScaleDecoder) buildStruct() {
	if s.TypeString != "" && string(s.TypeString[0]) == "(" && s.TypeString[len(s.TypeString)-1:] == ")" {

		var names, types []string
		reg := regexp.MustCompile(`[\<\(](.*?)[\>\)]`)
		typeString := s.TypeString[1 : len(s.TypeString)-1]
		typeParts := reg.FindAllString(typeString, -1)
		for _, part := range typeParts {
			typeString = strings.ReplaceAll(typeString, part, strings.ReplaceAll(part, ",", "#"))
		}

		for k, v := range strings.Split(typeString, ",") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			types = append(types, strings.ReplaceAll(strings.TrimSpace(v), "#", ","))
			names = append(names, fmt.Sprintf("col%d", k+1))
		}

		s.TypeMapping = &TypeMapping{Names: names, Types: types}
	}
}

func (s *ScaleDecoder) ProcessAndUpdateData(typeString string) interface{} {
	r := RuntimeType{Module: s.Module}

	if TypeRegistry == nil {
		r.Reg()
	}

	class, value, subType := r.DecoderClass(typeString, s.Spec)
	if class == nil {
		panic(fmt.Sprintf("Not found decoder class %s", typeString))
	}

	offsetStart := s.Data.Offset

	// init
	method, exist := class.MethodByName("Init")
	if !exist {
		panic(fmt.Sprintf("%s not implement init function", typeString))
	}
	option := ScaleDecoderOption{SubType: subType, Spec: s.Spec, Metadata: s.Metadata, Module: s.Module, TypeName: typeString}
	method.Func.Call([]reflect.Value{value, reflect.ValueOf(s.Data), reflect.ValueOf(&option)})

	// process
	value.MethodByName("Process").Call(nil)
	elementData := value.Elem().FieldByName("Data").Interface().(ScaleBytes)

	s.Data.Offset = elementData.Offset
	s.Data.Data = elementData.Data
	s.RawValue = utiles.BytesToHex(s.Data.Data[offsetStart:s.Data.Offset])

	return value.Elem().FieldByName("Value").Interface()
}

func Encode(typeString string, data interface{}) string {
	r := RuntimeType{}
	if TypeRegistry == nil {
		r.Reg()
	}
	class, value, _ := r.DecoderClass(typeString, -1)
	if class == nil {
		panic(fmt.Sprintf("Not found decoder class %s", typeString))
	}
	out := value.MethodByName("Encode").Call([]reflect.Value{reflect.ValueOf(data)})
	return out[0].String()
}
