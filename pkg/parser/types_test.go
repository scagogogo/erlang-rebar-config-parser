package parser

import (
	"testing"
)

// TestCompare tests the Compare method for all term types
func TestCompare(t *testing.T) {
	// Basic types
	atom1 := Atom{Value: "test", IsQuoted: false}
	atom2 := Atom{Value: "test", IsQuoted: true} // Quoting ignored for comparison
	atom3 := Atom{Value: "other"}
	string1 := String{Value: "hello"}
	string2 := String{Value: "hello"}
	string3 := String{Value: "world"}
	int1 := Integer{Value: 42}
	int2 := Integer{Value: 42}
	int3 := Integer{Value: 43}
	float1 := Float{Value: 3.14}
	float2 := Float{Value: 3.14}
	float3 := Float{Value: 3.15}

	// Complex types
	list1 := List{Elements: []Term{atom1, int1, string1}}
	list2 := List{Elements: []Term{atom2, int2, string2}} // Should compare equal to list1
	list3 := List{Elements: []Term{atom1, int3, string1}} // Different int value
	list4 := List{Elements: []Term{atom1, int1}}          // Different length
	list5 := List{Elements: []Term{}}
	list6 := List{Elements: []Term{}}

	tuple1 := Tuple{Elements: []Term{atom1, list1}}
	tuple2 := Tuple{Elements: []Term{atom2, list2}} // Should compare equal to tuple1
	tuple3 := Tuple{Elements: []Term{atom1, list3}} // Different list element
	tuple4 := Tuple{Elements: []Term{atom1}}        // Different length
	tuple5 := Tuple{Elements: []Term{}}
	tuple6 := Tuple{Elements: []Term{}}

	tests := []struct {
		a        Term
		b        Term
		expected bool
		name     string
	}{
		{atom1, atom2, true, "Atom comparison (ignore quote)"},
		{atom1, atom3, false, "Atom comparison (different value)"},
		{string1, string2, true, "String comparison (equal)"},
		{string1, string3, false, "String comparison (different)"},
		{int1, int2, true, "Integer comparison (equal)"},
		{int1, int3, false, "Integer comparison (different)"},
		{float1, float2, true, "Float comparison (equal)"},
		{float1, float3, false, "Float comparison (different)"},
		{list1, list2, true, "List comparison (equal)"},
		{list1, list3, false, "List comparison (different element)"},
		{list1, list4, false, "List comparison (different length)"},
		{list5, list6, true, "List comparison (empty)"},
		{tuple1, tuple2, true, "Tuple comparison (equal)"},
		{tuple1, tuple3, false, "Tuple comparison (different element)"},
		{tuple1, tuple4, false, "Tuple comparison (different length)"},
		{tuple5, tuple6, true, "Tuple comparison (empty)"},
		// Type mismatches
		{atom1, string1, false, "Type mismatch (Atom vs String)"},
		{int1, float1, false, "Type mismatch (Integer vs Float)"},
		{list1, tuple1, false, "Type mismatch (List vs Tuple)"},
		{list5, tuple5, false, "Type mismatch (Empty List vs Empty Tuple)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test both a.Compare(b) and b.Compare(a)
			if result1 := tt.a.Compare(tt.b); result1 != tt.expected {
				t.Errorf("Compare(%v, %v): expected %v, got %v", tt.a, tt.b, tt.expected, result1)
			}
			if result2 := tt.b.Compare(tt.a); result2 != tt.expected {
				t.Errorf("Compare(%v, %v): expected %v, got %v", tt.b, tt.a, tt.expected, result2)
			}
		})
	}
}

// TestTermStringMethods checks the String() method for various term types
func TestTermStringMethods(t *testing.T) {
	tests := []struct {
		term     Term
		expected string
	}{
		{Atom{Value: "simple", IsQuoted: false}, "simple"},
		{Atom{Value: "quoted atom", IsQuoted: true}, "'quoted atom'"},
		{String{Value: "hello"}, `"hello"`},
		{Integer{Value: 123}, "123"},
		{Integer{Value: -45}, "-45"},
		{Float{Value: 1.23}, "1.23"},
		{Float{Value: -0.5}, "-0.5"},
		{List{Elements: []Term{Integer{Value: 1}, Atom{Value: "a"}}}, "[1, a]"},
		{List{Elements: []Term{}}, "[]"},
		{Tuple{Elements: []Term{Atom{Value: "key"}, String{Value: "val"}}}, "{key, \"val\"}"},
		{Tuple{Elements: []Term{}}, "{}"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if str := tt.term.String(); str != tt.expected {
				t.Errorf("String() for %T: expected %q, got %q", tt.term, tt.expected, str)
			}
		})
	}
}
