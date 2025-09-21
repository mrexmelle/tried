package main

import (
	"fmt"

	"github.com/mrexmelle/tried/pkg/tried"
)

type Organization struct {
	Id string `json:"id"`
}

func (o *Organization) GetId() string {
	return o.Id
}

func (o *Organization) SetId(id string) {
	o.Id = id
}

func main() {
	o := tried.New[*Organization]("ABC", ".")
	o.Insert("ABC", &Organization{Id: "ABC"})
	fmt.Println(o.Root.Id)
	o.Insert("ABC.DEF", &Organization{Id: "DEF"})
	fmt.Println(o.Root.GetChildById("DEF").Id)
	_, err := o.Insert("ABC.DEF.GHI", &Organization{Id: "GHI"})
	if err != nil {
		fmt.Printf("Err: %s\n", err)
	} else {
		fmt.Println(o.Root.GetChildById("DEF").GetChildById("GHI").Id)
	}

	_, err = o.Insert("ABC.DEF.JKL", &Organization{Id: "JKL"})
	if err != nil {
		fmt.Printf("Err: %s\n", err)
	} else {
		fmt.Println(o.Root.GetChildById("DEF").GetChildById("JKL").Id)
	}

	_, err = o.Insert("ABC.MNO.GHI", &Organization{Id: "GHI"})
	if err != nil {
		fmt.Printf("Err: %s\n", err)
	} else {
		fmt.Println(o.Root.GetChildById("MNO").GetChildById("GHI").Id)
	}

	_, err = o.Insert("GHI", &Organization{Id: "GHI"})
	fmt.Println(o.Root.Id)
	fmt.Println(err)

	x := o.Root.GetChildById("MNO").GetChildById("GHI").Locate(".")
	fmt.Printf("x: %s\n", x)

	y, err := o.Root.Search("ABC.MNO.GHI", ".")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("y: %s\n", y.Id)
	}
}
