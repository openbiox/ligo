package parse

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// SraXML convert Sra XML to json
func SraXML(xmlPaths *[]string, stdin *[]byte, outfn string, keywords *[]string, thread int) {
	if len(*xmlPaths) == 1 {
		thread = 1
	}
	if len(*stdin) > 0 {
		*xmlPaths = append(*xmlPaths, "ParseSraXMLStdin")
	}
	sem := make(chan bool, thread)

	//|os.O_APPEND
	var of *os.File
	if outfn == "" {
		of = os.Stdout
	} else {
		of, err := os.OpenFile(outfn, os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatal(err)
		}
		defer of.Close()
	}

	var err error
	for i, xmlPath := range *xmlPaths {
		sem <- true
		go func(xmlPath string, i int) {
			defer func() {
				<-sem
			}()
			var sra = SraSets{}
			if xmlPath != "ParseSraXMLStdin" {
				xmlData, err := os.ReadFile(xmlPath)
				if err != nil {
					log.Warnln(err)
				}
				err = xml.Unmarshal(xmlData, &sra)
			} else if xmlPath == "ParseSraXMLStdin" && len(*stdin) > 0 {
				err = xml.Unmarshal(*stdin, &sra)
			}
			if err != nil {
				log.Warnf("%v", err)
				return
			}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			jsonData, _ := json.MarshalIndent(sra.EXPERIMENTPACKAGE, "", "  ")
			io.Copy(of, bytes.NewBuffer(jsonData))
		}(xmlPath, i)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

// SraSets defines extracted Sra fields
type SraSets struct {
	XMLName           xml.Name        `xml:"EXPERIMENT_PACKAGE_SET"`
	EXPERIMENTPACKAGE []ExperimentPkg `xml:"EXPERIMENT_PACKAGE"`
}

type ExperimentPkg struct {
	EXPERIMENT struct {
		Accession   string `xml:"accession,attr"`
		Alias       string `xml:"alias,attr"`
		IDENTIFIERS struct {
			PRIMARYID string `xml:"PRIMARY_ID"`
		} `xml:"IDENTIFIERS"`
		TITLE    string `xml:"TITLE"`
		STUDYREF struct {
			Accession   string `xml:"accession,attr"`
			IDENTIFIERS struct {
				PRIMARYID  string `xml:"PRIMARY_ID"`
				EXTERNALID struct {
					Namespace string `xml:"namespace,attr"`
				} `xml:"EXTERNAL_ID"`
			} `xml:"IDENTIFIERS"`
		} `xml:"STUDY_REF"`
		DESIGN struct {
			DESIGNDESCRIPTION string `xml:"DESIGN_DESCRIPTION"`
			SAMPLEDESCRIPTOR  struct {
				Accession   string `xml:"accession,attr"`
				IDENTIFIERS struct {
					PRIMARYID  string `xml:"PRIMARY_ID"`
					EXTERNALID struct {
						Namespace string `xml:"namespace,attr"`
					} `xml:"EXTERNAL_ID"`
				} `xml:"IDENTIFIERS"`
			} `xml:"SAMPLE_DESCRIPTOR"`
			LIBRARYDESCRIPTOR struct {
				LIBRARYNAME      string `xml:"LIBRARY_NAME"`
				LIBRARYSTRATEGY  string `xml:"LIBRARY_STRATEGY"`
				LIBRARYSOURCE    string `xml:"LIBRARY_SOURCE"`
				LIBRARYSELECTION string `xml:"LIBRARY_SELECTION"`
				LIBRARYLAYOUT    struct {
					PAIRED string `xml:"PAIRED"`
				} `xml:"LIBRARY_LAYOUT"`
			} `xml:"LIBRARY_DESCRIPTOR"`
		} `xml:"DESIGN"`
		PLATFORM struct {
			ILLUMINA struct {
				INSTRUMENTMODEL string `xml:"INSTRUMENT_MODEL"`
			} `xml:"ILLUMINA"`
		} `xml:"PLATFORM"`
	} `xml:"EXPERIMENT"`
	SUBMISSION struct {
		LabName     string `xml:"lab_name,attr"`
		CenterName  string `xml:"center_name,attr"`
		Accession   string `xml:"accession,attr"`
		Alias       string `xml:"alias,attr"`
		IDENTIFIERS struct {
			PRIMARYID string `xml:"PRIMARY_ID"`
		} `xml:"IDENTIFIERS"`
	} `xml:"SUBMISSION"`
	Organization struct {
		Type    string `xml:"type,attr"`
		Name    string `xml:"Name"`
		Address struct {
			PostalCode  string `xml:"postal_code,attr"`
			Department  string `xml:"Department"`
			Institution string `xml:"Institution"`
			Street      string `xml:"Street"`
			City        string `xml:"City"`
			Country     string `xml:"Country"`
		} `xml:"Address"`
		Contact struct {
			Email    string `xml:"email,attr"`
			SecEmail string `xml:"sec_email,attr"`
			Address  struct {
				PostalCode  string `xml:"postal_code,attr"`
				Department  string `xml:"Department"`
				Institution string `xml:"Institution"`
				Street      string `xml:"Street"`
				City        string `xml:"City"`
				Country     string `xml:"Country"`
			} `xml:"Address"`
			Name struct {
				First string `xml:"First"`
				Last  string `xml:"Last"`
			} `xml:"Name"`
		} `xml:"Contact"`
	} `xml:"Organization"`
	STUDY struct {
		CenterName  string `xml:"center_name,attr"`
		Alias       string `xml:"alias,attr"`
		Accession   string `xml:"accession,attr"`
		IDENTIFIERS struct {
			PRIMARYID  string `xml:"PRIMARY_ID"`
			EXTERNALID struct {
				Namespace string `xml:"namespace,attr"`
				Label     string `xml:"label,attr"`
			} `xml:"EXTERNAL_ID"`
		} `xml:"IDENTIFIERS"`
		DESCRIPTOR struct {
			STUDYTITLE string `xml:"STUDY_TITLE"`
			STUDYTYPE  struct {
				ExistingStudyType string `xml:"existing_study_type,attr"`
			} `xml:"STUDY_TYPE"`
			STUDYABSTRACT     string `xml:"STUDY_ABSTRACT"`
			CENTERPROJECTNAME string `xml:"CENTER_PROJECT_NAME"`
		} `xml:"DESCRIPTOR"`
	} `xml:"STUDY"`
	SAMPLE struct {
		Alias       string `xml:"alias,attr"`
		Accession   string `xml:"accession,attr"`
		IDENTIFIERS struct {
			PRIMARYID  string `xml:"PRIMARY_ID"`
			EXTERNALID struct {
				Namespace string `xml:"namespace,attr"`
			} `xml:"EXTERNAL_ID"`
		} `xml:"IDENTIFIERS"`
		TITLE      string `xml:"TITLE"`
		SAMPLENAME struct {
			TAXONID        string `xml:"TAXON_ID"`
			SCIENTIFICNAME string `xml:"SCIENTIFIC_NAME"`
		} `xml:"SAMPLE_NAME"`
		SAMPLELINKS struct {
			SAMPLELINK struct {
				XREFLINK struct {
					DB    string `xml:"DB"`
					ID    string `xml:"ID"`
					LABEL string `xml:"LABEL"`
				} `xml:"XREF_LINK"`
			} `xml:"SAMPLE_LINK"`
		} `xml:"SAMPLE_LINKS"`
		SAMPLEATTRIBUTES struct {
			SAMPLEATTRIBUTE []struct {
				TAG   string `xml:"TAG"`
				VALUE string `xml:"VALUE"`
			} `xml:"SAMPLE_ATTRIBUTE"`
		} `xml:"SAMPLE_ATTRIBUTES"`
	} `xml:"SAMPLE"`
	Pool struct {
		Member struct {
			MemberName  string `xml:"member_name,attr"`
			Accession   string `xml:"accession,attr"`
			SampleName  string `xml:"sample_name,attr"`
			SampleTitle string `xml:"sample_title,attr"`
			Spots       string `xml:"spots,attr"`
			Bases       string `xml:"bases,attr"`
			TaxID       string `xml:"tax_id,attr"`
			Organism    string `xml:"organism,attr"`
			IDENTIFIERS struct {
				PRIMARYID  string `xml:"PRIMARY_ID"`
				EXTERNALID struct {
					Namespace string `xml:"namespace,attr"`
				} `xml:"EXTERNAL_ID"`
			} `xml:"IDENTIFIERS"`
		} `xml:"Member"`
	} `xml:"Pool"`
	RUNSET struct {
		RUN struct {
			Accession           string `xml:"accession,attr"`
			Alias               string `xml:"alias,attr"`
			TotalSpots          string `xml:"total_spots,attr"`
			TotalBases          string `xml:"total_bases,attr"`
			Size                string `xml:"size,attr"`
			LoadDone            string `xml:"load_done,attr"`
			Published           string `xml:"published,attr"`
			IsPublic            string `xml:"is_public,attr"`
			ClusterName         string `xml:"cluster_name,attr"`
			StaticDataAvailable string `xml:"static_data_available,attr"`
			IDENTIFIERS         struct {
				PRIMARYID string `xml:"PRIMARY_ID"`
			} `xml:"IDENTIFIERS"`
			EXPERIMENTREF struct {
				Accession   string `xml:"accession,attr"`
				IDENTIFIERS string `xml:"IDENTIFIERS"`
			} `xml:"EXPERIMENT_REF"`
			Pool struct {
				Member struct {
					MemberName  string `xml:"member_name,attr"`
					Accession   string `xml:"accession,attr"`
					SampleName  string `xml:"sample_name,attr"`
					SampleTitle string `xml:"sample_title,attr"`
					Spots       string `xml:"spots,attr"`
					Bases       string `xml:"bases,attr"`
					TaxID       string `xml:"tax_id,attr"`
					Organism    string `xml:"organism,attr"`
					IDENTIFIERS struct {
						PRIMARYID  string `xml:"PRIMARY_ID"`
						EXTERNALID struct {
							Namespace string `xml:"namespace,attr"`
						} `xml:"EXTERNAL_ID"`
					} `xml:"IDENTIFIERS"`
				} `xml:"Member"`
			} `xml:"Pool"`
			SRAFiles   SraFns `xml:"SRAFiles"`
			CloudFiles struct {
				CloudFile []struct {
					Filetype string `xml:"filetype,attr"`
					Provider string `xml:"provider,attr"`
					Location string `xml:"location,attr"`
				} `xml:"CloudFile"`
			} `xml:"CloudFiles"`
			Statistics struct {
				Nreads string `xml:"nreads,attr"`
				Nspots string `xml:"nspots,attr"`
				Read   []struct {
					Index   string `xml:"index,attr"`
					Count   string `xml:"count,attr"`
					Average string `xml:"average,attr"`
					Stdev   string `xml:"stdev,attr"`
				} `xml:"Read"`
			} `xml:"Statistics"`
			Databases struct {
				Database struct {
					Table struct {
						Name       string `xml:"name,attr"`
						Statistics struct {
							Source string `xml:"source,attr"`
							Rows   struct {
								Count string `xml:"count,attr"`
							} `xml:"Rows"`
							Elements struct {
								Count string `xml:"count,attr"`
							} `xml:"Elements"`
						} `xml:"Statistics"`
					} `xml:"Table"`
				} `xml:"Database"`
			} `xml:"Databases"`
			Bases struct {
				CsNative string `xml:"cs_native,attr"`
				Count    string `xml:"count,attr"`
				Base     []struct {
					Value string `xml:"value,attr"`
					Count string `xml:"count,attr"`
				} `xml:"Base"`
			} `xml:"Bases"`
		} `xml:"RUN"`
	} `xml:"RUN_SET"`
}

type SraFns struct {
	SRAFile SRAFile `xml:"SRAFile"`
}

type SRAFile struct {
	Cluster      string `xml:"cluster,attr"`
	Filename     string `xml:"filename,attr"`
	URL          string `xml:"url,attr"`
	Size         string `xml:"size,attr"`
	Date         string `xml:"date,attr"`
	Md5          string `xml:"md5,attr"`
	SemanticName string `xml:"semantic_name,attr"`
	Supertype    string `xml:"supertype,attr"`
	Sratoolkit   string `xml:"sratoolkit,attr"`
	Alternatives []struct {
		URL        string `xml:"url,attr"`
		FreeEgress string `xml:"free_egress,attr"`
		AccessType string `xml:"access_type,attr"`
		Org        string `xml:"org,attr"`
	} `xml:"Alternatives"`
}

type ExperimentPkgJSON struct {
	EXPERIMENT struct {
		Accession   string `json:"Accession"`
		Alias       string `json:"Alias"`
		IDENTIFIERS struct {
			PRIMARYID string `json:"PRIMARYID"`
		} `json:"IDENTIFIERS"`
		TITLE    string `json:"TITLE"`
		STUDYREF struct {
			Accession   string `json:"Accession"`
			IDENTIFIERS struct {
				PRIMARYID  string `json:"PRIMARYID"`
				EXTERNALID struct {
					Namespace string `json:"Namespace"`
				} `json:"EXTERNALID"`
			} `json:"IDENTIFIERS"`
		} `json:"STUDYREF"`
		DESIGN struct {
			DESIGNDESCRIPTION string `json:"DESIGNDESCRIPTION"`
			SAMPLEDESCRIPTOR  struct {
				Accession   string `json:"Accession"`
				IDENTIFIERS struct {
					PRIMARYID  string `json:"PRIMARYID"`
					EXTERNALID struct {
						Namespace string `json:"Namespace"`
					} `json:"EXTERNALID"`
				} `json:"IDENTIFIERS"`
			} `json:"SAMPLEDESCRIPTOR"`
			LIBRARYDESCRIPTOR struct {
				LIBRARYNAME      string `json:"LIBRARYNAME"`
				LIBRARYSTRATEGY  string `json:"LIBRARYSTRATEGY"`
				LIBRARYSOURCE    string `json:"LIBRARYSOURCE"`
				LIBRARYSELECTION string `json:"LIBRARYSELECTION"`
				LIBRARYLAYOUT    struct {
					PAIRED string `json:"PAIRED"`
				} `json:"LIBRARYLAYOUT"`
			} `json:"LIBRARYDESCRIPTOR"`
		} `json:"DESIGN"`
		PLATFORM struct {
			ILLUMINA struct {
				INSTRUMENTMODEL string `json:"INSTRUMENTMODEL"`
			} `json:"ILLUMINA"`
		} `json:"PLATFORM"`
	} `json:"EXPERIMENT"`
	SUBMISSION struct {
		LabName     string `json:"LabName"`
		CenterName  string `json:"CenterName"`
		Accession   string `json:"Accession"`
		Alias       string `json:"Alias"`
		IDENTIFIERS struct {
			PRIMARYID string `json:"PRIMARYID"`
		} `json:"IDENTIFIERS"`
	} `json:"SUBMISSION"`
	Organization struct {
		Type    string `json:"Type"`
		Name    string `json:"Name"`
		Address struct {
			PostalCode  string `json:"PostalCode"`
			Department  string `json:"Department"`
			Institution string `json:"Institution"`
			Street      string `json:"Street"`
			City        string `json:"City"`
			Country     string `json:"Country"`
		} `json:"Address"`
		Contact struct {
			Email    string `json:"Email"`
			SecEmail string `json:"SecEmail"`
			Address  struct {
				PostalCode  string `json:"PostalCode"`
				Department  string `json:"Department"`
				Institution string `json:"Institution"`
				Street      string `json:"Street"`
				City        string `json:"City"`
				Country     string `json:"Country"`
			} `json:"Address"`
			Name struct {
				First string `json:"First"`
				Last  string `json:"Last"`
			} `json:"Name"`
		} `json:"Contact"`
	} `json:"Organization"`
	STUDY struct {
		CenterName  string `json:"CenterName"`
		Alias       string `json:"Alias"`
		Accession   string `json:"Accession"`
		IDENTIFIERS struct {
			PRIMARYID  string `json:"PRIMARYID"`
			EXTERNALID struct {
				Namespace string `json:"Namespace"`
				Label     string `json:"Label"`
			} `json:"EXTERNALID"`
		} `json:"IDENTIFIERS"`
		DESCRIPTOR struct {
			STUDYTITLE string `json:"STUDYTITLE"`
			STUDYTYPE  struct {
				ExistingStudyType string `json:"ExistingStudyType"`
			} `json:"STUDYTYPE"`
			STUDYABSTRACT     string `json:"STUDYABSTRACT"`
			CENTERPROJECTNAME string `json:"CENTERPROJECTNAME"`
		} `json:"DESCRIPTOR"`
	} `json:"STUDY"`
	SAMPLE struct {
		Alias       string `json:"Alias"`
		Accession   string `json:"Accession"`
		IDENTIFIERS struct {
			PRIMARYID  string `json:"PRIMARYID"`
			EXTERNALID struct {
				Namespace string `json:"Namespace"`
			} `json:"EXTERNALID"`
		} `json:"IDENTIFIERS"`
		TITLE      string `json:"TITLE"`
		SAMPLENAME struct {
			TAXONID        string `json:"TAXONID"`
			SCIENTIFICNAME string `json:"SCIENTIFICNAME"`
		} `json:"SAMPLENAME"`
		SAMPLELINKS struct {
			SAMPLELINK struct {
				XREFLINK struct {
					DB    string `json:"DB"`
					ID    string `json:"ID"`
					LABEL string `json:"LABEL"`
				} `json:"XREFLINK"`
			} `json:"SAMPLELINK"`
		} `json:"SAMPLELINKS"`
		SAMPLEATTRIBUTES struct {
			SAMPLEATTRIBUTE []struct {
				TAG   string `json:"TAG"`
				VALUE string `json:"VALUE"`
			} `json:"SAMPLEATTRIBUTE"`
		} `json:"SAMPLEATTRIBUTES"`
	} `json:"SAMPLE"`
	Pool struct {
		Member struct {
			MemberName  string `json:"MemberName"`
			Accession   string `json:"Accession"`
			SampleName  string `json:"SampleName"`
			SampleTitle string `json:"SampleTitle"`
			Spots       string `json:"Spots"`
			Bases       string `json:"Bases"`
			TaxID       string `json:"TaxID"`
			Organism    string `json:"Organism"`
			IDENTIFIERS struct {
				PRIMARYID  string `json:"PRIMARYID"`
				EXTERNALID struct {
					Namespace string `json:"Namespace"`
				} `json:"EXTERNALID"`
			} `json:"IDENTIFIERS"`
		} `json:"Member"`
	} `json:"Pool"`
	RUNSET struct {
		RUN struct {
			Accession           string `json:"Accession"`
			Alias               string `json:"Alias"`
			TotalSpots          string `json:"TotalSpots"`
			TotalBases          string `json:"TotalBases"`
			Size                string `json:"Size"`
			LoadDone            string `json:"LoadDone"`
			Published           string `json:"Published"`
			IsPublic            string `json:"IsPublic"`
			ClusterName         string `json:"ClusterName"`
			StaticDataAvailable string `json:"StaticDataAvailable"`
			IDENTIFIERS         struct {
				PRIMARYID string `json:"PRIMARYID"`
			} `json:"IDENTIFIERS"`
			EXPERIMENTREF struct {
				Accession   string `json:"Accession"`
				IDENTIFIERS string `json:"IDENTIFIERS"`
			} `json:"EXPERIMENTREF"`
			Pool struct {
				Member struct {
					MemberName  string `json:"MemberName"`
					Accession   string `json:"Accession"`
					SampleName  string `json:"SampleName"`
					SampleTitle string `json:"SampleTitle"`
					Spots       string `json:"Spots"`
					Bases       string `json:"Bases"`
					TaxID       string `json:"TaxID"`
					Organism    string `json:"Organism"`
					IDENTIFIERS struct {
						PRIMARYID  string `json:"PRIMARYID"`
						EXTERNALID struct {
							Namespace string `json:"Namespace"`
						} `json:"EXTERNALID"`
					} `json:"IDENTIFIERS"`
				} `json:"Member"`
			} `json:"Pool"`
			SRAFiles struct {
				SRAFile SRAFileJSON `json:"SRAFile"`
			} `json:"SRAFiles"`
			CloudFiles struct {
				CloudFile []struct {
					Filetype string `json:"Filetype"`
					Provider string `json:"Provider"`
					Location string `json:"Location"`
				} `json:"CloudFile"`
			} `json:"CloudFiles"`
			Statistics struct {
				Nreads string `json:"Nreads"`
				Nspots string `json:"Nspots"`
				Read   []struct {
					Index   string `json:"Index"`
					Count   string `json:"Count"`
					Average string `json:"Average"`
					Stdev   string `json:"Stdev"`
				} `json:"Read"`
			} `json:"Statistics"`
			Databases struct {
				Database struct {
					Table struct {
						Name       string `json:"Name"`
						Statistics struct {
							Source string `json:"Source"`
							Rows   struct {
								Count string `json:"Count"`
							} `json:"Rows"`
							Elements struct {
								Count string `json:"Count"`
							} `json:"Elements"`
						} `json:"Statistics"`
					} `json:"Table"`
				} `json:"Database"`
			} `json:"Databases"`
			Bases struct {
				CsNative string `json:"CsNative"`
				Count    string `json:"Count"`
				Base     []struct {
					Value string `json:"Value"`
					Count string `json:"Count"`
				} `json:"Base"`
			} `json:"Bases"`
		} `json:"RUN"`
	} `json:"RUNSET"`
}

type SRAFileJSON struct {
	Cluster      string `json:"Cluster"`
	Filename     string `json:"Filename"`
	URL          string `json:"URL"`
	Size         string `json:"Size"`
	Date         string `json:"Date"`
	Md5          string `json:"Md5"`
	SemanticName string `json:"SemanticName"`
	Supertype    string `json:"Supertype"`
	Sratoolkit   string `json:"Sratoolkit"`
	Alternatives []struct {
		URL        string `json:"URL"`
		FreeEgress string `json:"FreeEgress"`
		AccessType string `json:"AccessType"`
		Org        string `json:"Org"`
	} `json:"Alternatives"`
}
