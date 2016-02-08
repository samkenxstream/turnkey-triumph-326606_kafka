package proto

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

var _ = Suite(&MessagesSuite{})

func Test(t *testing.T) { TestingT(t) }

type MessagesSuite struct{}

type Request interface {
	Bytes() ([]byte, error)
	WriteTo(w io.Writer) (int64, error)
}

var _ Request = &MetadataReq{}
var _ Request = &ProduceReq{}
var _ Request = &FetchReq{}
var _ Request = &ConsumerMetadataReq{}
var _ Request = &OffsetReq{}
var _ Request = &OffsetCommitReq{}
var _ Request = &OffsetFetchReq{}

func testRequestSerialization(c *C, r Request) {
	var buf bytes.Buffer
	if n, err := r.WriteTo(&buf); err != nil {
		c.Fatalf("could not write request to buffer: %s", err)
	} else if n != int64(buf.Len()) {
		c.Fatalf("writer returned invalid number of bytes written %d != %d", n, buf.Len())
	}
	b, err := r.Bytes()
	if err != nil {
		c.Fatalf("could not convert request to bytes: %s", err)
	}
	if !bytes.Equal(b, buf.Bytes()) {
		c.Fatal("Bytes() and WriteTo() serialized request is of different form")
	}
}

func (s *MessagesSuite) TestMetadataRequest(c *C) {
	req1 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        nil,
	}
	testRequestSerialization(c, req1)
	b, _ := req1.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x0}

	if !bytes.Equal(b, expected) {
		c.Fatalf("expected different bytes representation: %v", b)
	}

	req2 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        []string{"foo", "bar"},
	}
	testRequestSerialization(c, req2)
	b, _ = req2.Bytes()
	expected = []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x2, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x3, 0x62, 0x61, 0x72}

	if !bytes.Equal(b, expected) {
		c.Fatalf("expected different bytes representation: %v", b)
	}

	r, _ := ReadMetadataReq(bytes.NewBuffer(expected))
	if !reflect.DeepEqual(r, req2) {
		c.Fatalf("malformed request: %#v", r)
	}
}

func (s *MessagesSuite) TestMetadataResponse(c *C) {
	msgb := []byte{0x0, 0x0, 0x1, 0xc7, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x10, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x12, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x11, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12}
	resp, err := ReadMetadataResp(bytes.NewBuffer(msgb))
	if err != nil {
		c.Fatalf("could not read metadata response: %s", err)
	}
	expected := &MetadataResp{
		CorrelationID: 123,
		Brokers: []MetadataRespBroker{
			{NodeID: 49168, Host: "172.17.42.1", Port: 49168},
			{NodeID: 49170, Host: "172.17.42.1", Port: 49170},
			{NodeID: 49169, Host: "172.17.42.1", Port: 49169},
			{NodeID: 49171, Host: "172.17.42.1", Port: 49171},
		},
		Topics: []MetadataRespTopic{
			{
				Name: "foo",
				Err:  error(nil),
				Partitions: []MetadataRespPartition{
					{Err: error(nil), ID: 2, Leader: 49171, Replicas: []int32{49171, 49168, 49169}, Isrs: []int32{49171, 49168, 49169}},
					{Err: error(nil), ID: 5, Leader: 49170, Replicas: []int32{49170, 49168, 49169}, Isrs: []int32{49170, 49168, 49169}},
					{Err: error(nil), ID: 4, Leader: 49169, Replicas: []int32{49169, 49171, 49168}, Isrs: []int32{49169, 49171, 49168}},
					{Err: error(nil), ID: 1, Leader: 49170, Replicas: []int32{49170, 49171, 49168}, Isrs: []int32{49170, 49171, 49168}},
					{Err: error(nil), ID: 3, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
					{Err: error(nil), ID: 0, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
				},
			},
			{
				Name: "test",
				Err:  error(nil),
				Partitions: []MetadataRespPartition{
					{Err: error(nil), ID: 1, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
					{Err: error(nil), ID: 0, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
				},
			},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		c.Fatalf("expected different message: %#v", resp)
	}

	if b, err := resp.Bytes(); err != nil {
		c.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb) {
			c.Fatalf("serialized representation different from expected: %#v", b)
		}
	}
}

func (s *MessagesSuite) TestProduceRequest(c *C) {
	req := &ProduceReq{
		CorrelationID: 241,
		ClientID:      "test",
		RequiredAcks:  RequiredAcksAll,
		Timeout:       time.Second,
		Topics: []ProduceReqTopic{
			{
				Name: "foo",
				Partitions: []ProduceReqPartition{
					{
						ID: 0,
						Messages: []*Message{
							{
								Offset: 0,
								Crc:    3099221847,
								Key:    []byte("foo"),
								Value:  []byte("bar"),
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		Compression Compression
		Expected    []byte
	}{
		{
			CompressionNone,
			[]byte{0x0, 0x0, 0x0, 0x49, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72},
		},
		{
			CompressionGzip,
			[]byte{0x0, 0x0, 0x0, 0x6d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x9d, 0x81, 0x74, 0xc4, 0x0, 0x1, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x2a, 0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0, 0xff, 0x62, 0x40, 0x0, 0x91, 0x1d, 0xbb, 0xe2, 0xc3, 0xc1, 0x2c, 0xe6, 0xb4, 0xfc, 0x7c, 0x10, 0x95, 0x94, 0x58, 0x4, 0x8, 0x0, 0x0, 0xff, 0xff, 0xa0, 0xbc, 0x10, 0xc2, 0x20, 0x0, 0x0, 0x0},
		},
		{
			CompressionSnappy,
			[]byte{0x0, 0x0, 0x0, 0x5c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x27, 0x2e, 0xd4, 0xed, 0xcd, 0x0, 0x2, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x19, 0x20, 0x0, 0x0, 0x19, 0x1, 0x10, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x5, 0xf, 0x28, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72},
		},
	}

	for _, tt := range tests {
		req.Compression = tt.Compression
		testRequestSerialization(c, req)
		b, _ := req.Bytes()

		if !bytes.Equal(b, tt.Expected) {
			c.Fatalf("expected different bytes representation: %#v", b)
		}

		r, _ := ReadProduceReq(bytes.NewBuffer(tt.Expected))
		req.Compression = CompressionNone // isn't set on deserialization
		if !reflect.DeepEqual(r, req) {
			c.Fatalf("malformed request: %#v", r)
		}
	}
}

func (s *MessagesSuite) TestProduceResponse(c *C) {
	msgb1 := []byte{0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x6, 0x66, 0x72, 0x75, 0x69, 0x74, 0x73, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5d, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	resp1, err := ReadProduceResp(bytes.NewBuffer(msgb1))
	if err != nil {
		c.Fatalf("could not read metadata response: %s", err)
	}
	expected1 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			{
				Name: "fruits",
				Partitions: []ProduceRespPartition{
					{
						ID:     93,
						Err:    ErrUnknownTopicOrPartition,
						Offset: -1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp1, expected1) {
		c.Fatalf("expected different message: %#v", resp1)
	}

	if b, err := resp1.Bytes(); err != nil {
		c.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb1) {
			c.Fatalf("serialized representation different from expected: %#v", b)
		}
	}

	msgb2 := []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	resp2, err := ReadProduceResp(bytes.NewBuffer(msgb2))
	if err != nil {
		c.Fatalf("could not read metadata response: %s", err)
	}
	expected2 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			{
				Name: "foo",
				Partitions: []ProduceRespPartition{
					{
						ID:     0,
						Err:    error(nil),
						Offset: 1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp2, expected2) {
		c.Fatalf("expected different message: %#v", resp2)
	}
	if b, err := resp2.Bytes(); err != nil {
		c.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb2) {
			c.Fatalf("serialized representation different from expected: %#v", b)
		}
	}
}

func (s *MessagesSuite) TestFetchRequest(c *C) {
	req := &FetchReq{
		CorrelationID: 241,
		ClientID:      "test",
		MaxWaitTime:   time.Second * 2,
		MinBytes:      12454,
		Topics: []FetchReqTopic{
			{
				Name: "foo",
				Partitions: []FetchReqPartition{
					{ID: 421, FetchOffset: 529, MaxBytes: 4921},
					{ID: 0, FetchOffset: 11, MaxBytes: 92},
				},
			},
		},
	}
	testRequestSerialization(c, req)
	b, _ := req.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x47, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x7, 0xd0, 0x0, 0x0, 0x30, 0xa6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x1, 0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x11, 0x0, 0x0, 0x13, 0x39, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x5c}

	if !bytes.Equal(b, expected) {
		c.Fatalf("expected different bytes representation: %#v", b)
	}

	r, _ := ReadFetchReq(bytes.NewBuffer(expected))
	if !reflect.DeepEqual(r, req) {
		c.Fatalf("malformed request: %#v", r)
	}
}

func (s *MessagesSuite) TestFetchResponse(c *C) {
	expected1 := &FetchResp{
		CorrelationID: 241,
		Topics: []FetchRespTopic{
			{
				Name: "foo",
				Partitions: []FetchRespPartition{
					{
						ID:        0,
						Err:       error(nil),
						TipOffset: 4,
						Messages: []*Message{
							{Offset: 2, Crc: 0xb8ba5f57, Key: []byte("foo"), Value: []byte("bar"), Topic: "foo", Partition: 0, TipOffset: 4},
							{Offset: 3, Crc: 0xb8ba5f57, Key: []byte("foo"), Value: []byte("bar"), Topic: "foo", Partition: 0, TipOffset: 4},
						},
					},
					{
						ID:        1,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}

	tests := []struct {
		Bytes     []byte
		RoundTrip bool // whether to compare re-serialized version
		Expected  *FetchResp
	}{
		{ // CompressionNone
			Bytes:     []byte{0x0, 0x0, 0x0, 0x75, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0},
			RoundTrip: true,
			Expected:  expected1,
		},
		{ // CompressionGzip
			Bytes:     []byte{0x0, 0x0, 0x0, 0x81, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x4c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x40, 0x7, 0x3c, 0x17, 0x35, 0x0, 0x1, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x32, 0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0, 0xff, 0x62, 0x80, 0x0, 0x26, 0x20, 0x16, 0xd9, 0xb1, 0x2b, 0x3e, 0x1c, 0xcc, 0x63, 0x4e, 0xcb, 0xcf, 0x7, 0x51, 0x49, 0x89, 0x45, 0x50, 0x79, 0x66, 0x5c, 0xf2, 0x80, 0x0, 0x0, 0x0, 0xff, 0xff, 0xab, 0xcc, 0x83, 0x80, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0},
			RoundTrip: false,
			Expected:  expected1,
		},
		{ // CompressionSnappy
			Bytes:     []byte{0x0, 0x0, 0x0, 0x75, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x34, 0x6, 0x8d, 0xfe, 0xe2, 0x0, 0x2, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x26, 0x40, 0x0, 0x0, 0x9, 0x1, 0x20, 0x2, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x5, 0xf, 0x28, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x5, 0x10, 0x8, 0x0, 0x0, 0x3, 0x5e, 0x20, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0},
			RoundTrip: false,
			Expected:  expected1,
		},
		{
			Bytes:     []byte{0x0, 0x0, 0x0, 0x48, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0},
			RoundTrip: true,
			Expected: &FetchResp{
				CorrelationID: 241,
				Topics: []FetchRespTopic{
					{
						Name: "test",
						Partitions: []FetchRespPartition{
							{
								ID:        0,
								Err:       ErrUnknownTopicOrPartition,
								TipOffset: -1,
								Messages:  []*Message{},
							},
							{
								ID:        1,
								Err:       ErrUnknownTopicOrPartition,
								TipOffset: -1,
								Messages:  []*Message{},
							},
							{
								ID:        8,
								Err:       ErrUnknownTopicOrPartition,
								TipOffset: -1,
								Messages:  []*Message{},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		resp, err := ReadFetchResp(bytes.NewBuffer(tt.Bytes))
		if err != nil {
			c.Fatalf("could not read fetch response: %s", err)
		}
		if !reflect.DeepEqual(resp, tt.Expected) {
			c.Fatalf("expected different message: %#v", resp)
		}
		if tt.RoundTrip {
			b, err := resp.Bytes()
			if err != nil {
				c.Fatalf("cannot serialize response: %s", err)
			}
			if !bytes.Equal(b, tt.Bytes) {
				c.Fatalf("serialized representation different from expected: %#v", b)
			}
		}
	}
}

func (s *MessagesSuite) TestSerializeEmptyMessageSet(c *C) {
	var buf bytes.Buffer
	messages := []*Message{}
	n, err := writeMessageSet(&buf, messages, CompressionNone)
	if err != nil {
		c.Fatalf("cannot serialize messages: %s", err)
	}
	if n != 0 {
		c.Fatalf("got n=%d result from writeMessageSet; want 0", n)
	}
	if l := len(buf.Bytes()); l != 0 {
		c.Fatalf("got len=%d for empty message set; should be 0", l)
	}
}

func (s *MessagesSuite) TestReadIncompleteMessage(c *C) {
	var buf bytes.Buffer
	_, err := writeMessageSet(&buf, []*Message{
		{Value: []byte("111111111111111")},
		{Value: []byte("222222222222222")},
		{Value: []byte("333333333333333")},
	}, CompressionNone)
	if err != nil {
		c.Fatalf("cannot serialize messages: %s", err)
	}

	b := buf.Bytes()
	// cut off the last bytes as kafka can do
	b = b[:len(b)-4]
	messages, err := readMessageSet(bytes.NewBuffer(b), int32(len(b)))
	if err != nil {
		c.Fatalf("cannot deserialize messages: %s", err)
	}
	if len(messages) != 2 {
		c.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Value[0] != '1' || messages[1].Value[0] != '2' {
		c.Fatal("expected different messages content")
	}
}

func BenchmarkProduceRequestMarshal(b *testing.B) {
	messages := make([]*Message, 100)
	for i := range messages {
		messages[i] = &Message{
			Offset: int64(i),
			Crc:    uint32(i),
			Key:    nil,
			Value:  []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.`),
		}

	}
	req := &ProduceReq{
		CorrelationID: 241,
		ClientID:      "test",
		Compression:   CompressionNone,
		RequiredAcks:  RequiredAcksAll,
		Timeout:       time.Second,
		Topics: []ProduceReqTopic{
			{
				Name: "foo",
				Partitions: []ProduceReqPartition{
					{
						ID:       0,
						Messages: messages,
					},
				},
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := req.Bytes(); err != nil {
			b.Fatalf("could not serialize messages: %s", err)
		}
	}
}

func BenchmarkProduceResponseUnmarshal(b *testing.B) {
	resp := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			{
				Name: "foo",
				Partitions: []ProduceRespPartition{
					{
						ID:     0,
						Err:    error(nil),
						Offset: 1,
					},
				},
			},
		},
	}
	raw, err := resp.Bytes()
	if err != nil {
		b.Fatalf("cannot serialize response: %s", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := ReadProduceResp(bytes.NewBuffer(raw)); err != nil {
			b.Fatalf("could not deserialize messages: %s", err)
		}
	}
}

func BenchmarkFetchRequestMarshal(b *testing.B) {
	req := &FetchReq{
		CorrelationID: 241,
		ClientID:      "test",
		MaxWaitTime:   time.Second * 2,
		MinBytes:      12454,
		Topics: []FetchReqTopic{
			{
				Name: "foo",
				Partitions: []FetchReqPartition{
					{ID: 421, FetchOffset: 529, MaxBytes: 4921},
					{ID: 0, FetchOffset: 11, MaxBytes: 92},
				},
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := req.Bytes(); err != nil {
			b.Fatalf("could not serialize messages: %s", err)
		}
	}
}

func BenchmarkFetchResponseUnmarshal(b *testing.B) {
	messages := make([]*Message, 100)
	for i := range messages {
		messages[i] = &Message{
			Offset: int64(i),
			Key:    nil,
			Value:  []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.`),
		}

	}
	resp := &FetchResp{
		CorrelationID: 241,
		Topics: []FetchRespTopic{
			{
				Name: "foo",
				Partitions: []FetchRespPartition{
					{
						ID:        0,
						TipOffset: 444,
						Messages:  messages,
					},
					{
						ID:        123,
						Err:       ErrBrokerNotAvailable,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}
	raw, err := resp.Bytes()
	if err != nil {
		b.Fatalf("cannot serialize response: %s", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := ReadFetchResp(bytes.NewBuffer(raw)); err != nil {
			b.Fatalf("could not deserialize messages: %s", err)
		}
	}
}

// vim has problem with coloring byte arrays in this file
// vim: set syntax=off:
