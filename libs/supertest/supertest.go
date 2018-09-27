package supertest

import "github.com/parnurzeal/gorequest"
import "github.com/pkg4go/urlx"
import "encoding/json"
import "net/http"
import "reflect"
import "strings"
import "testing"
import "errors"
import "fmt"

type Agent struct {
	host             string
	path             string
	method           string
	t                *testing.T
	asserts          [][]interface{}
	agent            *gorequest.SuperAgent
	structFieldsOnly bool
	debug            bool
}

func Request(host string, ts ...*testing.T) *Agent {
	r := &Agent{}
	r.host = host
	r.agent = gorequest.New()

	if len(ts) > 0 {
		r.t = ts[0]
	}
	return r
}

func (r *Agent) Get(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Get(host)
	return r
}

func (r *Agent) Post(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Post(host)
	return r
}

func (r *Agent) Put(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Put(host)
	return r
}

func (r *Agent) Delete(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Delete(host)
	return r
}

func (r *Agent) Patch(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Patch(host)
	return r
}

func (r *Agent) Head(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Head(host)
	return r
}

func (r *Agent) Options(path string) *Agent {
	host, _ := urlx.Resolve(r.host, path)
	r.agent.Options(host)
	return r
}

func (r *Agent) Set(param, value string) *Agent {
	r.agent.Set(param, value)
	return r
}

func (r *Agent) SetBasicAuth(username, password string) *Agent {
	r.agent.SetBasicAuth(username, password)
	return r
}

func (r *Agent) AddCookie(cookie *http.Cookie) *Agent {
	r.agent.AddCookie(cookie)
	return r
}

func (r *Agent) AddCookies(cookies []*http.Cookie) *Agent {
	r.agent.AddCookies(cookies)
	return r
}

func (r *Agent) Type(ts string) *Agent {
	r.agent.Type(ts)
	return r
}

func (r *Agent) Query(q interface{}) *Agent {
	r.agent.Query(q)
	return r
}

func (r *Agent) Send(data interface{}) *Agent {
	r.agent.Send(data)
	return r
}

func (r *Agent) Expect(args ...interface{}) *Agent {
	if (len(args) == 1 && getType(args[0]) == "struct") ||
		(len(args) == 2 && getType(args[1]) == "struct") {
		r.structFieldsOnly = true
	}
	r.asserts = append(r.asserts, args)
	return r
}

// IncludeStructFieldsOnly obtained response body should
// contain only fields in expected response struct
func (r *Agent) IncludeStructFieldsOnly() *Agent {
	r.structFieldsOnly = true
	return r
}

// Debug print out response if debug mode
func (r *Agent) Debug() *Agent {
	r.debug = true
	return r
}

func (r *Agent) throw(err error) {
	if r.t != nil {
		r.t.Error(err)
	} else {
		panic(err)
	}
}

func (r *Agent) End(cbs ...func(response gorequest.Response, body []byte, errors []error)) {
	r.agent.EndBytes(func(res gorequest.Response, body []byte, errs []error) {
		if r.structFieldsOnly {
			var derived []byte
			for _, assert := range r.asserts {
				if len(assert) == 1 && getType(assert[0]) == "struct" {
					derived = r.deriveBodyForStruct(assert[0], body)
					break
				} else if len(assert) == 2 && getType(assert[0]) == "struct" {
					derived = r.deriveBodyForStruct(assert[1], body)
					break
				}
			}
			r.checkAll(res, body, derived, errs)
			return
		}
		r.checkAll(res, body, []byte{}, errs)
		if len(cbs) > 0 {
			cbs[0](res, body, errs)
		}
	})
}

func (r *Agent) EndStruct(
	v interface{},
	cbs ...func(response gorequest.Response, v interface{}, body []byte, errors []error),
) {
	r.agent.EndStruct(
		v,
		func(res gorequest.Response, v interface{}, body []byte, errs []error) {
			derived, err := json.Marshal(&v)
			if err != nil {
				r.throw(err)
			}
			r.checkAll(res, body, derived, errs)

			if len(cbs) > 0 {
				cbs[0](res, v, body, errs)
			}
		},
	)
}

func (r *Agent) checkAll(res gorequest.Response, body, derived []byte, errs []error) {
	contentType := res.Header.Get("Content-Type")
	status := res.StatusCode

	for _, assert := range r.asserts {
		if len(assert) == 1 {
			v := assert[0]

			if getType(v) == "int" {
				r.checkStatus(v, status)
			} else {
				r.checkBody(v, body, derived, contentType)
			}
		} else if len(assert) == 2 {

			if getType(assert[0]) == "int" {
				// Expect(200, `body`)
				r.checkStatus(assert[0], status)
				r.checkBody(assert[1], body, derived, contentType)
			} else if getType(assert[0]) == "string" {
				// Expect("Content-Type", "application/json")
				r.checkHeader(res.Header, assert[0], assert[1])
			} else {
				r.throw(errors.New("Unknown Expect behavior"))
			}
		} else {
			r.throw(errors.New("Expect only accept one or two args"))
		}
	}
}

func (r *Agent) checkStatus(status interface{}, actual int) {
	expect := status.(int)
	if expect != actual {
		r.throw(fmt.Errorf("Expected status: [%d], but got: [%d]", expect, actual))
	}
}

func (r *Agent) checkHeader(header http.Header, key, val interface{}) {
	k := key.(string)
	actual := header.Get(k)
	expect := val.(string)
	if actual != expect {
		r.throw(fmt.Errorf("Expected header [%s] to equal: [%s], but got: [%s]", k, expect, actual))
	}
}

func (r *Agent) checkBody(tobe interface{}, body, derived []byte, contentType string) {
	var expect string
	var bodyString string

	if r.debug {
		outString := ".......\nDEBUG\n%v\n%v\n......."
		debugbody := fmt.Sprintf(
			outString,
			"BODY",
			string(body[0:len(body)]),
		)
		fmt.Println(debugbody)
		if len(derived) != 0 {
			debugderived := fmt.Sprintf(
				outString,
				"DERIVED",
				string(derived[0:len(derived)]),
			)
			fmt.Println(debugderived)
		}
	}

	if strings.HasPrefix(contentType, "application/json") {
		if getType(tobe) == "string" {
			expect = tobe.(string)
		} else {
			buf, err := json.Marshal(tobe)
			if err != nil {
				r.throw(err)
			}
			expect = string(buf[0:len(buf)])
		}
		if len(derived) != 0 {
			bodyString = string(derived[0:len(derived)])
		} else {
			bodyString = string(body[0:len(body)])
		}
		if trim(expect) != trim(bodyString) {
			r.throw(fmt.Errorf(
				"Expected body:\n%s\nbut got:\n%s",
				trim(expect),
				trim(string(body[0:len(body)])),
			))
		}
	} else {
		r.throw(fmt.Errorf("content-type: %s not supported", contentType))
	}
}

func getType(v interface{}) string {
	return reflect.ValueOf(v).Kind().String()
}

func trim(str string) string {
	return strings.Replace(strings.Replace(strings.Replace(str, "\n", "", -1), "\t", "", -1), " ", "", -1)
}

func (r *Agent) deriveBodyForStruct(tobe interface{}, body []byte) []byte {
	expectType := reflect.ValueOf(tobe).Type()
	bodyStruct := reflect.New(expectType)
	bodyInterface := bodyStruct.Interface()
	err := json.Unmarshal(body, &bodyInterface)
	if err != nil {
		r.throw(err)
	}
	derivedBody, err := json.Marshal(bodyInterface)
	if err != nil {
		r.throw(err)
	}
	return derivedBody
}
