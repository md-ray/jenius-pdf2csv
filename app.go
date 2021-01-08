// app.go

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type App struct {
	pdfFile string
	tmpFile string
}

func (a *App) Initialize(pdfFile string) {
	// Check camelot exist or not
	out, err := exec.Command("camelot", "--version").Output()
	if err != nil {
		log.Fatal("Cannot find Camelot script from command line. Please install it first")
	}
	if !strings.Contains(string(out), "version") {
		log.Fatal("Cannot find Camelot script from command line. Please install it first")
	}

	// Check file exist or not
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		log.Fatal("Cannot find specified PDF file")
	}

	a.pdfFile = pdfFile
	a.tmpFile = strconv.FormatInt(time.Now().UnixNano(), 10)

}

func (a *App) Run() {
	a.generateCSVFiles(a.pdfFile, a.tmpFile)
}

func (a *App) generateCSVFiles(pdfFile string, tmpFile string) {
	pathTmp := ".\\tmp\\"

	// Generate the first file
	cmdparams := fmt.Sprintf("-f csv -o %s.csv stream %s", pathTmp+tmpFile+".csv", pdfFile)
	out, err := exec.Command("camelot", "-f", "csv", "-o", pathTmp+tmpFile+".csv", "-p", "1-end", "stream", pdfFile).Output()
	if err != nil {
		log.Fatal("Failed to execute Camelot for conversion: " + cmdparams + " : " + string(out))
	}
	csvFile1 := pathTmp + tmpFile + "-page-1-table-1" + ".csv"

	// Load the first CSV, to get total page
	csvfile, err := os.Open(csvFile1)
	if err != nil {
		log.Fatalln("Couldn't open the csv file1", err)
	}

	// Parse the file
	// r := csv.NewReader(csvfile)
	r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	totalPage := 0
	for {

		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) >= 5 && strings.Contains(record[4], "1 of") {
			log.Println("total = " + record[4][5:])
			totalPage, _ = strconv.Atoi(record[4][5:])
		}

		fmt.Println(record[0] + record[1])

	}

	log.Println("Total page found = " + strconv.Itoa(totalPage))

	if totalPage == 0 {
		log.Fatal("Cannot get total page. Exiting....")
	}

	// Load and Parse CSV data into memory
	outputs := [][12]string{}
	for i := 1; i <= totalPage; i++ {
		csvFileN := ""
		if i == 1 {
			csvFileN = pathTmp + tmpFile + "-page-1-table-2" + ".csv"
		} else {
			csvFileN = fmt.Sprintf(pathTmp+tmpFile+"-page-%d-table-1"+".csv", i)
		}

		log.Println(csvFileN)

		csvfile, err := os.Open(csvFileN)
		if err != nil {
			log.Fatalln("Couldn't open the csv fileM", err)
		}

		// Parse the file
		// r := csv.NewReader(csvfile)
		r := csv.NewReader(bufio.NewReader(csvfile))

		// Iterate through the records
		counter := -1
		var lines [12]string
		isFirst := true
		for {

			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if counter == -1 {
				if strings.Contains(record[1], "Transaction ID") {
					counter = 0
				}
			} else {
				log.Println(record)

				match, _ := regexp.MatchString("^\\d\\d\\s...\\s\\d\\d\\d\\d", record[0])
				if match && !isFirst {
					counter = 0
					// log.Println(lines)
					lines[3] = strings.Replace(lines[3], ",", "", 10)
					outputs = append(outputs, lines)
				}

				if counter >= 3 {
					lines[10] = record[1]
					lines[11] = record[2]
				} else {
					lines[(counter * 4)] = string(record[0])
					lines[(counter*4)+1] = string(record[1])
					lines[(counter*4)+2] = string(record[2])
					if len(record) >= 5 {
						lines[(counter*4)+3] = string(record[4])
					} else {
						lines[(counter*4)+3] = string(record[3])
					}
				}

				counter++

				if isFirst {
					isFirst = false
				}

			}

		}
		lines[3] = strings.Replace(lines[3], ",", "", 10)
		outputs = append(outputs, lines)

		// Delete file here
		csvfile.Close()
		err2 := os.Remove(csvFileN)

		if err2 != nil {
			fmt.Println(err2)
		}
	}

	// Delete the first file as well
	os.Remove(pathTmp + tmpFile + "-page-1-table-1" + ".csv")
	fmt.Println("coba hapus = " + pathTmp + tmpFile + "-page-1-table-1" + ".csv")

	// Write final output
	file, err := os.Create(".\\tmp\\output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	for i := 0; i < len(outputs); i++ {
		writer.Write(outputs[i][0:11])
	}
	writer.Flush()

}
