package gui

import (
	"fix-SQ-scripts/core"
	"fix-SQ-scripts/logger"
	"sync"

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

	results := make([]core.PatchResult, 0)
	var resultsMutex sync.Mutex

	list := widget.NewList(
		func() int {
			resultsMutex.Lock()
			defer resultsMutex.Unlock()
			return len(results)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(nil),
				container.NewVBox(
					widget.NewLabel("File Path"),
					widget.NewLabel("Message"),
				),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			resultsMutex.Lock()
			result := results[i]
			resultsMutex.Unlock()

			hbox := o.(*fyne.Container)
			icon := hbox.Objects[0].(*widget.Icon)
			vbox := hbox.Objects[1].(*fyne.Container)
			filePath := vbox.Objects[0].(*widget.Label)
			message := vbox.Objects[1].(*widget.Label)

			icon.SetResource(statusIcon(result.Status))
			filePath.SetText(result.FilePath)
			filePath.TextStyle.Bold = true
			message.SetText(result.Message)
			message.Wrapping = fyne.TextWrapWord
		},
	)

	progressBar := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Processing files...")
	bottomBox := container.NewVBox(progressLabel, progressBar)

	startButton := widget.NewButton("Start", nil)
	closeButton := widget.NewButton("Close", func() {
		w.Close()
	})
	closeButton.Disable()

	startButton.OnTapped = func() {
		startButton.Disable()
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
				return
			}
			resultsChan := core.ProcessPaths(filesToProcess)
			filesProcessed := 0

			for result := range resultsChan {
				resultsMutex.Lock()
				results = append(results, result)
				resultsMutex.Unlock()
				list.Refresh()

				filesProcessed++
				progress := float64(filesProcessed) / float64(totalFiles)
				progressBar.SetValue(progress)
				list.ScrollToBottom()
			}
		}()
	}

	buttonContainer := container.NewHBox(startButton, closeButton)
	bottomBox.Add(buttonContainer)

	content := container.NewBorder(nil, bottomBox, nil, nil, list)

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
