package cmd

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new web project",
	Long:  `Create a new Golang web project based on Mercurius boilerplate`,
	Run: func(cmd *cobra.Command, args []string) {
		// checking number of arguments
		if len(args) > 0 {
			printColored(fmt.Sprintf("Abort: Too many arguments. Expected: 0 and Found: %v", len(args)), color.New(color.FgHiRed).PrintlnFunc())
			os.Exit(-1)
		}
		initGoPaths()
		setApplicationPath()
		copyNewAppFiles(confValues())
		if debug {
			printColored("2 - copyNewAppFiles -> OK!", color.New(color.FgGreen).PrintlnFunc())
		}
		packageStateCheck()
		if debug {
			printColored("3 - packageStateCheck -> OK!", color.New(color.FgGreen).PrintlnFunc())
		}
		// vendorize()
		printColored(fmt.Sprintf("Congratulations. Your Application is ready at: %s ", appPath), color.New(color.FgHiGreen, color.BgBlue).PrintlnFunc())
	},
}

const (
	mercuriusPath = "github.com/novatrixtech/mercurius"
	godepPath     = "github.com/tools/godep"
	debug         = true
)

var (
	// go related paths
	gopath  string
	srcRoot string

	// mercurius related paths
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
		printColored("Abort: GOPATH environment variable is not set. "+
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.", color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}

	// check for go executable
	var err error
	_, err = exec.LookPath("go")
	if err != nil {
		printColored("Go executable not found in PATH.", color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}

	// support relative path
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
	// set base project path
	// skeletonPath = filepath.Join(mercuriusPkg.Dir, "skeleton")
	if os.Getenv("GOPATH") == "" && os.Getenv("MERCURIUSPATH") == "" {
		printColored("Abort: neither GOPATH or MERCURIUSPATH are set. You need to define them to mercurius be able to access skeleton (template) files\n", color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}

	var err error
	appName = terminal("What is your application name?", "go-myapp")
	gitPath := terminal("What is your git source host? github.com, bitbucket.org or gitlab.com?", "github.com")
	gitUser := terminal("What is your git source host's username?", "")

	// check if gitUser is not empty to put gitUser between slashes
	if gitUser != "" {
		gitUser = fmt.Sprintf("/%s", gitUser)
	}

	// build import path
	importPath = fmt.Sprintf("%s%s/%s", gitPath, gitUser, appName)

	// check if import path is valid
	if importPath == "" {
		printColored("Abort: could not create a Mercurius application with empty application path.", color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}

	// validate relative path, we cannot use built-in functions
	// since Go import path is valid relative path too.
	// so check basic part of the path, which is "."
	if filepath.IsAbs(importPath) || strings.HasPrefix(importPath, ".") {
		printColored(fmt.Sprintf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.\n", importPath), color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		printColored(fmt.Sprintf("Alert: Import path %s already exists.\n", importPath), color.New(color.FgHiYellow).PrintlnFunc())
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

	if len(os.Getenv("GOPATH")) > 0 {
		skeletonPath = filepath.Join(os.Getenv("GOPATH"), "/src/", "github.com/novatrixtech/mercurius", "skeleton")
	} else {
		skeletonPath = filepath.Join(os.Getenv("MERCURIUSPATH"), "/", "skeleton")
	}

	if debug {
		color.Set(color.FgHiMagenta)
		defer color.Unset()
		fmt.Println("1 - Your runtime is: ", runtime.GOOS)
		// if runtime.GOOS == "windows" {
		fmt.Printf(" skeletonPath: %s \n", skeletonPath)
		fmt.Printf(" appName: %s \n", appName)
		fmt.Printf(" gitPath: %s \n", gitPath)
		fmt.Printf(" gitUser: %s \n", gitUser)
		fmt.Printf(" import-Path: %s \n", importPath)
		fmt.Printf(" app-Path: %s \n", appPath)
		fmt.Printf(" base-Path: %s \n", basePath)
		fmt.Printf(" gopath: %s\n", gopath)
		fmt.Printf(" srcRoot: %s\n\n", srcRoot)
		// }
	}
}

func copyNewAppFiles(cfgs map[string]interface{}) {
	err := os.MkdirAll(appPath, 0o777)
	if err != nil {
		printColored(fmt.Sprintf("Abort: Could not generate app %s\n", err), color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
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
	cfgs["DBType"] = strings.ToLower(terminal("What SQL Database do you want to use? MySQL or Postgres?", "mysql"))
	cfgs["DBUser"] = terminal("What is your database user?", "root")
	cfgs["DBPw"] = terminal("What is your database password?", "")
	cfgs["DBName"] = terminal("What is your database name?", "")
	cfgs["DBHost"] = terminal("What is your database host?", "localhost")
	cfgs["DBPort"] = terminal("What is your database port?", "3306")
	cfgs["MaxConn"] = terminal("What is your database max connection pool?", "10")
	cfgs["IdleConn"] = terminal("What is your database idle connection pool?", "10")
	cache := strings.ToLower(terminal("What cache do you want to use? Memory, File, Redis or Memcache?", "memory"))
	cfgs["CacheType"] = cache
	cfgs["CacheCfgs"] = terminal("What is your cache server address?", cacheMap[cache])
	if cache != "memory" {
		printColored("Don't forget to adjust cache config settings at app.go after the App being built.", color.New(color.FgHiYellow).PrintlnFunc())
		fmt.Println(" ")
	}
	cfgs["Key"] = terminal("What is your oauth key (key size must be 24 or 32)?", "12345678901234567890123456789012")
	cfgs["HttpPort"] = terminal("What is your HTTP port?", "8080")
	cfgs["MongoURI"] = terminal("What is your MongoDB URI?", "mongodb://localhost:27017/myMongoDb")
	cfgs["MongoDBName"] = terminal("What is your MongoDB database name?", "myMongoDb")
	return cfgs
}

func packageStateCheck() {
	cd(appPath)
	if debug {
		printColored(fmt.Sprintf("pkg.Dir: %q", appPath), color.New(color.FgHiMagenta).PrintlnFunc())
	}
	printColored("Generating go modules.", color.New(color.FgHiMagenta).PrintlnFunc())
	cmdGoMod := exec.Command("go", "mod", "tidy")
	cmdGoMod.Run()

	pkg := getGeneratedCode()
	getDependencies()
	printColored("Get dependencies finished. Let's build you application.", color.New(color.FgHiMagenta).PrintlnFunc())

	printColored("Preparing to build your app...", color.New(color.FgHiMagenta).PrintlnFunc())
	cmd := exec.Command("go", "build")

	printColored("Building your application...", color.New(color.FgHiYellow).PrintlnFunc())

	out, _ := cmd.CombinedOutput()
	msg := string(out)
	if debug {
		printColored(fmt.Sprintf("Msg is: %s", msg), color.New(color.FgHiMagenta).PrintlnFunc())
	}
	if msg != "" {
		if runtime.GOOS == "windows" {
			cmd = exec.Command("rd", "/s", "/q", pkg.Dir)
		} else {
			cmd = exec.Command("rm", "-Rf", pkg.Dir)
		}
		err := cmd.Run()
		if err != nil {
			printColored(fmt.Sprintf("Abort: %s\n", err), color.New(color.FgHiRed).PrintlnFunc())
			os.Exit(-1)
		}
	}
}

func getGeneratedCode() *build.Package {
	if debug {
		printColored(fmt.Sprintf("importPath: %s - absolutePath: %s\n", importPath, os.Getenv("GOPATH")+"/src/"+importPath), color.New(color.FgHiMagenta).PrintlnFunc())
	}

	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		printColored(fmt.Sprintf("Abort: Could not find generated app: %s - importPath: %s - absolutePath: %s\n", err.Error(), importPath, os.Getenv("GOPATH")+"/src/"+importPath), color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}
	return pkg
}

func getDependencies() {
	cmd := exec.Command("go", "get", "./...")
	printColored("Getting all dependencies and waiting Go finishes the job...", color.New(color.FgHiYellow).PrintlnFunc())
	err := cmd.Run()
	if err != nil {
		printColored(fmt.Sprintf("Error getting dependencies. Process aborted: %v - err.Error: %s", err, err.Error()), color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}
}

// TODO: update vendorize to usenew go mod's vendorize methodology
func vendorize() {
	v := terminal("Your App is ready to go. Do you also want to vendorize it using Godep?", "y")
	if v == "y" {
		printColored("Vendorizing...", color.New(color.FgHiYellow).PrintlnFunc())
		// getDependencies()
		// getGodep()
		pkg := getGeneratedCode()

		cd(pkg.Dir)

		printColored("Executing godep...", color.New(color.FgHiYellow).PrintlnFunc())
		cmd := exec.Command("godep", "save")
		printColored("Godep executed...", color.New(color.FgHiYellow).PrintlnFunc())
		err := cmd.Run()
		if err != nil {
			printColored(fmt.Sprintf("Abort: %s\n", err), color.New(color.FgHiRed).PrintlnFunc())
			os.Exit(-1)
		}
		return
	}
	fmt.Println("vendorization skipped")
}

func cd(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		printColored(fmt.Sprintf("Abort: %s\n", err), color.New(color.FgHiRed).PrintlnFunc())
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(newCmd)
}
