

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	firestoregorilla "github.com/GoogleCloudPlatform/firestore-gorilla-sessions"
	"github.com/gookit/color"
	"github.com/gorilla/sessions"
	"github.com/willf/pad"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)


type app struct {
	store sessions.Store
	tmpl  *template.Template
}

var (
	debug 			= false
	Yellow			= color.FgYellow.Render
	Blue			= color.FgBlue.Render
	defaultRoute	= "/go"
	greetings		= []string {
		"Hello World",
		"Hallo Welt",
		"Ciao Mondo",
		"Salut le Monde",
		"Hola Mundo",
	}
)


func main() {
	logPrefix := Yellow(pad.Right("main():", 15, " "))
	log.Printf("%sInitialize", logPrefix)

	defaultRoute	 = "/"+os.Getenv("NICKNAME")+"/";	if	defaultRoute	== "/"	{ defaultRoute	= "/uh-oh"	}
	port			:= os.Getenv("REMOTE_PORT");	if	port			== ""	{ port			= "8080"	}
	debug, _ 		 = strconv.ParseBool(os.Getenv("DEBUG"))

	log.Printf("%sListening on port: %s, and for route: %s", logPrefix, port, defaultRoute)

	projectID	:= os.Getenv("GOOGLE_CLOUD_PROJECT");	if	projectID	== ""	{ log.Fatalf("%sGOOGLE_CLOUD_PROJECT must be set", logPrefix) }
	a, err		:= newApp(projectID);						if	err			!= nil	{ log.Fatalf("%s%v", logPrefix, err) }

	if debug { log.Printf("%sNew app returned: %+v", logPrefix, a) }

	http.HandleFunc(defaultRoute, a.index);
	if	err	:= http.ListenAndServe(":"+port, nil); err != nil { log.Fatal(err) }
}


func newApp(projectID string) (*app, error) {
	logPrefix		:= Yellow(pad.Right("newApp():", 15, " "))
	ctx				:= context.Background()
	client,	err		:= firestore.NewClient(ctx, projectID); if err != nil { log.Fatalf("firestore.NewClient:		%v", err) }
	store,	err		:= firestoregorilla.New(ctx, client);	if err != nil { log.Fatalf("firestoregorilla.New:	%v", err) }
	tmpl,	err		:= template.New("Index").Parse(`<body>{{.views}} {{if eq .views 1.0}}view{{else}}views{{end}} for "{{.greeting}}"</body>`);
	if err != nil { return nil, fmt.Errorf("template.New: %v", err) }

	if debug { log.Printf("%sNew app client instantiated: %+v", logPrefix, client)}

	return &app{store: store, tmpl:  tmpl}, nil
}


func (a *app) index(w http.ResponseWriter, r *http.Request) {					// track view and assign random greeting
	logPrefix		:= Yellow(pad.Right("index():", 15, " "))
	log.Printf("%stest default route: %s", logPrefix, Blue(defaultRoute))

	if r.RequestURI != defaultRoute { return }

	name			:= os.Getenv("NICKNAME") + "-views"											// Firestore collection
	session, err	:= a.store.Get(r, name); if err != nil { log.Printf("store.Get: %v", err) }

	if session.IsNew {
		log.Printf("%sNew session required: %+v", logPrefix, session)
		session.Values["views"]		= float64(0)
		session.Values["greeting"]	= greetings[rand.Intn(len(greetings))]
	}

	session.Values["views"] = session.Values["views"].(float64) + 1				// float64s is a firestoregorilla thing

	log.Printf("%sSaving session: %+v", logPrefix, session)

	if err := session.Save(r, w);					err != nil { log.Printf("Save:		%v", err		) }
	if err := a.tmpl.Execute(w, session.Values);	err != nil { log.Printf("Execute:	%v", err	) }
}
