package util

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// ShowErrorDialog 显示错误对话框
func ShowErrorDialog(title string, err error, parent fyne.Window) {
	Error("%s: %v", title, err)
	dialog.ShowError(err, parent)
}

// ShowErrorDialogWithMessage 显示带自定义消息的错误对话框
func ShowErrorDialogWithMessage(title, message string, err error, parent fyne.Window) {
	Error("%s: %s - %v", title, message, err)
	fullMessage := fmt.Sprintf("%s\n\n错误详情: %v", message, err)
	dialog.ShowError(fmt.Errorf("%s", fullMessage), parent)
}

// ShowInfoDialog 显示信息对话框
func ShowInfoDialog(title, message string, parent fyne.Window) {
	Info("%s: %s", title, message)
	dialog.ShowInformation(title, message, parent)
}

// ShowConfirmDialog 显示确认对话框
func ShowConfirmDialog(title, message string, callback func(bool), parent fyne.Window) {
	dialog.ShowConfirm(title, message, callback, parent)
}

// ShowSuccessNotification 显示成功通知
func ShowSuccessNotification(message string, parent fyne.Window) {
	Info("成功: %s", message)
	dialog.ShowInformation("成功", message, parent)
}

// ShowWarningDialog 显示警告对话框
func ShowWarningDialog(title, message string, parent fyne.Window) {
	Warn("%s: %s", title, message)
	dialog.ShowInformation(title, message, parent)
}
