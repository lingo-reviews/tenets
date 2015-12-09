package astutil

import (
	"go/ast"

	"github.com/juju/errors"
)

// SameIdent returns true if a and b are the same.
func SameIdent(a, b *ast.Ident) bool {
	// TODO(waigani) Don't rely on name, it could change and still be the same
	// ident.
	if a.String() != b.String() {
		return false
	}

	// TODO(waigani) this happens if ident decl is outside of current
	// file. We need to use the FileSet to find it.
	if a.Obj == nil && b.Obj == nil {
		return true
	}

	pa, err := DeclPos(a)
	if err != nil {
		// TODO(waigani) log error
		return false
	}

	pb, err := DeclPos(b)
	if err != nil {
		// TODO(waigani) log error
		return false
	}

	if pa != pb {
		return false
	}

	return true
}

// IdentDeclExpr returns the expression that the identifier was declared with.
func IdentDeclExpr(ident *ast.Ident) (ast.Expr, error) {

	if ident.Obj == nil {
		return nil, errors.Errorf("ident object is nil for ident %q", ident.Name)
	}

	switch n := ident.Obj.Decl.(type) {
	case *ast.AssignStmt:

		// find position of ident on lhs of assignment
		var identPos int
		for i, exp := range n.Lhs {
			if a, ok := exp.(*ast.Ident); ok && SameIdent(a, ident) {
				identPos = i
				break
			}
		}

		return n.Rhs[identPos], nil
	case *ast.ValueSpec:

		// find position of ident on lhs of assignment
		var identPos int
		for i, name := range n.Names {
			if name.String() == ident.String() {
				identPos = i
			}
		}

		if n.Values != nil {
			// get the rhs counterpart
			return n.Values[identPos], nil
		}

	}

	return nil, errors.Errorf("no expr found for %T", ident.Name)
}

// DeclLhsPos returns the  position of the ident's variable on the left hand
// side of the assignment operator with which it was declared.
func DeclLhsPos(ident *ast.Ident) (int, error) {
	var identPos int
	switch n := ident.Obj.Decl.(type) {
	case *ast.AssignStmt:

		// find position of ident on lhs of assignment
		for i, exp := range n.Lhs {
			if a, ok := exp.(*ast.Ident); ok && SameIdent(a, ident) {
				identPos = i
				break
			}
		}
	case *ast.ValueSpec:

		// find position of ident on lhs of assignment
		for i, name := range n.Names {
			if name.String() == ident.String() {
				identPos = i
				break
			}
		}

	default:
		return 0, errors.New("could not get lhs position of ident: unknown decl type")
	}
	return identPos, nil

}

// DeclPos returns the possition the ident was declared.
func DeclPos(n *ast.Ident) (int, error) {
	if n.Obj == nil {
		// TODO(waigani) this happens if ident decl is outside of current
		// file. We need to use the FileSet to find it.
		return 0, errors.New("ident object is nil")
	}

	switch t := n.Obj.Decl.(type) {
	case *ast.AssignStmt:
		if t.TokPos.IsValid() {
			return int(t.TokPos), nil
		}
		return 0, errors.New("token not valid")
	case *ast.ValueSpec:
		if len(t.Names) == 0 {
			return 0, errors.New("decl statement has no names")
		}
		// Even if this is not the name of the ident, it is in the same decl
		// statement. We are interested in the pos of the decl statement.
		return int(t.Names[0].NamePos), nil
	default:
		// TODO(waigani) log
	}
	return 0, errors.New("could not get declaration position")
}

func DeclaredWithLit(ident *ast.Ident) (bool, error) {
	// TODO(waigani) decl is outside current file. Assume yes until we can be sure of an issue.
	if ident.Obj == nil {
		return true, nil
	}

	expr, err := IdentDeclExpr(ident)
	if err != nil {
		return false, errors.Trace(err)
	}

	switch n := expr.(type) {

	case *ast.CompositeLit:
		return true, nil

	case *ast.CallExpr:
		if f, ok := n.Fun.(*ast.Ident); ok {
			if f.Name == "make" {
				return true, nil
			}
		}
	}

	return false, nil
}

// Expects ident to have an object of type func and returns the returned vars
// of that func.
func FuncResults(ident *ast.Ident) ([]*ast.Field, error) {
	if ident.Obj == nil {
		return nil, errors.New("ident has no object")
	}

	if ident.Obj.Kind.String() != "func" {
		return nil, errors.New("expected type func")
	}

	if funcDecl, ok := ident.Obj.Decl.(*ast.FuncDecl); ok {
		return funcDecl.Type.Results.List, nil
	}

	return nil, errors.New("could not get func results")
}

// Returns a string representation of the identifier's type.
func TypeOf(ident *ast.Ident) (string, error) {

	switch n := ident.Obj.Decl.(type) {
	case *ast.AssignStmt:
		expr, err := IdentDeclExpr(ident)
		if err != nil {
			return "", errors.Trace(err)
		}

		switch n := expr.(type) {
		case *ast.CallExpr:

			pos, err := DeclLhsPos(ident)
			if err != nil {
				return "", errors.Trace(err)
			}

			if fun, ok := n.Fun.(*ast.Ident); ok {
				results, err := FuncResults(fun)
				if err != nil {
					return "", errors.Trace(err)
				}
				if p, ok := results[pos].Type.(*ast.Ident); ok {
					return p.Name, nil
				}
			}
		}
	case *ast.ValueSpec:

		if x, ok := n.Type.(*ast.Ident); ok {
			return x.Name, nil
		}

	default:
		// TODO(waigani) log
	}
	return "", errors.New("could not find type")
}
