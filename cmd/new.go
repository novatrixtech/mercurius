package cmd

import (
	"fmt"

	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
			fmt.Println("Too many arguments")
			os.Exit(-1)
		}
		initGoPaths()
		setApplicationPath()
		copyNewAppFiles(confValues())
		packageStateCheck()
		vendorize()
		fmt.Println("Congratulations. Your Application is ready at: ", appPath, appName)
	},
}

const (
	mercuriusPath = "github.com/novatrixtech/mercurius"
	godepPath     = "github.com/tools/godep"
)

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
	gitPath := terminal("What is your git source host? github.com, bitbucket.org or gitlab.com?", "github.com")
	gitUser := terminal("What is your git source host's username?", "")

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
		os.Exit(-1)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		fmt.Printf("Alert: Import path %s already exists.\n", importPath)
	}

	mercuriusPkg, err = build.Import(mercuriusPath, "", build.FindOnly)
	if err != nil {
		fmt.Printf("Abort: Could not find Mercurius source code: %s\n", err)
		os.Exit(-1)
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
		fmt.Printf("Abort: Could not generate app %s\n", err)
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
	cfgs["DBType"] = terminal("What SQL Database do you want to use? MySQL or PostgreSQL?", "mysql")
	cfgs["DBUser"] = terminal("What is your database user?", "root")
	cfgs["DBPw"] = terminal("What is your database password?", "")
	cfgs["DBName"] = terminal("What is your database name?", "")
	cfgs["DBHost"] = terminal("What is your database host?", "localhost")
	cfgs["DBPort"] = terminal("What is your database port?", "3306")
	cfgs["MaxConn"] = terminal("What is your database max connection pool?", "10")
	cfgs["IdleConn"] = terminal("What is your database idle connection pool?", "10")
	cache := terminal("What cache do you want to use? Memory, File, Redis or Memcache?", "memory")
	cfgs["CacheType"] = cache
	cfgs["CacheCfgs"] = terminal("What is your cache server address?", cacheMap[cache])
	if strings.ToLower(cache) != "memory" {
		fmt.Println("Don't forget to adjust cache config settings at app.go after the App being built.")
	}
	cfgs["Key"] = terminal("What is your oauth key (key size must be 24 or 32)?", "")
	cfgs["HttpPort"] = terminal("What is your HTTP port?", "8080")
	cfgs["MongoURI"] = terminal("What is your MongoDB URI?", "mongodb://localhost:27017/myMongoDb")
	cfgs["MongoDBName"] = terminal("What is your MongoDB database name?", "myMongoDb")
	return cfgs
}

func packageStateCheck() {
	pkg := getGeneratedCode()

	cd(pkg.Dir)

	cmd := exec.Command("go", "build")
	out, _ := cmd.CombinedOutput()
	msg := string(out)
	fmt.Println(msg)
	if msg != "" {
		if runtime.GOOS == "windows" {
			cmd = exec.Command("rd", "/s", "/q", pkg.Dir)
		} else {
			cmd = exec.Command("rm", "-Rf", pkg.Dir)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Abort: %s\n", err)
			os.Exit(-1)
		}
	}

}

func getGeneratedCode() *build.Package {
	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		fmt.Printf("Abort: Could not find generated app: %s\n", err)
		os.Exit(-1)
	}
	return pkg
}

func getGodep() {
	_, err := build.Import(godepPath, "", build.FindOnly)
	if err != nil {
		cmd := exec.Command("go", "get", godepPath)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Abort: %s\n", err)
			os.Exit(-1)
		}
	}
}

func getDependencies() {
	cmd := exec.Command("go", "get", "./...")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Abort: %s\n", err)
		os.Exit(-1)
	}
}

func vendorize() {
	v := terminal("Your App is ready to go. Do you also want to vendorize it using Godep?", "y")
	if v == "y" {
		fmt.Println("Vendorizing...")
		//getDependencies()
		//getGodep()
		pkg := getGeneratedCode()

		cd(pkg.Dir)

		fmt.Println("Executing godep...")
		cmd := exec.Command("godep", "save")
		fmt.Println("Godep executed...")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Abort: %s\n", err)
			os.Exit(-1)
		}
		return
	}
	fmt.Println("vendorization skipped")
}

func cd(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("Abort: %s\n", err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(newCmd)
}
