//Package sessions contains middleware for easy session management in Negroni
//This is inspired by martini-contrib/sessions

package sessions

import (
	//"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

const (
	errorFormat = "[negroni-sessions] Error! %s \n"
)

type Store interface {
	sessions.Store
}

type Options struct {
	Path 		string
	Domain 		string
	MaxAge		int
	Secure		bool
	HttpOnly	bool
}

type Session interface {
	Get(key interface{}) interface{}
	Set(key interface{}, val interface{})
	Delete(key interface{})
	Clear()
	AddFlash(value interface{}, vars ...string)
	Flashes(vars ...string) []interface{}
	Options(Options)
}


type NegroniHandler struct {
	name	string
	store	Store
}

func (s *NegroniHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("TODO: Actual handling of session data")
	next(rw, r)
}

func NewSessions(name string, store Store) *NegroniHandler {
	return &NegroniHandler {
		name: name,
		store: store,
	}
}

type session struct {
	name		string
	request		*http.Request
	logger		*log.Logger
	store		Store
	session 	*sessions.Session
	written		bool
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options Options) {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		check(err, s.logger)
	}

	return s.session
}

func (s *session) Written() bool {
	return s.written
}

func check(err error, l *log.Logger) {
	if err != nil {
		l.Printf(errorFormat, err)
	}
}

