package neko

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

const wildcard = "{}"

var pathTokenFinder = regexp.MustCompile("\\{[^{]*\\}")
var pathTokenValidator = regexp.MustCompile("^\\{([a-z_]+)(:([a-z]+))?\\}$")
var pathTokenReplacer = regexp.MustCompile("\\{[^{]+\\}")
var validTypeStrings = map[string]reflect.Kind{
	"int":    reflect.Int,
	"string": reflect.String,
}
var validMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}
var reEscapeChars = map[byte]string{
	'-': "\\-",
	'.': "\\.",
	'!': "\\!",
	'$': "\\$",
	'&': "\\&",
	'(': "\\(",
	')': "\\)",
	'*': "\\*",
	'+': "\\+",
	':': "\\:",
}

type PathToken struct {
	name string
	kind reflect.Kind
}

type Route struct {
	path        string
	pathTokens  []*PathToken
	pathRegex   *regexp.Regexp
	length      int
	handler     map[string]Handler
	router      *Router
	middlewares []Middleware
}

func (rt *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	var handler Handler
	var ok bool
	handler, ok = rt.handler[r.Method]
	if !ok {
		handler, ok = rt.handler["*"]
		if !ok {
			return NewStatusErrf(405, "Method \"%s\" not allowed by handler", r.Method)
		}
	}
	return handler.ServeHTTP(w, r)
}

func (rt *Route) Middleware(middlewares ...Middleware) *Route {
	if len(rt.handler) > 0 {
		panic("Middlewares must be added prior to adding handlers")
	}

	rt.middlewares = append(rt.middlewares, middlewares...)

	return rt
}

func (rt *Route) HandlerFunc(handlerFunc HandlerFunc, methods ...string) *Route {
	allMiddlewares := []Middleware{}
	allMiddlewares = append(allMiddlewares, rt.middlewares...)
	allMiddlewares = append(allMiddlewares, rt.router.middlewares...)
	if len(methods) == 0 {
		rt.handler["*"] = Cascade(MakeHandler(handlerFunc), allMiddlewares...)
	}
	for _, method := range methods {
		_, ok := validMethods[method]
		if !ok {
			panic(fmt.Sprintf("\"%s\" is not a valid http method", method))
		}
		rt.handler[method] = Cascade(MakeHandler(handlerFunc), allMiddlewares...)
	}

	return rt
}

func (rt *Route) Handler(handler Handler, methods ...string) *Route {
	allMiddlewares := []Middleware{}
	allMiddlewares = append(allMiddlewares, rt.middlewares...)
	allMiddlewares = append(allMiddlewares, rt.router.middlewares...)
	if len(methods) == 0 {
		rt.handler["*"] = Cascade(handler, allMiddlewares...)
	}
	for _, method := range methods {
		_, ok := validMethods[method]
		if !ok {
			panic(fmt.Sprintf("\"%s\" is not a valid http method", method))
		}
		rt.handler[method] = Cascade(handler, allMiddlewares...)
	}

	return rt
}

type Router struct {
	routerNode  *RouterNode
	routes      []*Route
	middlewares []Middleware
}

type RouterNode struct {
	lookup map[string]*RouterNode
	routes []*Route
}

func (rn *RouterNode) FindMatches(tokens []string, matches []*Route) []*Route {
	m := append(matches, rn.routes...)

	if len(tokens) > 0 {
		token := tokens[0]
		nextRouterNode, ok := rn.lookup[token]
		if ok {
			m = append(matches, nextRouterNode.FindMatches(tokens[1:], m)...)
		}

		nextRouterNode, ok = rn.lookup[wildcard]
		if ok {
			m = append(matches, nextRouterNode.FindMatches(tokens[1:], m)...)
		}
	}

	return m
}

func NewRouterNode() *RouterNode {
	return &RouterNode{
		lookup: map[string]*RouterNode{},
		routes: []*Route{},
	}
}

func NewRouter() *Router {
	return &Router{
		routerNode:  NewRouterNode(),
		routes:      []*Route{},
		middlewares: []Middleware{},
	}
}

func (r *Router) AddMiddleware(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) AddRoute(routePath string, isPrefix bool) *Route {
	if isPrefix && !strings.HasSuffix(routePath, "/") {
		panic(fmt.Sprintf("Route \"%s\" which is identified as a prefix must end in a \"/\" in order to qualify", routePath))
	}
	pathTokens := []*PathToken{}
	tokens := strings.Split(routePath, "/")
	for _, token := range tokens {
		match := pathTokenFinder.MatchString(token)
		if match {
			if token[0] != '{' || token[len(token)-1] != '}' {
				panic("Route tokens must occupy the entire span between path delimiters")
			}

			matches := pathTokenValidator.FindStringSubmatch(token)
			if matches == nil {
				panic(fmt.Sprintf("Route \"%s\" contains a malformed token \"%s\"", routePath, token))
			}

			tokenName := matches[1]
			tokenType := matches[3]
			if tokenType == "" {
				tokenType = "string"
			}
			kind, ok := validTypeStrings[tokenType]
			if !ok {
				panic(fmt.Sprintf("\"%s\" in route \"%s\" is not a valid token type (must be int or string)", tokenType, routePath))
			}

			pathTokens = append(pathTokens, &PathToken{
				name: tokenName,
				kind: kind,
			})
		}
	}

	var b strings.Builder
	inToken := false
	for i := 0; i < len(routePath); i++ {
		c := routePath[i]
		if c == '{' {
			inToken = true
		}
		if c == '}' {
			inToken = false
		}
		if inToken {
			b.WriteByte(c)
		} else {
			if rewrite, ok := reEscapeChars[c]; ok {
				b.WriteString(rewrite)
			} else {
				b.WriteByte(c)
			}
		}
	}
	regexSafeRoutePath := b.String()

	pathRegexStr := "^" + pathTokenReplacer.ReplaceAllString(regexSafeRoutePath, "([^/]+)")
	if !isPrefix {
		pathRegexStr = pathRegexStr + "$"
	}
	pathRegex, err := regexp.Compile(pathRegexStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to compile regular expression for route \"%s\"\n%v", routePath, err))
	}

	// Cut off the last token since in the case of routes ending in "/", this will add yet another blank
	// entry which is undesired here
	if len(tokens) > 0 && tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	route := &Route{
		path:        routePath,
		pathTokens:  pathTokens,
		pathRegex:   pathRegex,
		length:      len(tokens),
		handler:     map[string]Handler{},
		router:      r,
		middlewares: []Middleware{},
	}
	r.routes = append(r.routes, route)

	// Place the route in its correct location for lookup
	curNode := r.routerNode
	for _, token := range tokens {
		t := token
		if strings.HasPrefix(token, "{") && strings.HasSuffix(token, "}") {
			t = wildcard
		}
		nextNode, ok := curNode.lookup[t]
		if !ok {
			nextNode = NewRouterNode()
			curNode.lookup[t] = nextNode
		}
		curNode = nextNode
	}
	curNode.routes = append(curNode.routes, route)

	return route
}

// FindMatch looks through the routing entries for the most specific match against the candidate.  If a match
// cannot be established then nil is returned
func (r *Router) FindMatch(candidate string) *Route {
	tokens := strings.Split(candidate, "/")
	// Remove the final empty token if one exists
	if len(tokens) > 0 && tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	matches := r.routerNode.FindMatches(tokens, []*Route{})
	if len(matches) == 0 {
		return nil
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].length > matches[j].length
	})

	return matches[0]
}
