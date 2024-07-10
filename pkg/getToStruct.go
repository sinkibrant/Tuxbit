package pkg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aquasecurity/table"
)

var apiResponseProcessed []byte

// requestGet executa o Get na api e retorna o conteúdo do body
func requestGet(url string) []byte {

	apiResponse, err := http.Get(url)
	if err != nil {
		log.Fatal("\nerror: ", err)
	}

	apiResponseProcessed, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		log.Fatal("\nerror: ", err)
	}

	defer apiResponse.Body.Close()

	return apiResponseProcessed
}

func tableError(apiResponseProcessed *[]byte, errUnmarshal any) string {
	tableError := table.New(os.Stdout)
	defer tableError.Render()

	tableError.SetHeaderStyle(1)
	tableError.SetLineStyle(31)
	tableError.SetPadding(2)

	tableError.AddRow("ERROR", string(*apiResponseProcessed))

	if fmt.Sprint(errUnmarshal)[0:17] != "invalid character" {
		tableError.AddRow("error", fmt.Sprint(errUnmarshal))
		tableError.AddRow("error", "Contact the development team to report the error: https://github.com/xxx/xxx")
		return ""
	}

	return ""
}
