package resource

import "fmt"

type Currency struct {
	Code        string `json:"code"`
	Sign        string `json:"sign"`
	Description string `json:"description"`
}

func (c *Currency) String() string {
	return fmt.Sprintf("%s %s", c.Sign, c.Code)
}
