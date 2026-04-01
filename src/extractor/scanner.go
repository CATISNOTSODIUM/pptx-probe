package extractor

import (
	"html"
	"log"
	"os"
	"path/filepath"
	"ppt-probe/src/models"
	"regexp"
	"strconv"
	"strings"
)

const ARROW_KEYWORD = "cNvCxnSpPr"
const BOX_KEYWORD = "cNvCxnSpPr"
const BODY_KEYWORD = "txBody"
const PARAGRAPH_KEYWORD = "pPr"
const NAV_KEYWORD = "nvSpPr"
const TEXTBOX_KEYWORD = "cNvPr"

const LEAVES_TYPE = "p"

func ScanArrows(node models.Node) []models.Node {
	return Filter([]models.Node{node}, func(n models.Node) bool {
		return n.XMLName.Local == ARROW_KEYWORD
	})
}

func ScanLeaves(node models.Node) []models.Node {
	return Filter([]models.Node{node}, func(n models.Node) bool {
		return n.XMLName.Local == LEAVES_TYPE
	})
}

func GetLeafContent(node models.Node) string {
	if len(node.Nodes) == 0 {
		return string(node.Content)
	}
	var builder strings.Builder

	for _, n := range node.Nodes {
		content := GetLeafContent(n)
		if content == "" {
			continue
		}
		builder.WriteString(content)
	}
	return builder.String()
}

func ScanBody(node models.Node) []models.Node {
	return FilterWithParent([]models.Node{node}, func(n models.Node) bool {
		return n.XMLName.Local == BODY_KEYWORD
	}, nil)
}

func ExtractLevel(node models.Node) int {
	paragraphs := Filter([]models.Node{node}, func(n models.Node) bool {
		return n.XMLName.Local == PARAGRAPH_KEYWORD
	})

	if len(paragraphs) == 0 || paragraphs[0].Level == "" {
		return 0
	} else {
		parsed, err := strconv.Atoi(paragraphs[0].Level)
		if err != nil {
			log.Fatalf("Error try to parse level %s", paragraphs[0].Level)
			return 0
		} else {
			return parsed
		}
	}
}

func ParseBody(node models.Node) string {
	if node.XMLName.Local != BODY_KEYWORD {
		return ""
	}
	var builder strings.Builder
	for _, n := range node.Nodes {
		level := ExtractLevel(n)
		leaves := ScanLeaves(n)
		for _, leaf := range leaves {
			if len(leaf.Content) == 0 || leaf.Content == nil {
				continue
			}
			if level > 0 {
				tabs := strings.Repeat("\t", level)
				builder.WriteString(tabs)
			}
			builder.WriteString(GetLeafContent(leaf))
			builder.WriteString("\n")
		}
	}
	return html.UnescapeString(builder.String())
}

func ExtractArrowID(node models.Node) (string, string) {
	if node.XMLName.Local != ARROW_KEYWORD {
		return "", ""
	}
	xmlInput := string(node.Content)

	re := regexp.MustCompile(`<(?:a:)?(st|end)Cxn[^>]*id="(?P<id>\d+)"[^>]*idx="(?P<idx>\d+)"`)

	matches := re.FindAllStringSubmatch(xmlInput, -1)

	startID := ""
	endID := ""
	// Dirty implementation (to be replaced with proper XML)
	for _, match := range matches {
		// match[0] is the whole string
		// match[1] is 'st' or 'end'
		// match[2] is 'id'
		// match[3] is 'idx'
		if match[1] == "st" {
			startID = match[2]
		}
		if match[1] == "end" {
			endID = match[2]
		}
	}
	return startID, endID
}

func ExtractTextBoxID(node models.Node) string {
	if len(node.Nodes) == 0 || node.Nodes[0].XMLName.Local != NAV_KEYWORD {
		return ""
	}
	target := node.Nodes[0]
	if len(target.Nodes) == 0 || target.Nodes[0].XMLName.Local != TEXTBOX_KEYWORD {
		return ""
	}
	return target.Nodes[0].Id
}

func EnsureDirClean(dir string) error {

	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			if err := os.Remove(dir); err != nil {
				return err
			}
		}
	}
	return os.MkdirAll(dir, 0755)
}

// Entry point
// Return filename, content
func Parse(node models.Node, output string) {
	fileNameToID := make(map[string]string)
	IDToContent := make(map[string]string)
	IDToNextID := make(map[string]string)
	bodies := ScanBody(node)
	arrows := ScanArrows(node)
	err := EnsureDirClean(output)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	for _, arrow := range arrows {
		start, end := ExtractArrowID(arrow)
		IDToNextID[start] = end
	}

	for _, body := range bodies {
		id := ExtractTextBoxID(*body.Parent)
		content := FixPowerPointQuotes(ParseBody(body))
		fileName := ExtractFileName(content)

		if fileName != "" {
			fileNameToID[fileName] = id
			content = StripHeaderMarker(content)
		}
		IDToContent[id] = content
	}

	// Write the files
	for fileName, id := range fileNameToID {
		path := filepath.Join(output, strings.TrimSpace(fileName))
		err := EnsureDirClean(filepath.Dir(path))
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		// Perform DFS
		cur := id
		isVisited := make(map[string]bool)
		for {
			_, err := f.WriteString(IDToContent[cur])
			if err != nil {
				log.Fatalln(err.Error())
			}

			isVisited[cur] = true
			next, ok := IDToNextID[cur]
			if isVisited[next] || !ok {
				break
			}
			cur = next
		}
	}
}
