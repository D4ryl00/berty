// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: errcode.proto

package errcode

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ErrCode int32

const (
	Undefined                                  ErrCode = 0
	TODO                                       ErrCode = 666
	ErrNotImplemented                          ErrCode = 777
	ErrInternal                                ErrCode = 888
	ErrInvalidInput                            ErrCode = 100
	ErrInvalidRange                            ErrCode = 101
	ErrMissingInput                            ErrCode = 102
	ErrSerialization                           ErrCode = 103
	ErrDeserialization                         ErrCode = 104
	ErrStreamRead                              ErrCode = 105
	ErrStreamWrite                             ErrCode = 106
	ErrStreamTransform                         ErrCode = 110
	ErrStreamSendAndClose                      ErrCode = 111
	ErrStreamHeaderWrite                       ErrCode = 112
	ErrStreamHeaderRead                        ErrCode = 115
	ErrStreamSink                              ErrCode = 113
	ErrStreamCloseAndRecv                      ErrCode = 114
	ErrMissingMapKey                           ErrCode = 107
	ErrDBWrite                                 ErrCode = 108
	ErrDBRead                                  ErrCode = 109
	ErrCryptoRandomGeneration                  ErrCode = 200
	ErrCryptoKeyGeneration                     ErrCode = 201
	ErrCryptoNonceGeneration                   ErrCode = 202
	ErrCryptoSignature                         ErrCode = 203
	ErrCryptoSignatureVerification             ErrCode = 204
	ErrCryptoDecrypt                           ErrCode = 205
	ErrCryptoDecryptPayload                    ErrCode = 206
	ErrCryptoEncrypt                           ErrCode = 207
	ErrCryptoKeyConversion                     ErrCode = 208
	ErrCryptoCipherInit                        ErrCode = 209
	ErrCryptoKeyDerivation                     ErrCode = 210
	ErrMap                                     ErrCode = 300
	ErrForEach                                 ErrCode = 301
	ErrKeystoreGet                             ErrCode = 400
	ErrKeystorePut                             ErrCode = 401
	ErrOrbitDBInit                             ErrCode = 1000
	ErrOrbitDBOpen                             ErrCode = 1001
	ErrOrbitDBAppend                           ErrCode = 1002
	ErrOrbitDBDeserialization                  ErrCode = 1003
	ErrOrbitDBStoreCast                        ErrCode = 1004
	ErrIPFSAdd                                 ErrCode = 1050
	ErrIPFSGet                                 ErrCode = 1051
	ErrHandshakeOwnEphemeralKeyGenSend         ErrCode = 1100
	ErrHandshakePeerEphemeralKeyRecv           ErrCode = 1101
	ErrHandshakeRequesterAuthenticateBoxKeyGen ErrCode = 1102
	ErrHandshakeResponderAcceptBoxKeyGen       ErrCode = 1103
	ErrHandshakeRequesterHello                 ErrCode = 1104
	ErrHandshakeResponderHello                 ErrCode = 1105
	ErrHandshakeRequesterAuthenticate          ErrCode = 1106
	ErrHandshakeResponderAccept                ErrCode = 1107
	ErrHandshakeRequesterAcknowledge           ErrCode = 1108
	ErrContactRequestSameAccount               ErrCode = 1200
	ErrContactRequestContactAlreadyAdded       ErrCode = 1201
	ErrContactRequestContactBlocked            ErrCode = 1202
	ErrContactRequestContactUndefined          ErrCode = 1203
	ErrContactRequestIncomingAlreadyReceived   ErrCode = 1204
	ErrGroupMemberLogEventOpen                 ErrCode = 1300
	ErrGroupMemberLogEventSignature            ErrCode = 1301
	ErrGroupMemberUnknownGroupID               ErrCode = 1302
	ErrGroupSecretOtherDestMember              ErrCode = 1303
	ErrGroupSecretAlreadySentToMember          ErrCode = 1304
	ErrGroupInvalidType                        ErrCode = 1305
	ErrGroupMissing                            ErrCode = 1306
	ErrGroupActivate                           ErrCode = 1307
	ErrGroupDeactivate                         ErrCode = 1308
	ErrGroupInfo                               ErrCode = 1309
	// Event errors
	ErrEventListMetadata                   ErrCode = 1400
	ErrEventListMessage                    ErrCode = 1401
	ErrMessageKeyPersistencePut            ErrCode = 1500
	ErrMessageKeyPersistenceGet            ErrCode = 1501
	ErrBridgeInterrupted                   ErrCode = 1600
	ErrBridgeNotRunning                    ErrCode = 1601
	ErrMessengerInvalidDeepLink            ErrCode = 2000
	ErrDBEntryAlreadyExists                ErrCode = 2100
	ErrDBAddConversation                   ErrCode = 2101
	ErrDBAddContactRequestOutgoingSent     ErrCode = 2102
	ErrDBAddContactRequestOutgoingEnqueud  ErrCode = 2103
	ErrDBAddContactRequestIncomingReceived ErrCode = 2104
	ErrDBAddContactRequestIncomingAccepted ErrCode = 2105
	ErrDBAddGroupMemberDeviceAdded         ErrCode = 2106
	ErrDBMultipleRecords                   ErrCode = 2107
	ErrReplayProcessGroupMetadata          ErrCode = 2200
	ErrReplayProcessGroupMessage           ErrCode = 2201
	// API internals errors
	ErrAttachmentPrepare                 ErrCode = 2300
	ErrAttachmentRetrieve                ErrCode = 2301
	ErrProtocolSend                      ErrCode = 2302
	ErrCLINoTermcaps                     ErrCode = 3001
	ErrServicesAuth                      ErrCode = 4000
	ErrServicesAuthNotInitialized        ErrCode = 4001
	ErrServicesAuthWrongState            ErrCode = 4002
	ErrServicesAuthInvalidResponse       ErrCode = 4003
	ErrServicesAuthServer                ErrCode = 4004
	ErrServicesAuthCodeChallenge         ErrCode = 4005
	ErrServicesAuthServiceInvalidToken   ErrCode = 4006
	ErrServicesAuthServiceNotSupported   ErrCode = 4007
	ErrServicesAuthUnknownToken          ErrCode = 4008
	ErrServicesAuthInvalidURL            ErrCode = 4009
	ErrServiceReplication                ErrCode = 4100
	ErrServiceReplicationServer          ErrCode = 4101
	ErrServiceReplicationMissingEndpoint ErrCode = 4102
	ErrBertyAccount                      ErrCode = 5000
	ErrBertyAccountNoIDSpecified         ErrCode = 5001
	ErrBertyAccountAlreadyOpened         ErrCode = 5002
	ErrBertyAccountInvalidIDFormat       ErrCode = 5003
	ErrBertyAccountLoggerDecorator       ErrCode = 5004
	ErrBertyAccountGRPCClient            ErrCode = 5005
	ErrBertyAccountOpenAccount           ErrCode = 5006
	ErrBertyAccountDataNotFound          ErrCode = 5007
	ErrBertyAccountMetadataUpdate        ErrCode = 5008
	ErrBertyAccountManagerOpen           ErrCode = 5009
	ErrBertyAccountManagerClose          ErrCode = 5010
	ErrBertyAccountInvalidCLIArgs        ErrCode = 5011
	ErrBertyAccountFSError               ErrCode = 5012
	ErrBertyAccountAlreadyExists         ErrCode = 5013
	ErrBertyAccountNoBackupSpecified     ErrCode = 5014
	ErrBertyAccountIDGenFailed           ErrCode = 5015
	ErrBertyAccountCreationFailed        ErrCode = 5016
)

var ErrCode_name = map[int32]string{
	0:    "Undefined",
	666:  "TODO",
	777:  "ErrNotImplemented",
	888:  "ErrInternal",
	100:  "ErrInvalidInput",
	101:  "ErrInvalidRange",
	102:  "ErrMissingInput",
	103:  "ErrSerialization",
	104:  "ErrDeserialization",
	105:  "ErrStreamRead",
	106:  "ErrStreamWrite",
	110:  "ErrStreamTransform",
	111:  "ErrStreamSendAndClose",
	112:  "ErrStreamHeaderWrite",
	115:  "ErrStreamHeaderRead",
	113:  "ErrStreamSink",
	114:  "ErrStreamCloseAndRecv",
	107:  "ErrMissingMapKey",
	108:  "ErrDBWrite",
	109:  "ErrDBRead",
	200:  "ErrCryptoRandomGeneration",
	201:  "ErrCryptoKeyGeneration",
	202:  "ErrCryptoNonceGeneration",
	203:  "ErrCryptoSignature",
	204:  "ErrCryptoSignatureVerification",
	205:  "ErrCryptoDecrypt",
	206:  "ErrCryptoDecryptPayload",
	207:  "ErrCryptoEncrypt",
	208:  "ErrCryptoKeyConversion",
	209:  "ErrCryptoCipherInit",
	210:  "ErrCryptoKeyDerivation",
	300:  "ErrMap",
	301:  "ErrForEach",
	400:  "ErrKeystoreGet",
	401:  "ErrKeystorePut",
	1000: "ErrOrbitDBInit",
	1001: "ErrOrbitDBOpen",
	1002: "ErrOrbitDBAppend",
	1003: "ErrOrbitDBDeserialization",
	1004: "ErrOrbitDBStoreCast",
	1050: "ErrIPFSAdd",
	1051: "ErrIPFSGet",
	1100: "ErrHandshakeOwnEphemeralKeyGenSend",
	1101: "ErrHandshakePeerEphemeralKeyRecv",
	1102: "ErrHandshakeRequesterAuthenticateBoxKeyGen",
	1103: "ErrHandshakeResponderAcceptBoxKeyGen",
	1104: "ErrHandshakeRequesterHello",
	1105: "ErrHandshakeResponderHello",
	1106: "ErrHandshakeRequesterAuthenticate",
	1107: "ErrHandshakeResponderAccept",
	1108: "ErrHandshakeRequesterAcknowledge",
	1200: "ErrContactRequestSameAccount",
	1201: "ErrContactRequestContactAlreadyAdded",
	1202: "ErrContactRequestContactBlocked",
	1203: "ErrContactRequestContactUndefined",
	1204: "ErrContactRequestIncomingAlreadyReceived",
	1300: "ErrGroupMemberLogEventOpen",
	1301: "ErrGroupMemberLogEventSignature",
	1302: "ErrGroupMemberUnknownGroupID",
	1303: "ErrGroupSecretOtherDestMember",
	1304: "ErrGroupSecretAlreadySentToMember",
	1305: "ErrGroupInvalidType",
	1306: "ErrGroupMissing",
	1307: "ErrGroupActivate",
	1308: "ErrGroupDeactivate",
	1309: "ErrGroupInfo",
	1400: "ErrEventListMetadata",
	1401: "ErrEventListMessage",
	1500: "ErrMessageKeyPersistencePut",
	1501: "ErrMessageKeyPersistenceGet",
	1600: "ErrBridgeInterrupted",
	1601: "ErrBridgeNotRunning",
	2000: "ErrMessengerInvalidDeepLink",
	2100: "ErrDBEntryAlreadyExists",
	2101: "ErrDBAddConversation",
	2102: "ErrDBAddContactRequestOutgoingSent",
	2103: "ErrDBAddContactRequestOutgoingEnqueud",
	2104: "ErrDBAddContactRequestIncomingReceived",
	2105: "ErrDBAddContactRequestIncomingAccepted",
	2106: "ErrDBAddGroupMemberDeviceAdded",
	2107: "ErrDBMultipleRecords",
	2200: "ErrReplayProcessGroupMetadata",
	2201: "ErrReplayProcessGroupMessage",
	2300: "ErrAttachmentPrepare",
	2301: "ErrAttachmentRetrieve",
	2302: "ErrProtocolSend",
	3001: "ErrCLINoTermcaps",
	4000: "ErrServicesAuth",
	4001: "ErrServicesAuthNotInitialized",
	4002: "ErrServicesAuthWrongState",
	4003: "ErrServicesAuthInvalidResponse",
	4004: "ErrServicesAuthServer",
	4005: "ErrServicesAuthCodeChallenge",
	4006: "ErrServicesAuthServiceInvalidToken",
	4007: "ErrServicesAuthServiceNotSupported",
	4008: "ErrServicesAuthUnknownToken",
	4009: "ErrServicesAuthInvalidURL",
	4100: "ErrServiceReplication",
	4101: "ErrServiceReplicationServer",
	4102: "ErrServiceReplicationMissingEndpoint",
	5000: "ErrBertyAccount",
	5001: "ErrBertyAccountNoIDSpecified",
	5002: "ErrBertyAccountAlreadyOpened",
	5003: "ErrBertyAccountInvalidIDFormat",
	5004: "ErrBertyAccountLoggerDecorator",
	5005: "ErrBertyAccountGRPCClient",
	5006: "ErrBertyAccountOpenAccount",
	5007: "ErrBertyAccountDataNotFound",
	5008: "ErrBertyAccountMetadataUpdate",
	5009: "ErrBertyAccountManagerOpen",
	5010: "ErrBertyAccountManagerClose",
	5011: "ErrBertyAccountInvalidCLIArgs",
	5012: "ErrBertyAccountFSError",
	5013: "ErrBertyAccountAlreadyExists",
	5014: "ErrBertyAccountNoBackupSpecified",
	5015: "ErrBertyAccountIDGenFailed",
	5016: "ErrBertyAccountCreationFailed",
}

var ErrCode_value = map[string]int32{
	"Undefined":                                  0,
	"TODO":                                       666,
	"ErrNotImplemented":                          777,
	"ErrInternal":                                888,
	"ErrInvalidInput":                            100,
	"ErrInvalidRange":                            101,
	"ErrMissingInput":                            102,
	"ErrSerialization":                           103,
	"ErrDeserialization":                         104,
	"ErrStreamRead":                              105,
	"ErrStreamWrite":                             106,
	"ErrStreamTransform":                         110,
	"ErrStreamSendAndClose":                      111,
	"ErrStreamHeaderWrite":                       112,
	"ErrStreamHeaderRead":                        115,
	"ErrStreamSink":                              113,
	"ErrStreamCloseAndRecv":                      114,
	"ErrMissingMapKey":                           107,
	"ErrDBWrite":                                 108,
	"ErrDBRead":                                  109,
	"ErrCryptoRandomGeneration":                  200,
	"ErrCryptoKeyGeneration":                     201,
	"ErrCryptoNonceGeneration":                   202,
	"ErrCryptoSignature":                         203,
	"ErrCryptoSignatureVerification":             204,
	"ErrCryptoDecrypt":                           205,
	"ErrCryptoDecryptPayload":                    206,
	"ErrCryptoEncrypt":                           207,
	"ErrCryptoKeyConversion":                     208,
	"ErrCryptoCipherInit":                        209,
	"ErrCryptoKeyDerivation":                     210,
	"ErrMap":                                     300,
	"ErrForEach":                                 301,
	"ErrKeystoreGet":                             400,
	"ErrKeystorePut":                             401,
	"ErrOrbitDBInit":                             1000,
	"ErrOrbitDBOpen":                             1001,
	"ErrOrbitDBAppend":                           1002,
	"ErrOrbitDBDeserialization":                  1003,
	"ErrOrbitDBStoreCast":                        1004,
	"ErrIPFSAdd":                                 1050,
	"ErrIPFSGet":                                 1051,
	"ErrHandshakeOwnEphemeralKeyGenSend":         1100,
	"ErrHandshakePeerEphemeralKeyRecv":           1101,
	"ErrHandshakeRequesterAuthenticateBoxKeyGen": 1102,
	"ErrHandshakeResponderAcceptBoxKeyGen":       1103,
	"ErrHandshakeRequesterHello":                 1104,
	"ErrHandshakeResponderHello":                 1105,
	"ErrHandshakeRequesterAuthenticate":          1106,
	"ErrHandshakeResponderAccept":                1107,
	"ErrHandshakeRequesterAcknowledge":           1108,
	"ErrContactRequestSameAccount":               1200,
	"ErrContactRequestContactAlreadyAdded":       1201,
	"ErrContactRequestContactBlocked":            1202,
	"ErrContactRequestContactUndefined":          1203,
	"ErrContactRequestIncomingAlreadyReceived":   1204,
	"ErrGroupMemberLogEventOpen":                 1300,
	"ErrGroupMemberLogEventSignature":            1301,
	"ErrGroupMemberUnknownGroupID":               1302,
	"ErrGroupSecretOtherDestMember":              1303,
	"ErrGroupSecretAlreadySentToMember":          1304,
	"ErrGroupInvalidType":                        1305,
	"ErrGroupMissing":                            1306,
	"ErrGroupActivate":                           1307,
	"ErrGroupDeactivate":                         1308,
	"ErrGroupInfo":                               1309,
	"ErrEventListMetadata":                       1400,
	"ErrEventListMessage":                        1401,
	"ErrMessageKeyPersistencePut":                1500,
	"ErrMessageKeyPersistenceGet":                1501,
	"ErrBridgeInterrupted":                       1600,
	"ErrBridgeNotRunning":                        1601,
	"ErrMessengerInvalidDeepLink":                2000,
	"ErrDBEntryAlreadyExists":                    2100,
	"ErrDBAddConversation":                       2101,
	"ErrDBAddContactRequestOutgoingSent":         2102,
	"ErrDBAddContactRequestOutgoingEnqueud":      2103,
	"ErrDBAddContactRequestIncomingReceived":     2104,
	"ErrDBAddContactRequestIncomingAccepted":     2105,
	"ErrDBAddGroupMemberDeviceAdded":             2106,
	"ErrDBMultipleRecords":                       2107,
	"ErrReplayProcessGroupMetadata":              2200,
	"ErrReplayProcessGroupMessage":               2201,
	"ErrAttachmentPrepare":                       2300,
	"ErrAttachmentRetrieve":                      2301,
	"ErrProtocolSend":                            2302,
	"ErrCLINoTermcaps":                           3001,
	"ErrServicesAuth":                            4000,
	"ErrServicesAuthNotInitialized":              4001,
	"ErrServicesAuthWrongState":                  4002,
	"ErrServicesAuthInvalidResponse":             4003,
	"ErrServicesAuthServer":                      4004,
	"ErrServicesAuthCodeChallenge":               4005,
	"ErrServicesAuthServiceInvalidToken":         4006,
	"ErrServicesAuthServiceNotSupported":         4007,
	"ErrServicesAuthUnknownToken":                4008,
	"ErrServicesAuthInvalidURL":                  4009,
	"ErrServiceReplication":                      4100,
	"ErrServiceReplicationServer":                4101,
	"ErrServiceReplicationMissingEndpoint":       4102,
	"ErrBertyAccount":                            5000,
	"ErrBertyAccountNoIDSpecified":               5001,
	"ErrBertyAccountAlreadyOpened":               5002,
	"ErrBertyAccountInvalidIDFormat":             5003,
	"ErrBertyAccountLoggerDecorator":             5004,
	"ErrBertyAccountGRPCClient":                  5005,
	"ErrBertyAccountOpenAccount":                 5006,
	"ErrBertyAccountDataNotFound":                5007,
	"ErrBertyAccountMetadataUpdate":              5008,
	"ErrBertyAccountManagerOpen":                 5009,
	"ErrBertyAccountManagerClose":                5010,
	"ErrBertyAccountInvalidCLIArgs":              5011,
	"ErrBertyAccountFSError":                     5012,
	"ErrBertyAccountAlreadyExists":               5013,
	"ErrBertyAccountNoBackupSpecified":           5014,
	"ErrBertyAccountIDGenFailed":                 5015,
	"ErrBertyAccountCreationFailed":              5016,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4240057316120df7, []int{0}
}

type ErrDetails struct {
	Codes                []ErrCode `protobuf:"varint,1,rep,packed,name=codes,proto3,enum=berty.errcode.ErrCode" json:"codes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ErrDetails) Reset()         { *m = ErrDetails{} }
func (m *ErrDetails) String() string { return proto.CompactTextString(m) }
func (*ErrDetails) ProtoMessage()    {}
func (*ErrDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_4240057316120df7, []int{0}
}

func (m *ErrDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrDetails.Unmarshal(m, b)
}

func (m *ErrDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrDetails.Marshal(b, m, deterministic)
}

func (m *ErrDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrDetails.Merge(m, src)
}

func (m *ErrDetails) XXX_Size() int {
	return xxx_messageInfo_ErrDetails.Size(m)
}

func (m *ErrDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrDetails.DiscardUnknown(m)
}

var xxx_messageInfo_ErrDetails proto.InternalMessageInfo

func (m *ErrDetails) GetCodes() []ErrCode {
	if m != nil {
		return m.Codes
	}
	return nil
}

func init() {
	proto.RegisterEnum("berty.errcode.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterType((*ErrDetails)(nil), "berty.errcode.ErrDetails")
}

func init() { proto.RegisterFile("errcode.proto", fileDescriptor_4240057316120df7) }

var fileDescriptor_4240057316120df7 = []byte{
	// 1666 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x57, 0xc9, 0x73, 0x5c, 0x47,
	0x19, 0xcf, 0x24, 0xc4, 0x1a, 0xb5, 0xb7, 0x4f, 0x6d, 0xc7, 0x6b, 0xe2, 0x51, 0x4c, 0x9c, 0x09,
	0x06, 0xac, 0x2a, 0xb8, 0x71, 0x9b, 0x4d, 0xf2, 0x94, 0xb5, 0x4c, 0xcd, 0xc8, 0xa4, 0x8a, 0x5b,
	0xfb, 0xf5, 0xa7, 0x37, 0x8d, 0xde, 0x74, 0x3f, 0xf7, 0xeb, 0x51, 0x32, 0x9c, 0x81, 0x62, 0x27,
	0x81, 0x2c, 0x8e, 0x03, 0x55, 0xec, 0x4b, 0x15, 0x54, 0xb1, 0x84, 0x25, 0xe1, 0x02, 0x37, 0x96,
	0x2c, 0x4e, 0xe0, 0x08, 0x07, 0x6e, 0x6c, 0x7f, 0x80, 0xa9, 0x02, 0x8a, 0xea, 0x7e, 0xdf, 0x68,
	0xde, 0x48, 0x83, 0xc8, 0x49, 0xa3, 0xef, 0xfb, 0xf5, 0xb7, 0x6f, 0x8f, 0x1d, 0x45, 0x6b, 0x23,
	0x23, 0xf1, 0x4a, 0x6a, 0x8d, 0x33, 0xfc, 0xe8, 0x0d, 0xb4, 0x6e, 0x74, 0x85, 0x88, 0xe7, 0x4e,
	0xc6, 0x26, 0x36, 0x81, 0xb3, 0xe4, 0x7f, 0xe5, 0xa0, 0x8b, 0x1f, 0x60, 0xac, 0x65, 0x6d, 0x13,
	0x9d, 0x50, 0x49, 0xc6, 0xdf, 0xc3, 0xee, 0xf7, 0xd8, 0xec, 0x4c, 0x69, 0xf1, 0xbe, 0xc7, 0x8e,
	0xbd, 0xef, 0xd4, 0x95, 0x29, 0x11, 0x57, 0x5a, 0xd6, 0x36, 0x8c, 0xc4, 0x6e, 0x0e, 0xba, 0xfc,
	0xca, 0x79, 0x36, 0x47, 0x24, 0x7e, 0x94, 0xcd, 0x5f, 0xd7, 0x12, 0xb7, 0x94, 0x46, 0x09, 0xf7,
	0xf0, 0x79, 0xf6, 0x8e, 0xcd, 0x8d, 0xe6, 0x06, 0xdc, 0xbe, 0x9f, 0x9f, 0x62, 0x0b, 0x2d, 0x6b,
	0xd7, 0x8d, 0x6b, 0x0f, 0xd2, 0x04, 0x07, 0xa8, 0x1d, 0x4a, 0xf8, 0xe4, 0x21, 0x0e, 0xec, 0x70,
	0xcb, 0xda, 0xb6, 0x76, 0x68, 0xb5, 0x48, 0xe0, 0xee, 0x21, 0x7e, 0x82, 0x1d, 0x0f, 0x94, 0x1d,
	0x91, 0x28, 0xd9, 0xd6, 0xe9, 0xd0, 0x81, 0x9c, 0x26, 0x76, 0x85, 0x8e, 0x11, 0x90, 0x88, 0x6b,
	0x2a, 0xcb, 0x94, 0x8e, 0x73, 0xe4, 0x16, 0x3f, 0xc9, 0xa0, 0x65, 0x6d, 0x0f, 0xad, 0x12, 0x89,
	0xfa, 0x88, 0x70, 0xca, 0x68, 0x88, 0xf9, 0x29, 0xc6, 0x83, 0x83, 0xd9, 0x14, 0xbd, 0xcf, 0x17,
	0xd8, 0x51, 0x8f, 0x76, 0x16, 0xc5, 0xa0, 0x8b, 0x42, 0x82, 0xe2, 0x9c, 0x1d, 0xdb, 0x25, 0x3d,
	0x6e, 0x95, 0x43, 0xf8, 0x30, 0x3d, 0xcf, 0x69, 0x9b, 0x56, 0xe8, 0x6c, 0xcb, 0xd8, 0x01, 0x68,
	0x7e, 0x96, 0x3d, 0xb0, 0x4b, 0xef, 0xa1, 0x96, 0x35, 0x2d, 0x1b, 0x89, 0xc9, 0x10, 0x0c, 0x3f,
	0xc3, 0x4e, 0xee, 0xb2, 0xae, 0xa2, 0x90, 0x68, 0x73, 0x61, 0x29, 0x3f, 0xcd, 0x4e, 0xec, 0xe1,
	0x04, 0xcd, 0xd9, 0x94, 0x31, 0x3d, 0xa5, 0xb7, 0xe1, 0xe6, 0x94, 0x82, 0x20, 0xb9, 0xa6, 0x65,
	0x17, 0xa3, 0x1d, 0xb0, 0xe4, 0x28, 0x79, 0xbf, 0x26, 0xd2, 0x6b, 0x38, 0x82, 0x6d, 0x7e, 0x2c,
	0xcf, 0x64, 0x3d, 0x57, 0x96, 0xf8, 0x8c, 0x84, 0xff, 0x83, 0x8a, 0x01, 0xbf, 0xc0, 0xce, 0xfa,
	0x5c, 0xd9, 0x51, 0xea, 0x4c, 0x57, 0x68, 0x69, 0x06, 0x2b, 0xa8, 0xd1, 0xe6, 0xe1, 0xf8, 0x75,
	0x89, 0x9f, 0x67, 0xa7, 0x76, 0xf9, 0xd7, 0x70, 0x54, 0x60, 0xfe, 0xa6, 0xc4, 0x1f, 0x62, 0x67,
	0x76, 0x99, 0xeb, 0x46, 0x47, 0x58, 0x60, 0xff, 0xb6, 0xc4, 0x4f, 0x87, 0x20, 0xe5, 0xec, 0x9e,
	0x8a, 0xb5, 0x70, 0x43, 0x8b, 0xf0, 0xbb, 0x12, 0x7f, 0x27, 0xbb, 0xb0, 0x9f, 0xf1, 0x41, 0xb4,
	0x6a, 0x4b, 0x45, 0xf9, 0xeb, 0x57, 0x4b, 0xfc, 0x81, 0xe0, 0x4e, 0x0e, 0x6a, 0x62, 0xe4, 0xff,
	0xc2, 0x6b, 0x25, 0xfe, 0x20, 0x3b, 0xbd, 0x97, 0xdc, 0x11, 0xa3, 0xc4, 0x08, 0x09, 0xaf, 0x4f,
	0x3f, 0x6a, 0xe9, 0xfc, 0xd1, 0x1b, 0xfb, 0xbc, 0x68, 0x18, 0xbd, 0x83, 0x36, 0xf3, 0x8a, 0xee,
	0x94, 0xf8, 0x99, 0x10, 0xfe, 0x9c, 0xd9, 0x50, 0x69, 0x1f, 0x6d, 0x5b, 0x2b, 0x07, 0x6f, 0xee,
	0x7b, 0xd6, 0x44, 0xab, 0x76, 0x72, 0xfb, 0xde, 0x2a, 0xf1, 0xc3, 0xec, 0x90, 0x0f, 0xb7, 0x48,
	0xe1, 0x7b, 0xf7, 0xf2, 0xe3, 0x21, 0xca, 0xcb, 0xc6, 0xb6, 0x44, 0xd4, 0x87, 0xef, 0xdf, 0xcb,
	0x4f, 0x84, 0xa2, 0xb9, 0x86, 0xa3, 0xcc, 0x19, 0x8b, 0x2b, 0xe8, 0xe0, 0xa9, 0xfb, 0xf6, 0x10,
	0x3b, 0x43, 0x07, 0x4f, 0x8f, 0x89, 0x1b, 0xf6, 0x86, 0x72, 0xcd, 0x7a, 0xd0, 0xfc, 0x97, 0xb9,
	0x69, 0xe2, 0x46, 0x8a, 0x1a, 0xfe, 0x3a, 0x47, 0xce, 0x11, 0xb1, 0x96, 0xa6, 0xa8, 0x25, 0xfc,
	0x6d, 0x8e, 0x52, 0x48, 0xe4, 0xbd, 0x15, 0xfd, 0xf7, 0x39, 0xf2, 0x8f, 0xf8, 0x3d, 0xaf, 0xb9,
	0x21, 0x32, 0x07, 0xff, 0x98, 0x23, 0xab, 0xdb, 0x9d, 0xe5, 0x5e, 0x4d, 0x4a, 0xb8, 0x5d, 0x2e,
	0x10, 0xbc, 0xc5, 0x2f, 0x96, 0x79, 0x95, 0x5d, 0x6c, 0x59, 0x7b, 0x55, 0x68, 0x99, 0xf5, 0xc5,
	0x36, 0x6e, 0x3c, 0xa1, 0x5b, 0x69, 0x1f, 0x07, 0x68, 0x45, 0x92, 0x57, 0x83, 0x2f, 0x72, 0x78,
	0xb5, 0xcc, 0x2f, 0xb1, 0xc5, 0x22, 0xb0, 0x83, 0x68, 0x8b, 0xc8, 0x50, 0xa2, 0xaf, 0x95, 0xf9,
	0x12, 0xbb, 0x5c, 0x84, 0x75, 0xf1, 0xe6, 0x10, 0x33, 0x87, 0xb6, 0x36, 0x74, 0x7d, 0xd4, 0xce,
	0xa7, 0x1f, 0xeb, 0xe6, 0xc9, 0x5c, 0x36, 0xbc, 0x5e, 0xe6, 0xef, 0x62, 0x8f, 0x4c, 0x3f, 0xc8,
	0x52, 0xa3, 0x25, 0xda, 0x5a, 0x14, 0x61, 0xea, 0x26, 0xd0, 0x37, 0xca, 0xbc, 0xc2, 0xce, 0xcd,
	0x94, 0x7d, 0x15, 0x93, 0xc4, 0xc0, 0x9d, 0x19, 0x00, 0x92, 0x95, 0x03, 0xde, 0x2c, 0xf3, 0x47,
	0xd9, 0xc3, 0xff, 0xd7, 0x3a, 0x78, 0xab, 0xcc, 0x17, 0xd9, 0xf9, 0x03, 0x8c, 0x82, 0xdf, 0xef,
	0x0b, 0xc7, 0x44, 0x52, 0xb4, 0xad, 0xcd, 0x13, 0x09, 0xca, 0x18, 0xe1, 0x0f, 0x65, 0xfe, 0x30,
	0x7b, 0x30, 0x4c, 0x4a, 0xed, 0x44, 0xe4, 0x08, 0xd4, 0x13, 0x03, 0xac, 0x45, 0x91, 0x19, 0x6a,
	0x07, 0x3f, 0x98, 0xa7, 0x00, 0x4c, 0x43, 0xe8, 0xbf, 0x5a, 0x62, 0x51, 0xc8, 0x51, 0x4d, 0x4a,
	0x94, 0xf0, 0xc3, 0x79, 0xfe, 0x08, 0xab, 0xfc, 0x2f, 0x68, 0x3d, 0x31, 0xd1, 0x36, 0x4a, 0xf8,
	0xd1, 0x3c, 0x39, 0x39, 0x13, 0x35, 0x19, 0xd5, 0x3f, 0x9e, 0xe7, 0xef, 0x65, 0x8f, 0xed, 0xc3,
	0xb5, 0x75, 0x64, 0x06, 0x4a, 0xc7, 0xa4, 0xb9, 0x8b, 0x11, 0xaa, 0x1d, 0x94, 0xf0, 0xd2, 0x3c,
	0x05, 0x77, 0xc5, 0x9a, 0x61, 0xba, 0x86, 0x83, 0x1b, 0x68, 0x57, 0x4d, 0xdc, 0xda, 0x41, 0xed,
	0x42, 0xf5, 0x3e, 0xc3, 0xc8, 0xba, 0x19, 0x80, 0xc9, 0x68, 0x78, 0x96, 0x51, 0x44, 0x0a, 0xa8,
	0xeb, 0xda, 0x47, 0x4c, 0x07, 0x4a, 0xbb, 0x09, 0xcf, 0x31, 0x7e, 0x91, 0x3d, 0x34, 0x86, 0xf4,
	0x30, 0xb2, 0xe8, 0x36, 0x5c, 0x1f, 0xfd, 0x28, 0x77, 0xf9, 0x0b, 0x78, 0x9e, 0x91, 0x93, 0x05,
	0x0c, 0x59, 0xdc, 0x43, 0xed, 0x36, 0x0d, 0xe1, 0x6e, 0x31, 0xea, 0x8d, 0x5c, 0x78, 0xbe, 0x4b,
	0x36, 0x47, 0x29, 0xc2, 0x0b, 0x8c, 0x9f, 0x0c, 0xbb, 0x24, 0x37, 0x24, 0x1f, 0xa9, 0x70, 0x9b,
	0x51, 0x0b, 0x06, 0x6a, 0x2d, 0x72, 0x7e, 0x1a, 0x20, 0xbc, 0xc8, 0x68, 0xd2, 0x05, 0x72, 0x13,
	0xc5, 0x98, 0xf1, 0x25, 0xc6, 0x17, 0xd8, 0x91, 0x89, 0xfc, 0x2d, 0x03, 0x5f, 0x66, 0xfc, 0x6c,
	0xd8, 0x03, 0xc1, 0xf3, 0x55, 0xe5, 0x6d, 0x76, 0x42, 0x0a, 0x27, 0xe0, 0xee, 0xd8, 0x9a, 0x02,
	0x2b, 0xcb, 0x44, 0x8c, 0xf0, 0x4f, 0x46, 0x15, 0x47, 0x84, 0x6b, 0x38, 0xea, 0xf8, 0xf1, 0x95,
	0x39, 0xd4, 0x51, 0x18, 0x23, 0x7f, 0x3c, 0x7c, 0x10, 0xc2, 0xf7, 0xf2, 0x9f, 0x0e, 0x93, 0xe2,
	0xba, 0x55, 0x32, 0xc6, 0xb0, 0x5f, 0xed, 0x30, 0xf5, 0x4b, 0xf7, 0x97, 0x47, 0x48, 0x71, 0xce,
	0x5a, 0x37, 0xae, 0x3b, 0xd4, 0xda, 0x3b, 0xfc, 0xab, 0x23, 0x05, 0xb1, 0xa8, 0x63, 0x1c, 0x2f,
	0xdc, 0x26, 0x62, 0xba, 0xea, 0x17, 0xd2, 0x9d, 0xe3, 0x34, 0x90, 0x9b, 0xf5, 0x96, 0x76, 0x76,
	0x44, 0x61, 0x6e, 0x3d, 0xa9, 0x32, 0x97, 0xc1, 0x4b, 0x40, 0x4a, 0x9b, 0xf5, 0x9a, 0x94, 0x34,
	0x75, 0xf3, 0xb9, 0xf4, 0x13, 0xa0, 0xd9, 0x32, 0x66, 0x15, 0xaa, 0x6c, 0x63, 0xe8, 0x62, 0xa3,
	0x74, 0xec, 0x93, 0x05, 0x3f, 0x05, 0x7e, 0x99, 0x5d, 0x3a, 0x18, 0xd8, 0xd2, 0x37, 0x87, 0x38,
	0x94, 0xf0, 0x33, 0xe0, 0xef, 0x66, 0x8f, 0xce, 0xc6, 0x8e, 0x4b, 0x77, 0xb7, 0x66, 0x7f, 0xfe,
	0x36, 0xc0, 0x79, 0x43, 0xa3, 0x84, 0x97, 0x81, 0x96, 0x56, 0x00, 0x17, 0xca, 0xb3, 0x89, 0x3b,
	0x2a, 0xc2, 0xbc, 0x05, 0x5f, 0x99, 0xb8, 0xbb, 0x36, 0x4c, 0x9c, 0x4a, 0x13, 0xec, 0x62, 0x64,
	0xac, 0xcc, 0xe0, 0x17, 0x40, 0x65, 0xdb, 0xc5, 0x34, 0x11, 0xa3, 0x8e, 0x35, 0x11, 0x66, 0x19,
	0xc9, 0xa1, 0x02, 0xb8, 0xb5, 0x40, 0xd5, 0x3f, 0x0b, 0x93, 0x57, 0xc2, 0x0b, 0x0b, 0xa4, 0xa1,
	0xe6, 0x9c, 0x88, 0xfa, 0xfe, 0x6a, 0xea, 0x58, 0x4c, 0x85, 0x45, 0xf8, 0xd7, 0x02, 0x3f, 0x17,
	0x6e, 0x83, 0x09, 0xab, 0x8b, 0xce, 0x2a, 0xdc, 0x41, 0xf8, 0xf7, 0x02, 0x95, 0x73, 0xc7, 0x1f,
	0x77, 0x91, 0x49, 0xc2, 0xd4, 0xfe, 0xcf, 0xc2, 0x78, 0x5d, 0xae, 0xb6, 0xd7, 0xcd, 0x26, 0xda,
	0x41, 0x24, 0xd2, 0x0c, 0x5e, 0x3e, 0x4d, 0xe0, 0x1e, 0x5a, 0xef, 0x5b, 0xe6, 0xa7, 0x1f, 0x7c,
	0xa5, 0x42, 0x0e, 0x14, 0xa9, 0xfe, 0x7a, 0xd3, 0xca, 0x85, 0x75, 0x83, 0x12, 0xbe, 0x5a, 0xa1,
	0x5d, 0x54, 0xc4, 0x3c, 0x6e, 0x8d, 0x8e, 0x7b, 0xce, 0xf7, 0xc3, 0xd7, 0x2a, 0x14, 0xc4, 0x22,
	0x7f, 0x7c, 0xc2, 0x85, 0x19, 0x9a, 0x21, 0x7c, 0xbd, 0x42, 0x7e, 0x14, 0x41, 0xfe, 0x37, 0x5a,
	0xf8, 0x46, 0x85, 0x22, 0x54, 0xe4, 0xf9, 0x3b, 0xb3, 0xd1, 0x17, 0x49, 0xe2, 0x0b, 0x14, 0xbe,
	0x59, 0xa1, 0xba, 0xda, 0xfb, 0x5c, 0x45, 0x38, 0xee, 0x70, 0xb3, 0x8d, 0x1a, 0xbe, 0x75, 0x00,
	0x70, 0xdd, 0xb8, 0xde, 0x30, 0x4d, 0x8d, 0xf5, 0xa9, 0xff, 0x76, 0x85, 0x9a, 0xa0, 0x08, 0xa4,
	0xa9, 0x94, 0x8b, 0xfa, 0xce, 0x2c, 0xbf, 0x49, 0xd9, 0xf5, 0xee, 0x2a, 0x7c, 0x77, 0x8f, 0x4b,
	0x3e, 0xbf, 0xe3, 0x43, 0xe7, 0xa3, 0x8b, 0xd3, 0xd2, 0x0b, 0x3c, 0x72, 0xfa, 0x63, 0x8b, 0xb4,
	0x03, 0xf6, 0x23, 0x68, 0x30, 0xb5, 0xb4, 0x4c, 0x8d, 0xd2, 0x0e, 0x3e, 0xbe, 0x48, 0xa9, 0xab,
	0xfb, 0xfb, 0x7c, 0xbc, 0x44, 0x3e, 0x51, 0xa5, 0xa8, 0x15, 0xa9, 0xeb, 0xa6, 0xdd, 0xec, 0xa5,
	0x18, 0xa9, 0x2d, 0xe5, 0xef, 0xee, 0x59, 0x10, 0xea, 0x65, 0x3f, 0xbf, 0x51, 0xc2, 0xa7, 0xaa,
	0x94, 0xbc, 0x22, 0x64, 0x7c, 0x94, 0x37, 0x97, 0x8d, 0x1d, 0x08, 0x07, 0x9f, 0x9e, 0x05, 0x5a,
	0x35, 0x71, 0xec, 0x1b, 0x25, 0x32, 0x56, 0x38, 0x63, 0xe1, 0x33, 0x55, 0x0a, 0x57, 0x11, 0xb4,
	0xd2, 0xed, 0x34, 0x1a, 0x89, 0xf2, 0x1d, 0xff, 0xd9, 0x2a, 0x2d, 0x93, 0x22, 0xdf, 0x5b, 0x31,
	0x76, 0xe8, 0x73, 0x55, 0x8a, 0x59, 0x11, 0xd0, 0x14, 0x4e, 0xac, 0x1b, 0xb7, 0x6c, 0x86, 0x5a,
	0xc2, 0xe7, 0xab, 0x54, 0xad, 0x45, 0xc4, 0xb8, 0xd1, 0xae, 0xa7, 0xd2, 0x57, 0xe3, 0x53, 0xb3,
	0xd4, 0xac, 0x09, 0x2d, 0x62, 0xb4, 0x61, 0x67, 0x3d, 0x3d, 0x4b, 0x0d, 0x01, 0xf2, 0xa3, 0xfe,
	0x0b, 0xb3, 0xd4, 0x50, 0x4c, 0x1a, 0xab, 0xed, 0x9a, 0x8d, 0x33, 0xf8, 0x62, 0x95, 0xce, 0xc8,
	0x22, 0x66, 0xb9, 0xd7, 0xb2, 0xd6, 0x58, 0x78, 0xe6, 0x80, 0xb8, 0xd3, 0x0c, 0x7d, 0xb6, 0x4a,
	0xc7, 0xc4, 0x74, 0xf6, 0xea, 0x22, 0xda, 0x1e, 0xa6, 0x93, 0x0c, 0x3e, 0x37, 0xcb, 0x9b, 0x76,
	0x73, 0x05, 0xf5, 0xb2, 0x50, 0x09, 0x4a, 0x78, 0x7e, 0x96, 0xad, 0x0d, 0x8b, 0xa1, 0x90, 0x08,
	0x73, 0xab, 0x5a, 0xbf, 0x74, 0xe7, 0xcf, 0x17, 0xee, 0xf9, 0x50, 0x25, 0xff, 0xc0, 0x73, 0x18,
	0xf5, 0x97, 0xc2, 0xcf, 0xa5, 0xd8, 0x2c, 0xa5, 0xdb, 0xf1, 0x12, 0x7d, 0xf2, 0xdd, 0x38, 0x14,
	0x3e, 0x13, 0xdf, 0xff, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1d, 0x94, 0xb3, 0xfd, 0x5c, 0x0e,
	0x00, 0x00,
}