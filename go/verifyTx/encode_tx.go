package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"strconv"
	//"bytes"
	//math "math"
	math_bits "math/bits"
	// "github.com/cosmos/cosmos-sdk/codec"
	//types "github.com/cosmos/cosmos-sdk/codec/types"
	//github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	//tx "github.com/cosmos/cosmos-sdk/types/tx"
	signing "github.com/cosmos/cosmos-sdk/types/tx/signing"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	//"github.com/tendermint/tendermint/crypto/tmhash"
	//"github.com/google/go-cmp/cmp/internal/function"
	//cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	//"cosmossdk.io/api/cosmos/crypto/secp256k1"
	//"github.com/cosmos/cosmos-proto/anyutil"
	//txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	//"github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func StringToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func SovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}


func SovCoin(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func EncodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= SovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func EncodeVarintCoin(dAtA []byte, offset int, v uint64) int {
	offset -= SovCoin(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}



func (m *Transaction) Marshal() (dAtA []byte, err error) {
	size := m.Size()

	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Body) Marshal() (dAtA []byte, err error) {
	var s int
	_ = s
	s = m.Size()
		
	size := s + SovTx(uint64(s)) + 1

	dAtA = make([]byte, size)
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	fmt.Println(i)
			
	if &m != nil {
		{
			n, err := m.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return []byte{}, err
			}
			i -= n
			i = EncodeVarintTx(dAtA, i, uint64(n))
		}
		i--
		dAtA[i] = 0xa
	}
		return dAtA[i:], nil	
}

func (m *AuthInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size() + SovTx(uint64(m.Size())) + 1

	dAtA = make([]byte, size)
	i := len(dAtA)
	_ = i
	var l int
	_ = l

	dAtA = make([]byte, size)
	if &m != nil {
		{
			n, err := m.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return []byte{}, err
			}
			i -= n
			i = EncodeVarintTx(dAtA, i, uint64(n))
		}
		i--
		dAtA[i] = 0x12
	}
	return dAtA[i:], nil
}

func SignatureMarshal(m Transaction) (dAtA []byte, err error) {
	size := len(m.Signatures) * (len(First(base64.StdEncoding.DecodeString(m.Signatures[0]))) + 2)

	dAtA = make([]byte, size)
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	// fmt.Println(i)
	if len(m.Signatures) > 0 {
		for iNdEx := len(m.Signatures) - 1; iNdEx >= 0; iNdEx-- {
			signature := First(base64.StdEncoding.DecodeString(m.Signatures[iNdEx]))
			
			i -= len(signature)
			// fmt.Println(i)
			copy(dAtA[i:], signature)
			i = EncodeVarintTx(dAtA, i, uint64(len(signature)))
			// fmt.Println(i)
			i--
			// fmt.Println(i)
			dAtA[i] = 0x1a
		}
	}

	return dAtA[i:], nil
}

func (m *Transaction) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signatures) > 0 {
		for iNdEx := len(m.Signatures) - 1; iNdEx >= 0; iNdEx-- {
			signature := First(base64.StdEncoding.DecodeString(m.Signatures[iNdEx]))
			
			i -= len(signature)
			copy(dAtA[i:], signature)
			i = EncodeVarintTx(dAtA, i, uint64(len(signature)))
			i--
			dAtA[i] = 0x1a
		}
	}
	
	var test int
	_ = test
	test = i
	if &m.AuthInfo != nil {
		{
			size, err := m.AuthInfo.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	
	
	test = i
	if &m.Body != nil {
		{
			size, err := m.Body.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	
	test = i
	return len(dAtA) - i, nil
}

func (m *Transaction) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if &m.Body != nil {
		l = m.Body.Size()
		n += 1 + l + SovTx(uint64(l))
		
		
		
	}
	
	if &m.AuthInfo != nil {
		l = m.AuthInfo.Size()
		n += 1 + l + SovTx(uint64(l))
	}
	
	if len(m.Signatures) > 0 {
		for _, b := range m.Signatures {
			signature := First(base64.StdEncoding.DecodeString(b))
			l = len(signature)
			
			n += 1 + l + SovTx(uint64(l))
		}
	}
	
	return n
}

func (m *Body) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Messages) > 0 {
		for _, e := range m.Messages {
			l = e.Size()
			n += 1 + l + SovTx(uint64(l))
		}
		
	}
	l = len(m.Memo)
	if l > 0 {
		n += 1 + l + SovTx(uint64(l))
	}
	if First(StringToUint64(m.TimeoutHeight)) != 0 {
		n += 1 + SovTx(First(StringToUint64(m.TimeoutHeight)))
	}
	
	// if len(m.ExtensionOptions) > 0 {
	// 	for _, e := range m.ExtensionOptions {
	// 		l = e.Size()
	// 		n += 2 + l + SovTx(uint64(l))
	// 	}
	// }
	// if len(m.NonCriticalExtensionOptions) > 0 {
	// 	for _, e := range m.NonCriticalExtensionOptions {
	// 		l = e.Size()
	// 		n += 2 + l + SovTx(uint64(l))
	// 	}
	// }
	// n += 1 + SovTx(uint64(n))
	return n
}

func (m *Body) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	// if len(m.NonCriticalExtensionOptions) > 0 {
	// 	for iNdEx := len(m.NonCriticalExtensionOptions) - 1; iNdEx >= 0; iNdEx-- {
	// 		{
	// 			size, err := m.NonCriticalExtensionOptions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
	// 			if err != nil {
	// 				return 0, err
	// 			}
	// 			i -= size
	// 			i = EncodeVarintTx(dAtA, i, uint64(size))
	// 		}
	// 		i--
	// 		dAtA[i] = 0x7f
	// 		i--
	// 		dAtA[i] = 0xfa
	// 	}
	// }
	// if len(m.ExtensionOptions) > 0 {
	// 	for iNdEx := len(m.ExtensionOptions) - 1; iNdEx >= 0; iNdEx-- {
	// 		{
	// 			size, err := m.ExtensionOptions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
	// 			if err != nil {
	// 				return 0, err
	// 			}
	// 			i -= size
	// 			i = EncodeVarintTx(dAtA, i, uint64(size))
	// 		}
	// 		i--
	// 		dAtA[i] = 0x3f
	// 		i--
	// 		dAtA[i] = 0xfa
	// 	}
	// }
	if First(StringToUint64(m.TimeoutHeight)) != 0 {
		i = EncodeVarintTx(dAtA, i, uint64(First(StringToUint64(m.TimeoutHeight))))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Memo) > 0 {
		i -= len(m.Memo)
		copy(dAtA[i:], m.Memo)
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Memo)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Messages) > 0 {
		for iNdEx := len(m.Messages) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Messages[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = EncodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
			
		}
	}

	// i = EncodeVarintTx(dAtA, i, uint64(m.Size() - i))
	// i--
	// dAtA[i] = 0xa

	return len(dAtA) - i, nil
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l

	var l_value int
	_ = l_value
	if m.Type != "" {
		l = len(m.Type)
		n += 1 + l + SovTx(uint64(l))

		l_value = 1 + l + SovTx(uint64(l))
	}

	if m.Sender != "" {
		l = len(m.Sender)
		n += 1 + l + SovTx(uint64(l))
		
		l_value += (1 + l + SovTx(uint64(l)));
	}

	if m.Contract != "" {
		l = len(m.Contract)
		n += 1 + l + SovTx(uint64(l))
		
		l_value += (1 + l + SovTx(uint64(l)));
	}

	if m.Msg != "" {
		l = len(m.Msg)
		n += 1 + l + SovTx(uint64(l))

		l_value += (1 + l + SovTx(uint64(l)));
	}

	if len(m.Funds) > 0 {
		for _, e := range m.Funds {
			l = e.Size()
			n += 1 + l + SovTx(uint64(l))
			l_value += (1 + l + SovTx(uint64(l)));
		}
	}

	// if len(m.Funds) > 0 {
	// 	for _, b := range m.Funds {
	// 		if b != "" {
	// 			l = len(b)
	// 			n += 1 + l + SovTx(uint64(l))
	// 			l_value += (1 + l + SovTx(uint64(l)));
	// 		}
	// 	}
	// }

	n += 1 + SovTx(uint64(l_value))
	return n
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l

	var l_value int
	_ = l_value
	var idx = i;
	_ = idx
		idx = i
	if len(m.Funds) > 0 {
		for iNdEx := len(m.Funds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Funds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = EncodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
			
		}
	}

	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], []byte(m.Msg))
		
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x1a

		// l_value = 1 + len(m.Msg)
	}
	
	idx = i;

	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], []byte(m.Contract))
		fmt.Println("Contract", []byte(m.Contract))

		i = EncodeVarintTx(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x12

		// l_value += 1 +  len(m.Contract) 
	}
	
	idx = i;
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], []byte(m.Sender))
		fmt.Println("Sender", []byte(m.Sender))
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa

		// l_value += 1 + len(m.Sender) 
	}
	
	
	
	
	i = EncodeVarintTx(dAtA, i, uint64(len(dAtA) - i))
	i--
	dAtA[i] = 0x12 

	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], []byte(m.Type))
		
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
	}
	
	return len(dAtA) - i, nil
}

func (m *AuthInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if &m.Fee != nil {
		l = m.Fee.Size()
		n += 1 + l + SovTx(uint64(l))
		
	}
	
	if len(m.SignerInfos) > 0 {
		for _, e := range m.SignerInfos {
			l = e.Size()
			n += 1 + l + SovTx(uint64(l))
		}
		
	}
	// if m.Tip != nil {
	// 	l = m.Tip.Size()
	// 	n += 1 + l + SovTx(uint64(l))
	// }
	

	return n
}

func (m *AuthInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	// if m.Tip != nil {
	// 	{
	// 		size, err := m.Tip.MarshalToSizedBuffer(dAtA[:i])
	// 		if err != nil {
	// 			return 0, err
	// 		}
	// 		i -= size
	// 		i = EncodeVarintTx(dAtA, i, uint64(size))
	// 	}
	// 	i--
	// 	dAtA[i] = 0x1a
	// }
	if &m.Fee != nil {
		{
			size, err := m.Fee.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	var idx = i
	_ = idx

	
	
	if len(m.SignerInfos) > 0 {
		for iNdEx := len(m.SignerInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SignerInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = EncodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
			
			
		}
	}
	return len(dAtA) - i, nil
}

func (m *SignerInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if &m.PublicKey != nil {
		l = m.PublicKey.Size()
		n += 1 + l + SovTx(uint64(l))
		
	}
	if &m.ModeInfo != nil {
		l = m.ModeInfo.Size()
		n += 1 + l + SovTx(uint64(l))
		
	}
	if First(StringToUint64(m.Sequence)) != 0 {
		n += 1 + SovTx(First(StringToUint64(m.Sequence)))
		
	}
	return n
}

func (m *SignerInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if First(StringToUint64(m.Sequence)) != 0 {
		i = EncodeVarintTx(dAtA, i, First(StringToUint64(m.Sequence)))
		i--
		dAtA[i] = 0x18
	}
	
	var idx int
	_ = idx
	idx = i
	
	if &m.ModeInfo != nil {
		{
			size, err := m.ModeInfo.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	
	
	idx = i;
	if &m.PublicKey != nil {
		{
			size, err := m.PublicKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	
	
	return len(dAtA) - i, nil
}

func (m *PublicKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != "" {
		l = len(m.Type)
		n += 1 + l + SovTx(uint64(l))
	}
	if m.Key != "" {
		l = len(First(base64.StdEncoding.DecodeString(m.Key)))
		n += 1 + l + SovTx(uint64(l))
	}
	n += 1 + SovTx(uint64(1 + l + SovTx(uint64(l))))
	return n
}

func (m *PublicKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l

	if len(m.Key) > 0 {
		pubKey := First(base64.StdEncoding.DecodeString(m.Key))
		i -= len(pubKey)
		copy(dAtA[i:], []byte(pubKey))
		
		i = EncodeVarintTx(dAtA, i, uint64(len(pubKey)))
		i--
		dAtA[i] = 0xa

		i  = EncodeVarintTx(dAtA, i, uint64(len(dAtA) - i))
		i--
		dAtA[i] = 0x12
	}
	
	

	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], []byte(m.Type))
		
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}



func (m *ModeInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if &m.Single != nil {
		l = m.Single.Size()
		n += 1 + l + SovTx(uint64(l))
	}
	return n
}

func (m *ModeInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if &m.Single != nil {
		{
			
			
			size, err := m.Single.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}

			i -= size
			i = EncodeVarintTx(dAtA, i, uint64(size))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Single) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Mode != "" {
		l := GetMode(m.Mode)
		n += 1 + SovTx(uint64(l))
	}
	return n
}

func (m *Single) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	
	_ = i
	
	var l int
	_ = l
	if m.Mode != "" {
		{
			
			l := GetMode(m.Mode)
			//size := SovTx(uint64(l))
			i = EncodeVarintTx(dAtA, i, uint64(l))
			// i -= size
			i--
			dAtA[i] = 0x8
	
		}
	}
	return len(dAtA) - i, nil
}

func GetMode(mode string) signing.SignMode {
	if mode == "SIGN_MODE_DIRECT" {
		return signing.SignMode_SIGN_MODE_DIRECT
	} else if mode == "SIGN_MODE_DIRECT_AUX" {
		return signing.SignMode_SIGN_MODE_DIRECT_AUX
	} else if mode == "SIGN_MODE_LEGACY_AMINO_JSON" {
		return signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}
	return signing.SignMode_SIGN_MODE_EIP_191
}

func (m *Fee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + SovTx(uint64(l))
		}
	}
	
	if m.GasLimit != "" {
		n += 1 + SovTx(First(StringToUint64(m.GasLimit)))
	}

	l = len(m.Payer)
	if l > 0 {
		n += 1 + l + SovTx(uint64(l))
	}
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + SovTx(uint64(l))
	}
	return n
}

func (m *Fee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Payer) > 0 {
		i -= len(m.Payer)
		copy(dAtA[i:], m.Payer)
		i = EncodeVarintTx(dAtA, i, uint64(len(m.Payer)))
		i--
		dAtA[i] = 0x1a
	}
	if m.GasLimit != "" {
		i = EncodeVarintTx(dAtA, i, First(StringToUint64((m.GasLimit))))
		i--
		dAtA[i] = 0x10
	}
	
	
	var test int
	_ = test
	test = i
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = EncodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
			
			
		}
	}

	return len(dAtA) - i, nil
}

func (m *Amount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Denom != "" {	
		l = len(m.Denom)
		n += 1 + l + SovTx(uint64(l))
	}
	
	if m.Amount != "" {
		l = len(m.Amount)
		n += 1 + l + SovTx(uint64(l))
	}
	
	return n
}

func (m *Amount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := len(m.Amount)
		i -= size
		copy(dAtA[i:], m.Amount)
		i = EncodeVarintCoin(dAtA, i, uint64(size))
	}
	

	var idx = i;
	_ = idx;
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = EncodeVarintCoin(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	
	return len(dAtA) - i, nil
}



type Transaction struct {
	Body       Body     `json:"body"`
	AuthInfo   AuthInfo `json:"auth_info"`
	Signatures []string `json:"signatures"`
}

type Body struct {
	Messages                    []Message `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages"`
	Memo                        string    `protobuf:"bytes,2,opt,name=memo,proto3" json:"memo"`
	TimeoutHeight               string    `protobuf:"bytes,3,opt,name=timeout_height,proto3" json:"timeout_height"`
	ExtensionOptions            []string  `protobuf:"bytes,1023,rep,name=extensionOptions,proto3" json:"extension_options"`
	NonCriticalExtensionOptions []string  `protobuf:"bytes,2047,rep,name=non_critical_extension_options,proto3" json:"non_critical_extension_options"`
}

type Message struct {
	Type     string `protobuf:"bytes,1,opt,name=type,proto3" json:"@type"`
	Sender   string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender"`
	Contract string `protobuf:"bytes,3,opt,name=contract,proto3" json:"contract"`
	Msg      string    `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg"`
	Funds    []Amount `protobuf:"bytes,5,rep,name=funds,proto3" json:"funds"`
}

type AuthInfo struct {
	SignerInfos []SignerInfo `json:"signer_infos"`
	Fee         Fee          `json:"fee"`
}

type SignerInfo struct {
	PublicKey PublicKey `json:"public_key"`
	ModeInfo  ModeInfo  `json:"mode_info"`
	Sequence  string    `json:"sequence"`
}

type PublicKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ModeInfo struct {
	Single Single `json:"single"`
}

type Single struct {
	Mode string `json:"mode"`
}

type Fee struct {
	Amount   []Amount `json:"amount"`
	GasLimit string `json:"gas_limit"`
	Payer    string `json:"payer"`
	Granter  string `json:"granter"`
}

type Amount struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}

type Transaction_Json struct {
	Body       Body_Json     `json:"body"`
	AuthInfo   AuthInfo_Json `json:"auth_info"`
	Signatures []string `json:"signatures"`
}

type Body_Json struct {
	Messages                    []Message_Json `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages"`
	Memo                        string    `protobuf:"bytes,2,opt,name=memo,proto3" json:"memo"`
	TimeoutHeight               string    `protobuf:"bytes,3,opt,name=timeout_height,proto3" json:"timeout_height"`
	ExtensionOptions            []string  `protobuf:"bytes,1023,rep,name=extensionOptions,proto3" json:"extension_options"`
	NonCriticalExtensionOptions []string  `protobuf:"bytes,2047,rep,name=non_critical_extension_options,proto3" json:"non_critical_extension_options"`
}

type Message_Json struct {
	Type     string `protobuf:"bytes,1,opt,name=type,proto3" json:"@type"`
	Sender   string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender"`
	Contract string `protobuf:"bytes,3,opt,name=contract,proto3" json:"contract"`
	Msg      Msg_Json `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg"`
	Funds    []Amount_Json `protobuf:"bytes,5,rep,name=funds,proto3" json:"funds"`
}

type Msg_Json struct {
	Update_deposit_tree Update_deposit_tree_Json  `protobuf:"bytes,1,opt,name=update_deposit_tree,proto3" json:"update_deposit_tree"`
}

type Update_deposit_tree_Json struct {
	Root string  `protobuf:"bytes,2,opt,name=root,proto3" json:"root"`
	Proof []string  `protobuf:"int,1,opt,name=proof,proto3" json:"proof"`
}

type Preference_executor_fee_Json struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}

type AuthInfo_Json struct {
	SignerInfos []SignerInfo_Json `json:"signer_infos"`
	Fee         Fee_Json          `json:"fee"`
}

type SignerInfo_Json struct {
	PublicKey PublicKey_Json `json:"public_key"`
	ModeInfo  ModeInfo_Json  `json:"mode_info"`
	Sequence  string    `json:"sequence"`
}

type PublicKey_Json struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ModeInfo_Json struct {
	Single Single_Json `json:"single"`
}

type Single_Json struct {
	Mode string `json:"mode"`
}

type Fee_Json struct {
	Amount   []Amount_Json `json:"amount"`
	GasLimit string `json:"gas_limit"`
	Payer    string `json:"payer"`
	Granter  string `json:"granter"`
}

type Amount_Json struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}
func GetTx(transaction_path string) Transaction {
	var transaction Transaction_Json
	json.Unmarshal(GetBytesFileJson(transaction_path), &transaction)
	// transaction.Body.Messages[0].Msg = `{"request":{"threshold":2,"service":"orchai_price","preference_executor_fee":{"denom":"orai","amount":"0"}}}`
	
	var message []Message

	for i := 0; i < len(transaction.Body.Messages); i++ {
		//msg := string(First(json.Marshal(transaction.Body.Messages[i].Msg)))
		var fund []Amount
		for j := 0; j < len(transaction.Body.Messages[i].Funds); j++ {
			fund = append(fund, Amount{
				Denom: transaction.Body.Messages[i].Funds[j].Denom,
				Amount: transaction.Body.Messages[i].Funds[j].Amount,
			})
		}
		message = append(message, Message{
			Type: transaction.Body.Messages[i].Type,
			Sender: transaction.Body.Messages[i].Sender,
			Contract: transaction.Body.Messages[i].Contract,
			Msg: string(First(json.Marshal(transaction.Body.Messages[i].Msg))),
			Funds: fund,
		})
	}
	
	var signer_infos []SignerInfo
	for i := 0; i < len(transaction.AuthInfo.SignerInfos); i++ {
		signer_infos = append(signer_infos, SignerInfo{
			PublicKey: PublicKey{
				Type: transaction.AuthInfo.SignerInfos[i].PublicKey.Type,
				Key: transaction.AuthInfo.SignerInfos[i].PublicKey.Key,
			},
			ModeInfo: ModeInfo{
				Single: Single{
					Mode: transaction.AuthInfo.SignerInfos[i].ModeInfo.Single.Mode,
				},
			},
			Sequence: transaction.AuthInfo.SignerInfos[i].Sequence,
		})
	}

	var amount []Amount
	for i := 0; i < len(transaction.AuthInfo.Fee.Amount); i++ {
		amount = append(amount, Amount{
			Denom: transaction.AuthInfo.Fee.Amount[i].Denom,
			Amount: transaction.AuthInfo.Fee.Amount[i].Amount,
		})
	}

	var tx = Transaction{
		Body: Body{
			Messages: message,
			Memo: transaction.Body.Memo,
			TimeoutHeight: transaction.Body.TimeoutHeight,
			ExtensionOptions: transaction.Body.ExtensionOptions,
			NonCriticalExtensionOptions: transaction.Body.NonCriticalExtensionOptions,
		},
		AuthInfo: AuthInfo{
			SignerInfos: signer_infos,
			Fee: Fee{
				Amount: amount,
				GasLimit: transaction.AuthInfo.Fee.GasLimit,
				Payer: transaction.AuthInfo.Fee.Payer,
				Granter: transaction.AuthInfo.Fee.Granter,
			},
		},
		Signatures: transaction.Signatures,
	}

	
	return tx 
}

func GetTxByTxJson(transaction_path string) Transaction {
	var transaction Transaction_Json
	json.Unmarshal(GetBytesFileJson(transaction_path), &transaction)
	// transaction.Body.Messages[0].Msg = `{"request":{"threshold":2,"service":"orchai_price","preference_executor_fee":{"denom":"orai","amount":"0"}}}`
	
	var message []Message

	for i := 0; i < len(transaction.Body.Messages); i++ {
		//msg := string(First(json.Marshal(transaction.Body.Messages[i].Msg)))
		var fund []Amount
		for j := 0; j < len(transaction.Body.Messages[i].Funds); j++ {
			fund = append(fund, Amount{
				Denom: transaction.Body.Messages[i].Funds[j].Denom,
				Amount: transaction.Body.Messages[i].Funds[j].Amount,
			})
		}
		message = append(message, Message{
			Type: transaction.Body.Messages[i].Type,
			Sender: transaction.Body.Messages[i].Sender,
			Contract: transaction.Body.Messages[i].Contract,
			Msg: string(First(json.Marshal(transaction.Body.Messages[i].Msg))),
			Funds: fund,
		})
	}
	
	var signer_infos []SignerInfo
	for i := 0; i < len(transaction.AuthInfo.SignerInfos); i++ {
		signer_infos = append(signer_infos, SignerInfo{
			PublicKey: PublicKey{
				Type: transaction.AuthInfo.SignerInfos[i].PublicKey.Type,
				Key: transaction.AuthInfo.SignerInfos[i].PublicKey.Key,
			},
			ModeInfo: ModeInfo{
				Single: Single{
					Mode: transaction.AuthInfo.SignerInfos[i].ModeInfo.Single.Mode,
				},
			},
			Sequence: transaction.AuthInfo.SignerInfos[i].Sequence,
		})
	}

	var amount []Amount
	for i := 0; i < len(transaction.AuthInfo.Fee.Amount); i++ {
		amount = append(amount, Amount{
			Denom: transaction.AuthInfo.Fee.Amount[i].Denom,
			Amount: transaction.AuthInfo.Fee.Amount[i].Amount,
		})
	}

	var tx = Transaction{
		Body: Body{
			Messages: message,
			Memo: transaction.Body.Memo,
			TimeoutHeight: transaction.Body.TimeoutHeight,
			ExtensionOptions: transaction.Body.ExtensionOptions,
			NonCriticalExtensionOptions: transaction.Body.NonCriticalExtensionOptions,
		},
		AuthInfo: AuthInfo{
			SignerInfos: signer_infos,
			Fee: Fee{
				Amount: amount,
				GasLimit: transaction.AuthInfo.Fee.GasLimit,
				Payer: transaction.AuthInfo.Fee.Payer,
				Granter: transaction.AuthInfo.Fee.Granter,
			},
		},
		Signatures: transaction.Signatures,
	}

	
	return tx 
}


