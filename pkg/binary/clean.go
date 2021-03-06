package binary

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/loqutus/artifactory-replication/pkg/artifactory"
	"github.com/loqutus/artifactory-replication/pkg/credentials"
	"github.com/loqutus/artifactory-replication/pkg/helm"
	"github.com/loqutus/artifactory-replication/pkg/s3"
)

func removeStringFromSlice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func Clean(destinationRegistry string, destinationRegistryType string, sourceRegistry string, artifactFilterProd string, creds credentials.Creds, keepDays int, helmCdnDomain string, binaryCleanPrefix string) ([]string, error) {
	log.Println("Cleaning repo " + destinationRegistry + " from files older than " + strconv.Itoa(keepDays) + " days and not in repo " + sourceRegistry + "/" + artifactFilterProd)
	var filesToRemove []string
	if destinationRegistryType != "s3" {
		return nil, errors.New("Unknown destination registry type: " + destinationRegistryType)
	}
	log.Println("artifactory.ListAllFiles " + sourceRegistry + "/" + artifactFilterProd)
	sourceFilesProd, err := artifactory.ListAllFiles(sourceRegistry, artifactFilterProd, creds.SourceUser, creds.SourcePassword)
	if err != nil {
		return nil, err
	}
	log.Println("got " + string(strconv.Itoa(len(sourceFilesProd))) + " files from artifactory repo " + artifactFilterProd)
	log.Println("s3.GetFilesModificationDate: " + destinationRegistry)
	destinationFiles, err := s3.GetFilesModificationDate(destinationRegistry)
	if err != nil {
		return nil, err
	}
	log.Println("got " + string(strconv.Itoa(len(destinationFiles))) + " files with modification date from " + destinationRegistry)
	var destinationFilesFiltered = make(map[string]*time.Time)
	for destinationFileName, destinationFileModificationDate := range destinationFiles {
		if strings.HasPrefix(destinationFileName, binaryCleanPrefix) {
			destinationFilesFiltered[destinationFileName] = destinationFileModificationDate
		}
	}
	var excludedCounter int
	// O(n*n+n)
	for _, sourceFile := range sourceFilesProd {
		for destinationFile, _ := range destinationFilesFiltered {
			//log.Println(sourceFile, destinationFile)
			if sourceFile == destinationFile {
				excludedCounter++
				//log.Println("exluding " + sourceFile + " from removal")
				delete(destinationFilesFiltered, destinationFile)
			}
		}
	}
	log.Println("Excluded " + strconv.Itoa(excludedCounter) + " from removal")
	timeKeep := time.Now().AddDate(0, 0, -keepDays)
	for fileName, modificationDate := range destinationFilesFiltered {
		if modificationDate.Before(timeKeep) {
			filesToRemove = append(filesToRemove, fileName)
		}
	}
	log.Println("removing " + strconv.Itoa(len(filesToRemove)) + " files from " + destinationRegistry)
	removeFailed, err := s3.Delete(destinationRegistry, filesToRemove)
	if err != nil {
		return nil, err
	}
	if len(removeFailed) > 0 {
		log.Println("error removing files:")
		for _, file := range removeFailed {
			log.Println(file)
		}
	}
	var filesToReindex []string
	for fileName, _ := range destinationFiles {
		filesToReindex = append(filesToReindex, fileName)
	}
	if len(filesToRemove) > 0 {
		err := helm.Reindex(filesToRemove, destinationRegistry, filesToReindex, helmCdnDomain)
		if err != nil {
			log.Println("error regenerating index.yaml")
			panic(err)
		}
	} else {
		log.Println("Haven't found anything to remove, exiting...")
		return nil, nil
	}
	return filesToRemove, nil
}
