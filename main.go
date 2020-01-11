package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var store = sessions.NewCookieStore([]byte("secret"))

func main() {
	http.HandleFunc("/", RenderMainView)
	http.HandleFunc("/auth", RenderAuthView)
	http.HandleFunc("/auth/callback", Authenticate)
	http.HandleFunc("/logins", RenderloginView)

	log.Fatal(http.ListenAndServe(":8050", nil))
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(w, data)
}

// 메인 뷰 핸들러
func RenderMainView(w http.ResponseWriter, r *http.Request) {
}
var userNames = "최준용"
func RenderloginView(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "logins.html",userNames)
}
// 랜덤 state 값을 가진 구글 로그인 링크를 렌더링 해주는 뷰 핸들러
// 랜덤 state는 유저를 식별하는 용도로 사용된다
func RenderAuthView(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options = &sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	}
	state := RandToken()
	session.Values["state"] = state
	session.Save(r, w)
	//GetLoginURL 은 state에 만들어진 토큰값을 보내서 로그인 URL을 받아옴
	RenderTemplate(w, "auth.html", GetLoginURL(state))
}

// Google OAuth 인증 콜백 핸들러
func Authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	state := session.Values["state"]

	delete(session.Values, "state")
	session.Save(r, w)

	if state != r.FormValue("state") {
		http.Error(w, "Invalid session state", http.StatusUnauthorized)
		return
	}
	token, err := OAuthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := OAuthConf.Client(oauth2.NoContext, token)
	// UserInfoAPIEndpoint는 유저 정보 API URL을 담고 있음
	userInfoResp, err := client.Get(UserInfoAPIEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authUser User
	json.Unmarshal(userInfo, &authUser)

	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400,
	}
	session.Values["user"] = authUser.Email
	session.Values["username"] = authUser.Name
	session.Save(r, w)

	http.Redirect(w, r, "/logins", http.StatusFound)
}
