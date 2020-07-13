// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hashEpoch.proto

package fullHistory

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

//HashEpoch is used to store information about a hash to know which epoch it belongs to
type HashEpoch struct {
	Epoch uint32 `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
}

func (m *HashEpoch) Reset()      { *m = HashEpoch{} }
func (*HashEpoch) ProtoMessage() {}
func (*HashEpoch) Descriptor() ([]byte, []int) {
	return fileDescriptor_545a45a5af4a29bd, []int{0}
}
func (m *HashEpoch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HashEpoch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *HashEpoch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HashEpoch.Merge(m, src)
}
func (m *HashEpoch) XXX_Size() int {
	return m.Size()
}
func (m *HashEpoch) XXX_DiscardUnknown() {
	xxx_messageInfo_HashEpoch.DiscardUnknown(m)
}

var xxx_messageInfo_HashEpoch proto.InternalMessageInfo

func (m *HashEpoch) GetEpoch() uint32 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func init() {
	proto.RegisterType((*HashEpoch)(nil), "proto.HashEpoch")
}

func init() { proto.RegisterFile("hashEpoch.proto", fileDescriptor_545a45a5af4a29bd) }

var fileDescriptor_545a45a5af4a29bd = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x48, 0x2c, 0xce,
	0x70, 0x2d, 0xc8, 0x4f, 0xce, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52,
	0xba, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xe9, 0xf9, 0xe9, 0xf9,
	0xfa, 0x60, 0xe1, 0xa4, 0xd2, 0x34, 0x30, 0x0f, 0xcc, 0x01, 0xb3, 0x20, 0xba, 0x94, 0x14, 0xb9,
	0x38, 0x3d, 0x60, 0x06, 0x09, 0x89, 0x70, 0xb1, 0x82, 0x19, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xbc,
	0x41, 0x10, 0x8e, 0x93, 0xeb, 0x85, 0x87, 0x72, 0x0c, 0x37, 0x1e, 0xca, 0x31, 0x7c, 0x78, 0x28,
	0xc7, 0xd8, 0xf0, 0x48, 0x8e, 0x71, 0xc5, 0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c,
	0x92, 0x63, 0xbc, 0xf1, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x17, 0x8f, 0xe4, 0x18, 0x3e,
	0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18,
	0xa2, 0xb8, 0xd3, 0x4a, 0x73, 0x72, 0x3c, 0x32, 0x8b, 0x4b, 0xf2, 0x8b, 0x2a, 0x93, 0xd8, 0xc0,
	0x16, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xce, 0xde, 0x9d, 0x68, 0xb9, 0x00, 0x00, 0x00,
}

func (this *HashEpoch) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HashEpoch)
	if !ok {
		that2, ok := that.(HashEpoch)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Epoch != that1.Epoch {
		return false
	}
	return true
}
func (this *HashEpoch) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&fullHistory.HashEpoch{")
	s = append(s, "Epoch: "+fmt.Sprintf("%#v", this.Epoch)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringHashEpoch(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *HashEpoch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HashEpoch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HashEpoch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Epoch != 0 {
		i = encodeVarintHashEpoch(dAtA, i, uint64(m.Epoch))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintHashEpoch(dAtA []byte, offset int, v uint64) int {
	offset -= sovHashEpoch(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *HashEpoch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Epoch != 0 {
		n += 1 + sovHashEpoch(uint64(m.Epoch))
	}
	return n
}

func sovHashEpoch(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHashEpoch(x uint64) (n int) {
	return sovHashEpoch(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *HashEpoch) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HashEpoch{`,
		`Epoch:` + fmt.Sprintf("%v", this.Epoch) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringHashEpoch(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *HashEpoch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHashEpoch
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HashEpoch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HashEpoch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
			}
			m.Epoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHashEpoch
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Epoch |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHashEpoch(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHashEpoch
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHashEpoch
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipHashEpoch(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHashEpoch
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHashEpoch
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHashEpoch
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthHashEpoch
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHashEpoch
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHashEpoch
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHashEpoch        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHashEpoch          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHashEpoch = fmt.Errorf("proto: unexpected end of group")
)
