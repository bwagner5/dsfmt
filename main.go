package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Example:
// 259       0 nvme0n1 7807 1 757441 44250 3957 528 833976 12810 0 33516 43956 0 0 0 0
// 259       1 nvme0n1p1 7669 1 752537 44152 3955 520 833896 12804 0 33464 43932 0 0 0 0
// 259       2 nvme0n1p128 45 0 360 19 0 0 0 0 0 40 0 0 0 0 0

// man page:
// What:		/proc/diskstats
// Date:		February 2008
// Contact:	Jerome Marchand <jmarchan@redhat.com>
// Description:
// 		The /proc/diskstats file displays the I/O statistics
// 		of block devices. Each line contains the following 14
// 		fields:

// 		==  ===================================
// 		 1  major number
// 		 2  minor mumber
// 		 3  device name
// 		 4  reads completed successfully
// 		 5  reads merged
// 		 6  sectors read
// 		 7  time spent reading (ms)
// 		 8  writes completed
// 		 9  writes merged
// 		10  sectors written
// 		11  time spent writing (ms)
// 		12  I/Os currently in progress
// 		13  time spent doing I/Os (ms)
// 		14  weighted time spent doing I/Os (ms)
// 		==  ===================================

// 		Kernel 4.18+ appends four more fields for discard
// 		tracking putting the total at 18:

// 		==  ===================================
// 		15  discards completed successfully
// 		16  discards merged
// 		17  sectors discarded
// 		18  time spent discarding
// 		==  ===================================

// 		Kernel 5.5+ appends two more fields for flush requests:

// 		==  =====================================
// 		19  flush requests completed successfully
// 		20  time spent flushing
// 		==  =====================================

// 		For more details refer to Documentation/admin-guide/iostats.rst

var space = regexp.MustCompile(`\s+`)

type DiskStat struct {
	MajorNumber                   int
	MinorNumber                   int
	DeviceName                    string
	ReadsCompleted                int
	ReadsMerged                   int
	SectorsRead                   int
	TimeSpentReading              time.Duration
	WritesCompleted               int
	WritesMerged                  int
	SectorsWritten                int
	TimeSpentWriting              time.Duration
	IOsInProgress                 int
	TimeSpentDoingIOs             time.Duration
	WeightedTimeSpentDoingIOs     time.Duration
	DiscardsCompletedSuccessfully int
	DiscardsMerged                int
	SectorsDiscarded              int
	TimeSpentDiscarding           time.Duration
	FlushRequestsCompleted        int
	TimeSpentFlushing             time.Duration
}

func main() {
	short := flag.Bool("short", false, "Do not show Discard and Flush Stats")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	stats := map[string][]DiskStat{}
	colsLen := 0
	for scanner.Scan() {
		line := strings.TrimSpace(space.ReplaceAllString(scanner.Text(), " "))
		if line == "" {
			continue
		}
		cols := strings.Split(line, " ")
		colsLen = len(cols)
		if *short {
			colsLen = 14
		}
		diskStat := DiskStat{
			MajorNumber:               MustAtoi(cols[0]),
			MinorNumber:               MustAtoi(cols[1]),
			DeviceName:                cols[2],
			ReadsCompleted:            MustAtoi(cols[3]),
			ReadsMerged:               MustAtoi(cols[4]),
			SectorsRead:               MustAtoi(cols[5]),
			TimeSpentReading:          MustParseDuration(fmt.Sprintf("%sms", cols[6])),
			WritesCompleted:           MustAtoi(cols[7]),
			WritesMerged:              MustAtoi(cols[8]),
			SectorsWritten:            MustAtoi(cols[9]),
			TimeSpentWriting:          MustParseDuration(fmt.Sprintf("%sms", cols[10])),
			IOsInProgress:             MustAtoi(cols[11]),
			TimeSpentDoingIOs:         MustParseDuration(fmt.Sprintf("%sms", cols[12])),
			WeightedTimeSpentDoingIOs: MustParseDuration(fmt.Sprintf("%sms", cols[13])),
		}
		if len(cols) > 14 {
			diskStat.DiscardsCompletedSuccessfully = MustAtoi(cols[14])
			diskStat.DiscardsMerged = MustAtoi(cols[15])
			diskStat.SectorsDiscarded = MustAtoi(cols[16])
			diskStat.TimeSpentDiscarding = MustParseDuration(fmt.Sprintf("%sms", cols[17]))
		}
		if len(cols) > 18 {
			diskStat.FlushRequestsCompleted = MustAtoi(cols[18])
			diskStat.TimeSpentFlushing = MustParseDuration(fmt.Sprintf("%sms", cols[19]))
		}
		stats[diskStat.DeviceName] = append(stats[diskStat.DeviceName], diskStat)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Device",
		"Reads / Merged", "Sectors Read", "Read Time",
		"Writes / Merged", "Sectors Written", "Write Time",
		"IOs Now", "IOs Time", "Weighted IOs Time"}
	if colsLen > 14 {
		headers = append(headers, "Discards / Merged", "Sectors Discarded", "Discard Time")
	}
	if colsLen > 18 {
		headers = append(headers, "Flushes", "Time Flushing")
	}
	for _, stat := range stats {
		for _, dev := range stat {
			row := []string{
				fmt.Sprintf("%s %d/%d", dev.DeviceName, dev.MajorNumber, dev.MinorNumber),
				fmt.Sprintf("%d / %d", dev.ReadsCompleted, dev.ReadsMerged),
				fmt.Sprint(dev.SectorsRead),
				fmt.Sprint(dev.TimeSpentReading),
				fmt.Sprintf("%d / %d", dev.WritesCompleted, dev.WritesMerged),
				fmt.Sprint(dev.SectorsWritten),
				fmt.Sprint(dev.TimeSpentWriting),
				fmt.Sprint(dev.IOsInProgress),
				fmt.Sprint(dev.TimeSpentDoingIOs),
				fmt.Sprint(dev.WeightedTimeSpentDoingIOs),
			}
			if colsLen > 14 {
				row = append(row,
					fmt.Sprintf("%d / %d", dev.DiscardsCompletedSuccessfully, dev.DiscardsMerged),
					fmt.Sprint(dev.SectorsDiscarded),
					fmt.Sprint(dev.TimeSpentDiscarding))
			}
			if colsLen > 18 {
				row = append(row,
					fmt.Sprint(dev.FlushRequestsCompleted),
					fmt.Sprint(dev.TimeSpentFlushing))
			}
			data = append(data, row)
		}
	}
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}

func MustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("cannot convert \"%s\" from string to duration", s))
	}
	return d
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("cannot convert \"%s\" from string to int", s))
	}
	return i
}
