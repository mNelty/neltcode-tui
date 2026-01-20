package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/rivo/tview"
)

const emptyRepoMessage = "Bu depo boş görünüyor.\n\nSol panelde dosya ağacı, orta panelde kod görüntüleyici, altta durum/komut barı yer alacak."

func main() {
	app := tview.NewApplication()

	rootPath, err := os.Getwd()
	if err != nil {
		rootPath = "."
	}

	treeView := tview.NewTreeView()
	codeView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[::b]Durum:[-:-:-] Başlatılıyor...")

	emptyRepo := isRepoEmpty(rootPath)
	if emptyRepo {
		placeholderRoot := tview.NewTreeNode("(boş depo)").SetColor(tview.Colors.Gray)
		placeholderRoot.AddChild(tview.NewTreeNode("Henüz dosya yok").SetColor(tview.Colors.Gray))
		treeView.SetRoot(placeholderRoot).SetCurrentNode(placeholderRoot)
		codeView.SetText(emptyRepoMessage)
		statusBar.SetText("[yellow]Boş depo tespit edildi[-]")
	} else {
		treeRoot := buildFileTree(rootPath)
		treeView.SetRoot(treeRoot).SetCurrentNode(treeRoot)
		codeView.SetText("Bir dosya seçmek için sol paneli kullanın.")
		statusBar.SetText("[green]Hazır[-]")
	}

	leftPanel := tview.NewFrame(treeView).
		SetBorders(0, 0, 1, 1, 1, 1).
		SetBorder(true).
		AddText("Dosyalar", true, tview.AlignLeft, tview.Colors.White)
	centerPanel := tview.NewFrame(codeView).
		SetBorders(0, 0, 1, 1, 1, 1).
		SetBorder(true).
		AddText("Kod", true, tview.AlignLeft, tview.Colors.White)
	bottomPanel := tview.NewFrame(statusBar).
		SetBorders(1, 1, 1, 1, 1, 1).
		SetBorder(true).
		AddText("Durum / Komut", true, tview.AlignLeft, tview.Colors.White)

	mainPanels := tview.NewFlex().
		AddItem(leftPanel, 0, 1, true).
		AddItem(centerPanel, 0, 2, false)
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainPanels, 0, 1, true).
		AddItem(bottomPanel, 3, 0, false)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isRepoEmpty(root string) bool {
	entries, err := os.ReadDir(root)
	if err != nil {
		return true
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" || name == ".gitignore" {
			continue
		}
		return false
	}
	return true
}

func buildFileTree(root string) *tview.TreeNode {
	rootNode := tview.NewTreeNode(filepath.Base(root)).SetReference(root)
	addChildren(rootNode, root)
	return rootNode
}

func addChildren(node *tview.TreeNode, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() == entries[j].IsDir() {
			return entries[i].Name() < entries[j].Name()
		}
		return entries[i].IsDir()
	})

	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		fullPath := filepath.Join(path, entry.Name())
		child := tview.NewTreeNode(entry.Name()).SetReference(fullPath)
		if entry.IsDir() {
			child.SetColor(tview.Colors.Green)
			addChildren(child, fullPath)
		} else {
			child.SetColor(tview.Colors.White)
		}
		node.AddChild(child)
	}
}
