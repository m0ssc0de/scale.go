package scalecodec_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/m0ssc0de/scale.go"
	"github.com/m0ssc0de/scale.go/source"
	"github.com/m0ssc0de/scale.go/types"
	"github.com/m0ssc0de/scale.go/utiles"
)

func TestV14ExtrinsicDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(kusamaV14))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "network/kusama"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata, Spec: 9111}
	extrinsicRaw := "0x1d03840018c7717a3c5d2930f10eb5b0f67c191210e018bc55481bc44c1c1c7e810e945c01922c584c1c205b9747e76aad28cf2e73f729a9b3180072c47cd3abd205bb4b54f78a2627fa62a799f363fde25b5db74e5f8d8f3bde7a9a7f2ea8c95033d84e8585030800630301000400000000070088526a74080700000000070088526a74005ed0b200000000005ed0b20000000000000101000100000000010100707fd754e80e531ad411eb8b433548acbe669bf46a7e896e440feadc5ef3530800bca06501000000"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsicRaw)}, &option)
	e.Process()
}
