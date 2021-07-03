package main

import (
	"fmt"
	link "html-parser/cmd"
	"strings"
)

const exampleHTML = `
<html>
<body>
  	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
	<a href="/page">A link to some other 
		<span> something </span>
		page
	</a>
</body>
</html>`

func main() {
	r := strings.NewReader(exampleHTML)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
