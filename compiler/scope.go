package compiler

import "tiny-go/ast"

type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}

type Object struct {
	Name        string
	MangledName string
	Type        string
	ast.Node
}

func NewScope(outer *Scope) *Scope {
	return &Scope{outer, make(map[string]*Object)}
}

func (s *Scope) HasName(name string) bool {
	_, ok := s.Objects[name]
	return ok
}

func (s *Scope) Lookup(name string) (*Scope, *Object) {
	for ; s != nil; s = s.Outer {
		if obj := s.Objects[name]; obj != nil {
			return s, obj
		}
	}
	return nil, nil
}

func (s *Scope) Insert(obj *Object) (alt *Object) {
	if alt = s.Objects[obj.Name]; alt == nil {
		s.Objects[obj.Name] = obj
	}
	return
}
