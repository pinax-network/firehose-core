package json

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"slices"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/mr-tron/base58"
	"github.com/streamingfast/firehose-core/proto"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"
)

type kvList []*kv
type kv struct {
	key   string
	value any
}

type CustomEncoderFunc func(encoder *jsontext.Encoder, t []byte, options json.Options) error

type Marshaller struct {
	marshallers          *json.Marshalers
	includeUnknownFields bool
	registry             *proto.Registry
	bytesEncoderFunc     CustomEncoderFunc
}

type MarshallerOption func(*Marshaller)

func WithoutUnknownFields() MarshallerOption {
	return func(e *Marshaller) {
		e.includeUnknownFields = false
	}
}

func WithBytesEncoding(bytesEncoding string) MarshallerOption {
	return withBytesEncoderFunc(jsonEncoder(bytesEncoding))
}

// Deprecated: Use WithBytesEncoding instead passing the encoding directly.
var WithBytesEncoderFunc = withBytesEncoderFunc

func withBytesEncoderFunc(f CustomEncoderFunc) MarshallerOption {
	return func(e *Marshaller) {
		e.bytesEncoderFunc = f
	}
}

func NewMarshaller(registry *proto.Registry, options ...MarshallerOption) *Marshaller {
	m := &Marshaller{
		includeUnknownFields: true,
		registry:             registry,
		bytesEncoderFunc:     ToHex,
	}

	for _, opt := range options {
		opt(m)
	}

	m.marshallers = json.NewMarshalers(
		json.MarshalFuncV2(m.anypb),
		json.MarshalFuncV2(m.dynamicpbMessage),
		json.MarshalFuncV2(m.encodeKVList),
		json.MarshalFuncV2(m.bytesEncoderFunc),
	)

	return m
}

func (m *Marshaller) Marshal(in any) error {
	err := json.MarshalEncode(jsontext.NewEncoder(os.Stdout), in, json.WithMarshalers(m.marshallers))
	if err != nil {
		return fmt.Errorf("marshalling and encoding block to json: %w", err)
	}
	return nil
}

func (m *Marshaller) MarshalToString(in any, jsonEncoderOption ...json.Options) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := json.MarshalEncode(jsontext.NewEncoder(buf, jsonEncoderOption...), in, json.WithMarshalers(m.marshallers)); err != nil {
		return "", err
	}
	return buf.String(), nil

}

func (m *Marshaller) anypb(encoder *jsontext.Encoder, t *anypb.Any, options json.Options) error {
	msg, err := m.registry.Unmarshal(t)
	if err != nil {
		return fmt.Errorf("unmarshalling proto any: %w", err)
	}

	cnt, err := json.Marshal(msg, json.WithMarshalers(m.marshallers))
	if err != nil {
		return fmt.Errorf("json marshalling proto any: %w", err)
	}
	return encoder.WriteValue(cnt)
}

func (m *Marshaller) encodeKVList(encoder *jsontext.Encoder, t kvList, options json.Options) error {
	if err := encoder.WriteToken(jsontext.ObjectStart); err != nil {
		return err
	}
	for _, kv := range t {
		if err := encoder.WriteToken(jsontext.String(kv.key)); err != nil {
			return err
		}

		cnt, err := json.Marshal(kv.value, json.WithMarshalers(m.marshallers))
		if err != nil {
			return fmt.Errorf("json marshalling of value : %w", err)
		}

		if err := encoder.WriteValue(cnt); err != nil {
			return err
		}
	}
	return encoder.WriteToken(jsontext.ObjectEnd)
}

func (m *Marshaller) dynamicpbMessage(encoder *jsontext.Encoder, msg *dynamicpb.Message, options json.Options) error {
	var kvl kvList

	if m.includeUnknownFields {
		x := msg.GetUnknown()
		fieldNumber, ofType, l := protowire.ConsumeField(x)
		if l > 0 {
			unknownValue := x[:l]
			kvl = append(kvl, &kv{
				key:   fmt.Sprintf("__unknown_fields_%d_with_type_%d__", fieldNumber, ofType),
				value: hex.EncodeToString(unknownValue),
			})
		}
	}

	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.IsMap() {
			out := map[string]interface{}{}
			v.Map().Range(func(k protoreflect.MapKey, v2 protoreflect.Value) bool {
				key := k.String()
				out[key] = v2.Interface()
				return true
			})
			kvl = append(kvl, &kv{
				key:   string(fd.Name()),
				value: out,
			})
			return true
		}

		if fd.IsList() {
			out := make([]any, v.List().Len())
			for i := 0; i < v.List().Len(); i++ {
				out[i] = v.List().Get(i).Interface()
			}
			kvl = append(kvl, &kv{
				key:   string(fd.Name()),
				value: out,
			})
			return true
		}
		kvl = append(kvl, &kv{
			key:   string(fd.Name()),
			value: v.Interface(),
		})

		return true
	})

	slices.SortFunc(kvl, func(a, b *kv) int {
		if a.key < b.key {
			return -1
		}

		if a.key == b.key {
			return 0
		}

		return 1
	})

	cnt, err := json.Marshal(kvl, json.WithMarshalers(m.marshallers))
	if err != nil {
		return fmt.Errorf("json marshalling proto any: %w", err)
	}
	return encoder.WriteValue(cnt)
}

func ToBase58(encoder *jsontext.Encoder, t []byte, options json.Options) error {
	return encoder.WriteToken(jsontext.String(base58.Encode(t)))
}

func ToBase64(encoder *jsontext.Encoder, t []byte, options json.Options) error {
	return encoder.WriteToken(jsontext.String(base64.StdEncoding.EncodeToString(t)))
}

func ToHex(encoder *jsontext.Encoder, t []byte, options json.Options) error {
	return encoder.WriteToken(jsontext.String(hex.EncodeToString(t)))
}

func EncodeBase58(bytes []byte) string {
	return base58.Encode(bytes)
}

func EncodeBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func EncodeHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// Encode encodes the given bytes using the specified encoding.
func Encode(bytesEncoding string, bytes []byte) string {
	return Encoder(bytesEncoding)(bytes)
}

// Encoder returns the encoder function that will converts received bytes
// into a string of the specified encoding.
func Encoder(bytesEncoding string) func([]byte) string {
	switch {
	case strings.EqualFold(bytesEncoding, "base58"):
		return base58.Encode

	case strings.EqualFold(bytesEncoding, "base64"):
		return base64.StdEncoding.EncodeToString

	case strings.EqualFold(bytesEncoding, "hex"):
		return hex.EncodeToString
	}

	panic(fmt.Errorf("unsupported bytes encoding: %s", bytesEncoding))
}

func jsonEncoder(bytesEncoding string) CustomEncoderFunc {
	switch {
	case strings.EqualFold(bytesEncoding, "base58"):
		return ToBase58

	case strings.EqualFold(bytesEncoding, "base64"):
		return ToBase64

	case strings.EqualFold(bytesEncoding, "hex"):
		return ToHex
	}

	panic(fmt.Errorf("unsupported bytes encoding: %s", bytesEncoding))
}
