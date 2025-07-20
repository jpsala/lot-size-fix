package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("patrón de globbing inválido '%s': %v", path, err)
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, fmt.Errorf("no se pudo acceder a la ruta '%s': %v", match, err)
			}
			if info.IsDir() {
				filepath.Walk(match, func(walkPath string, walkInfo os.FileInfo, walkErr error) error {
					if walkErr == nil && !walkInfo.IsDir() && filepath.Ext(walkPath) == ".mq5" {
						filesToProcess = append(filesToProcess, walkPath)
					}
					return nil
				})
			} else if filepath.Ext(match) == ".mq5" {
				filesToProcess = append(filesToProcess, match)
			}
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

		for _, archivo := range filesToProcess {
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
			re := regexp.MustCompile(`double\s+PointValue\s*=\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_VALUE\s*\)\s*/\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_SIZE\s*\);`)

			if re.MatchString(contenidoString) {
				nuevoContenido := re.ReplaceAllString(contenidoString, lineaNueva)
				err = ioutil.WriteFile(archivo, []byte(nuevoContenido), 0644)
				if err != nil {
					results <- PatchResult{
						FilePath: archivo,
						Status:   "Error",
						Message:  fmt.Sprintf("Error al escribir en el archivo: %v", err),
					}
					continue
				}
				results <- PatchResult{
					FilePath: archivo,
					Status:   "Parcheado",
					Message:  "Archivo actualizado correctamente.",
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
