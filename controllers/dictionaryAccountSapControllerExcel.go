package controllers

import (
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

    "github.com/tealeg/xlsx"

	"sam-api/common"
	"sam-api/models"
	"sam-api/repository"
)

const (
	CELL_ID_P_NAME int = 22
	CELL_ID_DESC_EN_120 int = 45
	CELL_ID_ACCOUNT_MIN int = 28
	CELL_ID_ACCOUNT_MAX int = 39
)

// Load account -> desciprion mappings from Excel sheet
func dictionaryAccountSap(xfn string) (dictionary *map[string]string, err error) {
	// open Excel file
	log.Printf("Opening Excel file: %s\n", xfn)
    xf, err := xlsx.OpenFile(xfn)
    if err != nil {
		return  nil, fmt.Errorf("Error opening Excel file: %s - %s", xfn, err.Error())
    }
	log.Printf("Opened Excel file: %s", xfn)

	// load AC - AN, DESC EN(120) cells
	var d = map[string]string {}
    for _, sheet := range xf.Sheets {
		log.Printf("Found Excel sheet: %s\n", sheet.Name)
		if strings.HasPrefix(sheet.Name, "IKOS_AC_Master") {
			log.Printf("Processing Excel sheet: %s\n", sheet.Name)
			i := 0
			for _, row := range sheet.Rows {
				if i > 1 {
					desc := row.Cells[CELL_ID_DESC_EN_120].String()
					var account string
					for k := CELL_ID_ACCOUNT_MIN; k < CELL_ID_ACCOUNT_MAX; k++ {
						account += row.Cells[k].String()
					}
					if account != "" {
						d[account] = desc
					}
				}
				i++
			}
		} else {
			log.Printf("Skip Excel sheet: %s\n", sheet.Name)
		}
	}

	dictionary = &d
	log.Printf("Loaded Excel file: %s record: %d", xfn, len(d))

	return
}

// Parse Excel file loaded from request
func parseExcelPayload(r *http.Request, encoding string) (dictionary *map[string]string, err error) {
	log.Printf("Decoding Excel payload")

    tmp, err := ioutil.TempFile(os.TempDir(), "sam-api-")
    if err != nil {
        return nil, fmt.Errorf("Cannot create temporary file: %s", err.Error())
    }
    defer os.Remove(tmp.Name())
	log.Println("Created temporary file: " + tmp.Name())

	data, err := common.GetPayload(r, encoding)
    if _, err = tmp.Write(data); err != nil {
        return nil, fmt.Errorf("Failed to write to temporary file: %s", err.Error())
    }
	dictionary, err = dictionaryAccountSap(tmp.Name())
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Excel file: %s", err.Error())
	}

    if err = tmp.Close(); err != nil {
        return
    }

	log.Printf("Decoded Excel payload")
	return
}

//
// Create dictionary using the payload in XLSX format
//
func DictionaryAccountSapCreateExcel(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start processing request url: %s", r.URL.Path)
	
	encoding := r.Header.Get("Content-Encoding")
	d, err := parseExcelPayload(r, encoding)
	if err != nil {
		common.DisplayAppError(w, common.DecoderExcelError, "Error in parsing Excel payload - " + err.Error(), http.StatusInternalServerError)
		return		
	}
	log.Printf("Processed Excel payload")

	user := r.Header.Get("user")
	repo, err := repository.NewDictionaryAccountSapRepository(user)
	if err != nil {
		common.DisplayAppError(w, common.RepositoryNewError, "Error while creating repository - " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer repo.Close()
	
	_, err = repo.DeleteAll()
	if err != nil {
		common.DisplayAppError(w, common.RepositoryRunError, "Error in repository delete - " + err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Dictionary cleaned up")
	
	for k, v := range *d { 	
		e := &models.DictionaryAccountSap{
			Account: k,
			Name:    v,
		}
		if err := repo.Create(e); err != nil {
			common.DisplayAppError(w, common.RepositoryRunError, "Error while creating dictionary sap account - " + err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return creation result with headers and appropriate status
	WriteResponseJson(w, http.StatusCreated, nil)

	log.Printf("Created dictionary accounts sap, status: %d", http.StatusCreated)
}

