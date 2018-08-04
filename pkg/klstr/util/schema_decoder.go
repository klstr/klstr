package util

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

type SchemaDecoder struct {
	data []byte
}

func NewSchemaDecoder(data []byte) *SchemaDecoder {
	return &SchemaDecoder{data: data}
}

func (sc *SchemaDecoder) Decode(object runtime.Object) error {
	decoder := scheme.Codecs.UniversalDeserializer()
	gvk := &schema.GroupVersionKind{}
	_, _, err := decoder.Decode(sc.data, gvk, object)
	return err
}
