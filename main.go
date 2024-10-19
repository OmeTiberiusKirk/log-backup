package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
)

const PATH_OF_LOGS = "C:/Apache24/logs/"
const NAME_OF_ERR_LOG = "error"
const NAME_OF_ACC_LOG = "access"
const PATH_OF_ERR_LOG = PATH_OF_LOGS + NAME_OF_ERR_LOG + ".log"
const PATH_OF_ACC_LOG = PATH_OF_LOGS + NAME_OF_ACC_LOG + ".log"
const DIR_FORMAT = "%[1]d/%[2]d/%[3]d"

func main() {
	fmt.Println("Scheduling log backup.")
	// create a scheduler
	s, err := gocron.NewScheduler()

	handleLogs()
	os.Exit(0)

	if err != nil {
		log.Fatal(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.CronJob(
			"* 0 * * *",
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
	_, err = checkExistence(PATH_OF_ACC_LOG)
	check(err)

	fmt.Println(PATH_OF_ERR_LOG)
	compressFile(PATH_OF_ERR_LOG)

	// cmd := exec.Command("httpd", "-k", "stop")
	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("httpd stopped.")
	// }

	// createDateDirectories()
	// moveLogs()

	// cmd = exec.Command("httpd", "-k", "start")
	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("httpd started.")
	// }
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

func createDateDirectories() {
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

func moveLogs() {
	y, m, d := getDate()
	newLocation := fmt.Sprintf(PATH_OF_LOGS+DIR_FORMAT+"/"+NAME_OF_ERR_LOG+".log", y, m, d)
	err := os.Rename(PATH_OF_ERR_LOG, newLocation)
	if err != nil {
		log.Fatal(err)
	}

	newLocation = fmt.Sprintf(PATH_OF_LOGS+DIR_FORMAT+"/"+NAME_OF_ACC_LOG+".log", y, m, d)
	err = os.Rename(PATH_OF_ACC_LOG, newLocation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("log files were moved successfully.")
}

func compressFile(path string) {
	fmt.Println("creating zip archive.")

	archive, err := os.Create(PATH_OF_LOGS + "logs.zip")
	if err != nil {
		panic(err)
		// this is to catch errors if any
	}

	defer archive.Close()
	fmt.Println("archive file created successfully")

	//Create a new zip writer
	zipWriter := zip.NewWriter(archive)
	fmt.Println("opening first file")

	f1, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer f1.Close()

	fmt.Println("adding file to archive..")
	w1, err := zipWriter.Create(path)

	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(w1, f1); err != nil {
		panic(err)
	}

	fmt.Println("closing archive")
	zipWriter.Close()
}
