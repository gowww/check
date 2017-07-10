package check

import (
	"mime/multipart"
	"net/url"
)

var (
	testAlpha = "PeWGHxQBZfiMVQNbWGzkegklTRMHVmuO"
	testInt   = "-19001231"
	testFloat = "-19001231.558877"
	testURL   = "http://example.com"

	testChecker = &Checker{
		"email":     {Required, Email},
		"city":      {Alpha},
		"phone":     {Phone},
		"stars":     {Required, Range(3, 5)},
		"birthYear": {Max(2000)},
	}
	testCheckerData = &multipart.Form{Value: url.Values{
		"name":      {"foobar"},
		"phone":     {"0012345678901"},
		"stars":     {"2"},
		"birthYear": {"2001"},
	}}
	testCheckerWant = Errors{
		"email": {{Error: ErrRequired}},
		"stars": {{Error: ErrMin, Args: []interface{}{"3"}}},
	}
)

/*
func TestCheckerCheck(t *testing.T) {
	got := testChecker.Check(testCheckerData)
	if !reflect.DeepEqual(testCheckerWant, got) {
		t.Errorf("Checker.Check:\nwant %v\ngot  %v", testCheckerWant, got)
	}
}



func TestCheckerCheckRequest(t *testing.T) {
	got := testChecker.CheckRequest(&http.Request{MultipartForm: testCheckerData})
	if !reflect.DeepEqual(testCheckerWant, got) {
		t.Errorf("Checker.Check:\nwant %v\ngot  %v", testCheckerWant, got)
	}
}
*/
