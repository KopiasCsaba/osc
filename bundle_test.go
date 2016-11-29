package osc

import (
	"bytes"
	"testing"
)

func TestBundleBytes(t *testing.T) {
	for _, testcase := range []struct {
		Input  Bundle
		Output []byte
	}{
		{
			Input: Bundle{Timetag: 10},
			Output: bytes.Join([][]byte{
				ToBytes(BundleTag),
				[]byte{0, 0, 0, 0, 0, 0, 0, 0x0A},
			}, []byte{}),
		},
		{
			Input: Bundle{
				Timetag: 10,
				Packets: []Packet{
					Message{
						Address: "/foo",
						Arguments: Arguments{
							Int(2),
							Float(5.2314),
						},
					},
				},
			},
			Output: bytes.Join([][]byte{
				ToBytes(BundleTag),
				[]byte{0, 0, 0, 0, 0, 0, 0, 0x0A},
				[]byte{0, 0, 0, 0x14},
				[]byte{'/', 'f', 'o', 'o', 0, 0, 0, 0},
				[]byte{',', TypetagInt, TypetagFloat, 0},
				[]byte{0, 0, 0, 2},
				[]byte{0x40, 0xA7, 0x67, 0xA1},
			}, []byte{}),
		},
		// Bundle within a bundle.
		{
			Input: Bundle{
				Timetag: 10,
				Packets: []Packet{
					Bundle{
						Timetag: 20,
						Packets: []Packet{
							Message{
								Address: "/foobar",
								Arguments: Arguments{
									Float(1),
								},
							},
						},
					},
					Message{
						Address: "/foo",
						Arguments: Arguments{
							Int(2),
							Float(5.2314),
						},
					},
				},
			},
			Output: bytes.Join([][]byte{
				ToBytes(BundleTag),
				[]byte{0, 0, 0, 0, 0, 0, 0, 0x0A}, // Timetag
				[]byte{0, 0, 0, 0x24},             // Length of first bundle element
				ToBytes(BundleTag),
				[]byte{0, 0, 0, 0, 0, 0, 0, 0x14}, // Timetag
				[]byte{0, 0, 0, 0x10},             // Length of first element of bundle within bundle.
				[]byte{'/', 'f', 'o', 'o', 'b', 'a', 'r', 0},
				[]byte{TypetagPrefix, TypetagFloat, 0, 0},
				[]byte{0x3F, 0x80, 0x00, 0x00},
				[]byte{0, 0, 0, 0x14}, // Length of second bundle element
				[]byte{'/', 'f', 'o', 'o', 0, 0, 0, 0},
				[]byte{TypetagPrefix, TypetagInt, TypetagFloat, 0},
				[]byte{0, 0, 0, 2},
				[]byte{0x40, 0xA7, 0x67, 0xA1},
			}, []byte{}),
		},
	} {
		if expected, got := testcase.Output, testcase.Input.Bytes(); !bytes.Equal(expected, got) {
			t.Fatalf("expected %q\n                got %q", expected, got)
		}
	}
}

func TestBundleEqual(t *testing.T) {
	for _, testcase := range []struct {
		b  Bundle
		e  []Bundle
		ne []Packet
	}{
		{
			b: Bundle{Timetag: 5},
			e: []Bundle{
				{Timetag: 5},
			},
			ne: []Packet{
				Message{},
				Bundle{Timetag: 2},
				Bundle{
					Timetag: 5,
					Packets: []Packet{
						Message{Address: "/foo"},
					},
				},
			},
		},
		{
			b: Bundle{
				Timetag: 5,
				Packets: []Packet{
					Message{Address: "/bar"},
				},
			},
			ne: []Packet{
				Bundle{
					Timetag: 5,
					Packets: []Packet{
						Message{Address: "/foo"},
					},
				},
			},
		},
	} {
		b := testcase.b
		for i, e := range testcase.e {
			if !b.Equal(e) {
				t.Fatalf("(testcase %d) expected %q to equal %q", i, b, e)
			}
		}
		for i, ne := range testcase.ne {
			if b.Equal(ne) {
				t.Fatalf("(testcase %d) expected %q to not equal %q", i, b, ne)
			}
		}
	}
}

func TestParseBundle(t *testing.T) {
	type Output struct {
		bundle Bundle
		err    error
	}
	for _, testcase := range []struct {
		Input    []byte
		Expected Output
	}{
		{
			Input: []byte{},
			Expected: Output{
				bundle: Bundle{},
			},
		},
	} {
		b, err := ParseBundle(testcase.Input, nil)
		if testcase.Expected.err == nil {
			if err != nil {
				t.Fatal(err)
			}
			if expected, got := testcase.Expected.bundle, b; !expected.Equal(got) {
				t.Fatal("expected %q, got %q", expected, got)
			}
		}
	}
}
