package main

import (
	"fmt"
	"github.com/indigo-web/yajson"
	"github.com/indigo-web/yajson/flect"
)

type myJSONModel struct {
	Something string
	Nothing   string
	Lol       struct {
		Lorem string
		Ipsum int64
	}
}

func (m myJSONModel) String() string {
	return fmt.Sprintf(
		"myJSONModel{Something: %s, Nothing: %s, Lol.Lorem: %s}", m.Something, m.Nothing, m.Lol.Lorem,
	)
}

func main() {
	j := `
	{
		"Boah": "857",
		"Something": "okay, let it be",
		"Nothing": "Hello, world!",
		"Lol": {
			"Lorem": "ipsum!"
		}
	}
	`

	model := flect.NewModel[myJSONModel](nil)
	fmt.Println(model)

	parser := yajson.New[myJSONModel]()
	result, err := parser.Parse(j)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(result.String())
}
