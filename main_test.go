

package main

import (
	"github.com/gookit/color"
	"github.com/willf/pad"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testPink	= color.FgMagenta.Render
	testBlue	= color.FgBlue.Render
	testRed		= color.FgRed.Render
)


func TestIndex(t *testing.T) {
	logPrefix			:= testPink(pad.Right("main():", 20, " "))
	defaultRoute		 = "/" + os.Getenv("NICKNAME")			; if defaultRoute	== "/"	{ defaultRoute	= "/go" }
	port				:= os.Getenv("REMOTE_PORT")			; if port			== ""	{ port			= "8080" }
	projectID			:= os.Getenv("GOOGLE_CLOUD_PROJECT")	; if projectID		== ""	{ t.Skip("sGOOGLE_CLOUD_PROJECT not set") }

	log.Printf("%sListening for route: %s on port: %s", logPrefix, testPink(defaultRoute), testPink(port))

	a, err				:= newApp(projectID); if err != nil	{ t.Fatalf("newApp: %v", err)}

	http.HandleFunc(defaultRoute, a.index);

	r					:= httptest.NewRequest("GET", "/" + os.Getenv("NICKNAME"), nil)
	rr					:= httptest.NewRecorder()
	a.index(rr, r)
	if got, want		:= rr.Body.String(), "1 view"; !strings.Contains(got, want) { t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want) }

	r				 	 = httptest.NewRequest("GET", "/" + os.Getenv("NICKNAME"), nil)				// Include the cookie from first visit
	r.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))
	rr					 = httptest.NewRecorder()
	a.index(rr, r)
	if got, want		:= rr.Body.String(), "2 views"; !strings.Contains(got, want) { t.Errorf("index second visit got:\n----\n%v\n----\nWant to contain %q", got, want) }
}


func TestIndexCorrupted(t *testing.T) {
	logPrefix	:= testPink(pad.Right("main():", 20, " "))

	log.Printf("%sLooking for index corruption:", logPrefix)

	projectID	:= os.Getenv("GOOGLE_CLOUD_PROJECT");	if projectID	== ""	{ t.Skip("NICKNAME not set") }
	a, err		:= newApp(projectID);						if err			!= nil	{ t.Fatalf("%snewApp: %v", logPrefix, err) }
	r			:= httptest.NewRequest("GET", "/" + os.Getenv("NICKNAME"), nil)
	r.Header.Set("Cookie", "this is not a valid session ID")
	rr			:= httptest.NewRecorder()

	a.index(rr, r)

	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) { t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)}
}
