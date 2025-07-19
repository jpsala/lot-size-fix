package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

func mostrarAyuda() {
	fmt.Printf("%sModo de uso:%s emparchador %s<pattern>%s\n\n", ColorYellow, ColorReset, ColorCyan, ColorReset)
	fmt.Println("Corrige una línea de código específica en uno o más archivos que coincidan con el patrón.")
	fmt.Println()
	fmt.Printf("%sArgumentos:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  %s<pattern>%s   El patrón para encontrar los archivos (ej: \"*.mq5\", \"ruta/al/archivo.txt\").\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sEjemplo:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  emparchador.exe \"*.mq5\"\n")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--?" {
		mostrarAyuda()
		os.Exit(0)
	}

	patron := os.Args[1]
	archivos, err := filepath.Glob(patron)
	if err != nil {
		fmt.Printf("%sError al buscar archivos con el patrón: %s%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	if len(archivos) == 0 {
		fmt.Printf("%sNo se encontraron archivos con el patrón: %s%s\n", ColorYellow, patron, ColorReset)
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

	for _, archivo := range archivos {
		if filepath.Ext(archivo) != ".mq5" {
			fmt.Printf("%sOmitiendo archivo '%s' (no es .mq5).%s\n\n", ColorYellow, archivo, ColorReset)
			continue
		}

		fmt.Printf("Analizando archivo: %s%s%s\n", ColorCyan, archivo, ColorReset)

		contenido, err := ioutil.ReadFile(archivo)
		if err != nil {
			fmt.Printf("%sError al leer el archivo %s: %s%s\n", ColorRed, archivo, err, ColorReset)
			continue
		}

		contenidoString := string(contenido)

		// Use regex to find the line, ignoring whitespace variations
		re := regexp.MustCompile(`double\s+PointValue\s*=\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_VALUE\s*\)\s*/\s*SymbolInfoDouble\s*\(\s*correctedSymbol,\s*SYMBOL_TRADE_TICK_SIZE\s*\);`)

		if re.MatchString(contenidoString) {
			// Store the actual found line for logging purposes
			foundLine := re.FindString(contenidoString)

			nuevoContenido := re.ReplaceAllString(contenidoString, lineaNueva)
			err = ioutil.WriteFile(archivo, []byte(nuevoContenido), 0644)
			if err != nil {
				fmt.Printf("%sError al escribir en el archivo %s: %s%s\n", ColorRed, archivo, err, ColorReset)
				continue
			}
			fmt.Printf("%s✔ El archivo ha sido actualizado correctamente.%s\n", ColorGreen, ColorReset)
			fmt.Printf("  %s- Removido (Incorrecto): %s%s\n", ColorRed, foundLine, ColorReset)
			fmt.Printf("  %s+ Agregado (Correcto):   %s// Code block with verbose logging%s\n\n", ColorGreen, ColorYellow, ColorReset)
		} else {
			fmt.Printf("%sNo se encontró la línea a reemplazar. El archivo puede que ya estuviera actualizado.%s\n\n", ColorYellow, ColorReset)
		}
	}
}
