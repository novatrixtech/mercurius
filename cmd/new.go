package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new mercurius project",
	Long:  `Create a new mercurius project`,
	Run: func(cmd *cobra.Command, args []string) {
		// checking number of arguments
		if len(args) > 0 {
			fmt.Println("Too many arguments")
			os.Exit(-1)
		}
		initGoPaths()
		setApplicationPath()
		copyNewAppFiles(confValues())
	},
}

const mercuriusPath = "github.com/novatrixtech/mercurius"

var (
	// go related paths
	gopath  string
	srcRoot string

	// mercurius related paths
	mercuriusPkg *build.Package
	appPath      string
	appName      string
	basePath     string
	importPath   string
	skeletonPath string
)

// lookup and set Go related variables
func initGoPaths() {
	// lookup go path
	gopath = build.Default.GOPATH
	if gopath == "" {
		fmt.Println("Abort: GOPATH environment variable is not set. " +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
		os.Exit(-1)
	}

	// check for go executable
	var err error
	_, err = exec.LookPath("go")
	if err != nil {
		fmt.Println("Go executable not found in PATH.")
		os.Exit(-1)
	}

	//support relative path
	workingDir, _ := os.Getwd()
	goPathList := filepath.SplitList(gopath)
	for _, path := range goPathList {
		if strings.HasPrefix(strings.ToLower(workingDir), strings.ToLower(path)) {
			srcRoot = path
			break
		}

		path, _ = filepath.EvalSymlinks(path)
		if len(path) > 0 && strings.HasPrefix(strings.ToLower(workingDir), strings.ToLower(path)) {
			srcRoot = path
			break
		}
	}

	if len(srcRoot) == 0 {
		fmt.Println("Abort: could not create a Mercurius application outside of GOPATH.")
		os.Exit(-1)
	}

	// set go src path
	srcRoot = filepath.Join(srcRoot, "src")
}

func setApplicationPath() {
	var err error
	appName = terminal("What is your application name?", "")
	gitPath := terminal("What is your git or mercurial host?", "github.com")
	gitUser := terminal("What is your git or mercurial username?", "")

	//check if gitUser is not empty to put gitUser between slashes
	if gitUser != "" {
		gitUser = fmt.Sprintf("/%s", gitUser)
	}

	//build import path
	importPath = fmt.Sprintf("%s%s/%s", gitPath, gitUser, appName)

	//check if import path is valid
	if importPath == "" {
		fmt.Println("Abort: could not create a Mercurius application with empty application path.")
		os.Exit(-1)
	}

	// validate relative path, we cannot use built-in functions
	// since Go import path is valid relative path too.
	// so check basic part of the path, which is "."
	if filepath.IsAbs(importPath) || strings.HasPrefix(importPath, ".") {
		fmt.Printf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.\n", importPath)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		fmt.Printf("Abort: Import path %s already exists.\n", importPath)
	}

	mercuriusPkg, err = build.Import(mercuriusPath, "", build.FindOnly)
	if err != nil {
		fmt.Printf("Abort: Could not find Mercurius source code: %s\n", err)
	}

	appPath = filepath.Join(srcRoot, filepath.FromSlash(importPath))
	basePath = filepath.ToSlash(filepath.Dir(importPath))

	if basePath == "." {
		// we need to remove the a single '.' when
		// the app is in the $GOROOT/src directory
		basePath = ""
	} else {
		// we need to append a '/' when the app is
		// is a subdirectory such as $GOROOT/src/path/to/mercurius
		basePath += "/"
	}
	// set base project path
	skeletonPath = filepath.Join(mercuriusPkg.Dir, "skeleton")
}

func copyNewAppFiles(cfgs map[string]interface{}) {
	var err error
	err = os.MkdirAll(appPath, 0777)
	if err != nil {
		panic(err)
	}

	mustCopyDir(appPath, skeletonPath, cfgs)

	// Dotfiles are skipped by mustCopyDir, so we have to explicitly copy the .gitignore.
	gitignore := ".gitignore"
	mustCopyFile(filepath.Join(appPath, gitignore), filepath.Join(skeletonPath, gitignore))

}

func confValues() map[string]interface{} {
	// Define default value for cache types
	// See https://go-macaron.com/docs/middlewares/cache for details
	cacheMap := map[string]string{
		"memory":   "",
		"file":     "data/caches",
		"redis":    "addr=127.0.0.1:6379",
		"memcache": "127.0.0.1:11211",
	}

	// define configs
	cfgs := make(map[string]interface{})
	cfgs["AppPath"] = importPath
	cfgs["AppName"] = appName
	cfgs["DBType"] = terminal("What database do you want to use?", "mysql")
	cfgs["DBUser"] = terminal("What is your database user?", "root")
	cfgs["DBPw"] = terminal("What is your database password?", "")
	cfgs["DBName"] = terminal("What is your database name?", "")
	cfgs["DBHost"] = terminal("What is your database host?", "localhost")
	cfgs["DBPort"] = terminal("What is your database port?", "3306")
	cfgs["MaxConn"] = terminal("What is your database max connection pool?", "10")
	cfgs["IdleConn"] = terminal("What is your database idle connection pool?", "10")
	cache := terminal("What cache do you want to use?", "memory")
	cfgs["CacheType"] = cache
	cfgs["CacheCfgs"] = terminal("What is your cache server address?", cacheMap[cache])
	cfgs["Key"] = terminal("What is your oauth key (key size must be 16 or 32)?", "")
	return cfgs
}

func init() {
	RootCmd.AddCommand(newCmd)
}
