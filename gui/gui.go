package gui

import (
	"fix-SQ-scripts/core"
	"fix-SQ-scripts/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Start is the entry point for the GUI mode.
func Start(paths []string, debug bool) {
	if debug {
		logger.SetLogFile("debug.log")
	}

	a := app.New()
	w := a.NewWindow("File Patcher")

	resultsContainer := container.NewVBox()
	scrollableResults := container.NewScroll(resultsContainer)

	progressBar := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Processing files...")
	bottomBox := container.NewVBox(progressLabel, progressBar)

	closeButton := widget.NewButton("Close", func() {
		w.Close()
	})
	closeButton.Disable()

	go func() {
		defer func() {
			progressLabel.SetText("All tasks complete.")
			closeButton.Enable()
		}()

		filesToProcess, err := core.GetFilesToProcess(paths)
		if err != nil {
			progressLabel.SetText("Error: " + err.Error())
			return
		}

		totalFiles := len(filesToProcess)
		if totalFiles == 0 {
			progressLabel.SetText("No files found to process.")
			closeButton.Enable()
			return
		}
		availablePatches := []core.Patch{core.SQMMFixedAmount}
		resultsChan := core.ProcessPaths(filesToProcess, availablePatches)
		filesProcessed := 0

		for result := range resultsChan {
			item := createResultItem(result)
			resultsContainer.Add(item)

			filesProcessed++
			progress := float64(filesProcessed) / float64(totalFiles)
			progressBar.SetValue(progress)
			scrollableResults.ScrollToBottom()
		}
	}()

	buttonContainer := container.NewHBox(closeButton)
	bottomBox.Add(buttonContainer)

	content := container.NewBorder(nil, bottomBox, nil, nil, scrollableResults)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

// statusIcon returns an appropriate icon for the given status.
func statusIcon(status string) fyne.Resource {
	switch status {
	case "Patched":
		return theme.ConfirmIcon()
	case "Skipped":
		return theme.WarningIcon()
	case "Error":
		return theme.ErrorIcon()
	default:
		return theme.InfoIcon()
	}
}

func createResultItem(result core.PatchResult) fyne.CanvasObject {
	icon := widget.NewIcon(statusIcon(result.Status))
	filePath := widget.NewLabel(result.FilePath)
	filePath.TextStyle.Bold = true
	message := widget.NewLabel(result.Message)
	message.Wrapping = fyne.TextWrapWord

	vbox := container.NewVBox(filePath, message)
	border := container.NewBorder(nil, nil, icon, nil, vbox)
	return border
}
