package main

import (
	"flag"
	"fmt"
	"os"

	"fix-SQ-scripts/core"
	"fix-SQ-scripts/gui"
	"fix-SQ-scripts/logger"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

func mostrarAyuda() {
	fmt.Printf("%sUso:%s patcher %s<patrón>%s\n\n", ColorYellow, ColorReset, ColorCyan, ColorReset)
	fmt.Println("Corrige una línea de código específica en uno o más archivos que coincidan con el patrón.")
	fmt.Println()
	fmt.Printf("%sArgumentos:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  %s<patrón>%s   El patrón para encontrar los archivos (ej: \"*.mq5\", \"ruta/al/archivo.txt\").\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--gui%s      Lanzar en modo GUI.\n", ColorCyan, ColorReset)
	fmt.Printf("  %s--debug%s    Habilitar logging de depuración.\n", ColorCyan, ColorReset)
	fmt.Println()
	fmt.Printf("%sEjemplo:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  patcher.exe \"*.mq5\"\n")
}

func main() {
	guiFlag := flag.Bool("gui", false, "Lanzar en modo GUI")
	debugFlag := flag.Bool("debug", false, "Habilitar logging de depuración")
	flag.Parse()

	if *debugFlag {
		logger.SetLogFile("debug.log")
		logger.Logger.Println("Application started with args:", os.Args)
	}

	if *guiFlag {
		gui.Start(flag.Args(), *debugFlag)
	} else {
		if len(flag.Args()) < 1 {
			mostrarAyuda()
			os.Exit(0)
		}

		filesToProcess, err := core.GetFilesToProcess(flag.Args())
		if err != nil {
			fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

		availablePatches := []core.Patch{core.SQMMFixedAmount, core.LotSizeLogging}
		resultsChan := core.ProcessPaths(filesToProcess, availablePatches)

		for result := range resultsChan {
			var color string
			switch result.Status {
			case "Parcheado":
				color = ColorGreen
			case "Omitido":
				color = ColorYellow
			case "Error":
				color = ColorRed
			default:
				color = ColorReset
			}
			fmt.Printf("%s[%s]%s %s: %s\n", color, result.Status, ColorReset, result.FilePath, result.Message)
		}
	}
}
