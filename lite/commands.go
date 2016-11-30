package lite

import (
	"fmt"
	"os"

	// "github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
)

func WatchAction(c *cli.Context) error {
	fmt.Println(os.Getwd())
	fmt.Println("begin watch")
	return nil
}
