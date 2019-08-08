package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)


const (
	S3_REGION = "eu-central-1"
	S3_BUCKET = "cves"
)


type Cve struct {
	CVEDataType         string `json:"CVE_data_type"`
	CVEDataFormat       string `json:"CVE_data_format"`
	CVEDataVersion      string `json:"CVE_data_version"`
	CVEDataNumberOfCVEs string `json:"CVE_data_numberOfCVEs"`
	CVEDataTimestamp    string `json:"CVE_data_timestamp"`
	CVEItems            []struct {
		Cve struct {
			DataType    string `json:"data_type"`
			DataFormat  string `json:"data_format"`
			DataVersion string `json:"data_version"`
			CVEDataMeta struct {
				ID       string `json:"ID"`
				ASSIGNER string `json:"ASSIGNER"`
			} `json:"CVE_data_meta"`
			Affects struct {
				Vendor struct {
					VendorData []struct {
						VendorName string `json:"vendor_name"`
						Product    struct {
							ProductData []struct {
								ProductName string `json:"product_name"`
								Version     struct {
									VersionData []struct {
										VersionValue    string `json:"version_value"`
										VersionAffected string `json:"version_affected"`
									} `json:"version_data"`
								} `json:"version"`
							} `json:"product_data"`
						} `json:"product"`
					} `json:"vendor_data"`
				} `json:"vendor"`
			} `json:"affects"`
			Problemtype struct {
				ProblemtypeData []struct {
					Description []struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"description"`
				} `json:"problemtype_data"`
			} `json:"problemtype"`
			References struct {
				ReferenceData []struct {
					URL       string   `json:"url"`
					Name      string   `json:"name"`
					Refsource string   `json:"refsource"`
					Tags      []string `json:"tags"`
				} `json:"reference_data"`
			} `json:"references"`
			Description struct {
				DescriptionData []struct {
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"description_data"`
			} `json:"description"`
		} `json:"cve"`
		Configurations struct {
			CVEDataVersion string `json:"CVE_data_version"`
			Nodes          []struct {
				Operator string `json:"operator"`
				CpeMatch []struct {
					Vulnerable bool   `json:"vulnerable"`
					Cpe23URI   string `json:"cpe23Uri"`
				} `json:"cpe_match"`
			} `json:"nodes"`
		} `json:"configurations"`
		Impact struct {
			BaseMetricV2 struct {
				CvssV2 struct {
					Version               string  `json:"version"`
					VectorString          string  `json:"vectorString"`
					AccessVector          string  `json:"accessVector"`
					AccessComplexity      string  `json:"accessComplexity"`
					Authentication        string  `json:"authentication"`
					ConfidentialityImpact string  `json:"confidentialityImpact"`
					IntegrityImpact       string  `json:"integrityImpact"`
					AvailabilityImpact    string  `json:"availabilityImpact"`
					BaseScore             float64 `json:"baseScore"`
				} `json:"cvssV2"`
				Severity                string  `json:"severity"`
				ExploitabilityScore     float64 `json:"exploitabilityScore"`
				ImpactScore             float64 `json:"impactScore"`
				ObtainAllPrivilege      bool    `json:"obtainAllPrivilege"`
				ObtainUserPrivilege     bool    `json:"obtainUserPrivilege"`
				ObtainOtherPrivilege    bool    `json:"obtainOtherPrivilege"`
				UserInteractionRequired bool    `json:"userInteractionRequired"`
			} `json:"baseMetricV2"`
		} `json:"impact"`
		PublishedDate    string `json:"publishedDate"`
		LastModifiedDate string `json:"lastModifiedDate"`
	} `json:"CVE_Items"`
}














func main() {

	fileUrl := "https://nvd.nist.gov/feeds/json/cve/1.0/nvdcve-1.0-recent.json.zip"

	if err := DownloadFile("cve.zip", fileUrl); err != nil {
		panic(err)
	}


	files, err := Unzip("cve.zip", ".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	t := time.Now()
	var timestamp = t.Format("2006-01-02-15-04-05")


	/*dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}*/
	var filename = timestamp + "_cves.txt"
	//var s3file = dir + "\\" + timestamp + "_cves.txt"
	fmt.Println(filename)
	f, err := os.Create(filename)
	defer f.Close()
	//cveJson
	file, _ := ioutil.ReadFile("nvdcve-1.0-recent.json")
	data := Cve{}
	//var cve Cve
	_ = json.Unmarshal([]byte(file), &data)
	for i := 0; i < len(data.CVEItems); i++ {
		var desc = data.CVEItems[i].Cve.Description.DescriptionData[0].Value
		var id = data.CVEItems[i].Cve.CVEDataMeta.ID

		if (strings.Contains(desc, "AWS") || strings.Contains(desc, "Azure") || strings.Contains(desc, "Kubernetes")){
		//fmt.Println(data.CVEItems[i].Cve.CVEDataMeta.ID)
		fmt.Println(desc)
		f.WriteString(desc)
		f.Sync()
		f.WriteString("\n")
		f.Sync()
		f.WriteString("------------------------------------------------------")
		f.Sync()
		fmt.Println(id)
		f.WriteString(id)
		f.Sync()
		f.WriteString("\n")
		f.Sync()
		f.WriteString("------------------------------------------------------")
		f.Sync()
		f.WriteString("\n")
		f.Sync()
		w := bufio.NewWriter(f)
		w.Flush()
		}
	}


	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}

	// Upload
	err = AddFileToS3(s, filename)
	fmt.Println(filename)
	if err != nil {
		log.Fatal(err)
	}


}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		/*if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}*/

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}


func AddFileToS3(s *session.Session, fileDir string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(fileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

