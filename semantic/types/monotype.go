package types

import (
	"fmt"
	"sort"
	"strings"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/influxdata/flux/codes"
	"github.com/influxdata/flux/internal/errors"
	"github.com/influxdata/flux/semantic/internal/fbsemantic"
)

type fbTabler interface {
	Init(buf []byte, i flatbuffers.UOffsetT)
	Table() flatbuffers.Table
}

// MonoType represents a monotype.  This struct is a thin wrapper around
// Go code generated by the FlatBuffers compiler.
type MonoType struct {
	mt  fbsemantic.MonoType
	tbl fbTabler
}

// NewMonoType constructs a new monotype from a FlatBuffers table and the given kind of monotype.
func NewMonoType(tbl *flatbuffers.Table, t fbsemantic.MonoType) (*MonoType, error) {
	var tbler fbTabler
	switch t {
	case fbsemantic.MonoTypeNONE:
		return nil, errors.Newf(codes.Internal, "missing type, got type: %v", fbsemantic.EnumNamesMonoType[t])
	case fbsemantic.MonoTypeBasic:
		tbler = new(fbsemantic.Basic)
	case fbsemantic.MonoTypeVar:
		tbler = new(fbsemantic.Var)
	case fbsemantic.MonoTypeArr:
		tbler = new(fbsemantic.Arr)
	case fbsemantic.MonoTypeRow:
		tbler = new(fbsemantic.Row)
	case fbsemantic.MonoTypeFun:
		tbler = new(fbsemantic.Fun)
	default:
		return nil, errors.Newf(codes.Internal, "unknown type (%v)", t)
	}
	tbler.Init(tbl.Bytes, tbl.Pos)
	return &MonoType{mt: t, tbl: tbler}, nil
}

// Kind specifies a particular kind of monotype.
type Kind fbsemantic.MonoType

const (
	Unknown = Kind(fbsemantic.MonoTypeNONE)
	Basic   = Kind(fbsemantic.MonoTypeBasic)
	Var     = Kind(fbsemantic.MonoTypeVar)
	Arr     = Kind(fbsemantic.MonoTypeArr)
	Row     = Kind(fbsemantic.MonoTypeRow)
	Fun     = Kind(fbsemantic.MonoTypeFun)
)

// Kind returns what kind of monotype the receiver is.
func (mt *MonoType) Kind() Kind {
	return Kind(mt.mt)
}

// BasicKind specifies a basic type.
type BasicKind fbsemantic.Type

const (
	Bool     = BasicKind(fbsemantic.TypeBool)
	Int      = BasicKind(fbsemantic.TypeInt)
	Uint     = BasicKind(fbsemantic.TypeUint)
	Float    = BasicKind(fbsemantic.TypeFloat)
	String   = BasicKind(fbsemantic.TypeString)
	Duration = BasicKind(fbsemantic.TypeDuration)
	Time     = BasicKind(fbsemantic.TypeTime)
	Regexp   = BasicKind(fbsemantic.TypeRegexp)
	Bytes    = BasicKind(fbsemantic.TypeBytes)
)

func getBasic(tbl fbTabler) (*fbsemantic.Basic, error) {
	b, ok := tbl.(*fbsemantic.Basic)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a basic type")
	}
	return b, nil
}

// Basic returns the basic type for this monotype if it is a basic type,
// and an error otherwise.
func (mt *MonoType) Basic() (BasicKind, error) {
	b, err := getBasic(mt.tbl)
	if err != nil {
		return Bool, err
	}
	return BasicKind(b.T()), nil
}

func getVar(tbl fbTabler) (*fbsemantic.Var, error) {
	v, ok := tbl.(*fbsemantic.Var)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a type var")
	}
	return v, nil

}

// VarNum returns the type variable number if this monotype is a type variable,
// and an error otherwise.
func (mt *MonoType) VarNum() (uint64, error) {
	v, err := getVar(mt.tbl)
	if err != nil {
		return 0, err
	}
	return v.I(), nil
}

func monoTypeFromVar(v *fbsemantic.Var) *MonoType {
	return &MonoType{
		mt:  fbsemantic.MonoTypeVar,
		tbl: v,
	}
}

func getFun(tbl fbTabler) (*fbsemantic.Fun, error) {
	f, ok := tbl.(*fbsemantic.Fun)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a function")
	}
	return f, nil
}

// NumArguments returns the number of arguments if this monotype is a function,
// and an error otherwise.
func (mt *MonoType) NumArguments() (int, error) {
	f, err := getFun(mt.tbl)
	if err != nil {
		return 0, err
	}
	return f.ArgsLength(), nil
}

// Argument returns the argument give an ordinal position if this monotype is a function,
// and an error otherwise.
func (mt *MonoType) Argument(i int) (*Argument, error) {
	f, err := getFun(mt.tbl)
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= f.ArgsLength() {
		return nil, errors.Newf(codes.Internal, "request for out-of-bounds argument: %v of %v", i, f.ArgsLength())
	}
	a := new(fbsemantic.Argument)
	if !f.Args(a, i) {
		return nil, errors.New(codes.Internal, "missing argument")
	}
	return newArgument(a)
}

// SortedArguments returns a slice of function arguments,
// sorted by argument name, if this monotype is a function.
func (mt MonoType) SortedArguments() ([]*Argument, error) {
	nargs, err := mt.NumArguments()
	if err != nil {
		return nil, err
	}
	args := make([]*Argument, nargs)
	for i := 0; i < nargs; i++ {
		arg, err := mt.Argument(i)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}
	sort.Slice(args, func(i, j int) bool {
		return string(args[i].Name()) < string(args[j].Name())
	})
	return args, nil
}

func (mt *MonoType) ReturnType() (*MonoType, error) {
	f, ok := mt.tbl.(*fbsemantic.Fun)
	if !ok {
		return nil, errors.New(codes.Internal, "ReturnType() called on non-function MonoType")
	}
	tbl := new(flatbuffers.Table)
	if !f.Retn(tbl) {
		return nil, errors.New(codes.Internal, "missing return type")
	}
	return NewMonoType(tbl, f.RetnType())
}

func getArr(tbl fbTabler) (*fbsemantic.Arr, error) {
	arr, ok := tbl.(*fbsemantic.Arr)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not an array")
	}
	return arr, nil
}

// ElemType returns the element type if this monotype is an array, and an error otherise.
func (mt *MonoType) ElemType() (*MonoType, error) {
	arr, err := getArr(mt.tbl)
	if err != nil {
		return nil, err
	}
	tbl := new(flatbuffers.Table)
	if !arr.T(tbl) {
		return nil, errors.New(codes.Internal, "missing array type")
	}
	return NewMonoType(tbl, arr.TType())
}

func getRow(tbl fbTabler) (*fbsemantic.Row, error) {
	row, ok := tbl.(*fbsemantic.Row)
	if !ok {
		return nil, errors.New(codes.Internal, "MonoType is not a row")
	}
	return row, nil

}

// NumProperties returns the number of properties if this monotype is a row, and an error otherwise.
func (mt *MonoType) NumProperties() (int, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return 0, err
	}
	return row.PropsLength(), nil
}

// Property returns a property given its ordinal position if this monotype is a row, and an error otherwise.
func (mt *MonoType) Property(i int) (*Property, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= row.PropsLength() {
		return nil, errors.Newf(codes.Internal, "request for out-of-bounds property: %v of %v", i, row.PropsLength())
	}
	p := new(fbsemantic.Prop)
	if !row.Props(p, i) {
		return nil, errors.New(codes.Internal, "missing property")
	}
	return &Property{fb: p}, nil
}

// SortedProperties returns the properties for a Row monotype, sorted by
// key.  It's possible that there are duplicate keys with different types,
// in this case, this function preserves their order.
func (mt *MonoType) SortedProperties() ([]*Property, error) {
	nps, err := mt.NumProperties()
	if err != nil {
		return nil, err
	}
	ps := make([]*Property, nps)
	for i := 0; i < nps; i++ {
		ps[i], err = mt.Property(i)
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].Name() == ps[j].Name() {
			return i < j
		}
		return ps[i].Name() < ps[j].Name()
	})
	return ps, nil
}

// Extends returns the extending type variable if this monotype is a row, and an error otherwise.
func (mt *MonoType) Extends() (*MonoType, error) {
	row, err := getRow(mt.tbl)
	if err != nil {
		return nil, err
	}
	v := row.Extends(nil)
	if v == nil {
		return nil, nil
	}
	return monoTypeFromVar(v), nil
}

// Argument represents a function argument.
type Argument struct {
	*fbsemantic.Argument
}

func newArgument(fb *fbsemantic.Argument) (*Argument, error) {
	if fb == nil {
		return nil, errors.Newf(codes.Internal, "nil argument")
	}
	return &Argument{Argument: fb}, nil
}

// TypeOf returns the type of the function argument.
func (a *Argument) TypeOf() (*MonoType, error) {
	tbl := new(flatbuffers.Table)
	if !a.T(tbl) {
		return nil, errors.New(codes.Internal, "missing argument type")
	}
	argTy, err := NewMonoType(tbl, a.TType())
	if err != nil {
		return nil, err
	}
	return argTy, nil
}

// Property represents a property of a row.
type Property struct {
	fb *fbsemantic.Prop
}

// Name returns the name of the property.
func (p *Property) Name() string {
	return string(p.fb.K())
}

// TypeOf returns the type of the property.
func (p *Property) TypeOf() (*MonoType, error) {
	tbl := new(flatbuffers.Table)
	if !p.fb.V(tbl) {
		return nil, errors.Newf(codes.Internal, "missing property type")
	}
	return NewMonoType(tbl, p.fb.VType())
}

// String returns a string representation of this monotype.
func (mt *MonoType) String() string {
	switch tk := mt.Kind(); tk {
	case Unknown:
		return "<" + fbsemantic.EnumNamesMonoType[fbsemantic.MonoType(tk)] + ">"
	case Basic:
		b, err := mt.Basic()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		return strings.ToLower(fbsemantic.EnumNamesType[byte(b)])
	case Var:
		i, err := mt.VarNum()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		return fmt.Sprintf("t%d", i)
	case Arr:
		et, err := mt.ElemType()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		return "[" + et.String() + "]"
	case Row:
		var sb strings.Builder
		sb.WriteString("{")
		sprops, err := mt.SortedProperties()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		needBar := false
		for i := 0; i < len(sprops); i++ {
			if needBar {
				sb.WriteString(" | ")
			} else {
				needBar = true
			}
			prop := sprops[i]
			sb.WriteString(prop.Name() + ": ")
			ty, err := prop.TypeOf()
			if err != nil {
				return "<" + err.Error() + ">"
			}
			sb.WriteString(ty.String())
		}
		extends, err := mt.Extends()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		if extends != nil {
			if needBar {
				sb.WriteString(" | ")
			}
			sb.WriteString(extends.String())
		}
		sb.WriteString("}")
		return sb.String()
	case Fun:
		var sb strings.Builder
		sb.WriteString("(")
		needComma := false
		sargs, err := mt.SortedArguments()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		for _, arg := range sargs {
			if needComma {
				sb.WriteString(", ")
			} else {
				needComma = true
			}
			if arg.Optional() {
				sb.WriteString("?")
			} else if arg.Pipe() {
				sb.WriteString("<-")
			}
			sb.WriteString(string(arg.Name()) + ": ")
			argTyp, err := arg.TypeOf()
			if err != nil {
				return "<" + err.Error() + ">"
			}
			sb.WriteString(argTyp.String())
		}
		sb.WriteString(") -> ")
		rt, err := mt.ReturnType()
		if err != nil {
			return "<" + err.Error() + ">"
		}
		sb.WriteString(rt.String())
		return sb.String()
	default:
		return "<" + fmt.Sprintf("unknown monotype (%v)", tk) + ">"
	}
}

func updateTVarMap(counter *int, m map[uint64]int, tv uint64) {
	if _, ok := m[tv]; ok {
		return
	}
	m[tv] = *counter
	*counter++
}

func (mt *MonoType) getCanonicalMapping(counter *int, tvm map[uint64]int) error {
	switch tk := mt.Kind(); tk {
	case Var:
		tv, err := mt.VarNum()
		if err != nil {
			return err
		}
		updateTVarMap(counter, tvm, tv)
	case Arr:
		et, err := mt.ElemType()
		if err != nil {
			return err
		}
		if err := et.getCanonicalMapping(counter, tvm); err != nil {
			return err
		}
	case Row:
		n_props, err := mt.NumProperties()
		if err != nil {
			return err
		}
		for i := 0; i < n_props; i++ {
			p, err := mt.Property(i)
			if err != nil {
				return err
			}
			pt, err := p.TypeOf()
			if err != nil {
				return err
			}
			if err := pt.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
		evar, err := mt.Extends()
		if err != nil {
			return err
		}
		if evar != nil {
			if err := evar.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
	case Fun:
		nargs, err := mt.NumArguments()
		if err != nil {
			return err
		}
		for i := 0; i < nargs; i++ {
			arg, err := mt.Argument(i)
			if err != nil {
				return err
			}
			at, err := arg.TypeOf()
			if err != nil {
				return err
			}
			if err := at.getCanonicalMapping(counter, tvm); err != nil {
				return err
			}
		}
		rt, err := mt.ReturnType()
		if err != nil {
			return err
		}
		if err := rt.getCanonicalMapping(counter, tvm); err != nil {
			return err
		}
	}

	return nil
}
