package util

import (
	"bufio"
	"bytes"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

type SchemaDecoder struct {
	data []byte
}

func NewSchemaDecoder(data []byte) *SchemaDecoder {
	return &SchemaDecoder{data: data}
}

func (sc *SchemaDecoder) Decode() (runtime.Object, error) {
	decoder := scheme.Codecs.UniversalDeserializer()
	object, _, err := decoder.Decode(sc.data, nil, nil)
	return object, err
}

func (sc *SchemaDecoder) MultiDecode() ([]runtime.Object, error) {
	decoder := yaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(sc.data)))
	var objects []runtime.Object
	for {
		obj, err := decoder.Read()
		if len(obj) == 0 {
			break
		}
		if err != nil {
			log.Error("error reading from yaml ", err)
			return nil, err
		}
		sd := NewSchemaDecoder(obj)
		object, err := sd.Decode()
		if err != nil {
			log.Error("error decoding multi yaml ", err)
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}
