//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2021 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package lsmkv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapEncoderDecoderJourney(t *testing.T) {
	// this test first encodes the map pairs, then decodes them and replace
	// duplicates, remove tombstones, etc.
	type test struct {
		name string
		in   []MapPair
		out  []MapPair
	}

	tests := []test{
		{
			name: "single pair",
			in: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar"),
				},
			},
			out: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar"),
				},
			},
		},
		{
			name: "single pair, updated value",
			in: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar"),
				},
				{
					Key:   []byte("foo"),
					Value: []byte("bar2"),
				},
			},
			out: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar2"),
				},
			},
		},
		{
			name: "single pair, tombstone added",
			in: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar"),
				},
				{
					Key:       []byte("foo"),
					Tombstone: true,
				},
			},
			out: []MapPair{},
		},
		{
			name: "single pair, tombstone added, same value added again",
			in: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar"),
				},
				{
					Key:       []byte("foo"),
					Tombstone: true,
				},
				{
					Key:   []byte("foo"),
					Value: []byte("bar2"),
				},
			},
			out: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("bar2"),
				},
			},
		},
		{
			name: "multiple values, combination of updates and tombstones",
			in: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("never-updated"),
				},
				{
					Key:   []byte("foo1"),
					Value: []byte("bar1"),
				},
				{
					Key:   []byte("foo2"),
					Value: []byte("bar2"),
				},
				{
					Key:   []byte("foo2"),
					Value: []byte("bar2.2"),
				},
				{
					Key:       []byte("foo1"),
					Tombstone: true,
				},
				{
					Key:   []byte("foo2"),
					Value: []byte("bar2.3"),
				},
				{
					Key:   []byte("foo1"),
					Value: []byte("bar1.2"),
				},
			},
			out: []MapPair{
				{
					Key:   []byte("foo"),
					Value: []byte("never-updated"),
				},
				{
					Key:   []byte("foo1"),
					Value: []byte("bar1.2"),
				},
				{
					Key:   []byte("foo2"),
					Value: []byte("bar2.3"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := make([]value, len(test.in))
			for i, kv := range test.in {
				enc, err := newMapEncoder().Do(kv)
				require.Nil(t, err)
				encoded[i] = enc[0]
			}
			res, err := newMapDecoder().Do(encoded)
			require.Nil(t, err)
			// NOTE: we are accpeting that the order can be lost on updates
			assert.ElementsMatch(t, test.out, res)
		})
	}
}
