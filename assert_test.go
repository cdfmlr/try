package try

import (
	"reflect"
	"testing"
)

func TestAssert(t *testing.T) {
	// struct type for deep equal test
	type A struct {
		B int
		C []string
	}

	tests := []struct {
		name         string
		f            func() (err error)
		expectedErr  error
		expectPanics bool
	}{
		{
			"assertTrueWithoutHandle_badUsage",
			func() (err error) {
				Assert(true)
				return nil
			},
			nil,
			false,
		},
		{
			"assertTrueWithHandle",
			func() (err error) {
				defer Handle(&err)
				Assert(true)
				return nil
			},
			nil,
			false,
		},
		{
			"assertFalseWithoutHandle_badUsage",
			func() (err error) {
				Assert(false)
				return nil
			},
			nil,
			true,
		},
		{
			"assertFalseWithHandle",
			func() (err error) {
				defer Handle(&err)
				Assert(false)
				return nil
			},
			&ErrAssertion{
				Msg:      "",
				Got:      false,
				Expected: true,
			},
			false,
		},
		{
			"assertWithMsg",
			func() (err error) {
				defer Handle(&err)
				Assert(false, "this is a message", "with multiple parts")
				return nil
			},
			&ErrAssertion{
				Msg:      "this is a message with multiple parts",
				Got:      false,
				Expected: true,
			},
			false,
		},
		{
			"assertEqual_false",
			func() (err error) {
				defer Handle(&err)
				AssertEqual(1, 2)
				return nil
			},
			&ErrAssertion{
				Msg:      "",
				Got:      1,
				Expected: 2,
			},
			false,
		},
		{
			"assertEqual_true",
			func() (err error) {
				defer Handle(&err)
				AssertEqual("hello", "hello")
				return nil
			},
			nil,
			false,
		},
		{
			"assertDeepEqual_true",
			func() (err error) {
				defer Handle(&err)

				type A struct {
					B int
					C []string
				}

				AssertDeepEqual(A{B: 1, C: []string{"a", "b"}}, A{B: 1, C: []string{"a", "b"}})

				return nil
			},
			nil,
			false,
		},
		{
			"assertDeepEqual_false",
			func() (err error) {
				defer Handle(&err)

				AssertDeepEqual(A{B: 1, C: []string{"a", "b"}}, A{B: 1, C: []string{"a", "b", "c"}})

				return nil
			},
			&ErrAssertion{
				Msg:      "",
				Got:      A{B: 1, C: []string{"a", "b"}},
				Expected: A{B: 1, C: []string{"a", "b", "c"}},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			defer func() {
				r := recover()
				if (r != nil) != tt.expectPanics {
					t.Errorf("expect panics = %#v, got %#v", tt.expectPanics, r)
				}
				t.Logf("✅ recover() = %v (type %T)\n", r, r)
			}()

			err = tt.f()
			if !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("expect err = %#v, got %#v", tt.expectedErr, err)
			}
			t.Logf("✅ err = %#v\n", err)
		})
	}

	t.Run("ErrAssertError", func(t *testing.T) {
		err := assert(false, "this is a message", "with multiple parts")
		if err == nil {
			t.Errorf("expect err != nil, got %#v", err)
		}
		expected := "assert failed (expected true, got false): this is a message with multiple parts"
		if err.Error() != expected {
			t.Errorf("expect err.Error() = %#v, got %#v", expected, err.Error())
		}
	})
}
