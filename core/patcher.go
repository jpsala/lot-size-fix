package core

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// PatchResult holds the result of a patching operation for a single file.
type PatchResult struct {
	FilePath string
	Status   string // e.g., "Patched", "Skipped", "Error"
	Message  string
}

// Patch defines a single patching operation.
type Patch struct {
	Name        string
	Description string
	Apply       func(string) (string, error)
}

// GetFilesToProcess expands the given paths into a list of individual files to be processed.
func GetFilesToProcess(paths []string) ([]string, error) {
	var filesToProcess []string
	processed := make(map[string]bool)

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern '%s': %v", path, err)
		}

		for _, match := range matches {
			info, statErr := os.Stat(match)
			if statErr != nil {
				continue // Ignore if we can't stat it
			}

			if info.IsDir() {
				// If glob returns a directory, walk it for .mq5 files
				filepath.Walk(match, func(walkPath string, walkInfo os.FileInfo, walkErr error) error {
					if walkErr == nil && !walkInfo.IsDir() && filepath.Ext(walkPath) == ".mq5" {
						if !processed[walkPath] {
							filesToProcess = append(filesToProcess, walkPath)
							processed[walkPath] = true
						}
					}
					return nil
				})
			} else {
				// It's a file, check the extension
				if filepath.Ext(match) == ".mq5" {
					if !processed[match] {
						filesToProcess = append(filesToProcess, match)
						processed[match] = true
					}
				}
			}
		}
	}
	return filesToProcess, nil
}

// SQMMFixedAmount is a patch that fixes the lot size calculation for SQ-translated EAs.
var SQMMFixedAmount = Patch{
	Name:        "SQMMFixedAmount",
	Description: "Replaces the PointValue-based lot size calculation with a more reliable OrderCalcProfit method.",
	Apply: func(content string) (string, error) {
		originalContent := content
		var changes []string

		lineaNueva := `// --- FIX START ---
	// Calculate profit/loss for a 1-lot trade to determine the exact drawdown
	double oneLotSLDrawdown;
	if(!OrderCalcProfit(isLongOrder(orderType) ? ORDER_TYPE_BUY : ORDER_TYPE_SELL, correctedSymbol, 1.0, openPrice, sl, oneLotSLDrawdown)) {
		Print("OrderCalcProfit failed. Error: ", GetLastError());
		return 0;
	}
	oneLotSLDrawdown = MathAbs(oneLotSLDrawdown);
	Print(StringFormat("Money to risk: %.2f, One Lot SL Drawdown: %.2f", RiskedMoney, oneLotSLDrawdown));
	// --- FIX END ---`

		rePointValue := regexp.MustCompile(`double\s+PointValue\s*=\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_VALUE\s*\)\s*/\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_SIZE\s*\);`)
		reMagicNumber := regexp.MustCompile(`(input\s+int\s+MagicNumber\s*=\s*)\d+;`)
		reVerboseWithPointValue := regexp.MustCompile(`Verbose\s*\("Money to risk:.*?, DoubleToString\(PointValue\)\);`)
		reUseSQTickSize := regexp.MustCompile(`input\s+bool\s+UseSQTickSize\s*=\s*false;`)
		reOldDrawdown := regexp.MustCompile(`//Maximum drawdown of this order if we buy 1 lot\s*double\s+oneLotSLDrawdown\s*=\s*PointValue\s*\*\s*MathAbs\s*\(\s*openPrice\s*-\s*sl\s*\);`)

		if rePointValue.MatchString(content) {
			content = rePointValue.ReplaceAllString(content, "")
			content = reOldDrawdown.ReplaceAllString(content, lineaNueva)
			content = reVerboseWithPointValue.ReplaceAllString(content, `// $0`)
			changes = append(changes, "PointValue calculation replaced with OrderCalcProfit.")
		}

		if reUseSQTickSize.MatchString(content) {
			content = reUseSQTickSize.ReplaceAllString(content, "input bool UseSQTickSize = true;")
			changes = append(changes, "Set UseSQTickSize to true.")
		}

		var randomNumber int
		if reMagicNumber.MatchString(content) {
			randomNumber = rand.Intn(899999) + 100000
			replacementStr := fmt.Sprintf(`${1}%d; // Patched on %s`, randomNumber, time.Now().Format("2006-01-02"))
			content = reMagicNumber.ReplaceAllString(content, replacementStr)
			changes = append(changes, fmt.Sprintf("MagicNumber updated to %d.", randomNumber))
		}

		if originalContent != content {
			return content, nil
		}

		return content, fmt.Errorf("no changes applied")
	},
}

// ProcessPaths finds, reads, and patches .mq5 files based on the provided paths.
// It returns a channel of PatchResult to communicate the outcome of each operation.
func ProcessPaths(filesToProcess []string, patches []Patch) <-chan PatchResult {
	results := make(chan PatchResult)

	go func() {
		defer close(results)
		rand.Seed(time.Now().UnixNano())

		if len(filesToProcess) == 0 {
			results <- PatchResult{
				FilePath: "N/A",
				Status:   "Omitido",
				Message:  "No se encontraron archivos .mq5 para procesar.",
			}
			return
		}

		for _, archivo := range filesToProcess {
			ext := filepath.Ext(archivo)
			base := strings.TrimSuffix(archivo, ext)
			globPattern := fmt.Sprintf("%s-*%s", base, ext)
			matches, err := filepath.Glob(globPattern)
			if err != nil {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Error",
					Message:  fmt.Sprintf("Error checking for patched files: %v", err),
				}
				continue
			}
			if len(matches) > 0 {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Omitido",
					Message:  "An already patched version of this file exists, skipping.",
				}
				continue
			}

			contenido, err := ioutil.ReadFile(archivo)
			if err != nil {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Error",
					Message:  fmt.Sprintf("Error al leer el archivo: %v", err),
				}
				continue
			}

			contenidoString := string(contenido)

			if strings.Contains(contenidoString, "// Patched on") {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Omitido",
					Message:  "File already patched, skipping.",
				}
				continue
			}

			if strings.Contains(contenidoString, `Print("Lot size for ", _Symbol, " is ", DoubleToString(lotSize, 2));`) {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Omitido",
					Message:  "El archivo ya ha sido parcheado anteriormente.",
				}
				continue
			}

			originalContenido := contenidoString
			var appliedPatches []string
			currentContent := contenidoString

			for _, patch := range patches {
				newContent, err := patch.Apply(currentContent)
				if err == nil && newContent != currentContent {
					appliedPatches = append(appliedPatches, patch.Name)
					currentContent = newContent
				}
			}

			if originalContenido != currentContent {
				var newFilePath string
				reMagicNumber := regexp.MustCompile(`(input\s+int\s+MagicNumber\s*=\s*)\d+;`)
				if reMagicNumber.MatchString(currentContent) {
					matches := reMagicNumber.FindStringSubmatch(currentContent)
					if len(matches) > 1 {
						// This is a bit of a hack to get the new random number, assuming the patch set it.
						// A better approach would be for the patch to return the new number.
						var randomNumber int
						fmt.Sscanf(matches[0], "input int MagicNumber = %d;", &randomNumber)

						ext := filepath.Ext(archivo)
						base := archivo[:len(archivo)-len(ext)]
						newFilePath = fmt.Sprintf("%s-%s%s", base, strings.Split(matches[0], " ")[4], ext)
					} else {
						newFilePath = archivo
					}
				} else {
					newFilePath = archivo
				}

				err = ioutil.WriteFile(newFilePath, []byte(currentContent), 0644)
				if err != nil {
					results <- PatchResult{
						FilePath: archivo,
						Status:   "Error",
						Message:  fmt.Sprintf("Error al escribir en el archivo nuevo: %v", err),
					}
					continue
				}
				results <- PatchResult{
					FilePath: newFilePath,
					Status:   "Parcheado",
					Message:  fmt.Sprintf("File patched successfully with: %s", strings.Join(appliedPatches, ", ")),
				}
			} else {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Omitido",
					Message:  "No patches were applicable. The file might already be patched or doesn't need this fix.",
				}
			}
		}
	}()

	return results
}
