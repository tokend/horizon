package urlval

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDropMe(t *testing.T) {
	values, err := url.ParseQuery(`filter[a]=foobar&include=b&page[number]=2&page[size]=64&search=foo+bar&sort=a,-b,-c.d`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	type Request struct {
		AFilter    *string `filter:"a"`
		BInclude   bool    `include:"b"`
		CInclude   bool    `include:"c"`
		PageSize   uint64  `page:"size"`
		PageNumber uint64  `page:"number"`

		Sort   []Sort  `url:"sort"`
		Search *string `url:"search"`
	}

	var request Request

	if err = Decode(values, &request); err != nil {
		t.Fatalf("failed to decode request: %v", err)
	}

	require.NotNil(t, request.AFilter)
	require.Equal(t, "foobar", *request.AFilter)

	require.True(t, request.BInclude)
	require.False(t, request.CInclude)

	require.EqualValues(t, request.PageNumber, 2)
	require.EqualValues(t, request.PageSize, 64)

	require.NotNil(t, *request.Search)
	require.Equal(t, "foo bar", *request.Search)
	require.Equal(t, []Sort{"a", "-b", "-c.d"}, request.Sort)

	encoded, err := Encode(request)
	require.NoError(t, err)
	require.Equal(t, "filter%5Ba%5D=foobar&include=b&page%5Bnumber%5D=2&page%5Bsize%5D=64&search=foo+bar&sort=a%2C-b%2C-c.d", encoded)
}

func BenchmarkBasic(b *testing.B) {
	type Request struct {
		FQ *string `filter:"fq"`
		FW *string `filter:"fw"`
		FE *string `filter:"fe"`
		FR *string `filter:"fr"`
		FT *string `filter:"ft"`
		FY *uint64 `filter:"fy"`
		FU *uint64 `filter:"fu"`
		FI *uint64 `filter:"fi"`
		FO *uint64 `filter:"fo"`
		FP *uint64 `filter:"fp"`

		IA bool `include:"ia"`
		IS bool `include:"is"`
		ID bool `include:"id"`
		IF bool `include:"if"`
		IG bool `include:"ig"`
		IH bool `include:"ih"`
		IJ bool `include:"ij"`
		IK bool `include:"ik"`
		IL bool `include:"il"`

		PS uint64 `page:"size"`
		PN uint64 `page:"number"`
	}

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		values, _ := url.ParseQuery(`filter[fq]=foobar&filter[fw]=foobar&filter[fe]=foobar&filter[fr]=foobar&filter[ft]=foobar&filter[fy]=12&filter[fu]=1&include=ia,is,id,if,ig,ih,ij&page[number]=123&page[size]=444`)
		b.StartTimer()
		var r Request
		if err := Decode(values, &r); err != nil {
			b.Fatalf("failed to decode: %v", err)
		}
		_, err := Encode(r)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestSymmetry(t *testing.T) {
	type Includes struct {
		AInclude bool `include:"a"`
		BInclude bool `include:"b"`
	}
	type request struct {
		Includes
		AFilterString     *string  `filter:"a"`
		BFilterInt        *int     `filter:"b"`
		CFilterStrings    []string `filter:"c"`
		DFilterInts       []int    `filter:"d"`
		enexportedFilterF int      `filter:"f"`

		APage int    `page:"cursor"`
		BPage string `page:"size"`

		Currency *string `url:"currency"`
		Search   *string `url:"search"`

		DefaultParameter string `url:"parameter"`
	}
	t.Run("symmetric encoding and decoding with values", func(t *testing.T) {
		aFilter := "afilter"
		bFilter := 3
		currency := "EUR"
		search := "the world is gonna roll me"
		origin := request{
			AFilterString:     &aFilter,
			BFilterInt:        &bFilter,
			CFilterStrings:    []string{"somebody", "once", "told", "me"},
			DFilterInts:       []int{1, 2, 3},
			enexportedFilterF: 0,
			APage:             5,
			BPage:             "bla",
			Includes: Includes{
				AInclude: false,
				BInclude: true,
			},
			Currency:         &currency,
			Search:           &search,
			DefaultParameter: "default",
		}

		query, err := Encode(&origin)
		require.NoError(t, err)
		values, err := url.ParseQuery(query)
		require.NoError(t, err)
		result := request{}
		err = Decode(values, &result)
		require.NoError(t, err)
		require.Equal(t, origin, result)
	})
	t.Run("empty struct into empty query", func(t *testing.T) {
		query, err := Encode(request{})
		require.NoError(t, err)
		require.Equal(t, "page%5Bcursor%5D=0", query)
	})
	t.Run("empty struct from empty query", func(t *testing.T) {
		result := request{}
		err := Decode(url.Values{}, &result)
		require.NoError(t, err)
		require.Equal(t, result, request{})
	})
}
