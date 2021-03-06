package helm

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/loqutus/artifactory-replication/pkg/artifactory"
	"github.com/loqutus/artifactory-replication/pkg/s3"
	"github.com/loqutus/artifactory-replication/pkg/slack"
	"k8s.io/helm/pkg/repo"
)

func RegenerateIndexYaml(artifactsList []string, artifactsListProd []string, sourceRepoUrl string, destinationRepoUrl string, sourceRepo string, prodRepo string, helmCdnDomain string) error {
	log.Println("Regenerating index.yamls")
	files := make(map[string]string)
	replicatedArtifacts := append(artifactsList, artifactsListProd...)
	for _, fileName := range replicatedArtifacts {
		if strings.Contains(fileName, "/helm/") {
			s := strings.Split(fileName, "/")
			filePrefix := strings.Join(s[1:len(s)-1], "/")
			fileRepo := s[0]
			files[filePrefix] = fileRepo
		}
	}
	if len(files) == 0 {
		log.Println("No helm charts was copied, will not regenerate index.yamls")
		return nil
	}
	for filePrefix, fileRepo := range files {
		sourceFileLocalPath, err := artifactory.Download("https://"+sourceRepoUrl+"/artifactory/"+fileRepo+"/"+filePrefix+"/index.yaml", helmCdnDomain)
		if err != nil {
			err2 := slack.SendMessage(err.Error())
			if err2 != nil {
				log.Println(err)
			}
			return err
		}
		var sourceFileLocalPath2 string
		if fileRepo == sourceRepo {
			sourceFileLocalPath2, err = artifactory.Download("https://"+sourceRepoUrl+"/artifactory/"+prodRepo+"/"+filePrefix+"/index.yaml", helmCdnDomain)
			if err != nil {
				err2 := slack.SendMessage(err.Error())
				if err2 != nil {
					log.Println(err)
				}
				return err
			}
		} else if fileRepo == prodRepo {
			sourceFileLocalPath2, err = artifactory.Download("https://"+sourceRepoUrl+"/artifactory/"+sourceRepo+"/"+filePrefix+"/index.yaml", helmCdnDomain)
			if err != nil {
				err2 := slack.SendMessage(err.Error())
				if err2 != nil {
					log.Println(err)
				}
				return err
			}
		}
		sourceIndexFile, err := repo.LoadIndexFile(sourceFileLocalPath)
		if err != nil {
			err2 := slack.SendMessage(err.Error())
			if err2 != nil {
				log.Println(err)
			}
			return err
		}
		sourceIndexFile2, err := repo.LoadIndexFile(sourceFileLocalPath2)
		if err != nil {
			err2 := slack.SendMessage(err.Error())
			if err2 != nil {
				log.Println(err)
			}
			return err
		}
		sourceIndexFile.Merge(sourceIndexFile2)
		tempFile, err := ioutil.TempFile("", "index-yaml")
		if err != nil {
			err2 := slack.SendMessage(err.Error())
			if err2 != nil {
				log.Println(err)
			}
			return err
		}
		sourceIndexFile.WriteFile(tempFile.Name(), 0644)
		err = s3.Upload(destinationRepoUrl, filePrefix+"/index.yaml", tempFile.Name())
		if err != nil {
			err2 := slack.SendMessage(err.Error())
			if err2 != nil {
				log.Println(err)
			}
			return err
		}
	}
	return nil
}
