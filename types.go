package main

import (
	"log"
	"math"
	"strconv"
	"strings"
)

type PaymentMethod = string

const (
	Unknown PaymentMethod = "Unknown"
	Salary  PaymentMethod = "Salary"
	Wage    PaymentMethod = "Wage"
)

type ProgrammingJob struct {
	Title        string
	Company      string
	Location     string
	AveragePay   int
	PayMethod    PaymentMethod
	Technologies []string
}

func (job *ProgrammingJob) ParsePayment(payment string) {
	payment = strings.Split(payment, "\n")[0]
	payment = strings.ReplaceAll(payment, ",", "")
	str := strings.ReplaceAll(payment, "$", "")
	str = strings.TrimSpace(str)

	if strings.Contains(str, "-") {
		splitString := strings.Split(str, "-")
		for i, splitted := range splitString {
			splitString[i] = strings.TrimSpace(splitted)
		}

		minPay, err := strconv.ParseFloat(splitString[0], 64)
		if err != nil {
			log.Println("Failed to parse payment minimum in string:", payment, "Error:", err)
			return
		}

		job.AveragePay = int(math.Round(minPay))

		maxPay, err := strconv.ParseFloat(splitString[1], 64)
		if err != nil {
			log.Println("Failed to parse payment maximum in string:", payment, "Error:", err)
			return
		}

		job.AveragePay = int(math.Round((minPay + maxPay) / 2.0))
	} else {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Println("Failed to parse payment string:", payment, "Error:", err)
			return
		}

		job.AveragePay = int(math.Round(value))
	}
}

func (job *ProgrammingJob) ParseDescription(description string) {
	job.Technologies = []string{}

	words := strings.Split(description, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		lower = strings.ReplaceAll(lower, ",", "")
		lower = strings.TrimSpace(lower)
		if lower == "" {
			continue
		}

		for key, value := range technologies {
			if key == lower {
				job.Technologies = append(job.Technologies, value...)
			} else {
				keywords := strings.Split(key, " ")
				for keywordIndex, keyword := range keywords {
					if i+keywordIndex >= len(words) {
						break
					}

					if words[i+keywordIndex] != keyword {
						break
					}

					if keywordIndex == len(keywords)-1 {
						job.Technologies = append(job.Technologies, value...)
					}
				}
			}
		}
	}
}

var technologies = map[string][]string{
	"actionscript": {"ActionScript"},
	"asm":          {"Assembly"},
	"assembly":     {"Assembly"},

	// "b":     {"B"},
	"bash":  {"Bash"},
	"basic": {"BASIC"},
	"bat":   {"Batch"},
	"batch": {"Batch"},

	"c":       {"C"},
	"c3":      {"C3"},
	"c++":     {"C++"},
	"c/c++":   {"C", "C++"},
	"c#":      {"C#"},
	"c-sharp": {"C#"},
	"c sharp": {"C#"},
	"carbon":  {"Carbon"},
	"curl":    {"Curl"},

	// "d":    {"D"},
	"dart": {"Dart"},

	"elixir": {"Elixir"},
	"erlang": {"Erlang"},

	"fish":    {"fish"},
	"fortran": {"Fortran"},
	"zsh":     {"ZSH"},

	"gleam": {"Gleam"},

	"go":     {"Go"},
	"golang": {"Go"},

	"gml":        {"GameMaker"},
	"gamemaker":  {"GameMaker"},
	"game maker": {"GameMaker"},

	"godot":       {"Godot"},
	"gd":          {"Godot"},
	"gdextension": {"Godot"},
	"gdscript":    {"Godot"},

	"glsl":                       {"GLSL"},
	"opengl shading language":    {"GLSL"},
	"hlsl":                       {"HLSL"},
	"high level shader language": {"HLSL"},

	"haskell": {"Haskell"},
	"haxe":    {"Haxe"},

	"java":   {"Java"},
	"kotlin": {"Kotlin"},
	"scala":  {"Scala"},

	"js":         {"JavaScript"},
	"javascript": {"JavaScript"},
	"ts":         {"TypeScript"},
	"typescript": {"TypeScript"},

	"react":      {"React"},
	"next.js":    {"Next.js"},
	"next":       {"Next.js"},
	"vue.js":     {"Vue"},
	"vue":        {"Vue"},
	"angular":    {"Angular"},
	"angular.js": {"Angular"},
	"vite.js":    {"Vite"},
	"vite":       {"Vite"},
	"astro.js":   {"Astro"},
	"astro":      {"Astro"},

	"julia": {"Julia"},

	"lisp":   {"Lisp"},
	"lua":    {"Lua"},
	"luajit": {"Lua"},
	"luau":   {"LuaU"},

	"matlab": {"MATLAB"},
	"mojo":   {"Mojo"},

	"nim":    {"Nim"},
	"neko":   {"NekoVM"},
	"nekovm": {"NekoVM"},
	"nix":    {"Nix"},

	".net":   {".NET"},
	"dotnet": {".NET"},

	"npm":  {"NPM"},
	"pnpm": {"PNPM"},
	"yarn": {"Yarn"},
	"bun":  {"Bun"},

	"objc":        {"Objective-C"},
	"objective-c": {"Objective-C"},
	"objective c": {"Objective-C"},

	"swift":    {"Swift"},
	"swiftui":  {"SwiftUI"},
	"swift ui": {"SwiftUI"},
	"swift-ui": {"SwiftUI"},

	"ocaml":  {"OCaml"},
	"opengl": {"OpenGL"},
	"opencl": {"OpenCL"},

	"gles":      {"OpenGL ES"},
	"opengl es": {"OpenGL ES"},

	"d3d":     {"Direct3D"},
	"directx": {"DirectX"},
	"dx11":    {"DirectX11"},
	"dx12":    {"DirectX12"},

	"pascal": {"PASCAL"},
	"php":    {"PHP"},
	"cobalt": {"Cobalt"},
	"python": {"Python"},

	"powershell": {"Powershell"},
	"kali":       {"Kali Linux"},
	"linux":      {"Linux"},
	"git":        {"Git"},
	"mercurial":  {"Mercurial"},

	"r":             {"R"},
	"ruby":          {"Ruby"},
	"rubyonrails":   {"Ruby On Rails"},
	"ruby on rails": {"Ruby On Rails"},

	"rust":     {"Rust"},
	"rustlang": {"Rust"},
	"cargo":    {"Rust"},

	"electron":     {"Electron"},
	"react native": {"React Native"},

	"jdk": {"Java"},

	"unity":   {"Unity"},
	"unity3d": {"Unity"},

	"unreal": {"Unreal Engine"},

	// "v":     {"V"},
	"vlang": {"V"},

	"vim":        {"Vim"},
	"vim script": {"Vim Script"},

	"wasm":         {"Web Assembly"},
	"web assembly": {"Web Assembly"},

	"zig":     {"Zig"},
	"ziglang": {"Zig"},

	"vb":           {"VisualBasic"},
	"visual basic": {"VisualBasic"},
	"vb.net":       {"VisualBasic"},
	"visualbasic":  {"VisualBasic"},

	"sql":    {"SQL"},
	"sqlite": {"SQLite"},

	"html": {"HTML"},
	"xml":  {"XML"},
	"css":  {"CSS"},

	"jsx": {"JSX"},
	"tsx": {"TSX"},
}
