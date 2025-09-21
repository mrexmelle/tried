package node

import (
	"slices"
	"strings"
	"testing"
)

type AppendTestCase struct {
	name     string
	input    []string
	expected []string
}

type InsertTestCase struct {
	name     string
	input    InsertInput
	expected InsertExpected
}

type InsertInput struct {
	Text            string
	Nodes           []string
	ReplaceIfExists bool
}

type InsertExpected struct {
	Text  string
	Nodes []string
}

type LocateInput struct {
}

type GetChildByIdTestCase[T any] struct {
	name     string
	input    string
	expected *Class[T]
}

type Dummy struct {
	Text string
}

func TestAppendChild(t *testing.T) {
	tc := []AppendTestCase{
		{
			name:     "Appending 1 child",
			input:    []string{"DEF"},
			expected: []string{"DEF"},
		},
		{
			name:     "Appending 2 children",
			input:    []string{"DEF", "GHI"},
			expected: []string{"GHI", "DEF"},
		},
		{
			name:     "Appending 0 children",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, c := range tc {
		node := New(
			"ABC",
			&Dummy{Text: "ABC"},
			nil,
		)
		for _, in := range c.input {
			newNode := New(
				in,
				&Dummy{Text: "ABC"},
				nil,
			)
			node.AppendChild(*newNode)
		}

		cids := []string{}
		for _, child := range node.Children {
			cids = append(cids, child.Id)
		}

		if len(cids) != len(c.expected) {
			t.Errorf(
				"[%s]\nResult: %v\nExpected: %v\n",
				c.name,
				cids,
				c.expected,
			)
			return
		}

		slices.Sort(cids)
		slices.Sort(c.expected)
		if slices.Compare(c.expected, cids) != 0 {
			t.Errorf(
				"[%s]\nResult: %v\nExpected: %v\n",
				c.name,
				cids,
				c.expected,
			)
		}
	}
}

func TestInsert(t *testing.T) {
	tc := []InsertTestCase{
		{
			name: "Inserting a 1-level node, replacing data if exists",
			input: InsertInput{
				Text:            "Matt",
				Nodes:           []string{"ABC"},
				ReplaceIfExists: true,
			},
			expected: InsertExpected{
				Text:  "Matt",
				Nodes: []string{"ABC"},
			},
		},
		{
			name: "Inserting a 2-level node, replacing data if exists",
			input: InsertInput{
				Text:            "Matt",
				Nodes:           []string{"ABC", "DEF"},
				ReplaceIfExists: true,
			},
			expected: InsertExpected{
				Text:  "Matt",
				Nodes: []string{"ABC", "DEF"},
			},
		},
		{
			name: "Inserting a 3-level node, replacing data if exists",
			input: InsertInput{
				Text:            "Matt",
				Nodes:           []string{"ABC", "DEF", "GHI"},
				ReplaceIfExists: true,
			},
			expected: InsertExpected{
				Text:  "Matt",
				Nodes: []string{"ABC", "DEF", "GHI"},
			},
		},
		{
			name: "Inserting a 1-level node, not replacing data if exists",
			input: InsertInput{
				Text:            "Matt",
				Nodes:           []string{"ABC"},
				ReplaceIfExists: false,
			},
			expected: InsertExpected{
				Text:  "Jess",
				Nodes: []string{"ABC"},
			},
		},
	}

	var node *Class[Dummy] = nil
	for _, c := range tc {
		node = New("ABC", &Dummy{Text: "Jess"}, nil)

		path := strings.Join(c.input.Nodes, ".")
		_, err := node.Insert(path, ".", Dummy{Text: c.input.Text}, c.input.ReplaceIfExists)
		if (c.input.ReplaceIfExists && err != nil) || (!c.input.ReplaceIfExists && err != ErrAlreadyExists) {
			t.Errorf("[%s]\nFailed to insert\nReason: %s\n", c.name, err.Error())
			break
		}

		ptr := node
		i := 0
		for len(ptr.Children) > 0 {
			if ptr.Id != c.expected.Nodes[i] {
				t.Errorf(
					"[%s]\nResult: %s\nExpected: %s\n",
					c.name,
					ptr.Id,
					c.expected.Nodes[i],
				)
			}
			i++
			ptr = ptr.GetChildById(c.expected.Nodes[i])
		}

		if ptr.Data.Text != c.expected.Text {
			t.Errorf(
				"[%s]\nResult: %v\nExpected: %v\n",
				c.name,
				node.Data,
				c.expected,
			)
		}
	}
}

func TestGetChildById(t *testing.T) {
	tc := []GetChildByIdTestCase[Dummy]{
		{
			name:     "Getting a child",
			input:    "DEF",
			expected: New("DEF", &Dummy{Text: "DEF"}, nil),
		},
		{
			name:     "Getting a child",
			input:    "XYZ",
			expected: nil,
		},
	}

	node := New("ABC", &Dummy{}, nil)
	child := New("DEF", &Dummy{Text: "DEF"}, nil)
	node.AppendChild(*child)

	for _, c := range tc {
		exp := node.GetChildById(c.input)
		if exp != nil {
			if exp.Id != c.expected.Id {
				t.Errorf(
					"[%s]\nResult: %s\nExpected: %s\n",
					c.name,
					exp.Id,
					c.expected.Id,
				)
			}

			if exp.Data.Text != c.expected.Data.Text {
				t.Errorf(
					"[%s]\nResult: %s\nExpected: %s\n",
					c.name,
					exp.Data.Text,
					c.expected.Data.Text,
				)
			}
		} else {
			if c.expected != nil {
				t.Errorf(
					"[%s]\nResult: %v\nExpected: %s\n",
					c.name,
					nil,
					c.expected.Id,
				)
			}
		}
	}
}

func TestLocate(t *testing.T) {
	n := New("ABC", &Dummy{Text: "ABC"}, nil)
	o := New("DEF", &Dummy{Text: "DEF"}, nil)
	p := New("GHI", &Dummy{Text: "GHI"}, nil)
	n.AppendChild(*o)
	n.GetChildById("DEF").AppendChild(*p)
	l1 := n.GetChildById("DEF").GetChildById("GHI").Locate(".")
	if l1 != "ABC.DEF.GHI" {
		t.Errorf(
			"[%s]\nResult: %s\nExpected: %s\n",
			"Test Locating",
			l1,
			"ABC.DEF.GHI",
		)
	}
}

func TestSearch(t *testing.T) {
	n := New("ABC", &Dummy{Text: "ABC"}, nil)
	o := New("DEF", &Dummy{Text: "DEF"}, nil)
	p := New("GHI", &Dummy{Text: "GHI"}, nil)
	n.AppendChild(*o)
	n.GetChildById("DEF").AppendChild(*p)

	obj, err := n.Search("ABC.DEF.GHI", ".")
	if err == nil && obj.Locate(".") != "ABC.DEF.GHI" {
		t.Errorf(
			"[%s]\nResult: %s\nExpected: %s\n",
			"Searching an existing node",
			obj.Locate("."),
			"ABC.DEF.GHI",
		)
	}

	_, err = n.Search("ABC.DEF.JKL", ".")
	if err != ErrNotFound {
		t.Errorf(
			"[%s]\nResult: %s\nExpected: %s\n",
			"Searching an unexisting node",
			ErrNotFound.Error(),
			"nil",
		)
	}

	_, err = n.Search("", ".")
	if err != ErrBadPath {
		t.Errorf(
			"[%s]\nResult: %s\nExpected: %s\n",
			"Searching with an invalid path",
			ErrBadPath.Error(),
			"nil",
		)
	}
}
