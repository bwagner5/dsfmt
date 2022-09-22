package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	stats := map[string][]DiskStat{}
	for scanner.Scan() {
		line := strings.TrimSpace(space.ReplaceAllString(scanner.Text(), " "))
		fmt.Println(line)
		cols := strings.Split(line, " ")
		diskStat := DiskStat{
			MajorNumber:                   MustAtoi(cols[0]),
			MinorNumber:                   MustAtoi(cols[1]),
			DeviceName:                    cols[2],
			ReadsCompleted:                MustAtoi(cols[3]),
			SectorsRead:                   MustAtoi(cols[4]),
			TimeSpentReading:              MustParseDuration(fmt.Sprintf("%sms", cols[5])),
			WritesCompleted:               MustAtoi(cols[6]),
			WritesMerged:                  MustAtoi(cols[7]),
			SectorsWritten:                MustAtoi(cols[8]),
			TimeSpentWriting:              MustParseDuration(fmt.Sprintf("%sms", cols[9])),
			IOsInProgress:                 MustAtoi(cols[10]),
			TimeSpentDoingIOs:             MustParseDuration(fmt.Sprintf("%sms", cols[11])),
			WeightedTimeSpentDoingIOs:     MustParseDuration(fmt.Sprintf("%sms", cols[12])),
			DiscardsCompletedSuccessfully: MustAtoi(cols[13]),
			DiscardsMerged:                MustAtoi(cols[14]),
			SectorsDiscarded:              MustAtoi(cols[15]),
			TimeSpentDiscarding:           MustParseDuration(fmt.Sprintf("%sms", cols[16])),
		}
		stats[diskStat.DeviceName] = append(stats[diskStat.DeviceName], diskStat)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	fmt.Printf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | \n",
		"Major Number", "Minor Number", "Device Name", "Reads Completed", "Sectors Read", "Time Spent Reading",
		"Writes Completed", "Writes Merged", "Sectors Written", "Time Spent Writing", "IOs In-Progress", "Time Spent Doing IOs",
		"Weighted Time Spent Doing IOs", "Discards Completed Successfully", "Discards Merged", "Sectors Discarded", "Time Spent Discarding")
	fmt.Println("| ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- | ------- |")
	for _, stat := range stats {
		for _, dev := range stat {
			fmt.Printf("| %d | %d | %s | %d | %d | %s | %d | %d | %d | %s | %d | %s | %s | %d | %d | %d | %s | \n",
				dev.MajorNumber, dev.MinorNumber, dev.DeviceName, dev.ReadsCompleted, dev.SectorsRead, dev.TimeSpentReading,
				dev.WritesCompleted, dev.WritesMerged, dev.SectorsWritten, dev.TimeSpentWriting, dev.IOsInProgress, dev.TimeSpentDoingIOs,
				dev.WeightedTimeSpentDoingIOs, dev.DiscardsCompletedSuccessfully, dev.DiscardsMerged, dev.SectorsDiscarded, dev.TimeSpentDiscarding)
		}

	}

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
