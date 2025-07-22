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

// GetFilesToProcess expands the given paths into a list of individual files to be processed.
func GetFilesToProcess(paths []string) ([]string, error) {
	var filesToProcess []string
	for _, path := range paths {
		// Primero, intentar tratar el path como un archivo literal
		info, err := os.Stat(path)
		if err == nil {
			if info.IsDir() {
				// Si es un directorio, buscar archivos .mq5 recursivamente
				filepath.Walk(path, func(walkPath string, walkInfo os.FileInfo, walkErr error) error {
					if walkErr == nil && !walkInfo.IsDir() && filepath.Ext(walkPath) == ".mq5" {
						filesToProcess = append(filesToProcess, walkPath)
					}
					return nil
				})
			} else if filepath.Ext(path) == ".mq5" {
				// Si es un archivo .mq5, añadirlo directamente
				filesToProcess = append(filesToProcess, path)
			}
		} else if os.IsNotExist(err) {
			// Si no existe, entonces tratarlo como un patrón de glob
			matches, globErr := filepath.Glob(path)
			if globErr != nil {
				return nil, fmt.Errorf("patrón de globbing inválido '%s': %v", path, globErr)
			}
			for _, match := range matches {
				filesToProcess = append(filesToProcess, match)
			}
		} else {
			// Otro tipo de error al hacer Stat
			return nil, fmt.Errorf("no se pudo acceder a la ruta '%s': %v", path, err)
		}
	}
	return filesToProcess, nil
}

// ProcessPaths finds, reads, and patches .mq5 files based on the provided paths.
// It returns a channel of PatchResult to communicate the outcome of each operation.
func ProcessPaths(filesToProcess []string) <-chan PatchResult {
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
		lineaNueva := `// --- FIX START ---
	// Verbose PointValue calculation for debugging
	double incorrectPointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_VALUE) / SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_SIZE);
	double correctPointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE);
	
	Print(StringFormat("PointValue (Incorrect Method): %.5f", incorrectPointValue));
	Print(StringFormat("PointValue (Correct Method):   %.5f", correctPointValue));

	double PointValue = correctPointValue; // Use the correct value for the EA's logic
	// --- FIX END ---`

		rePointValue := regexp.MustCompile(`double\s+PointValue\s*=\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_VALUE\s*\)\s*/\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_SIZE\s*\);`)
		reMagicNumber := regexp.MustCompile(`(input\s+int\s+MagicNumber\s*=\s*)\d+;`)

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
			var changes []string

			if rePointValue.MatchString(contenidoString) {
				contenidoString = rePointValue.ReplaceAllString(contenidoString, lineaNueva)
				changes = append(changes, "PointValue actualizado.")
			}

			var randomNumber int
			magicNumberChanged := false
			if reMagicNumber.MatchString(contenidoString) {
				randomNumber = rand.Intn(899999) + 100000
				replacementStr := fmt.Sprintf(`${1}%d; // Patched on %s`, randomNumber, time.Now().Format("2006-01-02"))
				contenidoString = reMagicNumber.ReplaceAllString(contenidoString, replacementStr)
				changes = append(changes, fmt.Sprintf("MagicNumber actualizado a %d.", randomNumber))
				magicNumberChanged = true
			}

			if originalContenido != contenidoString {
				var newFilePath string
				if magicNumberChanged {
					ext := filepath.Ext(archivo)
					base := archivo[:len(archivo)-len(ext)]
					newFilePath = fmt.Sprintf("%s-%d%s", base, randomNumber, ext)
				} else {
					newFilePath = archivo
				}

				err = ioutil.WriteFile(newFilePath, []byte(contenidoString), 0644)
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
					Message:  fmt.Sprintf("Archivo actualizado y renombrado: %s", changes),
				}
			} else {
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Omitido",
					Message:  "No se encontró la línea a reemplazar. El archivo ya podría estar parcheado.",
				}
			}
		}
	}()

	return results
}
