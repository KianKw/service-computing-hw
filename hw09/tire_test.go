package resthttp

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var fakeHandlerValue string

func fakeHandler(val string) Handle {
	return func(http.ResponseWriter, *http.Request, Params) {
		fakeHandlerValue = val
	}
}

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
	ps         Params
}

func getParams() *Params {
	ps := make(Params, 0, 20)
	return &ps
}

func checkRequests(t *testing.T, tire *node, requests testRequests) {
	for _, request := range requests {
		handler, psp, _ := tire.getValue(request.path, getParams)

		switch {
		case handler == nil:
			if !request.nilHandler {
				t.Errorf("handle mismatch for route '%s': Expected non-nil handle", request.path)
			}
		case request.nilHandler:
			t.Errorf("handle mismatch for route '%s': Expected nil handle", request.path)
		default:
			handler(nil, nil, nil)
			if fakeHandlerValue != request.route {
				t.Errorf("handle mismatch for route '%s': Wrong handle (%s != %s)", request.path, fakeHandlerValue, request.route)
			}
		}

		var ps Params
		if psp != nil {
			ps = *psp
		}

		if !reflect.DeepEqual(ps, request.ps) {
			t.Errorf("Params mismatch for route '%s'", request.path)
		}
	}
}

func checkPriorities(t *testing.T, n *node) uint32 {
	var prio uint32
	for i := range n.children {
		prio += checkPriorities(t, n.children[i])
	}

	if n.handle != nil {
		prio++
	}

	if n.priority != prio {
		t.Errorf(
			"priority mismatch for node '%s': is %d, should be %d",
			n.path, n.priority, prio,
		)
	}

	return prio
}

func TestCountParams(t *testing.T) {
	if ParamNum("/path/:param1/static/*catch-all") != 2 {
		t.Fail()
	}
	if ParamNum(strings.Repeat("/:param", 256)) != 256 {
		t.Fail()
	}
}

func TestTireAddAndGet(t *testing.T) {
	tire := &node{}

	routes := [...]string{
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
	}
	for _, route := range routes {
		tire.addRoute(route, fakeHandler(route))
	}

	checkRequests(t, tire, testRequests{
		{"/a", false, "/a", nil},
		{"/", true, "", nil},
		{"/hi", false, "/hi", nil},
		{"/contact", false, "/contact", nil},
		{"/co", false, "/co", nil},
		{"/con", true, "", nil},
		{"/cona", true, "", nil},
		{"/no", true, "", nil},
		{"/ab", false, "/ab", nil},
		{"/α", false, "/α", nil},
		{"/β", false, "/β", nil},
	})

	checkPriorities(t, tire)
}

func catchPanic(testFunc func()) (recv interface{}) {
	defer func() {
		recv = recover()
	}()

	testFunc()
	return
}

type testRoute struct {
	path     string
	conflict bool
}

func testRoutes(t *testing.T, routes []testRoute) {
	tire := &node{}

	for i := range routes {
		route := routes[i]
		recv := catchPanic(func() {
			tire.addRoute(route.path, nil)
		})

		if route.conflict {
			if recv == nil {
				t.Errorf("no panic for conflicting route '%s'", route.path)
			}
		} else if recv != nil {
			t.Errorf("unexpected panic for route '%s': %v", route.path, recv)
		}
	}
}

func TestTireChildConflict(t *testing.T) {
	routes := []testRoute{
		{"/cmd/vet", false},
		{"/cmd/:tool/:sub", true},
		{"/src/AUTHORS", false},
		{"/src/*filepath", true},
		{"/user_x", false},
		{"/user_:name", true},
		{"/id/:id", false},
		{"/id:id", true},
		{"/:id", true},
		{"/*filepath", true},
	}
	testRoutes(t, routes)
}

func TestTireDupliatePath(t *testing.T) {
	tire := &node{}

	routes := [...]string{
		"/",
		"/doc/",
		"/src/*filepath",
		"/search/:query",
		"/user_:name",
	}
	for i := range routes {
		route := routes[i]
		recv := catchPanic(func() {
			tire.addRoute(route, fakeHandler(route))
		})
		if recv != nil {
			t.Fatalf("panic inserting route '%s': %v", route, recv)
		}
		recv = catchPanic(func() {
			tire.addRoute(route, nil)
		})
		if recv == nil {
			t.Fatalf("no panic while inserting duplicate route '%s", route)
		}
	}

	checkRequests(t, tire, testRequests{
		{"/", false, "/", nil},
		{"/doc/", false, "/doc/", nil},
		{"/src/some/file.png", false, "/src/*filepath", Params{Param{"filepath", "/some/file.png"}}},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query", Params{Param{"query", "someth!ng+in+ünìcodé"}}},
		{"/user_gopher", false, "/user_:name", Params{Param{"name", "gopher"}}},
	})
}

func TestTireCatchAllConflict(t *testing.T) {
	routes := []testRoute{
		{"/src/*filepath/x", true},
		{"/src2/", false},
		{"/src2/*filepath/x", true},
		{"/src3/*filepath", false},
		{"/src3/*filepath/x", true},
	}
	testRoutes(t, routes)
}

func TestTireCatchMaxParams(t *testing.T) {
	tire := &node{}
	var route = "/cmd/*filepath"
	tire.addRoute(route, fakeHandler(route))
}

func TestTireInvalidNodeType(t *testing.T) {
	const panicMsg = "invalid node type"

	tire := &node{}
	tire.addRoute("/", fakeHandler("/"))
	tire.addRoute("/:page", fakeHandler("/:page"))

	tire.children[0].nType = 42

	recv := catchPanic(func() {
		tire.getValue("/test", nil)
	})
	if rs, ok := recv.(string); !ok || rs != panicMsg {
		t.Fatalf("Expected panic '"+panicMsg+"', got '%v'", recv)
	}

	recv = catchPanic(func() {
		tire.findCaseInsensitivePath("/test", true)
	})
	if rs, ok := recv.(string); !ok || rs != panicMsg {
		t.Fatalf("Expected panic '"+panicMsg+"', got '%v'", recv)
	}
}
