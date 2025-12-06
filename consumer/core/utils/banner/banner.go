package banner

import "github.com/fatih/color"

func Banner() {
	color.New(color.FgRed, color.Bold).Println(" ____               __      ____                    __         ")
	color.New(color.FgRed, color.Bold).Println("/\\  _`\\            /\\ \\    /\\  _`\\                 /\\ \\        ")
	color.New(color.FgRed, color.Bold).Println("\\ \\ \\L\\ \\     __   \\_\\ \\   \\ \\ \\L\\ \\    ___     ___\\ \\ \\/'\\    ")
	color.New(color.FgRed, color.Bold).Println(" \\ \\ ,  /   /'__`\\ /'_` \\   \\ \\ ,  /   / __`\\  /'___\\ \\ , <    ")
	color.New(color.FgRed, color.Bold).Println("  \\ \\ \\\\ \\ /\\  __//\\ \\L\\ \\   \\ \\ \\\\ \\ /\\ \\L\\ \\/\\ \\__/\\ \\ \\\\`\\  ")
	color.New(color.FgRed, color.Bold).Println("   \\ \\_\\ \\_\\ \\____\\ \\___,_\\   \\ \\_\\ \\_\\ \\____/\\ \\____\\\\ \\_\\ \\_\\")
	color.New(color.FgRed, color.Bold).Println("    \\/_/\\/ /\\/____/\\/__,_ /    \\/_/\\/ /\\/___/  \\/____/ \\/_/\\/_/")
	color.New(color.FgRed, color.Bold).Println("                                                               ")
	color.New(color.FgRed, color.Bold).Println("                                                               ")
	color.New(color.FgCyan, color.Bold).Printf("Node: ")
	color.New(color.FgGreen, color.Bold).Println("Consumer")
}
