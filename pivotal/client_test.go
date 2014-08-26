package pivotal

import (
	"fmt"
	"testing"
)

func Test_buildURL(t *testing.T) {
	var expected string
	var p = make(map[string]string)

	expected = BASE_URL + "foo"
	if actual := buildURL("foo", p); actual != expected {
		fmt.Printf("no params - '%s' != '%s'\n", expected, actual)
		t.Fail()
	}

	p["bar"] = "banana"
	expected = BASE_URL + "foo?bar=banana"
	if actual := buildURL("foo", p); actual != expected {
		fmt.Printf("one param - '%s' != '%s'\n", expected, actual)
		t.Fail()
	}

	p["biz"] = "whiz"
	expected = BASE_URL + "foo?bar=banana&biz=whiz"
	if actual := buildURL("foo", p); actual != expected {
		fmt.Printf("one param - '%s' != '%s'\n", expected, actual)
		t.Fail()
	}
}
