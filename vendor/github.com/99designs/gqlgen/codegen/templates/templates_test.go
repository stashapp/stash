package templates

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/99designs/gqlgen/internal/code"

	"github.com/stretchr/testify/require"
)

func TestToGo(t *testing.T) {
	require.Equal(t, "ToCamel", ToGo("TO_CAMEL"))
	require.Equal(t, "ToCamel", ToGo("to_camel"))
	require.Equal(t, "ToCamel", ToGo("toCamel"))
	require.Equal(t, "ToCamel", ToGo("ToCamel"))
	require.Equal(t, "ToCamel", ToGo("to-camel"))
	require.Equal(t, "ToCamel", ToGo("-to-camel"))
	require.Equal(t, "ToCamel", ToGo("_to-camel"))
	require.Equal(t, "_", ToGo("_"))

	require.Equal(t, "RelatedURLs", ToGo("RelatedURLs"))
	require.Equal(t, "ImageIDs", ToGo("ImageIDs"))
	require.Equal(t, "FooID", ToGo("FooID"))
	require.Equal(t, "IDFoo", ToGo("IDFoo"))
	require.Equal(t, "FooASCII", ToGo("FooASCII"))
	require.Equal(t, "ASCIIFoo", ToGo("ASCIIFoo"))
	require.Equal(t, "FooUTF8", ToGo("FooUTF8"))
	require.Equal(t, "UTF8Foo", ToGo("UTF8Foo"))
	require.Equal(t, "JSONEncoding", ToGo("JSONEncoding"))

	require.Equal(t, "A", ToGo("A"))
	require.Equal(t, "ID", ToGo("ID"))
	require.Equal(t, "ID", ToGo("id"))
	require.Equal(t, "", ToGo(""))

	require.Equal(t, "RelatedUrls", ToGo("RelatedUrls"))
	require.Equal(t, "ITicket", ToGo("ITicket"))
	require.Equal(t, "FooTicket", ToGo("fooTicket"))
}

func TestToGoPrivate(t *testing.T) {
	require.Equal(t, "toCamel", ToGoPrivate("TO_CAMEL"))
	require.Equal(t, "toCamel", ToGoPrivate("to_camel"))
	require.Equal(t, "toCamel", ToGoPrivate("toCamel"))
	require.Equal(t, "toCamel", ToGoPrivate("ToCamel"))
	require.Equal(t, "toCamel", ToGoPrivate("to-camel"))

	require.Equal(t, "relatedURLs", ToGoPrivate("RelatedURLs"))
	require.Equal(t, "imageIDs", ToGoPrivate("ImageIDs"))
	require.Equal(t, "fooID", ToGoPrivate("FooID"))
	require.Equal(t, "idFoo", ToGoPrivate("IDFoo"))
	require.Equal(t, "fooASCII", ToGoPrivate("FooASCII"))
	require.Equal(t, "asciiFoo", ToGoPrivate("ASCIIFoo"))
	require.Equal(t, "fooUTF8", ToGoPrivate("FooUTF8"))
	require.Equal(t, "utf8Foo", ToGoPrivate("UTF8Foo"))
	require.Equal(t, "jsonEncoding", ToGoPrivate("JSONEncoding"))

	require.Equal(t, "relatedUrls", ToGoPrivate("RelatedUrls"))
	require.Equal(t, "iTicket", ToGoPrivate("ITicket"))

	require.Equal(t, "rangeArg", ToGoPrivate("Range"))

	require.Equal(t, "a", ToGoPrivate("A"))
	require.Equal(t, "id", ToGoPrivate("ID"))
	require.Equal(t, "id", ToGoPrivate("id"))
	require.Equal(t, "", ToGoPrivate(""))
	require.Equal(t, "_", ToGoPrivate("_"))
}

func Test_wordWalker(t *testing.T) {
	helper := func(str string) []*wordInfo {
		resultList := []*wordInfo{}
		wordWalker(str, func(info *wordInfo) {
			resultList = append(resultList, info)
		})
		return resultList
	}

	require.Equal(t, []*wordInfo{{Word: "TO"}, {Word: "CAMEL"}}, helper("TO_CAMEL"))
	require.Equal(t, []*wordInfo{{Word: "to"}, {Word: "camel"}}, helper("to_camel"))
	require.Equal(t, []*wordInfo{{Word: "to"}, {Word: "Camel"}}, helper("toCamel"))
	require.Equal(t, []*wordInfo{{Word: "To"}, {Word: "Camel"}}, helper("ToCamel"))
	require.Equal(t, []*wordInfo{{Word: "to"}, {Word: "camel"}}, helper("to-camel"))

	require.Equal(t, []*wordInfo{{Word: "Related"}, {Word: "URLs", HasCommonInitial: true}}, helper("RelatedURLs"))
	require.Equal(t, []*wordInfo{{Word: "Image"}, {Word: "IDs", HasCommonInitial: true}}, helper("ImageIDs"))
	require.Equal(t, []*wordInfo{{Word: "Foo"}, {Word: "ID", HasCommonInitial: true, MatchCommonInitial: true}}, helper("FooID"))
	require.Equal(t, []*wordInfo{{Word: "ID", HasCommonInitial: true, MatchCommonInitial: true}, {Word: "Foo"}}, helper("IDFoo"))
	require.Equal(t, []*wordInfo{{Word: "Foo"}, {Word: "ASCII", HasCommonInitial: true, MatchCommonInitial: true}}, helper("FooASCII"))
	require.Equal(t, []*wordInfo{{Word: "ASCII", HasCommonInitial: true, MatchCommonInitial: true}, {Word: "Foo"}}, helper("ASCIIFoo"))
	require.Equal(t, []*wordInfo{{Word: "Foo"}, {Word: "UTF8", HasCommonInitial: true, MatchCommonInitial: true}}, helper("FooUTF8"))
	require.Equal(t, []*wordInfo{{Word: "UTF8", HasCommonInitial: true, MatchCommonInitial: true}, {Word: "Foo"}}, helper("UTF8Foo"))

	require.Equal(t, []*wordInfo{{Word: "A"}}, helper("A"))
	require.Equal(t, []*wordInfo{{Word: "ID", HasCommonInitial: true, MatchCommonInitial: true}}, helper("ID"))
	require.Equal(t, []*wordInfo{{Word: "id", HasCommonInitial: true, MatchCommonInitial: true}}, helper("id"))
	require.Equal(t, []*wordInfo{}, helper(""))

	require.Equal(t, []*wordInfo{{Word: "Related"}, {Word: "Urls"}}, helper("RelatedUrls"))
	require.Equal(t, []*wordInfo{{Word: "ITicket"}}, helper("ITicket"))
}

func TestCenter(t *testing.T) {
	require.Equal(t, "fffff", center(3, "#", "fffff"))
	require.Equal(t, "##fffff###", center(10, "#", "fffff"))
	require.Equal(t, "###fffff###", center(11, "#", "fffff"))
}

func TestTemplateOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "gqlgen")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.RemoveAll(f.Name())
	err = Render(Options{Template: "hello", Filename: f.Name(), Packages: &code.Packages{}})
	if err != nil {
		t.Fatal(err)
	}
}
