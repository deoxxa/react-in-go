package main

import (
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"fknsrs.biz/p/ottoext/fetch"
	"fknsrs.biz/p/ottoext/loop"
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

func main() {
	logrus.StandardLogger().Level = logrus.DebugLevel

	logrus.Debug("starting up")

	baseVM := otto.New()

	jsRead := logDuration("read javascript")
	d, err := ioutil.ReadFile("../client/public/bundle.server.js")
	if err != nil {
		panic(err)
	}
	jsRead()

	jsCompiled := logDuration("compile javascript")
	s, err := baseVM.Compile("bundle.js", string(d))
	if err != nil {
		panic(err)
	}
	jsCompiled()

	jsInitialised := logDuration("initialise javascript")
	if _, err := baseVM.Run(s); err != nil {
		panic(err)
	}
	jsInitialised()

	tpl, err := template.ParseFiles("./index.template")
	if err != nil {
		panic(err)
	}

	apiRouter := mux.NewRouter()

	apiRouter.Path("/api/v1/posts").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts := `[
			{"id":1,"title":"hello1 ","content":"hi there one!"},
			{"id":2,"title":"hello 2","content":"hi there two!"}
		]`

		w.Header().Set("content-type", "application/json")
		if _, err := w.Write([]byte(posts)); err != nil {
			panic(err)
		}
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
		if _, err := fn.Call(otto.UndefinedValue(), r.URL.String(), v); err != nil {
			panic(err)
		}
		jsRun()

		jsFinished := logDuration("wait for js")
		if err := l.Run(); err != nil {
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
