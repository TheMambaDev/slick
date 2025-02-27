package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  slick run       - Run the Slick app")
	fmt.Println("  slick install   - Install the Slick project")
	fmt.Println("  slick new <option> <name> - Create a new <option>")
	fmt.Println("  - Options: project, handler, model, view")

	os.Exit(0)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
	}
	cmd := os.Args[1]

	switch cmd {
	case "run":
		if _, err := os.Stat("cmd/main.go"); err != nil {
			fmt.Println("not in slick app root: cmd/main.go not found")
			os.Exit(1)
		}
		if err := exec.Command("templ", "generate").Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exec.Command("go", "run", "cmd/main.go").Run()
	case "install":
		if err := installProject(); err != nil {
			fmt.Println(err)
		}
	case "new":

		if len(os.Args) < 4 {
			usage()
		}

		// handle new command
		if err := handleNewCommand(os.Args[2:]); err != nil {
			fmt.Println(err)
		}
	}
}

func handleNewCommand(args []string) error {
	option := args[0]
	name := args[1]

	_, isSlickProject := os.Stat("cmd/main.go")

	if option != "project" && isSlickProject != nil {
		fmt.Println("not in slick project root: cmd/main.go not found")
		os.Exit(1)
	}

	switch option {

	case "project":
		if err := generateProject(name); err != nil {
			return err
		}

	case "handler":
		// TODO: check if name is a folder and if it contains a .go extension
		if err := os.WriteFile("./handler/"+name+".go", writeHandlerContent(name), os.ModePerm); err != nil {
			return err
		}

	case "model":
		// INFO: this is not yet implemented

	case "view":
		// TODO: check if name is a folder and if it contains a .go extension
		if err := os.WriteFile("./view/"+name+".templ", writeViewContent(name), os.ModePerm); err != nil {
			return err
		}

	}

	return nil
}

func generateProject(name string) error {
	fmt.Println("creating new slick project:", name)
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		return err
	}

	folders := []string{"model", "handler", "view", "cmd", "public"}

	for _, folder := range folders {
		if err := os.Mkdir(name+"/"+folder, os.ModePerm); err != nil {
			return err
		}
	}

	if err := os.WriteFile(name+"/go.mod", writeGoModContents(name), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/.air.toml", writeAirTomlContents(), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/.env", writeEnvFileContents(), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/.gitignore", writeGitignore(), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/public/app.css", []byte(""), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/cmd/main.go", writeMainContents(name), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/handler/hello.go", writeHandlerContent(name), os.ModePerm); err != nil {
		return err
	}
	if err := os.Mkdir(name+"/view/hello", os.ModePerm); err != nil {
		return err
	}
	if err := os.Mkdir(name+"/view/layout", os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/view/layout/base.templ", writeBaseLayoutContent(), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/view/hello/hello.templ", writeViewContent(name), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func installProject() error {
	start := time.Now()
	fmt.Println("installing project...")
	if err := exec.Command("go", "get", "github.com/anthdm/slick@latest").Run(); err != nil {
		return err
	}
	if err := exec.Command("go", "get", "github.com/a-h/templ").Run(); err != nil {
		return err
	}
	if err := exec.Command("templ", "generate").Run(); err != nil {
		return err
	}
	fmt.Printf("done installing project in %v\n", time.Since(start))
	return nil
}

func writeEnvFileContents() []byte {
	return []byte(`
SLICK_HTTP_LISTEN_ADDR=:3000

SLICK_SQL_DB_NAME=
SLICK_SQL_DB_USER=
SLICK_SQL_DB_PASSWORD=
SLICK_SQL_DB_HOST=
SLICK_SQL_DB_PORT=
`)
}

func writeMainContents(mod string) []byte {
	c := fmt.Sprintf(`
package main

import (
	"github.com/anthdm/slick"
	"%s/handler"
)

func main() {
	app := slick.New()
	app.Get("/", handler.HandleHelloIndex)
	log.Fatal(app.Start())
}
`, mod)
	return []byte(c)
}

func writeGoModContents(mod string) []byte {
	buf := strings.Builder{}
	buf.WriteString("module " + mod)
	buf.WriteString("\n")
	buf.WriteString("\n")
	buf.WriteString("go 1.21.0")
	return []byte(buf.String())
}

func writeAirTomlContents() []byte {
	c := `
root = "."
testdata_dir = "testdata"
tmp_dir = ".build"

[build]
  args_bin = []
  bin = "./.build/main"
  cmd = "templ generate && go build -o ./.build/main ./cmd"
  delay = 1000
  exclude_dir = ["assets", ".build", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go", "_templ.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "templ"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true

`
	return []byte(c)
}

func writeViewContent(mod string) []byte {
	c := fmt.Sprintf(`
package hello

import (
	"%s/view/layout"
)

templ Index() {
	@layout.Base() {
		<h1>hello there</h1>
	}
}
`, mod)

	return []byte(c)
}

func writeHandlerContent(mod string) []byte {
	c := fmt.Sprintf(`
package handler

import (
	"github.com/anthdm/slick"
	"%s/view/hello"
)

func HandleHelloIndex(c *slick.Context) error {
	return c.Render(hello.Index())
}
`, mod)

	return []byte(c)
}

func writeGitignore() []byte {
	c := `
bin
.build
.env
`
	return []byte(c)
}

func writeBaseLayoutContent() []byte {
	c := `
package layout

templ Base() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>Slick Application</title>

			<script src="https://unpkg.com/htmx.org@1.9.9" defer></script>
		</head>
		<body>
			{ children... }
		</body>
	</html>
}
`
	return []byte(c)
}
