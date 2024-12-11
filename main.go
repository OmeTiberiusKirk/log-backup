package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/go-co-op/gocron/v2"
)

const PATH_OF_LOGS = "C:/Apache24/logs/"
const NAME_OF_ERR_LOG = "error"
const NAME_OF_REQUEST_LOG = "ssl_request"
const PATH_OF_ERR_LOG = PATH_OF_LOGS + NAME_OF_ERR_LOG + ".log"
const PATH_OF_REQUEST_LOG = PATH_OF_LOGS + NAME_OF_REQUEST_LOG + ".log"
const DIR_FORMAT = "%[1]d/%[2]d/%[3]d"

func main() {
	fmt.Println("Schedule logs backup.")

	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	// handleLogs()

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.CronJob(
			"0 0 * * *",
			false,
		),
		gocron.NewTask(
			func() {
				// do things
				handleLogs()
			},
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	// case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}

func handleLogs() {
	_, err := checkExistence(PATH_OF_ERR_LOG)
	check(err)
	_, err = checkExistence(PATH_OF_REQUEST_LOG)
	check(err)

	cmd := exec.Command("httpd", "-k", "stop")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("httpd stopped.")
	}

	createZip()
	createDirectories()
	moveZipFile()

	cmd = exec.Command("httpd", "-k", "start")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("httpd started.")
	}
}

func checkExistence(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	// if os.IsNotExist(err) {
	// 	fmt.Print(err)
	// 	return false, nil
	// }
	return false, err
}

func createDirectories() {
	y, m, d := getDate()

	err := os.MkdirAll(fmt.Sprintf(PATH_OF_LOGS+DIR_FORMAT, y, m, d), 0755)
	check(err)
}

func getDate() (int, int, int) {
	t := time.Now()
	return t.Year(), int(t.Month()), t.Day()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func moveZipFile() {
	y, m, d := getDate()
	newLocation := fmt.Sprintf(PATH_OF_LOGS+DIR_FORMAT+"/logs.zip", y, m, d)
	err := os.Rename(PATH_OF_LOGS+"logs.zip", newLocation)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(PATH_OF_ERR_LOG)
	os.Remove(PATH_OF_REQUEST_LOG)

	fmt.Println("move zip file.")
}

func createZip() {
	fmt.Println("create zip.")

	archive, err := os.Create(PATH_OF_LOGS + "logs.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	//Create a new zip writer
	zipWriter := zip.NewWriter(archive)

	file, err := os.Open(PATH_OF_ERR_LOG)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println("add error log to zip..")
	writerFile, err := zipWriter.Create("error.log")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(writerFile, file); err != nil {
		panic(err)
	}

	file, err = os.Open(PATH_OF_REQUEST_LOG)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Println("add request log to zip..")
	writerFile, err = zipWriter.Create("ssl_request.log")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(writerFile, file); err != nil {
		panic(err)
	}

	zipWriter.Close()
}
