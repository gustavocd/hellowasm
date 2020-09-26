package handler

// Original code cab be found here https://qvault.io/2020/07/01/running-go-in-the-browser-with-web-assembly-wasm/

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// CompileCodeHandler ...
func CompileCodeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	enableCors(&w)

	// Get code from params
	type parameters struct {
		Code string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// create file system location for compilation path
	usr, err := user.Current()
	if err != nil {
		respondWithError(w, 500, "Couldn't get system user")
		return
	}
	workingDir := filepath.Join(usr.HomeDir, ".wasm", uuid.New().String())
	err = os.MkdirAll(workingDir, os.ModePerm)
	if err != nil {
		respondWithError(w, 500, "Couldn't create directory for compilation")
		return
	}
	defer func() {
		err = os.RemoveAll(workingDir)
		if err != nil {
			respondWithError(w, 500, "Couldn't clean up code from compilation")
			return
		}
	}()
	f, err := os.Create(filepath.Join(workingDir, "main.go"))
	if err != nil {
		respondWithError(w, 500, "Couldn't create code file for compilation")
		return
	}
	defer f.Close()
	dat := []byte(params.Code)
	_, err = f.Write(dat)
	if err != nil {
		respondWithError(w, 500, "Couldn't write code to file for compilation")
		return
	}

	// compile the wasm
	const outputBinary = "main.wasm"
	os.Setenv("GOOS", "js")
	os.Setenv("GOARCH", "wasm")
	os.Setenv("GO111MODULE", "off")
	cmd := exec.Command("go", "build", "-o", outputBinary)
	cmd.Dir = workingDir
	stderr, err := cmd.StderrPipe()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	stdErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	stdErrString := string(stdErr)
	if stdErrString != "" {
		parts := strings.Split(stdErrString, workingDir)
		if len(parts) < 2 {
			respondWithError(w, 500, stdErrString)
			return
		}
		respondWithError(w, 400, parts[1])
		return
	}
	if err := cmd.Wait(); err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// write wasm binary to response
	dat, err = ioutil.ReadFile(filepath.Join(workingDir, outputBinary))
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	w.Write(dat)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func respondWithError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Write([]byte(err))
}
