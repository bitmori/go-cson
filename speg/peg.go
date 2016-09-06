package speg

import (
	"fmt"
	"regexp"
)

type ParseError struct {
	Msg    string
	Text   string
	Offset int
	Line   int
	Col    int
}

type _UnexpectedError struct {
	error
	state *pegState
	expr  string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("parse error %s on line %d, column %d, offset %d, content: %s", p.Msg, p.Line, p.Col, p.Offset, p.Text)
}

type strBoolPair struct {
	s string
	b bool
}

type pegState struct {
	pos      int
	line     int
	col      int
	vars     map[string]strBoolPair
	commited bool
}

type pegError struct {
	msg  string
	line int
	col  int
}

type peg struct {
	s       string
	states  []pegState
	errors  map[int]string
	reCache map[string]regexp.Regexp
}

func newPeg(p *peg) *peg {
	if p == nil {
		p = &peg{}
	}
	p.states = append(p.states, pegState{pos: 0, line: 1, col: 1, commited: false})
	return p
}

func (p *peg) _call(expr string) {
	var compiled Regexp
	if compiled, ok := p.reCache[expr]; !ok {
		compiled, err := regexp.Compile(expr)
		if err == nil {
			p.reCache[expr] = compiled
		} else {

		}
	}
}

func (p peg) _repr() string {
	pos := p.states[-1].pos
	vars := map[string]string{}
	for _, st := range p.states {
		for k, v := range st.vars {
			vars[k] = v
		}
	}
	return fmt.Sprintf("Peg(%s*%s, %v)", p.s[:pos], p.s[pos:], vars)
}

func (p *peg) _eof() {
	if p.states[len(p.states)-1].pos != len(p.s) {
		p._error(nil, nil)
	}
}

func (p *peg) _error(err string, expr string) {
	st := p.states[len(p.states)-1]
	if err == nil {
		err = fmt.Sprintf("expected %s, found %s", expr, p.s[st.pos:st.pos+4])
	}
	p.errors[st.pos] = &pegError{msg: err, line: st.line, col: st.col}
	panic(_UnexpectedError{state: &st, expr: expr})
}

func (p peg) get(key, def string) string {
	for i := len(p.states) - 1; i >= 0; i-- {
		if v, ok := p.states[i].vars[key]; ok {
			return v.s
		}
	}
	return def
}

func (p *peg) set(key, val string, global bool) {
	p.states[len(p.states)-1].vars[key] = strBoolPair{s: val, b: global}
}

func (p *peg) not(s) {
	p._call(s)
	p._error(nil, nil)
}

func (p *peg) _enter() {
	lastState := &(p.states[len(p.states)-1])
	lastState.commited = false

	p.states = append(p.states, pegState{pos: lastState.pos, line: lastState.line, col: lastState.col})
}

func (p *peg) _exit() {

}

func (p *peg) commit() {
	currentState := &(p.states[len(p.states)-1])
	previousState := &(p.states[len(p.states)-2])

	for key, v := range currentState.vars {
		if v.b {
			if previousStateVal, ok := previousState.vars[key]; ok {
				previousStateVal.s = v.s
			} else {
				previousState.vars[key] = strBoolPair{s: v.s, b: true}
			}
		}
	}

	previousState.pos = currentState.pos
	previousState.line = currentState.line
	previousState.col = currentState.col
	previousState.commited = true
}

func (p *peg) toBoolean() {
	return p.states[len(p.states)-1].commited
}
