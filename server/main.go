package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"fknsrs.biz/p/ottoext/fetch"
	"fknsrs.biz/p/ottoext/loop"
	"github.com/MathieuTurcotte/sourcemap"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/robertkrimen/otto"
)

func logDuration(s string) func() {
	before := time.Now()

	logrus.Debugf("%s [before]", s)

	return func() {
		logrus.WithField("duration", time.Now().Sub(before).String()).Debugf("%s [after]", s)
	}
}

func dumpDebugError(err error) string {
	if oerr, ok := err.(*otto.Error); ok {
		return oerr.String()
	}

	return ""
}

func main() {
	logrus.StandardLogger().Level = logrus.DebugLevel

	logrus.Debug("starting up")

	baseVM := otto.New()

	smRead := logDuration("read sourcemap")
	smData, err := ioutil.ReadFile("../client/public/bundle.server.js.map")
	if err != nil {
		panic(err)
	}
	smRead()

	smParse := logDuration("parse sourcemap")
	sm, err := sourcemap.Read(bytes.NewReader(smData))
	if err != nil {
		panic(err)
	}
	smParse()

	jsRead := logDuration("read javascript")
	jsData, err := ioutil.ReadFile("../client/public/bundle.server.js")
	if err != nil {
		panic(err)
	}
	jsRead()

	// lodash has a lookahead regex that we don't care about, and it makes otto
	// unhappy. this is a hack to get rid of it.
	jsData = bytes.Replace(jsData, []byte(`(?=`), []byte(`(`), -1)

	jsCompiled := logDuration("compile javascript")
	s, err := baseVM.CompileWithSourceMap("bundle.js", string(jsData), &sm)
	if err != nil {
		if st := dumpDebugError(err); st != "" {
			fmt.Println(st)
		}
		panic(err)
	}
	jsCompiled()

	jsInitialised := logDuration("initialise javascript")
	if _, err := baseVM.Run(s); err != nil {
		if st := dumpDebugError(err); st != "" {
			fmt.Println(st)
		}
		panic(err)
	}
	jsInitialised()

	tpl, err := template.ParseFiles("./index.template")
	if err != nil {
		panic(err)
	}

	apiRouter := mux.NewRouter()

	var cats = []Cat{
		Cat{Name: "fluffy", FavouriteFood: "marshmallows"},
		Cat{Name: "beans", FavouriteFood: "legume"},
		Cat{Name: "misty", FavouriteFood: "meat"},
		Cat{Name: "sara", FavouriteFood: "fish"},
	}

	apiRouter.Path("/api/v1/cats").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("search") == "" {
			w.Header().Set("content-type", "application/json")
			if err := json.NewEncoder(w).Encode(cats); err != nil {
				panic(err)
			}

			return
		}

		var d []Cat
		for _, c := range cats {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(r.URL.Query().Get("search"))) {
				d = append(d, c)
			}
		}

		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(d); err != nil {
			panic(err)
		}
	})

	apiRouter.Path("/api/v1/cats/{name}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		for _, c := range cats {
			if strings.ToLower(c.Name) == vars["name"] {
				w.Header().Set("content-type", "application/json")
				if err := json.NewEncoder(w).Encode(c); err != nil {
					panic(err)
				}

				return
			}
		}

		http.NotFound(w, r)
	})

	m := mux.NewRouter()

	m.PathPrefix("/api/v1").Handler(apiRouter)

	m.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vmCopied := logDuration("copy vm")
		vm := baseVM.Copy()
		vmCopied()

		l := loop.New(vm)
		if err := fetch.DefineWithHandler(vm, l, apiRouter); err != nil {
			panic(err)
		}

		fn, err := vm.Get("renderApplication")
		if err != nil {
			panic(err)
		}

		m := map[string]interface{}{
			"redirect": func(u string) {
				http.Redirect(w, r, u, 301)
			},
			"error": func(e string) {
				http.Error(w, e, http.StatusInternalServerError)
			},
			"notFound": func() {
				http.NotFound(w, r)
			},
			"success": func(d string) {
				w.Header().Set("content-type", "text/html")

				tpl.Execute(w, map[string]string{
					"html": d,
				})
			},
		}

		v, err := vm.ToValue(m)
		if err != nil {
			panic(err)
		}

		jsRun := logDuration("run javascript")
		if _, err := fn.Call(otto.UndefinedValue(), r.URL.String(), otto.NullValue(), v); err != nil {
			if st := dumpDebugError(err); st != "" {
				panic(fmt.Errorf("%s", st))
			}
			panic(err)
		}
		jsRun()

		jsFinished := logDuration("wait for js")
		if err := l.Run(); err != nil {
			if st := dumpDebugError(err); st != "" {
				panic(fmt.Errorf("%s", st))
			}
			panic(err)
		}
		jsFinished()
	})

	n := negroni.New()

	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	st := negroni.NewStatic(http.Dir("../client/public"))
	st.IndexFile = "NOT_USED_I_HOPE"
	n.Use(st)
	n.UseHandler(m)

	logrus.Info("listening")
	if err := http.ListenAndServe(":3000", n); err != nil {
		panic(err)
	}
}
