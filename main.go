package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		fmt.Printf("Error al buscar archivos con el patrón: %s\n", err)
		os.Exit(1)
	}

	if len(archivos) == 0 {
		fmt.Printf("No se encontraron archivos con el patrón: %s\n", patron)
		return
	}

	lineaVieja := "double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_VALUE) / SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_SIZE);"
	lineaNueva := "double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE);"

	for _, archivo := range archivos {
		if filepath.Ext(archivo) != ".mq5" {
			fmt.Printf("El archivo '%s' no es un archivo .mq5 y será omitido.\n", archivo)
			continue
		}

		contenido, err := ioutil.ReadFile(archivo)
		if err != nil {
			fmt.Printf("Error al leer el archivo %s: %s\n", archivo, err)
			continue
		}

		contenidoString := string(contenido)

		if strings.Contains(contenidoString, lineaVieja) {
			nuevoContenido := strings.Replace(contenidoString, lineaVieja, lineaNueva, -1)
			err = ioutil.WriteFile(archivo, []byte(nuevoContenido), 0644)
			if err != nil {
				fmt.Printf("Error al escribir en el archivo %s: %s\n", archivo, err)
				continue
			}
			fmt.Printf("El archivo '%s' ha sido actualizado correctamente.\n", archivo)
		} else {
			fmt.Printf("No se encontró la línea a reemplazar en '%s'. El archivo puede que ya estuviera actualizado.\n", archivo)
		}
	}
}
