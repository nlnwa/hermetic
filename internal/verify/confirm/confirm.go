package confirmImplmentation

import (
	"fmt"
)

func Verify(confirmTopicName string) error {
	fmt.Printf("Reading messages from topic '%s'\n", confirmTopicName)
	fmt.Printf("This command is not implemented yet\n" +
		"Aims to solve issue https://github.com/nlnwa/hermetic/issues/3 ")
	return nil
}
